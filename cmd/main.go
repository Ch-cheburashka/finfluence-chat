package main

import (
	"context"
	"errors"
	"finfluence-chat/internal/service"
	"finfluence-chat/internal/transport"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	store := service.NewInMemoryStore(50)
	chatService := service.NewChatService(store)
	handler := transport.NewTidioHandler(chatService)

	http.Handle("/message", handler)

	port := flag.String("port", "8080", "Port to listen on")
	flag.Parse()

	server := &http.Server{
		Addr:    ":" + *port,
		Handler: nil,
	}

	go func() {
		log.Printf("HTTP server listening on %s …", server.Addr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Kill)

	<-quit
	log.Println("Shutdown requested – draining connections…")

	ctx, cancel := context.WithTimeout(context.Background(), 25*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("server forced to shutdown: %v", err)
	}

	log.Println("Server exited cleanly")
}
