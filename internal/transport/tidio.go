package transport

import (
	"finfluence-chat/internal/service"
	"io"
	"log"
	"net/http"
)

func MessageHandler(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		log.Printf("Method \"%s\" Not Allowed\n", request.Method)
		http.Error(writer, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	data, err := io.ReadAll(request.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		http.Error(writer, "Error reading body", http.StatusInternalServerError)
		return
	}
	defer request.Body.Close()

	tidioRequest, err := TidioToDomain(data)
	if err != nil {
		log.Printf("Error reading tidio tidio_message: %v", err)
		http.Error(writer, "Error reading tidio tidio_message", http.StatusInternalServerError)
		return
	}

	chatService := service.NewChatService()

	chatGPTResponse, err := chatService.ProcessMessage(tidioRequest)
	if err != nil {
		log.Printf("Error creating chat: %v", err)
		http.Error(writer, "Error creating chat", http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)

	toTidio, err := ToTidio(chatGPTResponse)
	if err != nil {
		log.Printf("Error converting to tidio: %v", err)
		http.Error(writer, "Error converting to tidio", http.StatusInternalServerError)
		return
	}

	_, err = writer.Write(toTidio)
	if err != nil {
		log.Println("Failed to write response")
		http.Error(writer, "Internal server error", http.StatusInternalServerError)
		return
	}
}
