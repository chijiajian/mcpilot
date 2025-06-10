package llm

import (
	"context"
	"fmt"
	"os"

	openai "github.com/sashabaranov/go-openai"
)

type OpenAIClient struct {
	client *openai.Client
}

func NewOpenAIClient() *OpenAIClient {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		panic("OPENAI_API_KEY 未设置")
	}
	client := openai.NewClient(apiKey)
	return &OpenAIClient{client: client}
}

func (c *OpenAIClient) ChatCompletion(ctx context.Context, prompt string) (string, error) {
	resp, err := c.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: openai.GPT4,
		Messages: []openai.ChatCompletionMessage{
			{Role: "system", Content: "你是一个帮助用户部署ZStack云主机的助手。"},
			{Role: "user", Content: prompt},
		},
	})
	if err != nil {
		return "", err
	}
	return resp.Choices[0].Message.Content, nil
}

func Example() {
	cli := NewOpenAIClient()
	ctx := context.Background()
	output, err := cli.ChatCompletion(ctx, "帮我创建一台2核4G内存的ZStack虚拟机，名字叫testvm")
	if err != nil {
		panic(err)
	}
	fmt.Println("LLM响应:", output)
}
