//go:build !solution

package dupcall

import (
	"context"
	"reflect"
	"sync"
)

type cbResult struct {
	numActive int
	resultCh  chan struct{}
	result    interface{}
	err       error
	cancelFn  func()
	mu        *sync.Mutex
}

type Call struct {
	cbRes sync.Map
}

func (o *Call) Do(
	ctx context.Context,
	cb func(context.Context) (interface{}, error),
) (result interface{}, err error) {
	funcEntryPoint := reflect.ValueOf(cb).Pointer()

	val, exists := o.cbRes.LoadOrStore(funcEntryPoint, &cbResult{
		mu: &sync.Mutex{},
	})
	cbr := val.(*cbResult)
	cbr.mu.Lock()
	if !exists {
		innerCtx, cancel := context.WithCancel(context.Background())
		cbr.cancelFn = cancel
		cbr.resultCh = make(chan struct{})
		go func() {
			defer func() {
				o.cbRes.Delete(funcEntryPoint)
				close(cbr.resultCh)
			}()
			res, err := cb(innerCtx)
			cbr.result = res
			cbr.err = err
		}()
	}
	cbr.numActive++
	cbr.mu.Unlock()

	select {
	case <-ctx.Done():
		cbr.mu.Lock()
		cbr.numActive--
		if cbr.numActive == 0 {
			cbr.cancelFn()
		}
		cbr.mu.Unlock()
		return nil, ctx.Err()
	case <-cbr.resultCh:
		cbr.mu.Lock()
		cbr.numActive = 0
		cbr.mu.Unlock()
		return cbr.result, cbr.err
	}
}
