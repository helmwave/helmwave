package parallel

import (
	"fmt"
	"sync"

	"github.com/hashicorp/go-multierror"
)

// WaitGroup is an extension to sync.WaitGroup that provided channel for errors.
type WaitGroup struct {
	syncWG     *sync.WaitGroup
	errChan    chan error
	err        *multierror.Error
	closeMutex sync.Mutex
}

// ErrChan returns channel for errors.
func (wg *WaitGroup) ErrChan() chan<- error {
	return wg.errChan
}

// NewWaitGroup initializes new *WaitGroup and runs errors collection goroutine.
func NewWaitGroup() *WaitGroup {
	wg := &WaitGroup{
		syncWG:  &sync.WaitGroup{},
		errChan: make(chan error),
		err:     &multierror.Error{},
	}
	go wg.gatherErrors()

	return wg
}

func (wg *WaitGroup) gatherErrors() {
	wg.closeMutex.Lock()
	defer wg.closeMutex.Unlock()

	for err := range wg.errChan {
		wg.err = multierror.Append(wg.err, err)
	}
}

// Wait waits for all goroutines to exit (via Done method) and returns multierror for all errors or nil.
func (wg *WaitGroup) Wait() error {
	wg.syncWG.Wait()
	close(wg.errChan)
	wg.closeMutex.Lock()
	defer wg.closeMutex.Unlock()

	if err := wg.err.ErrorOrNil(); err != nil {
		return fmt.Errorf("one of goroutines in waitgroup sent error: %w", err)
	}

	return nil
}

// Add adds delta to WaitGroup counter.
func (wg *WaitGroup) Add(i int) {
	wg.syncWG.Add(i)
}

// Done decrements WaitGroup counter by 1.
func (wg *WaitGroup) Done() {
	wg.syncWG.Done()
}
