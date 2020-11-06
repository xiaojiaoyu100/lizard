package sonyflake

import (
	"fmt"
	"strconv"
	"time"
)

type ID int64

func IDFromString(s string) (ID, error) {
	id, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return ID(0), err
	}
	return ID(id), nil
}

func (i *ID) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%v\"", i)), nil
}

func (i *ID) UnmarshalJSON(b []byte) error {
	id, err := strconv.ParseInt(string(b[1:len(b)-1]), 10, 64)
	if err != nil {
		return err
	}
	*i = ID(id)
	return nil
}

func (i *ID) Int64() int64 {
	return int64(*i)
}

func (i *ID) String() string {
	return strconv.FormatInt(int64(*i), 10)
}

func (i *ID) Time() time.Time {
	return time.Unix(0, (toSonyflakeTime(DefaultStartTime)+(i.Int64()>>(BitLenSequence+BitLenMachineID)))*sonyflakeTimeUnit)
}
