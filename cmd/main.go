package main

import (
	"context"
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
	store := service.NewInMemoryStore()
	chatService := service.NewChatService(store)
	handler := http.TimeoutHandler(
		transport.NewTidioHandler(chatService),
		15*time.Second, `{"error":"timeout"}`,
	)

	http.Handle("/message", handler)

	port := flag.String("port", "8080", "Port to listen on")
	flag.Parse()

	server := &http.Server{
		Addr:         ":" + *port,
		Handler:      nil,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 20 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("HTTP server listening on %s …", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
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
