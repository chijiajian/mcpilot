package tool

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type RequireApprovalTool struct{}

func (t *RequireApprovalTool) Name() string {
	return "require_approval"
}

func (t *RequireApprovalTool) Description() string {
	return "该工具用于在关键步骤请求人工批准，输入应说明为什么需要批准。"
}

func (t *RequireApprovalTool) InputSchema() map[string]string {
	return map[string]string{
		"reason": "为什么需要审批",
	}
}

func (t *RequireApprovalTool) Run(params map[string]string) (string, error) {
	reason, ok := params["reason"]
	if !ok {
		return "", fmt.Errorf("缺少参数: reason")
	}

	fmt.Println("\n🔒 审批请求: " + reason)
	fmt.Print("是否批准操作？(yes/no) > ")

	reader := bufio.NewReader(os.Stdin)
	resp, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	resp = strings.TrimSpace(strings.ToLower(resp))
	if resp == "yes" || resp == "y" {
		return "approved", nil
	} else if resp == "no" || resp == "n" {
		return "rejected", nil
	} else {
		return "", fmt.Errorf("无效输入: 请输入 yes 或 no")
	}
}
