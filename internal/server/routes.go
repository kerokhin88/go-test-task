package server

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"beetest/pkg/model"
	"beetest/pkg/queue"
)

func (s *Server) RegisterRoutes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("PUT /queue/{queue}", s.PutQueueMessage)
	mux.HandleFunc("GET /queue/{queue}", s.GetQueueMessage)

	return mux
}

func (s *Server) PutQueueMessage(w http.ResponseWriter, r *http.Request) {
	queueName := r.PathValue("queue")
	if err := s.valideQueueName(queueName); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var message model.QueueMessage
	if err := json.NewDecoder(r.Body).Decode(&message); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := s.valideMessage(&message); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := s.queue.PutMessage(queueName, message); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (s *Server) GetQueueMessage(w http.ResponseWriter, r *http.Request) {
	queueName := r.PathValue("queue")
	if err := s.valideQueueName(queueName); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	timeoutParam := r.URL.Query().Get("timeout")
	timeout, err := strconv.Atoi(timeoutParam)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := s.valideTimeout(timeout); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	message, err := s.queue.GetMessage(queueName, time.Duration(timeout)*time.Second)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if errors.Is(err, queue.ErrNoMessage) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var b bytes.Buffer
	if err := json.NewEncoder(&b).Encode(message); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(b.Bytes())
}
