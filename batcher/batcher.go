//go:build !solution

package batcher

import (
	"gitlab.com/slon/shad-go/batcher/slow"
)

type Batcher struct {
	startLoadCh chan struct{}
	batchCh     chan struct{}
	batchSize   int
	waitLoadCh  chan any
	slowVal     *slow.Value
}

func (b *Batcher) clearWaiters() {
	lastVal := b.slowVal.Load()

	for range b.batchSize {
		b.waitLoadCh <- lastVal
	}
	b.batchSize = 0
}

func (b *Batcher) Load() any {
	select {
	case b.startLoadCh <- struct{}{}:
		defer func() {
			<-b.startLoadCh
		}()

		var firstVal any
		valChan := make(chan any)
		go func() {
			valChan <- b.slowVal.Load()
		}()
	accLoop:
		for {
			select {
			case firstVal = <-valChan:
				break accLoop
			case <-b.batchCh:
				b.batchSize++
			}
		}
		b.clearWaiters()
		return firstVal
	case b.batchCh <- struct{}{}:
		return <-b.waitLoadCh
	}
}

func NewBatcher(v *slow.Value) *Batcher {
	return &Batcher{
		startLoadCh: make(chan struct{}, 1),
		batchCh:     make(chan struct{}),
		waitLoadCh:  make(chan any),
		slowVal:     v,
	}
}
