package tool

import (
	"github.com/chijiajian/mcpilot/pkg/config"
)

type ShowReasoningTool struct{}

func (t *ShowReasoningTool) Name() string {
	return "show_reasoning"
}

func (t *ShowReasoningTool) Description() string {
	return "启用推理过程显示"
}

func (t *ShowReasoningTool) InputSchema() map[string]string {
	return map[string]string{} // 无参数
}

func (t *ShowReasoningTool) Run(params map[string]string) (string, error) {
	config.SetShowReasoning(true)
	return "推理输出已启用", nil
}

type HideReasoningTool struct{}

func (t *HideReasoningTool) Name() string {
	return "hide_reasoning"
}

func (t *HideReasoningTool) Description() string {
	return "关闭推理过程显示"
}

func (t *HideReasoningTool) InputSchema() map[string]string {
	return map[string]string{} // 无参数
}

func (t *HideReasoningTool) Run(params map[string]string) (string, error) {
	config.SetShowReasoning(false)
	return "推理输出已关闭", nil
}
