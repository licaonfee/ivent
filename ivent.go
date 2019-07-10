package ivent

import (
	"strconv"
	"time"
)

//Class interface to define custom event classes
type Class interface {
	Value() int64
	String() string
}

//Event ss
type Event struct {
	T    Class
	Time time.Time
	Tags map[string]string
	Data interface{}
}

//NewEvent create an event of type t qith data d and Time = Now
func NewEvent(t Class, tags map[string]string, d interface{}) Event {
	return Event{T: t, Data: d, Tags: tags, Time: time.Now()}
}

//Client this interface can be used by client libraries to send events
//final user must pass an implementation of Stream to this
//all implementations must be use default Noop Stream
type Client interface {
	//WithStream implementation must be concurrent safe
	WithStream(Stream)
}

//Stream this interface should be used for final users to get events from Client
type Stream interface {
	//Send all implementations of this method must be concurent safe
	Send(Event)
}

//Any default Class Interface implementation
type Any int

//Value just return int(Any)
func (a Any) Value() int64 {
	return int64(a)
}

//String return Any as a numeric string
func (a Any) String() string {
	return strconv.Itoa(int(a))
}
