package main

import (
	"encoding/json"
	"finfluence-chat/internal/client"
	"flag"
	"github.com/openai/openai-go"
	"io"
	"log"
	"net/http"
)

var history []openai.ChatCompletionMessageParamUnion

func messageHandler(writer http.ResponseWriter, request *http.Request) {
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

	tidioRequest, err := client.ReceiveTidioMessage(data)
	if err != nil {
		log.Printf("Error reading tidio message: %v", err)
		http.Error(writer, "Error reading tidio message", http.StatusInternalServerError)
		return
	}

	history = append(history, openai.UserMessage(tidioRequest.Message.Text))

	chatGPTResponse, err := client.ReceiveChatCompletion(history)
	if err != nil {
		log.Printf("Error creating chat: %v", err)
		http.Error(writer, "Error creating chat", http.StatusInternalServerError)
		return
	}

	history = append(history, openai.AssistantMessage(chatGPTResponse.Message.Text))

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	jsonData, err := json.Marshal(chatGPTResponse)
	if err != nil {
		log.Println("Failed to write create JSON object")
		http.Error(writer, "Internal server error", http.StatusInternalServerError)
		return
	}
	_, err = writer.Write(jsonData)
	if err != nil {
		log.Println("Failed to write response")
		http.Error(writer, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func main() {
	port := flag.String("port", "8080", "Port to listen on")
	flag.Parse()
	http.HandleFunc("/message", messageHandler)
	log.Printf("Server started at http://localhost:%s\n", *port)
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}
