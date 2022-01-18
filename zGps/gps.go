package zGps

import (
	"fmt"
	Math "math"
	"strconv"
)

/*
	WGS－84原始坐标系，一般用国际GPS纪录仪记录下来的经纬度  Google和高德地图(国外版)
	GCJ－02坐标系，又名“火星坐标系”，是我国国测局独创的坐标体系  Google和高德地图(国内版)
	百度坐标系:bd-09，百度坐标系是在GCJ－02坐标系的基础上再次加密偏移后形成的坐标系，只适用于百度地图
*/

var earthRadius float64 = 6378245.0     //地球半径（米）
var ee float64 = 0.00669342162296594323 //扁率

func TransformLat(x float64, y float64) float64 {
	ret := -100.0 + 2.0*x + 3.0*y + 0.2*y*y + 0.1*x*y + 0.2*Math.Sqrt(Math.Abs(x))
	ret += (20.0*Math.Sin(6.0*x*Math.Pi) + 20.0*Math.Sin(2.0*x*Math.Pi)) * 2.0 / 3.0
	ret += (20.0*Math.Sin(y*Math.Pi) + 40.0*Math.Sin(y/3.0*Math.Pi)) * 2.0 / 3.0
	ret += (160.0*Math.Sin(y/12.0*Math.Pi) + 320*Math.Sin(y*Math.Pi/30.0)) * 2.0 / 3.0
	return ret
}

func TransformLon(x float64, y float64) float64 {
	ret := 300.0 + x + 2.0*y + 0.1*x*x + 0.1*x*y + 0.1*Math.Sqrt(Math.Abs(x))
	ret += (20.0*Math.Sin(6.0*x*Math.Pi) + 20.0*Math.Sin(2.0*x*Math.Pi)) * 2.0 / 3.0
	ret += (20.0*Math.Sin(x*Math.Pi) + 40.0*Math.Sin(x/3.0*Math.Pi)) * 2.0 / 3.0
	ret += (150.0*Math.Sin(x/12.0*Math.Pi) + 300.0*Math.Sin(x/30.0*Math.Pi)) * 2.0 / 3.0
	return ret
}

func transform(lat float64, lon float64) (float64, float64) {
	if OutOfChina(lat, lon) {
		return lat, lon
	}
	dLat := TransformLat(lon-105.0, lat-35.0)
	dLon := TransformLon(lon-105.0, lat-35.0)
	radLat := lat / 180.0 * Math.Pi
	magic := Math.Sin(radLat)
	magic = 1 - ee*magic*magic
	SqrtMagic := Math.Sqrt(magic)
	dLat = (dLat * 180.0) / ((earthRadius * (1 - ee)) / (magic * SqrtMagic) * Math.Pi)
	dLon = (dLon * 180.0) / (earthRadius / SqrtMagic * Math.Cos(radLat) * Math.Pi)
	mgLat := lat + dLat
	mgLon := lon + dLon
	return mgLat, mgLon
}

func OutOfChina(lat float64, lon float64) bool {
	if lat < 0.8293 || lat > 55.8271 {
		return true
	}
	if lon < 72.004 || lon > 137.8347 {
		return true
	}
	return false
}

/**
 * 84 to 火星坐标系 (GCJ-02) World Geodetic System ==> Mars Geodetic System
 */

func Gps84ToGcj02(lat float64, lon float64) (float64, float64) {
	if OutOfChina(lat, lon) {
		return lat, lon
	}
	dLat := TransformLat(lon-105.0, lat-35.0)
	dLon := TransformLon(lon-105.0, lat-35.0)
	radLat := lat / 180.0 * Math.Pi
	magic := Math.Sin(radLat)
	magic = 1 - ee*magic*magic
	SqrtMagic := Math.Sqrt(magic)
	dLat = (dLat * 180.0) / ((earthRadius * (1 - ee)) / (magic * SqrtMagic) * Math.Pi)
	dLon = (dLon * 180.0) / (earthRadius / SqrtMagic * Math.Cos(radLat) * Math.Pi)
	mgLat := lat + dLat
	mgLon := lon + dLon
	return mgLat, mgLon
}

/*
 * 火星坐标系 (GCJ-02) to 84 * * @param lon * @param lat * @return
 */

func Gcj02ToGps84(lat float64, lon float64) (float64, float64) {
	mgLat, mgLon := transform(lat, lon)
	longitude := lon*2 - mgLon
	latitude := lat*2 - mgLat
	return latitude, longitude
}

// 火星坐标系 (GCJ-02) 与百度坐标系 (BD-09) 的转换算法 将 GCJ-02 坐标转换成 BD-09 坐标
var xPi float64 = Math.Pi * 3000.0 / 180.0

// 高德谷歌转为百度

func Gcj02ToBd09(lat float64, lon float64) (float64, float64) {
	x := lon
	y := lat
	z := Math.Sqrt(x*x+y*y) + 0.00002*Math.Sin(y*xPi)
	theta := Math.Atan2(y, x) + 0.000003*Math.Cos(x*xPi)
	tempLon := z*Math.Cos(theta) + 0.0065
	tempLat := z*Math.Sin(theta) + 0.006
	return tempLat, tempLon
}

// 火星坐标系 (GCJ-02) 与百度坐标系 (BD-09) 的转换算法 * * 将 BD-09 坐标转换成GCJ-02 坐标
// 百度坐标转为高德谷歌坐标

func Bd09ToGcj02(lat float64, lon float64) (float64, float64) {
	x := lon - 0.0065
	y := lat - 0.006
	z := Math.Sqrt(x*x+y*y) - 0.00002*Math.Sin(y*xPi)
	theta := Math.Atan2(y, x) - 0.000003*Math.Cos(x*xPi)
	tempLon := z * Math.Cos(theta)
	tempLat := z * Math.Sin(theta)
	return tempLat, tempLon
}

// GPS坐标转为百度坐标

func Gps84ToBd09(lat float64, lon float64) (float64, float64) {
	mgLat, mgLon := Gps84ToGcj02(lat, lon)
	return Gcj02ToBd09(mgLat, mgLon)
}

//百度坐标转成GPS坐标

func Bd09ToGps84(lat float64, lon float64) (float64, float64) {
	gcj02Lat, gcj02Lon := Bd09ToGcj02(lat, lon)
	gps84Lat, gps84Lon := Gcj02ToGps84(gcj02Lat, gcj02Lon)
	return Retain6(gps84Lat), Retain6(gps84Lon)
}

//保留小数点后六位

func Retain6(num float64) float64 {
	value, _ := strconv.ParseFloat(fmt.Sprintf("%.6f", num), 64)
	return value
}

// gga 原始数据转 gps数据

func GpsGgaToGps(lat float64, lon float64) (float64, float64) {
	tempLat := Math.Floor(lat/100) + (Math.Mod(lat, 100) / 60)
	tempLon := Math.Floor(lon/100) + (Math.Mod(lon, 100) / 60)
	return tempLat, tempLon
}

func rad(d float64) float64 {
	return d * Math.Pi / 180.0
}

func GetDistance(lat1 float64, lng1 float64, lat2 float64, lng2 float64) float64 {
	radLat1 := rad(lat1)
	radLat2 := rad(lat2)
	a := radLat1 - radLat2
	b := rad(lng1) - rad(lng2)
	s := 2 * Math.Asin(Math.Sqrt(Math.Pow(Math.Sin(a/2), 2)+
		Math.Cos(radLat1)*Math.Cos(radLat2)*Math.Pow(Math.Sin(b/2), 2)))
	return Retain6(s * earthRadius)
}
