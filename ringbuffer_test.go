package ringbuffer

import (
	"testing"
)

func Assert(t *testing.T, condition bool, reason interface{}) {
	if !condition {
		t.Fatal(reason)
	}
}

func Test_RingBuffer(t *testing.T) {
	data := make([]int, 2)
	rb := New(len(data))
	var err error

	err = rb.Put(func(i int) { data[i] = 1 })
	Assert(t, err == nil, err)

	err = rb.Put(func(i int) { data[i] = 3 })
	Assert(t, err == ErrFull, err)

	var got int
	err = rb.Get(func(i int) { got = data[i] })
	Assert(t, err == nil, err)
	Assert(t, got == 1, 1)

	err = rb.Get(func(i int) { got = data[i] })
	Assert(t, err == ErrNil, err)

}
