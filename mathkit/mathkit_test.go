package mathkit

import (
	"math"
	"testing"
)

func TestRound(t *testing.T) {
	var testCases = []struct {
		num float64
		dig int
		ret float64
	}{
		{89.99, 0, 90},
		{50.555, 2, 50.56},
		{50.554, 2, 50.55},
		{0.3333, 2, 0.33},
		{0.1, 20, 0.1},
		{-1.55, 1, -1.6},
		{-1.3333, 2, -1.33},
		{math.NaN(), 0, 0},
	}
	for _, tt := range testCases {
		if ret := Round(tt.num, tt.dig); ret != tt.ret {
			t.Logf("Round(%f, %d) got %f want %f", tt.num, tt.dig, ret, tt.ret)
		}
	}
}