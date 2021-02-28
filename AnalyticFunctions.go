package main

import (
	"github.com/binance-exchange/go-binance"
	"math"
)

//Simple Moving Average
// (a1 + ...  + an) / n
func SMA(candles []*binance.Kline, numberOfCandles int) float64 {
	close := 0.0
	for i := len(candles) - 1; i > -1; i-- {
		close += (candles[i].Close + candles[i].High + candles[i].Low) / 3
	}
	//Return SMA
	return close / float64(numberOfCandles)
}

func BOL(b binance.Binance, coin string, interval binance.Interval) (float64, float64, float64) {
	//20 should be also enough
	defaultNumberOfCandels := 20
	candles := getCandles(b, coin, interval, defaultNumberOfCandels)
	sma := SMA(candles, defaultNumberOfCandels)

	if sma == -1 {
		panic("FATAL BOL ERROR")
	}

	standardDivergence := 0.0
	for i := len(candles) - 1; i > -1; i-- {
		tmp := (candles[i].Close + candles[i].High + candles[i].Low) / 3
		standardDivergence += 2 * math.Abs(sma-tmp)
	}
	//Boolinger Band return
	div := standardDivergence / float64(defaultCandleNumber)
	return sma + div, sma, sma - div
}

func getNextSMA(b binance.Binance, currentPrice float64, interval binance.Interval, coin string) float64 {
	candles := getCandles(b, coin, interval, 240)
	sma9 := SMA(candles, 9)
	sma24 := SMA(candles, 24)
	sma60 := SMA(candles, 60)
	sma240 := SMA(candles, 240)

	smallestDist := 1000000000.0
	smallest := -1.0

	if sma9-currentPrice > 0 && sma9-currentPrice < smallestDist {
		smallestDist = sma9 - currentPrice
		smallest = sma9
	}

	if sma24 - currentPrice > 0 && sma24 - currentPrice < smallestDist {
		smallestDist = sma24 - currentPrice
		smallest = sma24
	}

	if sma60 - currentPrice > 0 && sma60 - currentPrice < smallestDist {
		smallestDist = sma60 - currentPrice
		smallest = sma60 - currentPrice
	}

	if sma240 - currentPrice > 0 && sma240 - currentPrice < smallestDist {
		return sma240
	}

	return smallest
}

/*
func EMA(candles []*binance.Kline, currentPrice float64, numberOfCandles int) float64 {
	candles = candles[len(candles) - numberOfCandles:]
	ema := currentPrice * (2 / (float64(numberOfCandles) + 1))

	for i := 0; i < numberOfCandles; i++ {

	}
}

func _EMA(candles []*binance.Kline, numberOfCandles int, currentN int) float64 {

}*/
