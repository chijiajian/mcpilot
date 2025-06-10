package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/chijiajian/mcpilot/pkg/agent"
	"github.com/chijiajian/mcpilot/pkg/llm"
	"github.com/chijiajian/mcpilot/pkg/rag"
	"github.com/chijiajian/mcpilot/pkg/tool"
)

func main() {

	ctx := context.Background()
	apiKey := os.Getenv("OPENAI_API_KEY")

	r, err := rag.NewRAG(ctx, "localhost", 6333, "zstack", apiKey)
	if err != nil {
		log.Fatalf("初始化 RAG 失败: %v", err)
	}
	defer r.Close()

	embedFunc := func(ctx context.Context, text string) ([]float32, error) {
		vec := make([]float32, 1536) // OpenAI 的嵌入向量维度
		for i := range vec {
			vec[i] = float32(i)
		}
		return vec, nil
	}

	ingestor := rag.NewIngestor(r, embedFunc, []string{".tf", ".md"})

	if err := ingestor.IngestDirectory(ctx, "./docs"); err != nil {
		log.Fatalf("文档导入失败: %v", err)
	}

	planner := agent.NewPlanner()

	// 注册工具
	planner.Register(&tool.ZStackVMTool{})
	planner.Register(&tool.HumanTool{})
	planner.Register(&tool.RequireApprovalTool{})
	planner.Register(&tool.ShowReasoningTool{})
	planner.Register(&tool.HideReasoningTool{})

	llmClient := llm.NewOpenAIClient()

	// 模拟 chat 输入
	userInput := "请帮我 create vm，配置2核CPU，4GB内存"

	result, err := planner.ExecuteWithLLM(ctx, userInput, "question", llmClient)
	if err != nil {
		log.Fatalf("执行失败: %v", err)
	}

	fmt.Println("执行结果:\n" + result)

}
