package tool

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type HumanTool struct{}

func (t *HumanTool) Name() string {
	return "human"
}

func (t *HumanTool) Description() string {
	return "当系统需要向用户提问或请求指导时使用。输入应为问题内容。"
}

func (t *HumanTool) InputSchema() map[string]string {
	return map[string]string{
		"question": "你要问用户的问题",
	}
}

func (t *HumanTool) Run(params map[string]string) (string, error) {
	question, ok := params["question"]
	if !ok {
		return "", fmt.Errorf("缺少参数: question")
	}

	fmt.Println("\n🤖 问题: " + question)
	fmt.Print("👤 请输入你的回答 > ")

	reader := bufio.NewReader(os.Stdin)
	answer, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(answer), nil
}
