package util

import "sync"

type ConcurrentQueue[T any] struct {
	items []T
	lock  sync.Mutex
	cond  *sync.Cond
}

func NewConcurrentQueue[T any]() *ConcurrentQueue[T] {
	q := &ConcurrentQueue[T]{}
	q.cond = sync.NewCond(&q.lock)
	return q
}

func (q *ConcurrentQueue[T]) Enqueue(item T) {
	q.lock.Lock()
	defer q.lock.Unlock()
	q.items = append(q.items, item)
	q.cond.Signal()
}

func (q *ConcurrentQueue[T]) Dequeue() T {
	q.lock.Lock()
	defer q.lock.Unlock()
	for len(q.items) == 0 {
		q.cond.Wait()
	}
	item := q.items[0]
	q.items = q.items[1:]
	return item
}

func (q *ConcurrentQueue[T]) IsEmpty() bool {
	return len(q.items) == 0
}

func (q *ConcurrentQueue[T]) Len() int {
	return len(q.items)
}
