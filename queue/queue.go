package queue

import (
	"container/list"
	"sync"
)

type queue struct {
	mtx        sync.Mutex
	innerChan  chan struct{}
	queueTasks *list.List
}

func newQueue() *queue {
	item := queue{}
	item.innerChan = make(chan struct{}, 1)
	item.queueTasks = list.New()
	return &item
}

func (item *queue) push(task *Task) {
	item.mtx.Lock()
	item.queueTasks.PushBack(task)
	item.mtx.Unlock()
}

func (item *queue) pop() *Task {
	item.mtx.Lock()
	defer item.mtx.Unlock()

	if item.queueTasks.Len() == 0 {
		return nil
	}
	elem := item.queueTasks.Front()
	item.queueTasks.Remove(elem)
	return elem.Value.(*Task)
}

func InpQueue(inp chan *Task) *queue {
	queue := newQueue()
	go inpProcess(inp, queue)
	return queue
}

func inpProcess(inp chan *Task, queue *queue) {
	for value := range inp {
		queue.push(value)

		select {
		case queue.innerChan <- struct{}{}:
		default:
		}
	}
	close(queue.innerChan)
}

type Task struct {
	ID   int
	Data string
}