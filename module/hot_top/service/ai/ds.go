package ai

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-deepseek/deepseek"
	"github.com/go-deepseek/deepseek/request"
)

type DeepSeek struct {
	Client deepseek.Client
}

func NewAiDeepSeek(token string) (*DeepSeek, error) {
	client, err := deepseek.NewClient(token)
	if err != nil {
		return nil, err
	}
	return &DeepSeek{
		Client: client,
	}, nil
}

func (a *DeepSeek) Call(question string) (answer string, err error) {
	if question == "" {
		return "", errors.New("question is empty")
	}
	chatReq := &request.ChatCompletionsRequest{
		Model:  deepseek.DEEPSEEK_CHAT_MODEL,
		Stream: false,
		Messages: []*request.Message{
			{
				Role:    "user",
				Content: question, // set your input message
			},
		},
	}

	chatResp, err := (a.Client).CallChatCompletionsChat(context.Background(), chatReq)
	if err != nil {
		fmt.Println("Error =>", err)
		return "", err
	}
	if len(chatResp.Choices) == 0 {
		return "", errors.New("no answer")
	}
	answer = chatResp.Choices[0].Message.Content
	return answer, nil
}
