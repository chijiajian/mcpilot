package rag

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	openai "github.com/sashabaranov/go-openai"
)

// Ingestor handles reading files, generating embeddings, and storing in Qdrant.
type Ingestor struct {
	rag       *RAG
	embedFunc func(ctx context.Context, text string) ([]float32, error)
	exts      []string
}

// NewIngestor creates a new Ingestor instance.embedFunc: 将文本转成向量的嵌入函数（可插拔，默认可用 OpenAI）。
func NewIngestor(r *RAG, embedFunc func(ctx context.Context, text string) ([]float32, error), exts []string) *Ingestor {
	return &Ingestor{
		rag:       r,
		embedFunc: embedFunc,
		exts:      exts,
	}
}

// IngestFile reads a single file, splits it, embeds each chunk, and stores in vector DB.
func (i *Ingestor) IngestFile(ctx context.Context, filePath string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	chunks := ChunkText(string(content))
	var docs []Document
	for idx, chunk := range chunks {
		if len(strings.TrimSpace(chunk)) == 0 {
			continue
		}

		vector, err := i.embedFunc(ctx, chunk)

		if err != nil {
			return fmt.Errorf("embedding failed on chunk %d: %w", idx, err)
		}
		fmt.Printf("👉 Embedding for chunk %d has dimension: %d\n", idx, len(vector))

		docs = append(docs, Document{
			ID:     fmt.Sprintf("%s#%d", filepath.Base(filePath), idx),
			Vector: vector,
			Payload: map[string]string{
				"content": chunk,
				"source":  filePath,
			},
		})
	}

	return i.rag.AddDocument(ctx, docs)
}

// IngestDirectory recursively ingests all .md/.tf/.txt files in a directory.
func (i *Ingestor) IngestDirectory(ctx context.Context, dir string) error {
	return filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}
		if !i.shouldIngest(path) {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read file %s: %w", path, err)
		}

		vec, err := i.embedFunc(ctx, string(content))
		if err != nil {
			return fmt.Errorf("embed failed for %s: %w", path, err)
		}

		doc := Document{
			ID:     uuid.NewString(),
			Vector: vec,
			Payload: map[string]string{
				"path": path,
				"text": string(content), // 如果内容较大可裁剪或做摘要
			},
		}

		if err := i.rag.AddDocument(ctx, []Document{doc}); err != nil {
			return fmt.Errorf("add document failed for %s: %w", path, err)
		}
		return nil
	})
}

func (i *Ingestor) shouldIngest(path string) bool {
	ext := filepath.Ext(path)
	for _, e := range i.exts {
		if ext == e {
			return true
		}
	}
	return false
}

// ChunkText splits text into logical chunks (paragraphs here). 暂未实现基于 token 数的 chunker
func ChunkText(text string) []string {
	scanner := bufio.NewScanner(strings.NewReader(text))
	scanner.Split(bufio.ScanLines)

	var chunks []string
	var currentChunk strings.Builder

	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			if currentChunk.Len() > 0 {
				chunks = append(chunks, currentChunk.String())
				currentChunk.Reset()
			}
			continue
		}
		currentChunk.WriteString(line + "\n")
	}
	if currentChunk.Len() > 0 {
		chunks = append(chunks, currentChunk.String())
	}
	return chunks
}

// EmbedOpenAI is a default embedding function using OpenAI.
func EmbedOpenAI(apiKey string) func(ctx context.Context, text string) ([]float32, error) {
	client := openai.NewClient(apiKey)
	return func(ctx context.Context, text string) ([]float32, error) {
		if len(text) > 8191 {
			text = text[:8191]
		}

		resp, err := client.CreateEmbeddings(ctx, openai.EmbeddingRequest{
			Input: []string{text},
			Model: openai.AdaEmbeddingV2,
		})
		if err != nil {
			return nil, err
		}
		if len(resp.Data) == 0 {
			return nil, errors.New("no embedding returned")
		}
		return resp.Data[0].Embedding, nil
	}
}
