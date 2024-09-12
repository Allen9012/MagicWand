package fanout

import (
	"context"
	"sync"
)

type Middleware func(name string, opts *Options, handler FanoutHandler) FanoutHandler

type FanoutHandler interface {
	Do(ctx context.Context, f func(ctx context.Context)) (err error)
	SyncDo(ctx context.Context, f func(ctx context.Context)) (err error)
	Close() error
}

var _globalHandlers []Middleware
var _globalMu sync.RWMutex

func RegisterGlobalMiddleware(ms ...Middleware) {
	_globalMu.Lock()
	defer _globalMu.Unlock()
	_globalHandlers = append(_globalHandlers, ms...)
}

func newFanoutWarpper(name string, opts *Options, h FanoutHandler) FanoutHandler {
	_globalMu.RLock()
	defer _globalMu.RUnlock()
	var mergedHandler = h
	for i := len(_globalHandlers) - 1; i >= 0; i-- {
		mergedHandler = _globalHandlers[i](name, opts, mergedHandler)
	}
	return mergedHandler
}

func init() {
	RegisterGlobalMiddleware(newTraceHandler)
}
