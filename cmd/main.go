package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/chijiajian/mcpilot/pkg/rag"
	"github.com/chijiajian/mcpilot/pkg/source"
	"github.com/qdrant/go-client/qdrant"
)

// DummyEmbedder returns a fixed-length vector for testing
func DummyEmbedder(ctx context.Context, text string) ([]float32, error) {
	vector := make([]float32, 1536)
	for i := range vector {
		vector[i] = float32(i%10) * 0.1 // 简单模式：0.0, 0.1, ..., 0.9, 重复
	}
	return vector, nil
}

func main() {
	ctx := context.Background()

	// Create temp dir for downloaded source
	tmpDir, err := os.MkdirTemp("", "rag-demo-*")
	if err != nil {
		log.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// 1. Init Git source (any public repo with .tf files)
	gitSrc := source.NewGitSource("https://github.com/ZStack-Robot/terraform-provider-zstack.git", "main")

	fmt.Println("📦 Downloading repo to:", tmpDir)
	if err := gitSrc.Fetch(ctx, tmpDir); err != nil {
		log.Fatalf("❌ Failed to download source: %v", err)
	}

	// 2. Init RAG
	ragClient, err := rag.NewRAG(ctx, "localhost", 6334, "terraform", "")
	if err != nil {
		log.Fatalf("❌ Failed to create RAG client: %v", err)
	}

	// 🔧 Ensure collection exists before ingesting
	err = rag.EnsureCollection(ctx, ragClient.Client(), "terraform", 1536, qdrant.Distance_Cosine)
	if err != nil {
		log.Fatalf("❌ Failed to ensure collection: %v", err)
	}

	defer ragClient.Close()

	// 3. Init Ingestor
	ingestor := rag.NewIngestor(ragClient, DummyEmbedder, []string{".tf", ".md"})

	// 4. Ingest downloaded directory
	fmt.Println("🧠 Ingesting downloaded knowledge...")
	if err := ingestor.IngestDirectory(ctx, tmpDir); err != nil {
		log.Fatalf("❌ Failed to ingest: %v", err)
	}

	fmt.Println("✅ Done. Files ingested from:", tmpDir)
}
