package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/chijiajian/mcpilot/pkg/llm"
	"github.com/chijiajian/mcpilot/pkg/tool"
)

type Planner struct {
	Tools map[string]tool.Tool
}

type LLMOutput struct {
	Tool   string            `json:"tool"`
	Params map[string]string `json:"params"`
}

func NewPlanner() *Planner {
	return &Planner{
		Tools: map[string]tool.Tool{},
	}
}

func (p *Planner) Register(t tool.Tool) {
	p.Tools[strings.ToLower(t.Name())] = t

}

func ParseLLMOutput(content string) (*LLMOutput, error) {
	var out LLMOutput
	if err := json.Unmarshal([]byte(content), &out); err != nil {
		return nil, fmt.Errorf("解析 LLM 输出失败: %v", err)
	}
	if out.Tool == "" {
		return nil, fmt.Errorf("输出中未包含工具名")
	}
	return &out, nil
}

// 执行LLM的结构化输出
func (p *Planner) ExecuteFromLLMOutput(content string) (string, error) {
	llmOut, err := ParseLLMOutput(content)
	if err != nil {
		return "", fmt.Errorf("解析 LLM 输出失败: %v", err)
	}

	toolImpl, ok := p.Tools[strings.ToLower(llmOut.Tool)]
	if !ok {
		return "", fmt.Errorf("未知工具: %s", llmOut.Tool)
	}
	return toolImpl.Run(llmOut.Params)
}

func (p *Planner) ExecuteWithLLM(ctx context.Context, userInput string, question string, llmClient *llm.OpenAIClient) (string, error) {
	prompt := llm.BuildRAGPrompt(userInput, question)
	content, err := llmClient.ChatCompletion(ctx, prompt)
	if err != nil {
		return "", err
	}
	return p.ExecuteFromLLMOutput(content)
}

func PlanFromInput(userInput string) (map[string]interface{}, error) {
	// 简单规则 mock，后续你可以接入 llm + rag
	plan := make(map[string]interface{})

	if strings.Contains(userInput, "2C4G") {
		plan["cpu"] = 2
		plan["memory"] = 4096
	}

	if strings.Contains(userInput, "云主机") {
		plan["resource_type"] = "zstack_vm_instance"
		plan["image"] = "CentOS-7.9"
	}

	return plan, nil
}
