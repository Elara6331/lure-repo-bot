// Package queue provides an unbounded FIFO queue
// for payloads received via webhooks.
package queue

import (
	"container/list"
)

type Queue[T any] struct {
	buf      *list.List
	out      chan T
	valAdded chan struct{}
}

func New[T any]() *Queue[T] {
	q := &Queue[T]{
		buf:      list.New().Init(),
		out:      make(chan T),
		valAdded: make(chan struct{}),
	}

	go func() {
		for {
			if q.buf.Len() == 0 {
				<-q.valAdded
			}

			e := q.buf.Front()
			q.out <- e.Value.(T)
			q.buf.Remove(e)
		}
	}()

	return q
}

func (q *Queue[T]) Add(val T) {
	q.buf.PushBack(val)
	select {
	case q.valAdded <- struct{}{}:
	default:
	}
}

func (q *Queue[T]) Channel() <-chan T {
	return q.out
}
