package main

import (
	"errors"
	"sync"
)

type Queue struct {
	m sync.Mutex
	q []string
}

func NewQueue() *Queue {
	return &Queue{
		m: sync.Mutex{},
		q: make([]string, 0),
	}
}

func (q *Queue) BulkPush(elms []string) {
	q.m.Lock()
	q.q = append(q.q, elms...)
	q.m.Unlock()
}

func (q *Queue) Push(elm string) {
	q.m.Lock()
	q.q = append(q.q, elm)
	q.m.Unlock()
}

func (q *Queue) Pop() (string, error) {
	q.m.Lock()
	defer q.m.Unlock()
	if len(q.q) == 0 {
		return "", errors.New("")
	}
	p := q.q[0]
	q.q = q.q[1:]
	return p, nil
}

func (q *Queue) IsEmpty() bool {
	q.m.Lock()
	defer q.m.Unlock()
	return len(q.q) == 0
}
