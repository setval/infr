package math

import "math"

func deg2rad(deg float32) float64 {
	return float64(deg * (math.Pi / 180))
}

func FindDistanceBeetwenCities(lat1, lon1, lat2, lon2 float32) int {
	var R = 6371.210                // Radius of the earth in km
	var dLat = deg2rad(lat2 - lat1) // deg2rad below
	var dLon = deg2rad(lon2 - lon1)
	var a = math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(deg2rad(lat1))*math.Cos(deg2rad(lat2))*
			math.Sin(dLon/2)*math.Sin(dLon/2)
	var c = 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	var d = float64(R) * c // Distance in km
	return int(d)
}
