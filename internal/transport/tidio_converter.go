package transport

import (
	"encoding/json"
	"finfluence-chat/internal/model"
)

type tidioVisitor struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type tidioMessage struct {
	ID        string `json:"id"`
	Text      string `json:"text"`
	Timestamp int64  `json:"timestamp"`
}

type tidioWebhook struct {
	Event     string       `json:"event"`
	Timestamp int64        `json:"timestamp"`
	Visitor   tidioVisitor `json:"visitor"`
	Message   tidioMessage `json:"message"`
}

func TidioToDomain(data []byte) (model.Message, error) {
	webhook := tidioWebhook{}

	if err := json.Unmarshal(data, &webhook); err != nil {
		return model.Message{}, err
	}

	return model.Message{
		UserID:    webhook.Visitor.ID,
		Text:      webhook.Message.Text,
		Role:      model.RoleUser,
		Timestamp: webhook.Message.Timestamp,
	}, nil
}

func ToTidio(message model.Message) ([]byte, error) {
	webhook := tidioWebhook{
		Event:     "chat_message",
		Timestamp: message.Timestamp,
		Visitor: tidioVisitor{
			ID:   message.UserID,
			Name: "assistant",
		},
		Message: tidioMessage{
			ID:        message.UserID,
			Text:      message.Text,
			Timestamp: message.Timestamp,
		},
	}
	data, err := json.Marshal(webhook)
	if err != nil {
		return nil, err
	}
	return data, nil
}
