package rag

import (
	"context"
	"fmt"
	"os"

	"github.com/chijiajian/mcpilot/pkg/source"
)

// Pipeline connects a Source with an Ingestor to fetch + embed + store documents.
type Pipeline struct {
	Source   source.Source // Source interface for fetching documents
	Ingestor *Ingestor
}

// Run executes the full fetch → embed → store pipeline.
func (p *Pipeline) Run(ctx context.Context) error {
	// Create a temporary directory to store downloaded or cloned files
	tmpDir, err := os.MkdirTemp("", "rag_pipeline_")
	if err != nil {
		return fmt.Errorf("failed to create temp dir: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	// Step 1: Fetch documents from source
	if err := p.Source.Fetch(ctx, tmpDir); err != nil {
		return fmt.Errorf("failed to fetch source: %w", err)
	}

	// Step 2: Ingest documents into vector DB
	if err := p.Ingestor.IngestDirectory(ctx, tmpDir); err != nil {
		return fmt.Errorf("failed to ingest documents: %w", err)
	}

	return nil
}
