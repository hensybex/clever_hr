// repository/embedding.go

package repository

import (
	"context"
	"fmt"
	"log"

	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"github.com/milvus-io/milvus-sdk-go/v2/entity"
)

type ResumeMatch struct {
	ResumeID uint
	Score    float64
}

type EmbeddingRepository interface {
	CreateResumeEmbedding(resumeID uint, embedding []float32) error
	CreateVacancyEmbedding(vacancyID uint, embedding []float32) error
	GetVacancyEmbedding(vacancyID uint) ([]float32, error)
	GetResumeEmbedding(resumeID uint) ([]float32, error)
	FindMatchingResumes(vacancyEmbedding []float32, limit int) ([]ResumeMatch, error)
}

type embeddingRepositoryImpl struct {
	milvusClient client.Client
}

func NewEmbeddingRepository(milvusClient client.Client) EmbeddingRepository {
	return &embeddingRepositoryImpl{milvusClient}
}

func (repo *embeddingRepositoryImpl) CreateResumeEmbedding(resumeID uint, embedding []float32) error {
	ctx := context.Background()
	idColumn := entity.NewColumnInt64("resume_id", []int64{int64(resumeID)})
	embeddingColumn := entity.NewColumnFloatVector("embedding", 1024, [][]float32{embedding})

	_, err := repo.milvusClient.Insert(ctx, "resumes", "", idColumn, embeddingColumn)
	if err != nil {
		return fmt.Errorf("failed to insert resume embedding: %w", err)
	}

	// Flush to ensure data is persisted
	if err := repo.milvusClient.Flush(ctx, "resumes", false); err != nil {
		return fmt.Errorf("failed to flush resumes collection: %w", err)
	}

	return nil
}

func (repo *embeddingRepositoryImpl) CreateVacancyEmbedding(vacancyID uint, embedding []float32) error {
	ctx := context.Background()
	idColumn := entity.NewColumnInt64("vacancy_id", []int64{int64(vacancyID)})
	embeddingColumn := entity.NewColumnFloatVector("embedding", 1024, [][]float32{embedding})

	_, err := repo.milvusClient.Insert(ctx, "vacancies", "", idColumn, embeddingColumn)
	if err != nil {
		return fmt.Errorf("failed to insert vacancy embedding: %w", err)
	}

	// Flush to ensure data is persisted
	if err := repo.milvusClient.Flush(ctx, "vacancies", false); err != nil {
		return fmt.Errorf("failed to flush vacancies collection: %w", err)
	}

	return nil
}

func (repo *embeddingRepositoryImpl) GetResumeEmbedding(resumeID uint) ([]float32, error) {
	ctx := context.Background()
	expr := fmt.Sprintf("resume_id == %d", resumeID)
	projections := []string{"embedding"}

	res, err := repo.milvusClient.Query(ctx, "resumes", nil, expr, projections)
	if err != nil {
		return nil, fmt.Errorf("failed to query resume embedding: %w", err)
	}

	if len(res) == 0 {
		return nil, fmt.Errorf("no embedding found for resume_id %d", resumeID)
	}

	// Access the embedding field using FieldByIndex (column index is 0 in this case)
	embeddingField := res[0].(*entity.ColumnFloatVector) // Assuming embedding is the first field
	return embeddingField.Data()[0], nil
}

func (repo *embeddingRepositoryImpl) GetVacancyEmbedding(vacancyID uint) ([]float32, error) {
	ctx := context.Background()
	expr := fmt.Sprintf("vacancy_id == %d", vacancyID)
	projections := []string{"embedding"}

	log.Printf("Executing query on Milvus with expression: %s", expr)

	res, err := repo.milvusClient.Query(ctx, "vacancies", nil, expr, projections)
	if err != nil {
		log.Printf("Error during query: %v", err)
		return nil, fmt.Errorf("failed to query vacancy embedding: %w", err)
	}

	if len(res) == 0 {
		log.Printf("No embedding found for vacancy_id %d", vacancyID)
		return nil, fmt.Errorf("no embedding found for vacancy_id %d", vacancyID)
	}

	log.Printf("Query result for vacancy ID %d: %v", vacancyID, res)

	embeddingField := res[0].(*entity.ColumnFloatVector)
	log.Printf("Extracted embedding: %v", embeddingField.Data())

	return embeddingField.Data()[0], nil
}

func (repo *embeddingRepositoryImpl) FindMatchingResumes(vacancyEmbedding []float32, limit int) ([]ResumeMatch, error) {
	ctx := context.Background()

	// Convert vacancy embedding to Milvus Vector format
	queryEmbedding := entity.FloatVector(vacancyEmbedding)
	metricType := entity.L2
	topK := limit
	sp, err := entity.NewIndexFlatSearchParam() // Flat search param (for IVF index, modify as per need)
	if err != nil {
		return nil, fmt.Errorf("failed to create search param: %w", err)
	}

	results, err := repo.milvusClient.Search(
		ctx,
		"resumes",                       // Collection name
		nil,                             // Partitions, not used in this case
		"",                              // Expression
		[]string{"resume_id"},           // Fields to return
		[]entity.Vector{queryEmbedding}, // Search vectors
		"embedding",                     // Field containing embeddings
		metricType,                      // Similarity metric
		topK,                            // Top K results
		sp,                              // Search param
	)
	if err != nil {
		return nil, fmt.Errorf("failed to search for matching resumes: %w", err)
	}

	var matches []ResumeMatch
	for _, result := range results {
		// Access the resume IDs correctly
		resumeIDField := result.Fields[0].(*entity.ColumnInt64) // Access the first field which should be resume_id
		for i := 0; i < int(result.ResultCount); i++ {
			matches = append(matches, ResumeMatch{
				ResumeID: uint(resumeIDField.Data()[i]),
				Score:    float64(result.Scores[i]), // Cast score to float64
			})
		}
	}

	return matches, nil
}
