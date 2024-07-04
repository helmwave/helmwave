package fileref

import (
	"context"
	"errors"
	"fmt"
	"github.com/containerd/log"
	"github.com/helmwave/helmwave/pkg/template"
)

//type RenderedFiles struct {
//	bufs map[string]fmt.Stringer
//	cond *sync.Cond
//}
//
//func NewRenderedFiles() *RenderedFiles {
//	r := &RenderedFiles{
//		bufs: make(map[string]fmt.Stringer),
//	}
//	r.cond = sync.NewCond(&sync.Mutex{})
//
//	return r
//}
//
//func (r *RenderedFiles) Add(name string, buf fmt.Stringer) {
//	r.cond.L.Lock()
//	r.bufs[name] = buf
//	r.cond.Broadcast()
//	r.cond.L.Unlock()
//}
//
//func (r *RenderedFiles) Get(name string) fmt.Stringer {
//	r.cond.L.Lock()
//	for r.bufs[name] == nil {
//		r.cond.Wait()
//	}
//	r.cond.L.Unlock()
//
//	return r.bufs[name]
//}

func (v *Config) tplOpts() []template.TemplaterOptions {
	return []template.TemplaterOptions{
		template.SetDelimiters(v.DelimiterLeft, v.DelimiterRight),
	}
}

// Render renders the fileref
func (v *Config) Render(ctx context.Context, data any, l *log.Entry) (err error) {
	// Try to get the file
	err = v.fetch(ctx, l)
	if !v.Strict && errors.Is(ErrNotExist, err) {
		return nil
	}

	// Rendering
	if v.isURL() {
		err = template.Tpl2yml(ctx, v.Dst, v.Dst, data, v.Renderer, v.tplOpts()...)
	} else {
		err = template.Tpl2yml(ctx, v.Src, v.Dst, data, v.Renderer, v.tplOpts()...)
	}

	if err != nil {
		return fmt.Errorf("failed to render %q file: %w", v.Src, err)
	}

	return nil

}
