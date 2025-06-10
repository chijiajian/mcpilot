package rag

import (
	"context"
	"fmt"

	qdrant "github.com/qdrant/go-client/qdrant"
)

type RAG struct {
	client     *qdrant.Client
	collection string
}

type Document struct {
	ID      string            `json:"id"`
	Vector  []float32         `json:"vector"`
	Payload map[string]string `json:"payload"`
	Score   float32           `json:"score,omitempty"`
}

func (r *RAG) Client() *qdrant.Client {
	return r.client
}

func NewRAG(ctx context.Context, host string, port int, collection string, apiKey string) (*RAG, error) {
	// Create proper config based on the actual Config struct
	config := &qdrant.Config{
		Host:   host,
		Port:   port,
		APIKey: apiKey,
		// Set UseTLS to true if using HTTPS/TLS
		UseTLS: false, // Change as needed
	}

	// Create the Qdrant client with proper error handling
	client, err := qdrant.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create Qdrant client: %w", err)
	}

	return &RAG{
		client:     client,
		collection: collection,
	}, nil
}

// Add a Close method to properly clean up resources
func (r *RAG) Close() error {
	if closer, ok := any(r.client).(interface{ Close() error }); ok {
		return closer.Close()
	}
	return nil
}

func (r *RAG) AddDocument(ctx context.Context, docs []Document) error {
	var points []*qdrant.PointStruct
	for _, doc := range docs {
		points = append(points, &qdrant.PointStruct{
			Id:      &qdrant.PointId{PointIdOptions: &qdrant.PointId_Uuid{Uuid: doc.ID}},
			Vectors: &qdrant.Vectors{VectorsOptions: &qdrant.Vectors_Vector{Vector: &qdrant.Vector{Data: doc.Vector}}},
			Payload: stringMapToPayload(doc.Payload),
		})
	}

	_, err := r.client.Upsert(ctx, &qdrant.UpsertPoints{
		CollectionName: r.collection,
		Points:         points,
	})
	return err
}

func (r *RAG) Search(ctx context.Context, queryVector []float32, topK int) ([]Document, error) {
	limit := uint64(topK)

	resp, err := r.client.Query(ctx, &qdrant.QueryPoints{
		CollectionName: r.collection,
		Query:          qdrant.NewQuery(queryVector...),
		Limit:          &limit,
		WithPayload:    qdrant.NewWithPayload(true),
		WithVectors:    qdrant.NewWithVectors(true),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to query points: %w", err)
	}

	var docs []Document
	for _, point := range resp {
		var vector []float32
		if point.Vectors != nil && point.Vectors.GetVector() != nil {
			vector = point.Vectors.GetVector().Data
		}
		docs = append(docs, Document{
			ID:      point.Id.GetUuid(),
			Vector:  vector,
			Payload: payloadToStringMap(point.Payload),
			Score:   point.Score,
		})
	}

	return docs, nil
}

func stringMapToPayload(input map[string]string) map[string]*qdrant.Value {
	payload := make(map[string]*qdrant.Value)
	for k, v := range input {
		payload[k] = &qdrant.Value{Kind: &qdrant.Value_StringValue{StringValue: v}}
	}
	return payload
}

func payloadToStringMap(input map[string]*qdrant.Value) map[string]string {
	result := make(map[string]string)
	for k, v := range input {
		if s := v.GetStringValue(); s != "" {
			result[k] = s
		}
	}
	return result
}
