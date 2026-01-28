//go:build !solution

package rwmutex

// A RWMutex is a reader/writer mutual exclusion lock.
// The lock can be held by an arbitrary number of readers or a single writer.
// The zero value for a RWMutex is an unlocked mutex.
//
// If a goroutine holds a RWMutex for reading and another goroutine might
// call Lock, no goroutine should expect to be able to acquire a read lock
// until the initial read lock is released. In particular, this prohibits
// recursive read locking. This is to ensure that the lock eventually becomes
// available; a blocked Lock call excludes new readers from acquiring the
// lock.

type LockType int

const (
	Reader LockType = iota
	Writer
)

type RWMutex struct {
	globalCh        chan struct{}
	releaseWriterCh chan struct{}
	writerLockCh    chan struct{}
	counter         int
	writerWaiting   bool
}

// New creates *RWMutex.
func New() *RWMutex {
	return &RWMutex{
		globalCh:        make(chan struct{}, 1),
		releaseWriterCh: make(chan struct{}),
		writerLockCh:    make(chan struct{}, 1),
		counter:         0,
		writerWaiting:   false,
	}
}

// RLock locks rw for reading.
//
// It should not be used for recursive read locking; a blocked Lock
// call excludes new readers from acquiring the lock. See the
// documentation on the RWMutex type.
func (rw *RWMutex) RLock() {
	rw.writerLockCh <- struct{}{}
	rw.globalCh <- struct{}{}
	<-rw.writerLockCh
	rw.counter++
	<-rw.globalCh
}

// RUnlock undoes a single RLock call;
// it does not affect other simultaneous readers.
// It is a run-time error if rw is not locked for reading
// on entry to RUnlock.
func (rw *RWMutex) RUnlock() {
	rw.globalCh <- struct{}{}
	rw.counter--
	switch {
	case rw.counter == 0 && rw.writerWaiting:
		rw.writerWaiting = false
		rw.releaseWriterCh <- struct{}{}
	case rw.counter < 0:
		panic("reader unlock: tried to unlock already unlocked read lock\n")
	}
	<-rw.globalCh
}

// Lock locks rw for writing.
// If the lock is already locked for reading or writing,
// Lock blocks until the lock is available.
func (rw *RWMutex) Lock() {
	rw.writerLockCh <- struct{}{}
	rw.globalCh <- struct{}{}
	if rw.counter > 0 {
		rw.writerWaiting = true
		<-rw.globalCh
		<-rw.releaseWriterCh
		rw.globalCh <- struct{}{}
	}
	<-rw.globalCh
}

// Unlock unlocks rw for writing. It is a run-time error if rw is
// not locked for writing on entry to Unlock.
//
// As with Mutexes, a locked RWMutex is not associated with a particular
// goroutine. One goroutine may RLock (Lock) a RWMutex and then
// arrange for another goroutine to RUnlock (Unlock) it.
func (rw *RWMutex) Unlock() {
	rw.globalCh <- struct{}{}
	<-rw.writerLockCh
	<-rw.globalCh
}
