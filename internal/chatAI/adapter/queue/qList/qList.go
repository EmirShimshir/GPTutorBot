package qList

import (
"container/list"
"errors"
	"github.com/EmirShimshir/tasker-bot/internal/chatAI/port"
)

var ErrValueFound = errors.New("value not found in queue")
var ErrOneEl = errors.New("one el in queue")

type QList struct {
	List *list.List
}

func NewQList() *QList {
	return &QList{
		List: list.New(),
	}
}

func (q *QList) Get() *port.Token {
	return q.List.Back().Value.(*port.Token)

}

func (q *QList) Add(v *port.Token) {
	q.List.PushFront(v)
}

func (q *QList) Next() *port.Token {
	back := q.List.Back()
	q.List.MoveToFront(back)
	return q.Get()
}

func (q *QList) GetAll() []*port.Token {
	res := make([]*port.Token, 0, 1)
	for e := q.List.Front(); e != nil; e = e.Next() {
		res = append(res, e.Value.(*port.Token))
	}

	return res
}

func (q *QList) Remove(token string) error {
	if q.List.Len() < 2 {
		return ErrOneEl
	}
	for e := q.List.Front(); e != nil; e = e.Next() {
		if e.Value.(*port.Token).Value == token {
			q.List.Remove(e)
			return nil
		}
	}

	return ErrValueFound
}
