package parallel

import (
	"sync"

	"github.com/hashicorp/go-multierror"
)

type WaitGroup struct {
	syncWG     *sync.WaitGroup
	errChan    chan error
	err        *multierror.Error
	closeMutex sync.Mutex
}

func (wg *WaitGroup) ErrChan() chan<- error {
	return wg.errChan
}

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

func (wg *WaitGroup) Wait() error {
	wg.syncWG.Wait()
	close(wg.errChan)
	wg.closeMutex.Lock()
	defer wg.closeMutex.Unlock()

	return wg.err.ErrorOrNil()
}

func (wg *WaitGroup) Add(i int) {
	wg.syncWG.Add(i)
}

func (wg *WaitGroup) Done() {
	wg.syncWG.Done()
}
