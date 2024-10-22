// internal/model/milvus_schema.go

package model

import (
	"context"
	"fmt"
	"log"

	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"github.com/milvus-io/milvus-sdk-go/v2/entity"
)

func CreateResumeCollection(milvusClient client.Client) error {
	ctx := context.Background()

	// Check if collection exists
	hasCollection, err := milvusClient.HasCollection(ctx, "resumes")
	if err != nil {
		return err
	}

	// If collection exists, skip creation
	if hasCollection {
		log.Println("Resumes collection already exists, skipping creation.")
		return nil
	}

	// Define schema
	resumeSchema := &entity.Schema{
		CollectionName: "resumes",
		Description:    "Resume embeddings",
		Fields: []*entity.Field{
			{
				Name:       "resume_id",
				DataType:   entity.FieldTypeInt64,
				PrimaryKey: true,
				AutoID:     false,
			},
			{
				Name:     "embedding",
				DataType: entity.FieldTypeFloatVector,
				TypeParams: map[string]string{
					"dim": "1024", // Dimension of your embeddings
				},
			},
		},
	}

	// Create the collection
	err = milvusClient.CreateCollection(ctx, resumeSchema, 2)
	if err != nil {
		return err
	}

	log.Println("Resume collection created successfully")

	// Create index on the 'embedding' field
	index, err := entity.NewIndexIvfFlat(entity.L2, 1024)
	if err != nil {
		return fmt.Errorf("failed to create index params: %w", err)
	}

	err = milvusClient.CreateIndex(ctx, "resumes", "embedding", index, false)
	if err != nil {
		return fmt.Errorf("failed to create index on 'resumes' collection: %w", err)
	}

	log.Println("Index created successfully on 'resumes' collection")

	// Load the collection into memory
	err = milvusClient.LoadCollection(ctx, "resumes", false)
	if err != nil {
		return fmt.Errorf("failed to load 'resumes' collection: %w", err)
	}

	log.Println("'resumes' collection loaded into memory")

	return nil
}

// CreateVacancyCollection creates the Milvus collection schema for vacancy embeddings.
func CreateVacancyCollection(milvusClient client.Client) error {
	ctx := context.Background()

	// Check if collection exists
	hasCollection, err := milvusClient.HasCollection(ctx, "vacancies")
	if err != nil {
		return err
	}

	// If collection exists, skip creation
	if hasCollection {
		log.Println("Vacancies collection already exists, skipping creation.")
		return nil
	}

	// Define schema
	vacancySchema := &entity.Schema{
		CollectionName: "vacancies",
		Description:    "Vacancy embeddings",
		Fields: []*entity.Field{
			{
				Name:       "vacancy_id",
				DataType:   entity.FieldTypeInt64,
				PrimaryKey: true,
				AutoID:     false,
			},
			{
				Name:     "embedding",
				DataType: entity.FieldTypeFloatVector,
				TypeParams: map[string]string{
					"dim": "1024", // Dimension of your embeddings
				},
			},
		},
	}

	// Create the collection
	err = milvusClient.CreateCollection(ctx, vacancySchema, 2)
	if err != nil {
		return err
	}

	log.Println("Vacancy collection created successfully")

	// Create index on the 'embedding' field
	index, err := entity.NewIndexIvfFlat(entity.L2, 1024)
	if err != nil {
		return fmt.Errorf("failed to create index params: %w", err)
	}

	err = milvusClient.CreateIndex(ctx, "vacancies", "embedding", index, false)
	if err != nil {
		return fmt.Errorf("failed to create index on 'vacancies' collection: %w", err)
	}

	log.Println("Index created successfully on 'vacancies' collection")

	// Load the collection into memory
	err = milvusClient.LoadCollection(ctx, "vacancies", false)
	if err != nil {
		return fmt.Errorf("failed to load 'vacancies' collection: %w", err)
	}

	log.Println("'vacancies' collection loaded into memory")

	return nil
}

// CreateMilvusCollections initializes both resume and vacancy collections in Milvus.
func CreateMilvusCollections(milvusClient client.Client) error {
	// Create Resume Collection
	if err := CreateResumeCollection(milvusClient); err != nil {
		return err
	}

	// Create Vacancy Collection
	if err := CreateVacancyCollection(milvusClient); err != nil {
		return err
	}

	log.Println("All Milvus collections created successfully")
	return nil
}
