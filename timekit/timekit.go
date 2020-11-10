package timekit

import "time"

// DurationToMillis converts duration to milliseconds.
func DurationToMillis(d time.Duration) int64 {
	return int64(d / time.Millisecond)
}

// NowInMillis returns timestamp in milliseconds.
func NowInMillis() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

// NowInSecs returns timestamp in seconds.
func NowInSecs() int64 {
	return time.Now().Unix()
}

// UTCNowTime returns current time in utc.
func UTCNowTime() time.Time {
	return time.Now().UTC()
}

/*
从 octopus/async/utils/time_util.go 中迁移过来的函数
 */
// CurrentTimeStamp returns current time in utc.
func CurrentTimeStamp() int64 {
	return NowInSecs()
}

// TimeStamp2UTCTime timestamp int64 in seconds to utc time
func TimeStamp2UTCTime(t int64) time.Time {
	return time.Unix(t, 0).UTC()
}

// UTCTime2Timestamp  time.Time to timestamp in seconds
func UTCTime2Timestamp(t time.Time) int64 {
	if t.IsZero() {
		return 0
	} else {
		return t.UTC().Unix()
	}
}

// UnixMSec time.Time 转为 毫秒
func UnixMSec(t time.Time) int64 {
	return t.UnixNano() / 1e6
}

// CurrentUTCTime returns utc time.Time
func CurrentUTCTime() time.Time {
	return UTCNowTime()
}

// ZeroTime returns 当前北京时间的当天的0点0时0分的 utc的时间
func ZeroTime() time.Time {
	now := CurrentUTCTime()
	secondsEastOfUTC := int((8 * time.Hour).Seconds())
	bjZone := time.FixedZone("Beijing", secondsEastOfUTC)
	bjTime := now.In(bjZone)

	zeroBjTime := time.Date(bjTime.Year(), bjTime.Month(), bjTime.Day(), 0, 0, 0, 0, bjZone)
	return zeroBjTime.In(time.UTC)
}

// e.g. 20170109031201 2017-01-09 03:12:01
func CurrentUTCTimeStrInSecs() string {
	return CurrentUTCTime().Format("20060102150405")
}

// e.g. 20170109 2017-01-09
func CurrentUTCDateStr() string {
	return CurrentUTCTime().Format("20060102")
}

// e.g. 2017-01-09 03:12:01
func CurrentUTCDateDetailStr() string {
	return CurrentUTCTime().Format("2006-01-02 15:04:05")
}

// e.g. 2017-01-09 03:12:01
func Int642CurrentUTCDateDetailStr(t int64) string {
	tmpTime := time.Unix(t, 0)
	return tmpTime.UTC().Format("2006-01-02 15:04:05")
}

func DescBetweenTime(t1, t2 time.Time) float64 {
	return t1.Sub(t2).Seconds()
}

// e.g. 2020-11-10 00:00:00 +0800 CST
func Timestamp2DateTime(t int64) time.Time {
	tmpTime := time.Unix(t, 0)
	year, month, day := tmpTime.Date()
	dateTime := time.Date(year, month, day, 0, 0, 0, 0, tmpTime.Location())
	return dateTime
}

// MaxTime returns maxTime
func MaxTime() time.Time {
	t, _ := time.Parse("2006-01-02 15:04:05", "9999-12-31 00:00:00")
	return t
}

// CheckTimeBetween returns targetTime is between startTime and endTime
func CheckTimeBetween(targetTime, startTime, endTime time.Time) bool {
	targetStamp := UTCTime2Timestamp(targetTime)
	startStamp := UTCTime2Timestamp(startTime)
	endStamp := UTCTime2Timestamp(endTime)
	return targetStamp <= endStamp && targetStamp >= startStamp
}

// AsyncCountCostms 从 start开始算过了多少秒
func AsyncCountCostms(start time.Time) int64 {
	return time.Since(start).Nanoseconds() / 1e6
}

// BeforeMonthUtcDateSepStr
func BeforeMonthUtcDateSepStr(t time.Time, monthNum int) string {
	return t.AddDate(0, -monthNum, 0).Format("2006-01-02")
}

// AfterDayUtcDateSepStr
func AfterDayUtcDateSepStr(t time.Time, dayNum int) string {
	return t.AddDate(0, 0, dayNum).Format("2006-01-02")
}

// BeforeDayUtcDate
func BeforeDayUtcDate(t time.Time, dayNum int) time.Time {
	return t.AddDate(0, 0, -dayNum)
}

// AfterDayUtcDate
func AfterDayUtcDate(t time.Time, dayNum int) time.Time {
	return t.AddDate(0, 0, dayNum)
}