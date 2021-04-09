package sonyflake

import (
	"strconv"
	"time"
)

// ID sonyflake id
type ID int64

// IDFromString Transfer format from string to ID.
func IDFromString(s string) (ID, error) {
	id, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return ID(0), err
	}
	return ID(id), nil
}

// MarshalText ...
func (i ID) MarshalText() ([]byte, error) {
	return []byte(strconv.FormatInt(int64(i), 10)), nil
}

// UnmarshalText ...
func (i *ID) UnmarshalText(b []byte) error {
	if len(string(b)) == 0 {
		*i = ID(0)
		return nil
	}
	id, err := strconv.ParseInt(string(b), 10, 64)
	if err != nil {
		return err
	}
	*i = ID(id)
	return nil
}

// Int64 Get int64 from ID
func (i *ID) Int64() int64 {
	return int64(*i)
}

// String Get string from ID
func (i *ID) String() string {
	return strconv.FormatInt(int64(*i), 10)
}

// Time Get time from ID
func (i *ID) Time() time.Time {
	return time.Unix(0, (toSonyflakeTime(DefaultStartTime)+(i.Int64()>>(BitLenSequence+BitLenMachineID)))*sonyflakeTimeUnit)
}
