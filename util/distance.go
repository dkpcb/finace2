package util

import "math"

const (
	TokyoStationLat = 35.6809591
	TokyoStationLon = 139.7673068
	EarthRadius     = 6371.0 // 地球の半径 [km]
)

func CalculateDistance(lat1, lon1, lat2, lon2 float64) float64 {
	// 緯度差と経度差
	dLat := lat2 - lat1 //緯度差
	dLon := lon2 - lon1 //経度差

	// 平均緯度
	meanLat := (lat1 + lat2) / 2.0

	// 計算式
	x := dLon * math.Cos(meanLat*math.Pi/180.0)
	y := dLat
	distance := math.Sqrt(x*x+y*y) * EarthRadius * math.Pi / 180.0

	// 小数点第1位に丸める
	return math.Round(distance*10) / 10
}
