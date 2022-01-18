package zTime

import (
	"fmt"
	"strconv"
	"time"
)

var TimeLayoutStr = "2006-01-02 15:04:05"

func GetNowSeconds() int64 {
	return time.Now().Unix()
}

func GetNano(t time.Time) int64 {
	return t.UnixNano()
}

func GetMillisecond(t time.Time) int64 {
	return t.UnixNano() / 1e6
}

func GetNowString() string {
	return Time2String(time.Now())
}

func GetToday0ClockSeconds() int64 {
	ts := GetTodaySeconds()
	sStr := strconv.Itoa(ts)
	s, _ := strconv.ParseInt(sStr, 10, 64)
	return GetNowSeconds() - s
}

func GetFormatYMDHMS(t time.Time) string {
	return t.Format(TimeLayoutStr)
}

func SecondsToTime(s int64) time.Time {
	return time.Unix(s, 0)
}

func GetTodaySeconds() int {
	return time.Now().Hour()*60*60 + time.Now().Minute()*60 + time.Now().Second()
}

func GetLastWeek0ClockSeconds() int64 {
	n := time.Now()
	s := GetToday0ClockSeconds()
	t := SecondsToTime(s)
	return t.AddDate(0, 0, 1-int(n.Weekday())).Unix()
}

func GetLastMonth0ClockSeconds() int64 {
	n := time.Now()
	str := fmt.Sprintf("%04d-%02d-%02d 00:00:00", n.Year(), n.Month(), n.Day())
	t := String2Time(str)
	return t.AddDate(0, 0, 1-n.Day()).Unix()
}

func GetNextMonth0ClockSeconds() int64 {
	n := time.Now()
	str := fmt.Sprintf("%04d-%02d-%02d 00:00:00", n.Year(), n.Month(), n.Day())
	t := String2Time(str)
	return t.AddDate(0, 1, 1-n.Day()).Unix()
}

func GetThisMonthRestDays() int {
	s := GetNextMonth0ClockSeconds() - time.Now().Unix()
	return int(s / 60 / 60 / 24)
}

func String2Time(s string) time.Time {
	loc, _ := time.LoadLocation("Local")
	theTime, _ := time.ParseInLocation(TimeLayoutStr, s, loc) //使用模板在对应时区转化为time.time类型

	return theTime
}
func StringFormat2Time(formatStr string, s string) time.Time {
	loc, _ := time.LoadLocation("Local")
	theTime, _ := time.ParseInLocation(formatStr, s, loc) //使用模板在对应时区转化为time.time类型

	return theTime
}

func Time2String(t time.Time) string {
	return t.Format(TimeLayoutStr)
}

func Time2FormatString(formatStr string, t time.Time) string {
	return t.Format(formatStr)
}

func Timestamp2String(ts int64) string {
	return Time2String(SecondsToTime(ts))
}

func Timestamp2Time(ts int64) time.Time {
	return time.Unix(ts, 0)
}
