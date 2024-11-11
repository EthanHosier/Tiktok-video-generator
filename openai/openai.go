package openai

import (
	"context"
	"os"

	"github.com/sashabaranov/go-openai"
)

type OpenaiHandler interface{}

type OpenaiClient struct {
	client *openai.Client
}

func NewOpenaiClient() *OpenaiClient {
	client := openai.NewClient(os.Getenv("OPENAI_KEY"))

	return &OpenaiClient{client: client}
}

func (oc *OpenaiClient) CreateChatCompletion(ctx context.Context, prompt string) (*string, error) {
	resp, err := oc.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
		},
	)

	if err != nil {
		return nil, err
	}

	return &resp.Choices[0].Message.Content, nil
}
