package once

import (
	"sync"
	"sync/atomic"
)

// Once 支持可重入、Init and Close多次（历史sdk问题）
type Once struct {
	num atomic.Int32
	m   sync.Mutex
}

// Do Init once
func (o *Once) Do(f func()) {
	if o.num.Add(1) == 1 {
		o.m.Lock()
		defer o.m.Unlock()
		f()
	}
}

// UnDo for close
func (o *Once) UnDo(f func()) {
	if o.num.Add(-1) == 0 {
		o.m.Lock()
		defer o.m.Unlock()
		f()
	}
}
