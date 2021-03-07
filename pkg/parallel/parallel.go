package parallel

import (
	"github.com/hashicorp/go-multierror"
	"sync"
)

type WaitGroup struct {
	syncWG  *sync.WaitGroup
	errChan chan error
}

func (wg *WaitGroup) ErrChan() chan<- error {
	return wg.errChan
}

func NewWaitGroup() *WaitGroup {
	return &WaitGroup{
		syncWG:  &sync.WaitGroup{},
		errChan: make(chan error),
	}
}

func (wg *WaitGroup) Wait() error {
	wg.syncWG.Wait()

	result := &multierror.Error{}

	for err := range wg.errChan {
		if err != nil {
			result = multierror.Append(result, err)
		}
	}

	return result.ErrorOrNil()
}

func (wg *WaitGroup) Add(i int) {
	wg.syncWG.Add(i)
}

func (wg *WaitGroup) Done() {
	wg.syncWG.Done()
}
