package queue

import (
	"container/list"
	"errors"
)

var ErrValueFound = errors.New("value not found in queue")

type Queue struct {
	List *list.List
}

func NewQueue() *Queue {
	return &Queue{
		List: list.New(),
	}
}

func (q *Queue) Get() any {
	return q.List.Back().Value
}

func (q *Queue) Add(v any) {
	q.List.PushFront(v)
}

func (q *Queue) Next() any {
	back := q.List.Back()
	q.List.MoveToFront(back)
	return q.Get()
}

func (q *Queue) GetAll() []any {
	res := make([]any, 0, 1)
	for e := q.List.Front(); e != nil; e = e.Next() {
		res = append(res, e.Value)
	}

	return res
}

func (q *Queue) Remove(v any) error {
	for e := q.List.Front(); e != nil; e = e.Next() {
		if e.Value == v {
			q.List.Remove(e)
			return nil
		}
	}

	return ErrValueFound
}