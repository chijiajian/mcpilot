package rag

import (
	"context"
	"testing"

	"github.com/chijiajian/mcpilot/pkg/llm"
)

func DummyEmbedFunc(ctx context.Context, text string) ([]float32, error) {
	vector := make([]float32, 1536)
	for i := range vector {
		vector[i] = float32(i%10) * 0.1 // 和 main.go 中一致
	}
	return vector, nil
}

func DummyLLMChatFunc(ctx context.Context, prompt string) (string, error) {
	return "模拟回答：" + prompt, nil
}

func TestRAGPipeline_Ask(t *testing.T) {
	ctx := context.Background()

	ragClient, err := NewRAG(ctx, "localhost", 6334, "terraform", "")
	if err != nil {
		t.Fatalf("连接 Qdrant 失败: %v", err)
	}
	defer ragClient.Close()

	// 不需要 EnsureCollection，假设已由 main.go 创建

	pipeline := RAGPipeline{
		RAG:        ragClient,
		EmbedFunc:  DummyEmbedFunc,
		LLMChat:    DummyLLMChatFunc,
		TopK:       2,
		PromptFunc: llm.DefaultPromptTemplate,
	}

	answer, err := pipeline.Ask(ctx, "ZStack 有哪些 Terraform 资源？")
	if err != nil {
		t.Fatalf("调用 Ask 失败: %v", err)
	}

	t.Logf("返回回答：%s", answer)
}
