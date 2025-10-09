package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	httpserver "github.com/Prathamesh314/http_server_from_scrath_learning/internal/server"
)

const port = 42069

func main() {
	server, err := httpserver.Serve(port)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}