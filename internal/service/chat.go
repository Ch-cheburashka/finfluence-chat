package service

import (
	"context"
	"finfluence-chat/internal/model"
	"fmt"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"os"
)

type ChatService struct {
	OpenaiClient openai.Client
	History      HistoryStore
}

func NewChatService(historyStore HistoryStore) *ChatService {
	return &ChatService{
		OpenaiClient: openai.NewClient(
			option.WithAPIKey(os.Getenv("OPENAI_API_KEY")),
		),
		History: historyStore,
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
	s.History.Add(message.UserID, message)
	chatCompletion, err := s.OpenaiClient.Chat.Completions.New(context.TODO(), openai.ChatCompletionNewParams{
		Messages: convertToUnion(s.History.Get(message.UserID)),
		Model:    openai.ChatModelGPT4,
	})
	if err != nil {
		return model.Message{}, err
	}
	if len(chatCompletion.Choices) == 0 {
		return model.Message{}, fmt.Errorf("openai returned no choices")
	}
	aiMessage := model.Message{
		UserID:    message.UserID,
		Text:      chatCompletion.Choices[0].Message.Content,
		Role:      model.RoleAssistant,
		Timestamp: chatCompletion.Created,
	}
	s.History.Add(aiMessage.UserID, aiMessage)
	return aiMessage, nil
}
