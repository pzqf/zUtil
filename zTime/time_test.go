package zTime

import (
	"fmt"
	"testing"
)

func Test(t *testing.T) {

	fmt.Println("Now Seconds:", Now().Time().Unix())
	fmt.Println("Now String:", Time2String(Now().Time()))

	fmt.Println("now Nano:", Now().Time().UnixNano())
	fmt.Println("GetMillisecond:", Now().Time().UnixMilli())

	fmt.Println("SecondsToTime:", Seconds2Time(Now().Time().Unix()))
	fmt.Println("GetTodaySeconds:", int(Now().Time().Sub(Now().BeginOfDay()).Seconds()))
	fmt.Println("GetTodaySeconds:", int(Now().Time().Sub(Now().SetZone("dd", -3).BeginOfDay()).Seconds()))

	fmt.Println("Seconds2String:", Seconds2String(Now().Time().Unix()))

	fmt.Println("GetThisMonthRestDays:", int(Now().EndOfMonth().Sub(Now().Time()).Seconds())/24/60/60)
	fmt.Println("GetThisYearRestDays:", int(Now().EndOfYear().Sub(Now().Time()).Seconds())/24/60/60)

	fmt.Println("BeginOfDay:", Now().BeginOfDay())
	fmt.Println("EndOfDay:", Now().EndOfDay())

	fmt.Println("BeginOfWeek:", Now().BeginOfWeek())
	fmt.Println("EndOfWeek:", Now().EndOfWeek())

	fmt.Println("BeginOfMonth:", Now().BeginOfMonth())
	fmt.Println("EndOfMonth:", Now().EndOfMonth())

	fmt.Println("BeginOfYear:", Now().BeginOfYear())
	fmt.Println("EndOfYear:", Now().EndOfYear())
}
