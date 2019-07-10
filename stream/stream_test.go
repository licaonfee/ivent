package stream_test

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/licaonfee/ivent"
	"github.com/licaonfee/ivent/stream"
)

func TestAsync_Send(t *testing.T) {

	nilEvent := func() ivent.Event {
		return ivent.NewEvent(ivent.Any(0), nil, nil)
	}

	type fields struct {
		timeout time.Duration
	}
	type args struct {
		e []ivent.Event
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		wantTimeout bool
	}{
		{"Without timeout",
			fields{time.Duration(0)},
			args{[]ivent.Event{nilEvent()}},
			false},
		{"With no expired timeout",
			fields{time.Second},
			args{[]ivent.Event{nilEvent()}},
			false},
		{"With expired timeout",
			fields{time.Second},
			args{[]ivent.Event{nilEvent()}},
			true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//This to ensure context not expire when timeout ==0
			newTout := (time.Millisecond * 500) + tt.fields.timeout*2
			ctx, c := context.WithTimeout(context.Background(), newTout)
			defer c()
			s := stream.NewAsync(ctx, tt.fields.timeout)
			defer s.Close()
			for _, e := range tt.args.e {
				s.Send(e)
			}
			if tt.wantTimeout {
				time.Sleep(tt.fields.timeout + time.Second)
			}
			got := make([]ivent.Event, 0, 0)
		EXIT_LOOP:
			for {
				select {
				case e := <-s.Get():
					got = append(got, e)
				case <-ctx.Done():
					break EXIT_LOOP
				}
			}
			if tt.wantTimeout && len(got) > 0 {
				t.Errorf("Stream expect timeout but results where given: %v", got)
			} else if !tt.wantTimeout && !reflect.DeepEqual(got, tt.args.e) {
				t.Errorf("Stream = %v , want %v", got, tt.args.e)
			}
		})
	}
}

func Benchmark_AsyncStream_WithTimeout(b *testing.B) {
	b.StopTimer()
	e := ivent.NewEvent(ivent.Any(1), map[string]string{"name": "name", "item": "1"}, []byte("foo"))
	async := stream.NewAsync(context.Background(), time.Second)
	go func() {
		for range async.Get() {

		}
	}()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		async.Send(e)
	}
}

func Benchmark_AsyncStream_NOTimeout(b *testing.B) {
	b.StopTimer()
	e := ivent.NewEvent(ivent.Any(1), map[string]string{"name": "name", "item": "1"}, []byte("foo"))
	async := stream.NewAsync(context.Background(), 0)
	go func() {
		for range async.Get() {

		}
	}()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		async.Send(e)
	}
}

func Benchmark_SyncStream_unbuffered(b *testing.B) {
	b.StopTimer()
	e := ivent.NewEvent(ivent.Any(1), map[string]string{"name": "name", "item": "1"}, []byte("foo"))
	sync := stream.NewSync(context.Background(), 0)
	go func() {
		for range sync.Get() {

		}
	}()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		sync.Send(e)
	}
}

func Benchmark_SyncStream_buffer_1000(b *testing.B) {
	b.StopTimer()
	e := ivent.NewEvent(ivent.Any(1), map[string]string{"name": "name", "item": "1"}, []byte("foo"))
	sync := stream.NewSync(context.Background(), 1000)
	go func() {
		for range sync.Get() {

		}
	}()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		sync.Send(e)
	}
}
