// Package backoff exposes exponential backoff algorithm with jitter.
package backoff

import (
	"errors"
	"testing"
)

func TestOn(t *testing.T) {
	invocationCounter := 0
	errFake := errors.New("fake error")
	type args struct {
		fn      func() error
		retries int
	}
	tests := []struct {
		name                string
		expectedInvocations int
		args                args
		wantErr             bool
	}{
		{"expect an error",
			1,
			args{func() error {
				invocationCounter++
				return errFake
			}, 1},
			true,
		},
		{"expect invocations",
			3,
			args{
				func() error {
					invocationCounter++
					if invocationCounter < 3 {
						return errFake
					}
					return nil
				},
				3},
			false,
		},
	}
	for _, tt := range tests {
		invocationCounter = 0
		t.Log(tt.name)
		if err := On(tt.args.fn, tt.args.retries); (err != nil) != tt.wantErr {
			t.Errorf("On() error = %v, wantErr %v", err, tt.wantErr)
		}
		if invocationCounter != tt.expectedInvocations {
			t.Errorf("expected %d invocations, got %d ", tt.expectedInvocations, invocationCounter)
		}
	}
}
