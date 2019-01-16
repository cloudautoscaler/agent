// Package backoff exposes exponential backoff algorithm with jitter.
package backoff

import (
	"math/rand"
	"time"
)

// Magnitude specifies the order of magnitude of the delay.
var Magnitude = time.Second

// On executes passed function applying exponential backoff with jitter on
// passed function. After specified number of attempts it will give up, returning
// the last error calling fn.
func On(fn func() error, retries int) error {
	var err error
	attempts := 0
	for attempts < retries {
		err = fn()
		if err == nil {
			return nil
		}
		attempts++
		if attempts == retries {
			break
		}
		time.Sleep(Magnitude * time.Duration(jitter(2^attempts)))
	}
	return err
}

func jitter(d int) int {
	return rand.Intn(d+1) + 1
}
