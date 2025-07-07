package model

type Role string

const (
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
)

type Message struct {
	UserID    string
	Text      string
	Role      Role
	Timestamp int64
}
