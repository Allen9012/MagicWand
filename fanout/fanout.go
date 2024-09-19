package fanout

import (
	xcontext "MagicWand/library/context"
	"context"
	"errors"
	"fmt"
	"os"
	"runtime"
	"sync"
)

var (
	// ErrFull chan full.
	ErrFull = errors.New("fanout: chan full")
)

type Options struct {
	worker int
	buffer int
}

// Option fanout option
type Option func(*Options)

// Worker specifies the worker of fanout
func Worker(n int) Option {
	if n <= 0 {
		panic("fanout: worker should > 0")
	}
	return func(o *Options) {
		o.worker = n
	}
}

// Buffer specifies the buffer of fanout
func Buffer(n int) Option {
	if n <= 0 {
		panic("fanout: buffer should > 0")
	}
	return func(o *Options) {
		o.buffer = n
	}
}

// Fanout async consume data from chan.
type Fanout struct {
	handler FanoutHandler
}

// New new a fanout struct.
func New(name string, opts ...Option) *Fanout {
	if name == "" {
		name = "anonymous"
	}
	o := &Options{
		worker: 1,
		buffer: 1024,
	}
	for _, op := range opts {
		op(o)
	}
	return &Fanout{handler: newFanoutWarpper(name, o, newFanout(name, o))}
}

// Do save a callback func with channel full err
func (c *Fanout) Do(ctx context.Context, f func(ctx context.Context)) (err error) {
	if f == nil {
		return nil
	}
	return c.handler.Do(ctx, f)
}

// SyncDo save a callback func no channel full err
func (c *Fanout) SyncDo(ctx context.Context, f func(ctx context.Context)) (err error) {
	if f == nil {
		return nil
	}
	return c.handler.SyncDo(ctx, f)
}

// Close close fanout
func (c *Fanout) Close() error {
	return c.handler.Close()
}

type item struct {
	f   func(c context.Context)
	ctx context.Context
}

type fanout struct {
	name    string
	ch      chan *item
	options *Options
	waiter  sync.WaitGroup

	ctx    context.Context
	cancel func()
}

func newFanout(name string, opts *Options) FanoutHandler {
	c := &fanout{
		ch:      make(chan *item, opts.buffer),
		name:    name,
		options: opts,
	}
	c.ctx, c.cancel = context.WithCancel(context.Background())
	c.waiter.Add(opts.worker)
	for i := 0; i < opts.worker; i++ {
		go c.proc()
	}
	_metricChanCap.Set(float64(opts.buffer), name)
	return c
}

func (c *fanout) proc() {
	defer c.waiter.Done()
	for {
		t := <-c.ch
		if t == nil {
			return
		}
		wrapFunc(t.f)(t.ctx)
		_metricChanSize.Set(float64(len(c.ch)), c.name)
		_metricCount.Inc(c.name)
	}
}

func wrapFunc(f func(c context.Context)) (res func(context.Context)) {
	res = func(ctx context.Context) {
		defer func() {
			if r := recover(); r != nil {
				buf := make([]byte, 64*1024)
				buf = buf[:runtime.Stack(buf, false)]
				fmt.Fprintf(os.Stderr, "fanout: panic recovered: %s\n%s\n", r, buf)
				log.Errorc(ctx, "panic in fanout proc, err: %s, stack: %s", r, buf)
			}
		}()
		f(ctx)
		if tr, ok := trace.FromContext(ctx); ok {
			tr.Finish(nil)
		}
	}
	return
}

// Do save a callback func with channel full err
func (c *fanout) Do(ctx context.Context, f func(ctx context.Context)) (err error) {
	if c.ctx.Err() != nil {
		return c.ctx.Err()
	}
	select {
	case c.ch <- &item{f: f, ctx: xcontext.Detach(ctx)}:
	default:
		err = ErrFull
		_metricChanFullCount.Inc(c.name)
	}
	_metricChanSize.Set(float64(len(c.ch)), c.name)
	return
}

// SyncDo save a callback func no channel full err
func (c *fanout) SyncDo(ctx context.Context, f func(ctx context.Context)) (err error) {
	if c.ctx.Err() != nil {
		return c.ctx.Err()
	}
	select {
	case c.ch <- &item{f: f, ctx: xcontext.Detach(ctx)}:
	case <-ctx.Done():
		err = ctx.Err()
	}
	_metricChanSize.Set(float64(len(c.ch)), c.name)
	return
}

// Close close fanout
func (c *fanout) Close() error {
	if err := c.ctx.Err(); err != nil {
		return err
	}
	c.cancel()
	for i := 0; i < c.options.worker; i++ {
		c.ch <- nil
	}
	c.waiter.Wait()
	return nil
}
