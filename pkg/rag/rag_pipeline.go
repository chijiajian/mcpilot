package rag

import (
	"context"
	"fmt"
)

type RAGPipeline struct {
	RAG        *RAG
	EmbedFunc  func(ctx context.Context, text string) ([]float32, error)
	LLMChat    func(ctx context.Context, prompt string) (string, error)
	TopK       int
	PromptFunc func(question string, contexts []string) string
}

// Ask 执行嵌入 → 检索 → 组装 prompt → 调用LLM
func (p *RAGPipeline) Ask(ctx context.Context, question string) (string, error) {
	// 1. 问题 → 向量
	vec, err := p.EmbedFunc(ctx, question)
	if err != nil {
		return "", fmt.Errorf("embed question failed: %w", err)
	}

	// 2. 向量检索
	docs, err := p.RAG.Search(ctx, vec, p.TopK)
	if err != nil {
		return "", fmt.Errorf("vector search failed: %w", err)
	}

	// 3. 抽取上下文
	var contexts []string
	for _, doc := range docs {
		if text, ok := doc.Payload["text"]; ok {
			contexts = append(contexts, text)
		} else if content, ok := doc.Payload["content"]; ok {
			contexts = append(contexts, content)
		}
	}

	// 4. 构建提示词
	prompt := p.PromptFunc(question, contexts)

	// 5. 调用LLM
	answer, err := p.LLMChat(ctx, prompt)
	if err != nil {
		return "", fmt.Errorf("LLM chat failed: %w", err)
	}
	return answer, nil
}
