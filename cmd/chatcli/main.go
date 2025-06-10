package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/chijiajian/mcpilot/pkg/agent"
)

func main() {
	fmt.Println("💬 欢迎使用 Appliot Chat CLI。输入任务，输入 exit 退出。")

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("🧑‍💻 你：")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("读取输入失败: %v\n", err)
			continue
		}

		input = strings.TrimSpace(input)
		if input == "exit" || input == "退出" {
			fmt.Println("👋 再见！")
			break
		}

		plan, err := agent.PlanFromInput(input)
		if err != nil {
			fmt.Printf("❌ 生成计划失败: %v\n", err)
			continue
		}

		fmt.Println("✅ 生成的资源计划如下：")
		for k, v := range plan {
			fmt.Printf("  - %s: %v\n", k, v)
		}
	}
}
