package Features

import (
	"math"
)

func GetFibRetracements(x1 float64, x2 float64, d int) []float64 {
	//standard := (1 + math.Sqrt(5)) / 2
	dist := math.Abs(x1 - x2)
	dist = math.Sqrt(dist * dist + float64(d * d))
	retracements := []float64{0.236, 0.382, 0.5, 0.618, 0.764, 1, 1.382}
	extensions := []float64{0.618, 1.382, 1.618, 2, 2.618}

	var answer []float64
	if x1 > x2 {
		//Downtrend
		for _, el := range retracements {
			answer = append(answer, x2 + (dist * el))
		}
		for _, el := range extensions {
			answer = append(answer, x2 - (dist * el))
		}
	} else {
		//Uptrend
		for _, el := range retracements {
			answer = append(answer, x1 - (dist * el))
		}
		for _,el := range extensions {
			answer = append(answer, x1 + (dist * el))
		}
	}
	return answer
}
