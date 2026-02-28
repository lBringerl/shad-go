//go:build !solution

package keylock

import (
	"sync"
)

type writerData struct {
	keys         []string
	permissionCh chan bool
}

type KeyLock struct {
	keysMutex sync.Mutex
	keys      map[string]struct{}
	startCh   chan struct{}
	waiters   chan struct{}
	writers   chan writerData
}

func (l *KeyLock) processWaiters() {
	waitersNumber := 1 // at least 1 waiter started processing
iterWaiters:
	for {
		select {
		case <-l.waiters:
			waitersNumber++
		default:
			break iterWaiters
		}
	}

	for range waitersNumber {
		wrData := <-l.writers
		if l.tryLockKeys(wrData.keys) {
			wrData.permissionCh <- true
		} else {
			wrData.permissionCh <- false
		}
	}
}

func (l *KeyLock) tryLockKeys(keys []string) bool {
	l.keysMutex.Lock()
	defer l.keysMutex.Unlock()
	for _, key := range keys {
		_, exists := l.keys[key]
		if exists {
			return false
		}
	}
	for _, key := range keys {
		l.keys[key] = struct{}{}
	}
	return true
}

func New() *KeyLock {
	lock := &KeyLock{
		keys:    make(map[string]struct{}),
		startCh: make(chan struct{}, 1),
		waiters: make(chan struct{}),
		writers: make(chan writerData),
	}
	return lock
}

func (l *KeyLock) LockKeys(keys []string, cancel <-chan struct{}) (canceled bool, unlock func()) {
	unlock = func() {
		l.keysMutex.Lock()
		for _, key := range keys {
			delete(l.keys, key)
		}
		l.keysMutex.Unlock()

		select {
		case <-l.startCh:
		default:
		}
	}

	keysLocked := l.tryLockKeys(keys)
	if keysLocked {
		return
	}

	permissionCh := make(chan bool, 1)

	for {
		select {
		case l.startCh <- struct{}{}:
			go l.processWaiters()
		case l.waiters <- struct{}{}:
		case <-cancel:
			return true, nil
		}

		l.writers <- writerData{
			keys:         keys,
			permissionCh: permissionCh,
		}

		select {
		case permission := <-permissionCh:
			if !permission {
				continue
			}
			return
		case <-cancel:
			return true, nil
		}
	}
}
