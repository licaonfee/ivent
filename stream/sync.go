package stream

import (
	"context"

	"github.com/licaonfee/ivent"
)

//Sync allow to get Event in a syncronous way, optionally allows to set a buffer
type Sync struct {
	c      chan ivent.Event
	ctx    context.Context
	cancel context.CancelFunc
}

func (s *Sync) Send(e ivent.Event) {
	select {
	case s.c <- e:
	case <-s.ctx.Done():
		return
	}
}

func (s *Sync) Get() <-chan ivent.Event {
	return s.c
}

func (s *Sync) Close() {
	s.cancel()
	close(s.c)
}

func NewSync(ctx context.Context, buffer int) *Sync {
	s := &Sync{}
	s.c = make(chan ivent.Event, buffer)
	s.ctx, s.cancel = context.WithCancel(ctx)
	return s
}
