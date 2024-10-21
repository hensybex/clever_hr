// internal/service/similarity.go

package service

import (
	"fmt"
	"math"
)

// CalculateCosineSimilarity computes the cosine similarity between two vectors
func CalculateCosineSimilarity(vec1, vec2 []float32) (float64, error) {
	if len(vec1) != len(vec2) {
		return 0, fmt.Errorf("vectors must be of the same length")
	}

	var dotProduct float64
	var normA float64
	var normB float64

	for i := 0; i < len(vec1); i++ {
		dotProduct += float64(vec1[i]) * float64(vec2[i])
		normA += float64(vec1[i]) * float64(vec1[i])
		normB += float64(vec2[i]) * float64(vec2[i])
	}

	if normA == 0 || normB == 0 {
		return 0, fmt.Errorf("one of the vectors is zero")
	}

	return dotProduct / (math.Sqrt(normA) * math.Sqrt(normB)), nil
}
