package zTime

import (
	"time"
)

// TimeLayoutStr such as time.Layout  "2006-01-02 15:04:05.000"
var TimeLayoutStr = "2006-01-02 15:04:05"

func Seconds2Time(s int64) time.Time {
	return NewFromSeconds(s).Time()
}

func Time2Seconds(t time.Time) int64 {
	return t.Unix()
}

func String2Time(s string) time.Time {
	return StringFormat2Time(TimeLayoutStr, s)
}

func StringFormat2Time(layout string, s string) time.Time {
	return NewFromString(layout, s).Time()
}

func Time2String(t time.Time) string {
	return t.Format(TimeLayoutStr)
}

func Time2FormatString(layout string, t time.Time) string {
	return t.Format(layout)
}

func Seconds2String(ts int64) string {
	return Time2String(Seconds2Time(ts))
}

func String2Seconds(s string) int64 {
	return String2Time(s).Unix()
}
