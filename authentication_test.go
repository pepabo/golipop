package lolp

import "testing"

func TestMask(t *testing.T) {
	tests := []struct {
		in  string
		out string
	}{
		{"", "***[masked]"},
		{"helloworld!", "hel***[masked]"},
		{"yo", "***[masked]"},
	}
	for _, tt := range tests {
		s := mask(tt.in)
		if s != tt.out {
			t.Errorf("expected %s to be %s", tt.out, s)
		}
	}
}
