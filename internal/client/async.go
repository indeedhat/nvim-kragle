package client

import (
	"time"
)

// race for either the callback to finish or the time limit to be reached
func race(callback func(), timeLimit time.Duration) {
	ch := make(chan bool)

	go func() {
		callback()
		ch <- true
	}()

	select {
	case <-ch:
		// done

	case <-time.After(timeLimit):
		// done
	}
}
