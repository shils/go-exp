package channels

import (
	"gotest.tools/v3/assert"
	"testing"
)

func TestSendUntilBlocked(t *testing.T) {
	ch := make(chan int, 1)
	in := []int{5, 4, 3, 2, 1}
	n := SendUntilBlocked(ch, in)
	assert.Equal(t, 1, n)
	assert.Equal(t, 5, <-ch)
}

func TestSend(t *testing.T) {
	ch := make(chan int, 1)
	in := []int{5, 4, 3, 2, 1}

	var received []int

	go func() {
		for i := range ch {
			received = append(received, i)
		}
	}()
	n := Send(ch, in)
	assert.Equal(t, 5, n)
	assert.DeepEqual(t, in, received)
}
