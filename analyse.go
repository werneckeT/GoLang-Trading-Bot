package main

import (
	"github.com/binance-exchange/go-binance"
	"math"
	"strconv"
	"time"
)

/**
Getting Data for the Coin in a given timeframe
Returns nil if there was an error, e.g. Coin does not exist
*/

var defaultCandleNumber = 10
var defaultDivergence = 0.0005

var chartAlertsList []float64

func analyseWrapper(b binance.Binance, sleepTime int, intervals []binance.Interval) {
	for true {
		coins := getCoins("config.lev")
		for _, element := range coins {
			go analyse(b, element, intervals)
			println("Started analyse for ", element)
		}
		time.Sleep(time.Duration(sleepTime))
	}
}

func analyse(b binance.Binance, coin string, intervals []binance.Interval) {
	for _, element := range intervals {
		go _analyse(b, coin, element)
	}
}
func _analyse(b binance.Binance, coin string, interval binance.Interval) {
	candles := getCandles(b, coin, interval, defaultCandleNumber)
	//Divergence in relation to the Price -> Bigger divergence for more expensive coins
	currentPrice := getCurrentPrice(b, coin)
	divergence := currentPrice * defaultDivergence * getMultiplicator(interval)

	highChartlines := getHighLines(candles, divergence)
	lowChartlines := getLowLines(candles, divergence)

	println(coin, " ", " Interval: ", interval, "\n")

	for _, chartline := range highChartlines {
		expectedResistanceHigh := chartline.M*float64(defaultCandleNumber) + chartline.N
		if currentPrice > expectedResistanceHigh+(divergence*2) {
			if !arrayContains(chartAlertsList, expectedResistanceHigh) {
				nextSMA := getNextSMA(b, currentPrice, interval, coin)
				msg := "BREAKOUT\n " + coin + " " + string(interval) + " LONG\n" +
					"\nExpected Value: $" + strconv.FormatFloat(expectedResistanceHigh, 'f', 4, 64) +
					"\nCurrent  Value: $" + strconv.FormatFloat(currentPrice, 'f', 4, 64) +
					"\n\nNext SMA Resistance: " + strconv.FormatFloat(nextSMA, 'f', 4, 64)
				newMessages = append(newMessages, msg)
				println("ADDED ALERT")
				chartAlertsList = append(chartAlertsList, expectedResistanceHigh)
				break
			}
		}
	}

	for _, chartline := range lowChartlines {
		expectedSupportResistance := chartline.M*float64(defaultCandleNumber) + 1 + chartline.N
		if currentPrice < expectedSupportResistance-divergence {
			if !arrayContains(chartAlertsList, expectedSupportResistance) {
				//msg := "BREAKOUT\n " + coin + " " + string(interval) + "\n$" + strconv.FormatFloat(currentPrice, 'f', 4, 64) + "\n\n ENTER SHORT\n" + strconv.FormatFloat(expectedSupportResistance, 'f', 4, 64)
				//println(strconv.FormatFloat(chartline.M, 'f', 4, 64), " ", strconv.FormatFloat(chartline.N, 'f', 4, 64))
				//newMessages = append(newMessages, msg)
				//chartAlertsList = append(chartAlertsList, expectedSupportResistance)
				break
			}
		}
	}
}

func getCandles(b binance.Binance, coin string, time binance.Interval, numberOfCandles int) []*binance.Kline {
	request, err := b.Klines(binance.KlinesRequest{
		Symbol:   coin,
		Interval: time,
		Limit:    numberOfCandles,
	})

	if err != nil {
		println(err)
		return nil
	}
	return request
}

func getHighLinesTest(arr []float64, divergence float64) []Chartline {
	var chartlines []Chartline

	for i := 0; i < len(arr)-2; i++ {
		for j := i + 1; j < len(arr)-1; j++ {

			m := getMovementByValues(float64(i), float64(j), arr[i], arr[j])
			n := getNByValues(float64(i), float64(j), arr[i], arr[j])

			ch := Chartline{
				M:          m,
				N:          n,
				StartPoint: i,
			}
			verificationCounter := 0
			for k := j + 1; k < len(arr); k++ {
				val := math.Abs(arr[k] - (float64(k)*m + n))
				if val < divergence {
					verificationCounter++
				}
			}
			if verificationCounter > 1 {
				chartlines = append(chartlines, ch)
			}
		}
	}
	return chartlines
}

func getHighLines(candles []*binance.Kline, divergence float64) []Chartline {
	if candles == nil {
		println("Error at Candle parsing!")
		return nil
	}
	var chartlines []Chartline

	for i := 0; i < len(candles)-3; i++ {

		for j := i + 1; j < len(candles)-2; j++ {
			verificationCounter := 0
			m := getMovementByValues(float64(i), float64(j), candles[i].High, candles[j].High)
			n := getNByValues(float64(i), float64(j), candles[i].High, candles[j].High)
			ch := Chartline{
				M:          m,
				N:          n,
				StartPoint: i,
				Points:     nil,
			}
			ch.Points = []float64{candles[i].High, candles[j].High}

			for k := j + 1; k < len(candles)-1; k++ {
				if math.Abs((m*float64(k)+n)-candles[k].High) <= divergence {
					ch.Points = append(ch.Points, candles[k].High)
					verificationCounter++
				}
				/*
					if candles[k].High > (m*float64(k)+n) && math.Abs((m*float64(k)+n)-candles[k].High) > divergence {
						//useable = false
						//k = len(candles) - 2
					} else if math.Abs((m*float64(k)+n)-candles[k].High) <= divergence {
						verificationCounter++
						ch.Points = append(ch.Points, candles[k].High)
					}*/
			}
			if verificationCounter > 1 {
				chartlines = append(chartlines, ch)
			}
		}
	}
	tmp := chartlines
	tmp = checkElementsHigh(candles, chartlines, 0)
	return tmp
}

func checkElementsHigh(candles []*binance.Kline, arr []Chartline, index int) []Chartline {
	if index >= len(arr) {
		return arr
	}
	divergence := candles[0].High * defaultDivergence
	for i := arr[index].StartPoint; i < len(candles); i++ {
		if candles[i].High-divergence > arr[index].M*float64(i)+arr[index].N {
			return checkElementsHigh(candles, removeElement(arr, index), index)
		}
	}
	println(arr[index].M, " ", arr[index].N, " ", arr[index].StartPoint)
	for _, element := range arr[index].Points {
		print(strconv.FormatFloat(element, 'f', 4, 64), " ")
	}
	println("\n", arr[index].M*float64(defaultCandleNumber)+arr[index].N)
	print(" ", strconv.FormatFloat(arr[index].M, 'f', 2, 64), " ", strconv.FormatFloat(arr[index].N, 'f', 2, 64))
	print("\n\n")
	return checkElementsHigh(candles, arr, index+1)
}

func checkElementsLow(candles []*binance.Kline, arr []Chartline, index int) []Chartline {
	if index >= len(arr) {
		return arr
	}
	divergence := candles[0].Low * defaultDivergence
	for i := arr[index].StartPoint; i < len(candles); i++ {
		if candles[i].Low+divergence < arr[index].M*float64(i)+arr[index].N {
			return checkElementsLow(candles, removeElement(arr, index), index)
		}
	}
	/*println(arr[index].M, " ", arr[index].N, " ", arr[index].StartPoint)
	for _, element := range arr[index].Points {
		print(strconv.FormatFloat(element, 'f', 4, 64), " ")
	}
	print("\n")*/
	return checkElementsLow(candles, arr, index+1)
}

func removeElement(arr []Chartline, index int) []Chartline {
	var newArr []Chartline
	if index > 0 {
		newArr = append(newArr, arr[0:index]...)
		newArr = append(newArr, arr[index+1:]...)
	} else {
		newArr = append(newArr, arr[index+1:]...)
	}
	return newArr
}

func removeElementInt(arr []int, index int) []int {
	var newArr []int
	if index > 0 {
		newArr = append(newArr, arr[0:index]...)
		newArr = append(newArr, arr[index+1:]...)
	} else {
		newArr = append(newArr, arr[index+1:]...)
	}
	return newArr
}

func getLowLines(candles []*binance.Kline, divergence float64) []Chartline {
	if candles == nil {
		println("Error at Candle parsing!")
		return nil
	}
	var chartlines []Chartline

	for i := 0; i < len(candles)-3; i++ {

		for j := i + 1; j < len(candles)-2; j++ {
			verificationCounter := 0
			m := getMovementByValues(float64(i), float64(j), candles[i].Low, candles[j].Low)
			n := getNByValues(float64(i), float64(j), candles[i].Low, candles[j].Low)

			ch := Chartline{
				M:          m,
				N:          n,
				StartPoint: i,
				Points:     nil,
			}

			ch.Points = []float64{candles[i].High, candles[j].High}

			for k := j + 1; k < len(candles)-1; k++ {
				if math.Abs((m*float64(k)+n)-candles[k].Low) <= divergence {
					verificationCounter++
					ch.Points = append(ch.Points, candles[k].Low)
				}

			}
			if verificationCounter > 1 {
				chartlines = append(chartlines, ch)
			}
		}
	}

	tmp := chartlines
	tmp = checkElementsLow(candles, chartlines, 0)
	return tmp
}

func getMovementByValues(x1 float64, x2 float64, y1 float64, y2 float64) float64 {
	return (y2 - y1) / (x2 - x1)
}

func getNByValues(x1 float64, x2 float64, y1 float64, y2 float64) float64 {
	return (x2*y1 - x1*y2) / (x2 - x1)
}

var alerts []priceAlert

func alertLoop(b binance.Binance) {
	println("Started AlertLoop!")
	for true {
		if len(alerts) > 0 {
			for index, alert := range alerts {
				if coinExists(b, alert.coin) {
					if getCurrentPrice(b, alert.coin) >= alert.price {
						newMessages = append(newMessages, alert.coin + " just crossed $" + strconv.FormatFloat(alert.price, 'f', 4, 64) + "\nAdded by: " + alert.username)
						alerts = removeAlert(alerts, index)
					}
				}
			}
		}
		time.Sleep(10000000000)
	}
}

func removeAlert(alerts []priceAlert, index int) []priceAlert {
	if index == 0 {
		return alerts[1:]
	}
	alerts[len(alerts)-1], alerts[index] = alerts[index], alerts[len(alerts) - 1]
	return alerts[:len(alerts) - 1]

}
