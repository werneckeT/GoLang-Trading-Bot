package main

import (
	"github.com/binance-exchange/go-binance"
	"time"
)

type Chartline struct {
	M          float64
	N          float64
	StartPoint int
	timeframe  binance.Interval
	Points     []float64
}

type resistanceLevel struct {
	price      float64
	breakTests int
}

type chartAlert struct {
	coin           string
	candleOpenTime time.Time
	timeInterval   binance.Interval
	targetPrice    float64
	breakout       string
}

type priceAlert struct {
	coin       string
	price      float64
	candleType string
	username   string
}
