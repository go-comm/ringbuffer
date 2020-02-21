package ringbuffer

import "testing"

func Test_RingBuffer(t *testing.T) {
	rb := New(3)

	t.Log(rb.Get())

	t.Log(rb.Put("1"))
	t.Log(rb.Put("2"))
	t.Log(rb.Put("3"))
	t.Log(rb.Put("4"))

	t.Log(rb.Get())
	t.Log(rb.Get())
	t.Log(rb.Get())
}
