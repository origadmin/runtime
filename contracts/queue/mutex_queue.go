/*
 * Copyright (c) 2024 OrigAdmin. All rights reserved.
 */

package queue

import (
	"sync"
)

// MutexQueue is a mutex based queue implementation
type MutexQueue[E any] struct {
	mu    sync.Mutex
	items []E
}

// NewMutexQueue creates a new MutexQueue
func NewMutexQueue[E any]() *MutexQueue[E] {
	return &MutexQueue[E]{
		items: make([]E, 0),
	}
}

func (q *MutexQueue[E]) Offer(e E) bool {
	if q.mu.TryLock() {
		defer q.mu.Unlock()
		q.items = append(q.items, e)
		return true
	}
	return false
}

func (q *MutexQueue[E]) Poll() (E, bool) {
	var zero E
	if !q.mu.TryLock() {
		return zero, false
	}
	defer q.mu.Unlock()
	if len(q.items) == 0 {
		return zero, false
	}
	e := q.items[0]
	q.items = q.items[1:]
	return e, true
}

func (q *MutexQueue[E]) Peek() (E, bool) {
	var zero E
	if !q.mu.TryLock() {
		return zero, false
	}
	defer q.mu.Unlock()
	if len(q.items) == 0 {
		return zero, false
	}
	return q.items[0], true
}

func (q *MutexQueue[E]) Size() int64 {
	q.mu.Lock()
	defer q.mu.Unlock()
	return int64(len(q.items))
}

func (q *MutexQueue[E]) IsEmpty() bool {
	q.mu.Lock()
	defer q.mu.Unlock()
	return len(q.items) == 0
}

func (q *MutexQueue[E]) Clear() {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.items = q.items[:0]
}

func (q *MutexQueue[E]) ToSlice() []E {
	q.mu.Lock()
	defer q.mu.Unlock()
	result := make([]E, len(q.items))
	copy(result, q.items)
	return result
}

func (q *MutexQueue[E]) Iterator() Iterator[E] {
	return &mutexQueueIterator[E]{
		queue: q,
		index: 0,
	}
}

type mutexQueueIterator[E any] struct {
	queue *MutexQueue[E]
	index int
	items []E
}

func (it *mutexQueueIterator[E]) Next() bool {
	if it.items == nil {
		it.queue.mu.Lock()
		it.items = make([]E, len(it.queue.items))
		copy(it.items, it.queue.items)
		it.queue.mu.Unlock()
	}
	return it.index < len(it.items)
}

func (it *mutexQueueIterator[E]) Value() E {
	return it.items[it.index]
}

var _ Queue[any] = (*MutexQueue[any])(nil)
