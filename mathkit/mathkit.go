package mathkit

import "math"

func Round(val float64, places int) float64 {
	pow := math.Pow(10, float64(places))
	newVal := pow * val
	trunc := math.Trunc(newVal)
	if math.Abs(newVal-trunc) >= 0.5 {
		trunc += math.Copysign(1, trunc)
	}
	round := trunc / pow
	if math.IsNaN(round) {
		return 0.0
	}
	return round
}

func IntDivRound2Digit(i, j int) float64 {
	if j == 0 {
		return 0
	}
	return Round(float64(i)/float64(j), 2)
}

func IntDivRound4Digit(i, j int) float64 {
	if j == 0 {
		return 0
	}
	return Round(float64(i)/float64(j), 4)
}

func Int64DivRound4Digit(i, j int64) float64 {
	if j == 0 {
		return 0
	}
	return Round(float64(i)/float64(j), 4)
}

func IntDiv(i, j int) float64 {
	if j == 0 {
		return 0
	}
	return float64(i) / float64(j)
}

func IntDivWithRound(i, j int, places int) float64 {
	return Round(IntDiv(i, j), places)
}

func Float64DivWithRound(i, j float64, places int) float64 {
	if j == 0 {
		return 0
	}
	return Round(i/j, places)
}

// 整数百分比
func Float64DivWithRoundPercentage(i, j float64, places int) float64 {
	if j == 0 {
		return 0
	}
	// 进行乘法运算后回因为浮点数精度造成微小的误差
	// 再进行一次Round运算可以的到精确的整数百分比
	return Round(Round(i/j, places)*100, 0)
}
