package model

type TidioMessage struct {
	ID   string `json:"id"`
	Text string `json:"text"`
}

type visitor struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type TidioRequest struct {
	Timestamp int64        `json:"timestamp"`
	Message   TidioMessage `json:"message"`
	Visitor   visitor      `json:"visitor"`
	Role      Role         `json:"role"`
}
