package Features

import (
	"github.com/binance-exchange/go-binance"
	"math"
)

type resistanceLevel struct {
	priceLevel	float64
	interval    binance.Interval
	verificationCounter int
}

func GetResistances(candles []*binance.Kline, coin string, interval binance.Interval, divergence float64) []resistanceLevel {
	var levels []resistanceLevel
	for i := 0; i < len(candles)-3; i++ {

		verificationCounter := []int{0, 0, 0}
		val := candles[i].Close
		valHigh := candles[i].High
		valLow := candles[i].Low

		for j := i + 1; j < len(candles); j++ {
			//for Close Value
			verificationCounter[0] += checkValues(val, candles[j], divergence)
			//for candleHigh
			verificationCounter[1] += checkValues(valHigh, candles[j], divergence)
			//for candleLow
			verificationCounter[1] += checkValues(valLow, candles[j], divergence)
		}
		maxIndex := maxIndex(verificationCounter)
		if verificationCounter[maxIndex]  > 4 {
			var res = resistanceLevel{
				priceLevel:          0,
				interval:            interval,
				verificationCounter: verificationCounter[maxIndex],
			}
			switch maxIndex {
			case 0:
				res.priceLevel = val
				break
			case 1:
				res.priceLevel = valHigh
				break
			case 2:
				res.priceLevel = valLow
				break
			default:
				println("WASN HIER LOS")
				break
			}
			levels = append(levels, res)
		}
	}
	return levels
}

func checkValues(value float64, candle *binance.Kline, divergence float64) int{
	if math.Abs(value - candle.Close) < divergence {
		return 1
	}else if math.Abs(value - candle.High) < divergence {
		return 1
	}else if math.Abs(value - candle.Low) < divergence {
		return 1
	}
	return 0
}

func maxIndex(arr []int) int {
	max := arr[0]
	maxIndex := 0
	for i := 1; i < len(arr); i++ {
		if arr[i] > max {
			max = arr[i]
			maxIndex = i
		}
	}
	return maxIndex
}