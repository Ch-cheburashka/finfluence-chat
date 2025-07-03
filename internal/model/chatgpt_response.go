package model

type Role string

const (
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
)

type ChatMessage struct {
	Role Role   `json:"role"`
	Text string `json:"content"`
}

type ChatGPTResponse struct {
	Timestamp int64       `json:"created"`
	Model     string      `json:"model"`
	Message   ChatMessage `json:"message"`
}
