package mass

import (
	"math"
)

type Mass struct {
	batchSize    int
	count        int
	currentCount int
	times        int
	currentTimes int
}

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
