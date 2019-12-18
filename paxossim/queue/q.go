package queue

import (
	"container/list"
	"log"
	"sync"
)

// Queue - An interface to queues
type Queue interface {
	// EnQueue an element to the queue
	Enqueue(item interface{})
	// DeQueue an element from the queue
	Dequeue() (interface{}, bool)
	// Return the length of the underlying queue
	Len() int
	// Wait until an item in the underlying queue is ready
	// Once an item is available. Dequeue this item and return it
	WaitForItem() interface{}
}

// ConcurrentQueue - A queue which is go-routine safe
type ConcurrentQueue struct {
	// q is the backing linked list for this queue
	q *list.List
	// A cond variable used for th WaitForItem call
	c *sync.Cond
	// An underlying mutex to provide go-routine safety
	mu *sync.Mutex
}

func NewQueue() Queue {
	mu := &sync.Mutex{}
	return &ConcurrentQueue{
		q:  list.New(),
		mu: mu,
		c:  sync.NewCond(mu),
	}
}

func (mq *ConcurrentQueue) Enqueue(item interface{}) {
	mq.c.L.Lock()
	defer mq.c.L.Unlock()
	mq.q.PushBack(item)
	// Broadcast all waiting go-routines
	// ToDo: figure out if signal is sufficient
	mq.c.Broadcast()
}

func (mq *ConcurrentQueue) Dequeue() (interface{}, bool) {
	mq.c.L.Lock()
	defer mq.c.L.Unlock()
	return mq.dequeue()
}

func (mq *ConcurrentQueue) dequeue() (interface{}, bool) {
	if mq.q.Len() == 0 {
		return nil, false
	}
	elem := mq.q.Front()
	mq.q.Remove(elem)
	return elem.Value, true
}

func (mq *ConcurrentQueue) Len() int {
	mq.c.L.Lock()
	defer mq.c.L.Lock()
	return mq.q.Len()
}

func (mq *ConcurrentQueue) WaitForItem() interface{} {
	mq.c.L.Lock()
	for mq.q.Len() == 0 {
		// Refer the comment for cond.Wait
		// Wait atomically unlocks c.L and suspends execution
		// of the calling goroutine. After later resuming execution,
		// Wait locks c.L before returning. Unlike in other systems,
		// Wait cannot return unless awoken by Broadcast or Signal.
		//
		// Because c.L is not locked when Wait first resumes, the caller
		// typically cannot assume that the condition is true when
		// Wait returns. Instead, the caller should Wait in a loop:
		mq.c.Wait()
	}
	defer mq.c.L.Unlock()
	v, ok := mq.dequeue()
	if !ok {
		log.Fatalf("Fatal error expected to have at least one element")
	}
	return v
}
