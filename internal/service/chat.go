package service

import (
	"context"
	"finfluence-chat/internal/model"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"os"
)

type ChatService struct {
	OpenaiClient openai.Client
	History      []model.Message
}

func NewChatService() *ChatService {
	return &ChatService{
		OpenaiClient: openai.NewClient(
			option.WithAPIKey(os.Getenv("OPENAI_API_KEY")),
		),
		History: make([]model.Message, 0),
	}
}

func convertToUnion(messages []model.Message) []openai.ChatCompletionMessageParamUnion {
	converted := make([]openai.ChatCompletionMessageParamUnion, len(messages))
	for i, m := range messages {
		if m.Role == model.RoleAssistant {
			converted[i] = openai.AssistantMessage(m.Text)
		} else {
			converted[i] = openai.UserMessage(m.Text)
		}
	}
	return converted
}

func (s *ChatService) ProcessMessage(message model.Message) (model.Message, error) {
	s.History = append(s.History, message)
	client := ChatService{
		OpenaiClient: s.OpenaiClient,
		History:      s.History,
	}
	chatCompletion, err := client.OpenaiClient.Chat.Completions.New(context.TODO(), openai.ChatCompletionNewParams{
		Messages: convertToUnion(s.History),
		Model:    openai.ChatModelGPT4,
	})
	if err != nil {
		return model.Message{}, err
	}
	return model.Message{
		UserID:    message.UserID,
		Text:      chatCompletion.Choices[0].Message.Content,
		Role:      model.RoleAssistant,
		Timestamp: chatCompletion.Created,
	}, nil
}
