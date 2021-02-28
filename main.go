package main

import (
	"github.com/binance-exchange/go-binance"
	_ "github.com/binance-exchange/go-binance"
	"os"
)

func main() {
	if !pathExists("config.lev") {
		createFiles()
		println("Please update your Data in config.lev")
		os.Exit(0)
	}

	var (
		apiKey    = getData(0, "config.lev")
		apiSecret = getData(1, "config.lev")
		tgKey     = getData(2, "config.lev")
	)

	if apiKey == "123456789" || apiSecret == "987654321" || tgKey == "abcdefghiklmnopqrstuvwxyz" {
		println("Please update your Data in config.lev")
		os.Exit(0)
	}
	b := createConnection(apiKey, apiSecret)

	//intervals := []binance.Interval{binance.FifteenMinutes, binance.ThirtyMinutes, binance.Hour, binance.FourHours, binance.Day}
	intervals := []binance.Interval{binance.FifteenMinutes}
	go analyseWrapper(b, 30000000000, intervals)
	//go ResistanceWrapper(b, 60000000000)

	startTelegramBot(b, tgKey)
	//getSpotCoins(b)
}
