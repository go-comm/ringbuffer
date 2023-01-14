package ringbuffer

import (
	"errors"
	"sync/atomic"
)

var (
	ErrFull = errors.New("ringbuffer: data full")
	ErrNil  = errors.New("ringbuffer: no data")
)

type sequence struct {
	cur int32
	_   [7]uint64
}

func (s *sequence) Request() int32 {
	return atomic.LoadInt32(&s.cur)
}

func (s *sequence) Commit(req int32, new int32) (ok bool) {
	ok = atomic.CompareAndSwapInt32(&s.cur, req, new)
	return ok
}

type RingBuffer interface {
	Put(f func(i int)) error
	Get(f func(i int)) error
}

func New(n int) RingBuffer {
	return &ringBuffer{
		head: &sequence{},
		tail: &sequence{},
		size: int32(n),
	}
}

type ringBuffer struct {
	head *sequence
	tail *sequence
	size int32
}

func (rb *ringBuffer) Put(f func(i int)) error {
	for {
		headReq := rb.head.Request()
		tailReq := rb.tail.Request()
		if headReq-tailReq == -1 || headReq-tailReq == rb.size-1 {
			return ErrFull
		}
		next := (headReq + 1) % rb.size
		if ok := rb.head.Commit(headReq, next); !ok {
			continue
		}
		f(int(next))
		break
	}
	return nil
}

func (rb *ringBuffer) Get(f func(i int)) error {
	for {
		headReq := rb.head.Request()
		tailReq := rb.tail.Request()
		if headReq == tailReq {
			return ErrNil
		}
		next := (tailReq + 1) % rb.size
		if ok := rb.tail.Commit(tailReq, next); !ok {
			continue
		}
		f(int(next))
		break
	}
	return nil
}
