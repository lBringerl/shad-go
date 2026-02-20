//go:build !solution

package tparallel

type T struct {
	parallelBarrier chan *T
	parallelInvoked chan struct{}
	done            chan struct{}
	mut             chan struct{}

	parent *T
}

func (t *T) Parallel() {
	t.mut <- struct{}{}
	<-t.mut
	t.parallelInvoked <- struct{}{}
	t.parent.parallelBarrier <- t
}

func processParallel(t *T) {
	t.mut <- struct{}{}
	subtests := make([]*T, 0)
clearLoop:
	for {
		select {
		case parallelT := <-t.parallelBarrier:
			subtests = append(subtests, parallelT)
		default:
			break clearLoop
		}
	}
	<-t.mut
	for _, subtest := range subtests {
		<-subtest.done
	}
}

func (t *T) Run(subtest func(t *T)) {
	innerT := &T{
		parallelInvoked: make(chan struct{}),
		parallelBarrier: t.parallelBarrier,
		done:            make(chan struct{}),
		mut:             t.mut,

		parent: t,
	}

	go func() {
		defer func() {
			close(innerT.done)
		}()
		subtest(innerT)
		innerT.done <- struct{}{}
	}()

	select {
	case <-innerT.done:
		processParallel(t)
	case <-innerT.parallelInvoked:
		close(innerT.parallelInvoked)
	}
}

func Run(topTests []func(t *T)) {
	t := &T{
		parallelBarrier: make(chan *T),
		parallelInvoked: nil,
		done:            nil,
		mut:             make(chan struct{}, 1),

		parent: nil,
	}
	defer func() {
		close(t.parallelBarrier)
		close(t.mut)
	}()

	for _, test := range topTests {
		t.Run(test)
	}
	processParallel(t)
}
