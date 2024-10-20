package memory

import (
	"context"
	"sync"
	"time"

	"beetest/pkg/model"
	queuepkg "beetest/pkg/queue"
)

type MemoryQueueController struct {
	mx       sync.RWMutex
	queueMap map[string]*queueImpl
}

func NewMemoryQueueController() *MemoryQueueController {
	mc := MemoryQueueController{
		queueMap: make(map[string]*queueImpl),
	}
	return &mc
}

func (mq *MemoryQueueController) PutMessage(queue string, message model.QueueMessage) error {
	mq.mx.Lock()
	q, ok := mq.queueMap[queue]
	if !ok {
		q = newQueue()
		mq.queueMap[queue] = q
	}
	mq.mx.Unlock()

	q.getInChan() <- message

	return nil
}

func (mq *MemoryQueueController) GetMessage(queue string, timeout time.Duration) (model.QueueMessage, error) {
	mq.mx.Lock()
	q, ok := mq.queueMap[queue]
	if !ok {
		q = newQueue()
		mq.queueMap[queue] = q
	}
	mq.mx.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	select {
	case msg := <-q.getOutChan():
		return msg, nil
	case <-ctx.Done():
		return model.QueueMessage{}, queuepkg.ErrNoMessage
	}
}
