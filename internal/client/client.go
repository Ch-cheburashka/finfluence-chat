package client

import (
	"context"
	"encoding/json"
	"finfluence-chat/internal/model"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"os"
)

func ReceiveTidioMessage(data []byte) (model.TidioRequest, error) {
	tidioRequest := model.TidioRequest{}
	err := json.Unmarshal(data, &tidioRequest)
	if err != nil {
		return model.TidioRequest{}, err
	}
	tidioRequest.Role = model.RoleUser
	return tidioRequest, nil
}

func ReceiveChatCompletion(history []openai.ChatCompletionMessageParamUnion) (model.ChatGPTResponse, error) {
	client := openai.NewClient(
		option.WithAPIKey(os.Getenv("OPENAI_API_KEY")),
	)
	chatCompletion, err := client.Chat.Completions.New(context.TODO(), openai.ChatCompletionNewParams{
		Messages: history,
		Model:    openai.ChatModelGPT4_1Mini,
	})
	if err != nil {
		return model.ChatGPTResponse{}, err
	}
	chatGPTResponse := model.ChatGPTResponse{
		Timestamp: chatCompletion.Created,
		Model:     chatCompletion.Model,
		Message: model.ChatMessage{
			Role: model.RoleAssistant,
			Text: chatCompletion.Choices[0].Message.Content,
		},
	}
	return chatGPTResponse, nil
}
