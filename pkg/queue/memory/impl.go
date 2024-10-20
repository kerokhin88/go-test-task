package memory

import (
	"beetest/pkg/model"
)

type queueImpl struct {
	in, out chan model.QueueMessage
	list    []model.QueueMessage
}

func newQueue() *queueImpl {
	q := &queueImpl{
		in:   make(chan model.QueueMessage),
		out:  make(chan model.QueueMessage),
		list: make([]model.QueueMessage, 0),
	}
	go q.run()
	return q
}

func (q *queueImpl) run() {
	for {
		if len(q.list) > 0 {
			select {
			case q.out <- q.list[0]:
				q.list = q.list[1:]
			case msg := <-q.in:
				q.list = append(q.list, msg)
			}
		} else {
			msg := <-q.in
			q.list = append(q.list, msg)
		}
	}
}

func (q *queueImpl) getInChan() chan model.QueueMessage {
	return q.in
}

func (q *queueImpl) getOutChan() chan model.QueueMessage {
	return q.out
}
