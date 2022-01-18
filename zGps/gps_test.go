package zGps

import (
	"fmt"
	"testing"
)

func TestGps(t *testing.T) {
	GGALat, GGALon := 3108.15, 10404.9
	gps84Lat, gps84Lon := GpsGgaToGps(GGALat, GGALon)
	fmt.Println(gps84Lat, gps84Lon)
	gcj02Lat, gcj02Lon := Gps84ToGcj02(gps84Lat, gps84Lon)
	fmt.Println(gcj02Lat, gcj02Lon)
}
