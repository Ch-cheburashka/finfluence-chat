package transport

import (
	"finfluence-chat/internal/service"
	"io"
	"log"
	"net/http"
)

func NewTidioHandler(chatService *service.ChatService) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodPost {
			http.Error(writer, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		body, err := io.ReadAll(request.Body)
		defer request.Body.Close()
		if err != nil {
			log.Printf("read error: %v", err)
			http.Error(writer, "Internal Error", http.StatusInternalServerError)
			return
		}

		incomingMessage, err := TidioToDomain(body)
		if err != nil {
			log.Printf("decode error: %v", err)
			http.Error(writer, "Bad Request", http.StatusBadRequest)
			return
		}

		AIResponse, err := chatService.ProcessMessage(incomingMessage)
		if err != nil {
			log.Printf("chatService error: %v", err)
			http.Error(writer, "Internal Error", http.StatusInternalServerError)
			return
		}

		outgoingMessage, err := ToTidio(AIResponse)
		if err != nil {
			log.Printf("encode error: %v", err)
			http.Error(writer, "Internal Error", http.StatusInternalServerError)
			return
		}

		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)
		writer.Write(outgoingMessage)
	}
}
