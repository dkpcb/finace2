package util

import (
	"math"
	"testing"
)

func TestCalculateDistance(t *testing.T) {
	tests := []struct {
		name                   string
		lat1, lon1, lat2, lon2 float64
		expected               float64
	}{
		{"東京駅から岐阜市柳津町までの距離", TokyoStationLat, TokyoStationLon, 35.355743, 136.725408, 278.3},
		{"東京駅から東京駅まで（自己距離）", TokyoStationLat, TokyoStationLon, TokyoStationLat, TokyoStationLon, 0.0},
		{"東京駅から大阪駅までの距離", TokyoStationLat, TokyoStationLon, 34.702485, 135.495951, 403.3},
		{"東京駅から札幌駅までの距離", TokyoStationLat, TokyoStationLon, 43.068661, 141.350755, 831.7},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalculateDistance(tt.lat1, tt.lon1, tt.lat2, tt.lon2)

			// 許容誤差を広げてテスト
			diff := math.Abs(got - tt.expected)
			if diff > 1.0 { // 許容誤差を 1.0 km に設定
				t.Errorf("CalculateDistance() = %v, want %v (diff %v)", got, tt.expected, diff)
			}
		})
	}
}
