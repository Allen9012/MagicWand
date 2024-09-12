package fanout

import "context"

var _ FanoutHandler = &TraceHandler{}

var traceTags = []trace.Tag{
	{Key: trace.TagSpanKind, Value: "background"},
	{Key: trace.TagComponent, Value: "sync/pipeline/fanout"},
}

type TraceHandler struct {
	FanoutHandler
	name string
}

func newTraceHandler(name string, _ *Options, h FanoutHandler) FanoutHandler {
	return &TraceHandler{FanoutHandler: h, name: name}
}

func (t *TraceHandler) Do(ctx context.Context, f func(ctx context.Context)) (err error) {
	if tr, ok := trace.FromContext(ctx); ok {
		tr = tr.Fork("", "Fanout:Do").SetTag(traceTags...)
		tr.SetTag(trace.String("fanout.name", t.name))
		ctx = trace.NewContext(ctx, tr)
	}
	return t.FanoutHandler.Do(ctx, f)
}

func (t *TraceHandler) SyncDo(ctx context.Context, f func(ctx context.Context)) (err error) {
	if tr, ok := trace.FromContext(ctx); ok {
		tr = tr.Fork("", "Fanout:SyncDo").SetTag(traceTags...)
		tr.SetTag(trace.String("fanout.name", t.name))
		ctx = trace.NewContext(ctx, tr)
	}
	return t.FanoutHandler.SyncDo(ctx, f)
}

func (t *TraceHandler) Close() (err error) {
	return t.FanoutHandler.Close()
}
