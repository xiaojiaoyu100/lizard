// Package sonyflake implements Sonyflake, a distributed unique ID generator inspired by Twitter's Snowflake.
//
// A Sonyflake ID is composed of
//     39 bits for time in units of 10 msec
//     10 bits for a sequence number (102400num/s)
//     14 bits for a machine id (0 - 16383)
package sonyflake

import (
	"errors"
	"strconv"
	"sync"
	"time"
)

// These constants are the bit lengths of Sonyflake ID parts.
const (
	BitLenTime      = 39                               // bit length of time
	BitLenSequence  = 10                               // bit length of sequence number
	BitLenMachineID = 63 - BitLenTime - BitLenSequence // bit length of machine id

	sonyflakeTimeUnit = 1e7 // nsec, i.e. 10 msec
)

var DefaultStartTime = time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC)

// Sonyflake is a distributed unique ID generator.
type Sonyflake struct {
	mutex       *sync.Mutex
	startTime   int64
	elapsedTime int64
	sequence    uint16
	machineID   uint16
}

// NewSonyflake returns a new Sonyflake configured with the given Settings.
func NewSonyflake(machineID uint16) (*Sonyflake, error) {
	maxMachineID := -1 ^ (-1 << BitLenMachineID)
	if machineID > uint16(maxMachineID) {
		return nil, errors.New("Node number must be between 0 and " + strconv.FormatInt(int64(maxMachineID), 10))
	}

	sf := new(Sonyflake)
	sf.mutex = new(sync.Mutex)
	sf.sequence = uint16(1<<BitLenSequence - 1)
	sf.startTime = toSonyflakeTime(DefaultStartTime)
	sf.machineID = machineID

	return sf, nil
}

// NextID generates a next unique ID.
// After the Sonyflake time overflows, NextID returns an error.
func (sf *Sonyflake) NextID() (ID, error) {
	const maskSequence = uint16(1<<BitLenSequence - 1)

	sf.mutex.Lock()
	defer sf.mutex.Unlock()

	current := currentElapsedTime(sf.startTime)
	if sf.elapsedTime < current {
		sf.elapsedTime = current
		sf.sequence = 0
	} else { // sf.elapsedTime >= current
		sf.sequence = (sf.sequence + 1) & maskSequence
		if sf.sequence == 0 {
			sf.elapsedTime++
			overtime := sf.elapsedTime - current
			time.Sleep(sleepTime(overtime))
		}
	}

	return sf.toID()
}

func GetTimeByID(id ID) time.Time {
	return time.Unix(0, (toSonyflakeTime(DefaultStartTime)+(id.ToInt64()>>(BitLenSequence+BitLenMachineID)))*sonyflakeTimeUnit)
}

func toSonyflakeTime(t time.Time) int64 {
	return t.UTC().UnixNano() / sonyflakeTimeUnit
}

func currentElapsedTime(startTime int64) int64 {
	return toSonyflakeTime(time.Now()) - startTime
}

func sleepTime(overtime int64) time.Duration {
	return time.Duration(overtime)*10*time.Millisecond -
		time.Duration(time.Now().UTC().UnixNano()%sonyflakeTimeUnit)*time.Nanosecond
}

func (sf *Sonyflake) toID() (ID, error) {
	if sf.elapsedTime >= 1<<BitLenTime {
		return 0, errors.New("over the time limit")
	}

	return ID(uint64(sf.elapsedTime)<<(BitLenSequence+BitLenMachineID) |
		uint64(sf.sequence)<<BitLenMachineID |
		uint64(sf.machineID)), nil
}

type IDComposition struct {
	ID        int64 `json:"id"`
	Msb       int64 `json:"msb"`
	Time      int64 `json:"time"`
	Sequence  int64 `json:"sequence"`
	MachineId int64 `json:"machine_id"`
}

// Decompose returns a set of Sonyflake ID parts.
func Decompose(id int64) IDComposition {
	const maskSequence = (1<<BitLenSequence - 1) << BitLenMachineID
	const maskMachineID = 1<<BitLenMachineID - 1

	msb := id >> 63
	decomposeTime := id >> (BitLenSequence + BitLenMachineID)
	sequence := id & maskSequence >> BitLenMachineID
	machineID := id & maskMachineID
	return IDComposition{
		ID:        id,
		Msb:       msb,
		Time:      decomposeTime,
		Sequence:  sequence,
		MachineId: machineID,
	}
}
