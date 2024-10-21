// internal/config/milvus_client.go

package config

import (
	"context"
	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"log"
)

func InitMilvusClient() client.Client {
	ctx := context.Background()
	milvusClient, err := client.NewGrpcClient(ctx, "milvus:19530")
	if err != nil {
		log.Fatalf("Failed to connect to Milvus: %v", err)
	}
	return milvusClient
}
