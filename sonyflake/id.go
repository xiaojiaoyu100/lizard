package sonyflake

import (
	"strconv"
	"time"
)

type ID int64

func (i *ID) MarshalJSON() ([]byte, error) {
	b := []byte(strconv.FormatInt(int64(*i), 10))
	return b, nil
}

func (i *ID) UnmarshalJSON(b []byte) error {
	id, err := strconv.ParseInt(string(b), 10, 64)
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

func (i *ID) GetTimeByID() time.Time {
	return time.Unix(0, (toSonyflakeTime(DefaultStartTime)+(i.Int64()>>(BitLenSequence+BitLenMachineID)))*sonyflakeTimeUnit)
}
