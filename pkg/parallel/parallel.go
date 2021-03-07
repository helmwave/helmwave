package parallel

import (
	"github.com/hashicorp/go-multierror"
	"sync"
)

type WaitGroup struct {
	syncWG  *sync.WaitGroup
	errChan chan error
	err     *multierror.Error
}

func (wg *WaitGroup) ErrChan() chan<- error {
	return wg.errChan
}

func NewWaitGroup() *WaitGroup {
	return &WaitGroup{
		syncWG:  &sync.WaitGroup{},
		errChan: make(chan error),
		err:     &multierror.Error{},
	}
}

func (wg *WaitGroup) gatherErrors() {
	for err := range wg.errChan {
		if err != nil {
			wg.err = multierror.Append(wg.err, err)
		}
	}
}

func (wg *WaitGroup) Wait() error {
	go wg.gatherErrors()
	wg.syncWG.Wait()
	close(wg.errChan)

	return wg.err.ErrorOrNil()
}

func (wg *WaitGroup) Add(i int) {
	wg.syncWG.Add(i)
}

func (wg *WaitGroup) Done() {
	wg.syncWG.Done()
}
