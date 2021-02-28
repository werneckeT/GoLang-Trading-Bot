package main

import (
	"github.com/binance-exchange/go-binance"
	"math"
	"strconv"
	"time"
)
var defaultNumberOfCandels = 75
func ResistanceWrapper(b binance.Binance, sleepTime int) {
	intervals := []binance.Interval{binance.Hour, binance.FourHours, binance.Day}
	for true {
		for _, el := range intervals {
			coins := getCoins("config.lev")
			for _, elm := range coins {
				go _ResistanceWrapper(b, elm, el)
			}
		}
		time.Sleep(time.Duration(sleepTime))
	}
}

func _ResistanceWrapper(b binance.Binance, coin string, interval binance.Interval) {
	resistances := findHighReistances(b, coin, interval, defaultNumberOfCandels)
	currentPrice := getCurrentPrice(b, coin)

	for _, el := range resistances {
		if math.Abs((el.price/ currentPrice) - 100) < 3 {
			msg := "RESISTANCE\nPossible resistance for Coin: " + coin + " @" + strconv.FormatFloat(el.price, 'f', 4, 64) + "\nCurrent Price: " +
				strconv.FormatFloat(currentPrice, 'f', 4, 64)
			newMessages = append(newMessages, msg)
		}
	}
}

func findHighReistances(b binance.Binance, coin string, interval binance.Interval, numberOfCandles int) []resistanceLevel {
	candles := getCandles(b, coin, interval, numberOfCandles)
	divergence := 0.0
	var resistances []resistanceLevel
	for i := 0; i < len(candles); i++ {
		resistanceCounter := 0
		for j := 0; j < len(candles); j++ {
			if j == i {
				continue
			}
			if closeTo(candles[i].Close, candles[j].Close, divergence) {
				resistanceCounter++
			} else if closeTo(candles[i].Close, candles[j].Low, divergence) {
				resistanceCounter++
			} else if closeTo(candles[i].Close, candles[j].High, divergence) {
				resistanceCounter++
			}
		}
		if resistanceCounter > 3 {
			resLevel := resistanceLevel{
				price:      candles[i].Close,
				breakTests: resistanceCounter,
			}
			resistances = append(resistances, resLevel)
		}
	}
	return resistances
}

func closeTo(base float64, value float64, divergence float64) bool {
	if math.Abs(base-value) <= divergence {
		return true
	}
	return false
}