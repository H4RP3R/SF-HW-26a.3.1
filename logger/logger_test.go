package logger

import (
	"errors"
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		target  string
		wantErr error
	}{
		{"valid param", "none", nil},
		{"valid param", "file", nil},
		{"valid param", "console", nil},
		{"invalid param", "qwerty", ErrUnsupportedDestination},
		{"empty string param", "", ErrUnsupportedDestination},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := New(tc.target)
			if !errors.Is(err, tc.wantErr) {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}
