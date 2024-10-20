package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"beetest/pkg/queue"
)

type Server struct {
	port  int
	queue queue.Queue
}

func NewServer(queue queue.Queue) *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT")) //move to flags
	NewServer := &Server{
		port:  port,
		queue: queue,
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
