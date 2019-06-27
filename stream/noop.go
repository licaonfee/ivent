package stream

import "github.com/licaonfee/ivent"

//Noop default ivent implementation , does nothing
type Noop struct {
}

//Send this implementation just discard the event
func (n *Noop) Send(e ivent.Event) {}

//Get return a closed chanel
func Get() <-chan ivent.Event {
	c := make(chan ivent.Event)
	defer close(c)
	return c
}

//NewNoop create a dummy ivent.Stream
func NewNoop() ivent.Stream {
	return &Noop{}
}
