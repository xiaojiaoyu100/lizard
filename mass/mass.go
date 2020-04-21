package mass

import (
	"math"
)

// Mass 分批
type Mass struct {
	batchSize    int
	count        int
	currentCount int
	times        int
	currentTimes int
}

// New 生成分批
func New(count, batchSize int) *Mass {
	mass := new(Mass)

	if count < 0 {
		count = 0
	}

	if batchSize <= 0 {
		batchSize = 1
	}

	mass.count = count
	mass.batchSize = batchSize
	mass.times = int(math.Ceil(float64(count) / float64(batchSize)))
	return mass
}

// Iter 迭代分批
func (mass *Mass) Iter(start, length *int) bool {
	roundSize := mass.batchSize

	if mass.currentTimes != 0 {
		*start += *length
	}

	if *start+roundSize > mass.count {
		roundSize = mass.count - *start
	}

	mass.currentCount += roundSize
	*length = roundSize

	mass.currentTimes++

	return *start < mass.count
}
