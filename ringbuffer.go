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
	cap := roundUp(int32(n))
	return &ringBuffer{
		head: &sequence{},
		tail: &sequence{},
		mask: cap - 1,
		cap:  cap,
	}
}

type ringBuffer struct {
	head *sequence
	tail *sequence
	mask int32
	cap  int32
}

func (rb *ringBuffer) Put(f func(i int)) error {
	for {
		headReq := rb.head.Request()
		tailReq := rb.tail.Request()
		if headReq-tailReq == -1 || headReq-tailReq == rb.cap-1 {
			return ErrFull
		}
		next := (headReq + 1) & rb.mask
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
		next := (tailReq + 1) & rb.mask
		if ok := rb.tail.Commit(tailReq, next); !ok {
			continue
		}
		f(int(next))
		break
	}
	return nil
}

func roundUp(v int32) int32 {
	v--
	v |= v >> 1
	v |= v >> 2
	v |= v >> 4
	v |= v >> 8
	v |= v >> 16
	v++
	return v
}
