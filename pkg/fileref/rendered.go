package fileref

import (
	"fmt"
	"sync"
)

type renderedFiles struct {
	bufs map[string]fmt.Stringer
	cond *sync.Cond
}

func NewRenderedFiles() *renderedFiles {
	r := &renderedFiles{
		bufs: make(map[string]fmt.Stringer),
	}
	r.cond = sync.NewCond(&sync.Mutex{})

	return r
}

func (r *renderedFiles) Add(name string, buf fmt.Stringer) {
	r.cond.L.Lock()
	r.bufs[name] = buf
	r.cond.Broadcast()
	r.cond.L.Unlock()
}

func (r *renderedFiles) Get(name string) fmt.Stringer {
	r.cond.L.Lock()
	for r.bufs[name] == nil {
		r.cond.Wait()
	}
	r.cond.L.Unlock()

	return r.bufs[name]
}
