package utils

import (
	"container/list"
	"sync"
)

type Buffer interface {
	Push(v interface{}) *list.Element
	Pop() interface{}
	Find(f func(interface{}) bool) *list.Element
	Filter(f func(interface{}) bool) []interface{}
}

type BufferList struct {
	sync.RWMutex
	List *list.List
}

func NewBufferList() *BufferList {
	return &BufferList{List: list.New()}
}

func (this *BufferList) Push(v interface{}) *list.Element {
	this.Lock()
	e := this.List.PushFront(v)
	this.Unlock()
	return e
}

func (this *BufferList) Pop() interface{} {
	this.Lock()
	if el := this.List.Back(); el != nil {
		item := this.List.Remove(el)
		this.Unlock()
		return item
	}
	this.Unlock()
	return nil
}

func (this *BufferList) Find(f func(interface{}) bool) *list.Element {
	this.RLock()
	defer this.RUnlock()
	for e := this.List.Front(); e != nil; e = e.Next() {
		if f(e.Value) {
			return e
		}
	}
	return nil
}

func (this *BufferList) Filter(f func(interface{}) bool) []interface{} {
	this.RLock()
	defer this.RUnlock()
	var r = make([]interface{}, 0)
	for e := this.List.Front(); e != nil; e = e.Next() {
		if f(e.Value) {
			r = append(r, e.Value)
		}
	}
	return r
}
