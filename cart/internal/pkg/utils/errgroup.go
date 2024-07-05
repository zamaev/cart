package utils

import (
	"context"
	"sync"
)

type ErrGroup struct {
	sync.RWMutex
	wg     sync.WaitGroup
	cancel context.CancelFunc
	err    error
}

func NewErrGroup(ctx context.Context) (*ErrGroup, context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	return &ErrGroup{
		cancel: cancel,
	}, ctx
}

func (e *ErrGroup) Wait() error {
	e.wg.Wait()
	return e.err
}

func (e *ErrGroup) Go(f func() error) {
	e.RLock()
	defer e.RUnlock()
	if e.err != nil {
		return
	}
	e.wg.Add(1)
	go func() {
		defer e.wg.Done()
		if err := f(); err != nil {
			e.Lock()
			defer e.Unlock()
			if e.err == nil {
				e.cancel()
				e.err = err
			}
		}
	}()
}
