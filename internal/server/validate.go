package server

import (
	"beetest/pkg/model"
	"fmt"
)

func (s *Server) valideQueueName(queue string) error {
	if queue == "" {
		return fmt.Errorf("invalid queue name")
	}
	return nil
}

func (s *Server) valideMessage(message *model.QueueMessage) error {
	if message == nil {
		return fmt.Errorf("empty message")
	}
	if message.Message == "" {
		return fmt.Errorf("invalid messsage payload")
	}
	return nil
}

func (s *Server) valideTimeout(timeout int) error {
	if timeout <= 0 {
		return fmt.Errorf("invalid timeout")
	}
	return nil
}
