package rag

import (
	"context"
	"fmt"

	qdrant "github.com/qdrant/go-client/qdrant"
)

func EnsureCollection(ctx context.Context, client *qdrant.Client, name string, vectorSize int, distance qdrant.Distance) error {
	exists, err := client.CollectionExists(ctx, name)
	if err != nil {
		return fmt.Errorf("failed to check if collection exists: %w", err)
	}
	if exists {
		return nil
	}

	req := &qdrant.CreateCollection{
		CollectionName: name,
		VectorsConfig: &qdrant.VectorsConfig{
			Config: &qdrant.VectorsConfig_Params{
				Params: &qdrant.VectorParams{
					Size:     uint64(vectorSize),
					Distance: distance,
				},
			},
		},
	}

	if err := client.CreateCollection(ctx, req); err != nil {
		return fmt.Errorf("failed to create collection: %w", err)
	}
	return nil
}

func DeleteCollection(ctx context.Context, client *qdrant.Client, name string) error {
	return client.DeleteCollection(ctx, name)
}

func ListCollections(ctx context.Context, client *qdrant.Client) ([]string, error) {
	return client.ListCollections(ctx)
}

func CountCollectionPoints(ctx context.Context, client *qdrant.Client, name string) (uint64, error) {
	return client.Count(ctx, &qdrant.CountPoints{CollectionName: name})
}
