package sonyflake

import (
	"errors"
	"time"
)

// NewIdByTime generates a unique ID by time.
// After the Sonyflake time overflows, NewIdByTime returns an error.
// Less than Sonyflake min time, NewIdByTime returns an error.
func (sf *Sonyflake) NewIdByTime(t time.Time) (int64, error) {
	if t.UnixNano() < DefaultStartTime.UnixNano() {
		return 0, errors.New("min date is 2015/01/01")
	}
	const maskSequence = uint16(1<<BitLenSequence - 1)

	sf.mutex.Lock()
	defer sf.mutex.Unlock()

	current := currentElapsedTimeWithNow(t, sf.startTime)
	if sf.elapsedTime < current {
		sf.elapsedTime = current
		sf.sequence = 0
	} else { // sf.elapsedTime >= current
		sf.sequence = (sf.sequence + 1) & maskSequence
		if sf.sequence == 0 {
			sf.elapsedTime++
			overtime := sf.elapsedTime - current
			time.Sleep(sleepTimeWithNow(t, overtime))
		}
	}

	return sf.toID()
}

func sleepTimeWithNow(t time.Time, overtime int64) time.Duration {
	return time.Duration(overtime)*10*time.Millisecond -
		time.Duration(t.UTC().UnixNano()%sonyflakeTimeUnit)*time.Nanosecond
}

func currentElapsedTimeWithNow(t time.Time, startTime int64) int64 {
	return toSonyflakeTime(t) - startTime
}
