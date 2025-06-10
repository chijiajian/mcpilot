package llm

import (
	"fmt"
	"strings"
)

const ragSystemPrompt = `你是一个基础设施专家助手，请根据下面的知识内容，回答用户提出的问题。
如果无法从中找到答案，请直接说明 “知识库中未找到相关内容”。

知识内容如下：
%s

问题：%s
请用简洁中文回答：`

func BuildRAGPrompt(context string, question string) string {
	return fmt.Sprintf(ragSystemPrompt, context, question)
}

func DefaultPromptTemplate(question string, contexts []string) string {
	contextBlock := strings.Join(contexts, "\n\n---\n\n")
	return fmt.Sprintf(`你是一个基础设施专家，请基于以下上下文回答用户问题：

上下文：
%s

问题：
%s

请使用中文回答，如果上下文不足以回答，可以说明你不确定。`, contextBlock, question)
}
