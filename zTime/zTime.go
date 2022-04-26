package zTime

import (
	"fmt"
	"time"
)

type ZTime struct {
	t   time.Time
	loc *time.Location
}

func (zt *ZTime) SetZone(name string, offset int) *ZTime {
	zt.loc = time.FixedZone(name, offset*3600)
	zt.t = zt.t.In(zt.loc)
	return zt
}

func (zt *ZTime) GetZone() *time.Location {
	if zt.loc == nil {
		zt.loc = time.FixedZone("CST", 8*3600)
	}
	return zt.loc
}

func Now() *ZTime {
	zt := ZTime{
		t: time.Now(),
	}
	zt.loc = zt.GetZone()
	return &zt
}
func New(t time.Time) *ZTime {
	zt := ZTime{
		t: time.Now(),
	}
	zt.loc = zt.GetZone()
	return &zt
}

func NewFromSeconds(s int64) *ZTime {
	return New(time.Unix(s, 0))
}

func NewFromString(layout, s string) *ZTime {
	zt := ZTime{}
	zt.loc = zt.GetZone()
	theTime, _ := time.ParseInLocation(layout, s, zt.GetZone())
	zt.t = theTime
	return &zt
}

func (zt *ZTime) Time() time.Time {
	return zt.t.In(zt.GetZone())
}

func (zt *ZTime) BeginOfDay() time.Time {
	tz := zt.t.In(zt.GetZone())
	theTime, _ := time.ParseInLocation(TimeLayoutStr,
		fmt.Sprintf("%04d-%02d-%02d 00:00:00", tz.Year(), tz.Month(), tz.Day()),
		zt.GetZone())
	return theTime
}

func (zt *ZTime) EndOfDay() time.Time {
	tz := zt.t.In(zt.GetZone())
	theTime, _ := time.ParseInLocation(TimeLayoutStr,
		fmt.Sprintf("%04d-%02d-%02d 23:59:59", tz.Year(), tz.Month(), tz.Day()),
		zt.GetZone())
	return theTime
}

func (zt *ZTime) BeginOfWeek() time.Time {
	tz := zt.BeginOfDay()
	offDay := -int(tz.Weekday()) + 1
	if tz.Weekday() == time.Sunday {
		offDay = -7 + 1
	}
	return tz.AddDate(0, 0, offDay)
}

func (zt *ZTime) EndOfWeek() time.Time {
	nextWeekBegin := zt.BeginOfWeek().AddDate(0, 0, 7)
	return time.Unix(nextWeekBegin.Unix()-1, 0)
}

func (zt *ZTime) BeginOfMonth() time.Time {
	tz := zt.t.In(zt.GetZone())
	theTime, _ := time.ParseInLocation(TimeLayoutStr,
		fmt.Sprintf("%04d-%02d-01 00:00:00", tz.Year(), tz.Month()),
		zt.GetZone())
	return theTime
}
func (zt *ZTime) EndOfMonth() time.Time {
	nextMonthBegin := zt.BeginOfMonth().AddDate(0, 1, 0)
	return time.Unix(nextMonthBegin.Unix()-1, 0)
}

func (zt *ZTime) BeginOfYear() time.Time {
	tz := zt.t.In(zt.GetZone())
	theTime, _ := time.ParseInLocation(TimeLayoutStr,
		fmt.Sprintf("%04d-01-01 00:00:00", tz.Year()),
		zt.GetZone())
	return theTime
}
func (zt *ZTime) EndOfYear() time.Time {
	nextYearBegin := zt.BeginOfYear().AddDate(1, 0, 0)
	return time.Unix(nextYearBegin.Unix()-1, 0)
}
