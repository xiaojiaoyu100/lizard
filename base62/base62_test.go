package base62

import (
	"math"
	"testing"
)

func TestEncode(t *testing.T) {
	tests := [...]struct {
		In   int64
		Want string
	}{
		0: {
			0,
			"0",
		},
		1: {
			1,
			"1",
		},
		2: {
			61,
			"z",
		},
		3: {
			62,
			"10",
		},
		4: {
			math.MaxInt64,
			"AzL8n0Y58m7",
		},
	}
	for _, test := range tests {
		got := Encode(test.In)
		if got != test.Want {
			t.Errorf("value: %d, want: %s, got: %s", test.In, test.Want, got)
		}
	}
}

func TestDecode(t *testing.T) {
	tests := [...]struct {
		In   string
		Want int64
	}{
		0: {

			"0",
			0,
		},
		1: {
			"1",
			1,
		},
		2: {
			"z",
			61,
		},
		3: {
			"10",
			62,
		},
		4: {
			"AzL8n0Y58m7",
			math.MaxInt64,
		},
	}
	for _, test := range tests {
		got := Decode(test.In)
		if got != test.Want {
			t.Errorf("value: %s, want: %d, got: %d", test.In, test.Want, got)
		}
	}
}
