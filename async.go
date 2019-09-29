package main

import (
	"time"
)

func setTimeout(callback func(), after time.Duration) {
	ch := make(chan bool)

	go func() {
		callback()
		ch <- true
	}()

	select {
	case <-ch:
		// done

	case <-time.After(after):
		// done
	}
}
