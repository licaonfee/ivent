package stream

import (
	"context"
	"time"

	"github.com/licaonfee/ivent"
)

//Async allow to send in a Non-blocking unbufered chanel
type Async struct {
	c       chan ivent.Event
	timeout time.Duration
	ctx     context.Context
	cancel  context.CancelFunc
}

func (s *Async) sendWithTimeout(e ivent.Event) {
	ctx, c := context.WithTimeout(s.ctx, s.timeout)
	defer c()
	select {
	case s.c <- e:
		return
	case <-ctx.Done():
		return
	}
}

func (s *Async) sendWithContext(e ivent.Event) {
	ctx, c := context.WithCancel(s.ctx)
	defer c()
	select {
	case s.c <- e:
	case <-ctx.Done():
		return
	}
}

//Send put an event in queue using a new goroutine whith expiration
func (s *Async) Send(e ivent.Event) {
	if s.timeout > time.Duration(0) {
		go s.sendWithTimeout(e)
		return
	}
	go s.sendWithContext(e)
}

//Get return a read only chanel
func (s *Async) Get() <-chan ivent.Event {
	return s.c
}

//Close chan
func (s *Async) Close() {
	s.cancel()
	close(s.c)
}

//NewAsync create a stream of events if timeout is 0 then messages never expire
func NewAsync(ctx context.Context, timeout time.Duration) *Async {
	s := &Async{}
	s.c = make(chan ivent.Event)
	s.timeout = timeout
	s.ctx, s.cancel = context.WithCancel(ctx)
	return s
}
