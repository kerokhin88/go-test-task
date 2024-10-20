package queue

import (
	"time"

	"beetest/pkg/model"
)

type Queue interface {
	PutMessage(queue string, message model.QueueMessage) error
	GetMessage(queue string, timeout time.Duration) (model.QueueMessage, error)
}
