package main

import (
	"context"
	"github.com/binance-exchange/go-binance"
	"github.com/go-kit/kit/log"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

func createConnection(apiKey string, apiSecret string) binance.Binance {
	var logger log.Logger
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.With(logger, "time", log.DefaultTimestampUTC, "caller", log.DefaultCaller)

	hmacSigner := &binance.HmacSigner{
		Key: []byte(apiSecret),
	}
	ctx, _ := context.WithCancel(context.Background())

	binanceService := binance.NewAPIService(
		"https://www.binance.com",
		apiKey,
		hmacSigner,
		logger,
		ctx,
	)
	b := binance.NewBinance(binanceService)
	println(" => Finished Binance Authorization")
	return b
}

func pathExists(path string) bool {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func createFiles() {
	var fileName = "config.lev"
	if !pathExists(fileName) {
		os.Remove(fileName)
	}
	os.Create("config.lev")
	var insertData = "Binance API Key{123456789}\nBinance Secret Key{987654321}\nTelegram Bot Token{abcdefghiklmnopqrstuvwxyz}"
	fileWriteErr := ioutil.WriteFile(fileName, []byte(insertData), 0644)
	if fileWriteErr != nil {
		panic(fileWriteErr)
	}
}

func getData(line int, file string) string {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}
	var split = strings.Split(string(data), "\n")
	return strings.Replace(strings.Split(split[line], "{")[1], "}", "", 1)
}

func getCoins(path string) []string {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		println(err)
		return nil
	}
	tmp := strings.Split(string(data), "\n")
	return tmp[4:]
}

func addCoin(coin string) {
	if contains(getCoins("config.lev"), coin) {
		newMessages = append(newMessages, coin+" does already exist in the Database!")
	}

	f, err := os.OpenFile("config.lev", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		println(err)
		defer f.Close()
		return
	}

	defer f.Close()

	if _, err = f.WriteString("\n" + coin); err != nil {
		println(err)
	}
}

func getCurrentPrice(b binance.Binance, coin string) float64 {
	response, err := b.Ticker24(binance.TickerRequest{
		Symbol: coin,
	})

	if err != nil {
		println(err)
		return -1
	}
	return response.LastPrice
}

func arrayContains(arr []float64, element float64) bool {
	for _, el := range arr {
		if el == element {
			return true
		}
	}
	return false
}

func getMultiplicator(interval binance.Interval) float64 {
	var multiplicator float64
	switch interval {
	case binance.Minute:
		multiplicator = 1
		break
	case binance.FiveMinutes:
		multiplicator = 1.5
		break
	case binance.FifteenMinutes:
		multiplicator = 2
		break
	case binance.ThirtyMinutes:
		multiplicator = 2.75
		break
	case binance.Hour:
		multiplicator = 3
		break
	case binance.FourHours:
		multiplicator = 3.2
		break
	case binance.Day:
		multiplicator = 3.4
		break
	default:
		multiplicator = 1
		break
	}
	return multiplicator
}

func coinExists(b binance.Binance, coin string) bool {
	if getCurrentPrice(b, coin) != -1 {
		return true
	}
	return false
}

func contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}

func alertContains(alerts []priceAlert, alert priceAlert) bool {
	for i := 0; i < len(alerts); i++ {
		if alerts[i].coin == alert.coin && alerts[i].price == alert.price {
			return true
		}
	}
	return false
}

func getPinnedMessageID(path string) string {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
		return ""
	}
	return string(data)
}

func getSpotCoins(b binance.Binance) []string {
	binanceAcc, err := b.Account(binance.AccountRequest{
		RecvWindow: 5 * time.Second,
		Timestamp:  time.Now(),
	})
	if err != nil {
		panic(err)
	}
	var spotCoins []string
	for _, el := range binanceAcc.Balances {
		if !(el.Free > 0) {
			continue
		}
		if !strings.HasPrefix(el.Asset, "USDT") && !strings.HasPrefix(el.Asset, "LD") {
			if coinExists(b, el.Asset+"USDT") {
				spotCoins = append(spotCoins, el.Asset+"USDT")
			}
		}
	}
	return spotCoins
}
