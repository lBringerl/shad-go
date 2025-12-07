//go:build !solution

package lrucache

import (
	"container/list"
)

type CacheDataPoint struct {
	value       int
	listElement *list.Element
}

type CacheImpl struct {
	keysList *list.List
	mem      map[int]*CacheDataPoint
	cap      int
}

func New(cap int) Cache {
	return CacheImpl{
		keysList: list.New(),
		mem:      make(map[int]*CacheDataPoint, cap),
		cap:      cap,
	}
}

func (c CacheImpl) setNotExists(key int, val int) *CacheDataPoint {
	if c.cap > 0 && c.keysList.Len() >= c.cap {
		c.removeLastKey()
	}

	element := c.keysList.PushFront(key)
	dataPoint := &CacheDataPoint{
		value:       val,
		listElement: element,
	}
	c.mem[key] = dataPoint

	return dataPoint
}

func (c CacheImpl) removeLastKey() {
	removedKey := c.keysList.Remove(c.keysList.Back()).(int)
	delete(c.mem, removedKey)
}

func (c CacheImpl) Get(key int) (int, bool) {
	if c.cap == 0 {
		return 0, false
	}

	val, exists := c.mem[key]
	if !exists {
		val = c.setNotExists(key, 0)
	} else {
		c.keysList.MoveBefore(val.listElement, c.keysList.Front())
	}

	return val.value, exists
}

func (c CacheImpl) Set(key, value int) {
	if c.cap == 0 {
		return
	}

	dataPoint, exists := c.mem[key]
	if !exists {
		c.setNotExists(key, value)
	} else {
		c.keysList.MoveBefore(dataPoint.listElement, c.keysList.Front())
		dataPoint.value = value
	}
}

func (c CacheImpl) Range(f func(key, value int) bool) {
	curElement := c.keysList.Back()
	for curElement != nil {
		if !f(curElement.Value.(int), c.mem[curElement.Value.(int)].value) {
			break
		}
		curElement = curElement.Prev()
	}
}

func (c CacheImpl) Clear() {
	clear(c.mem)
	c.keysList.Init()
}
