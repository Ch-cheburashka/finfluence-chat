package transport

import (
	"context"
	"finfluence-chat/internal/service"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

func NewTidioHandler(chatService *service.ChatService) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodPost {
			http.Error(writer, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		ctx, cancel := context.WithTimeout(request.Context(), 60*time.Second)
		defer cancel()

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

		AIResponse, err := chatService.ProcessMessage(ctx, incomingMessage)
		if err != nil {
			log.Printf("chatService error: %v", err)
			http.Error(writer, fmt.Sprint(err), http.StatusInternalServerError)
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
		_, err = writer.Write(outgoingMessage)
		if err != nil {
			log.Println("Failed to write response")
			http.Error(writer, "Internal server error", http.StatusInternalServerError)
			return
		}
	}
}
