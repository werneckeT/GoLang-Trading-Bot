package Features

import (
	"github.com/binance-exchange/go-binance"
	"math"
)

//Three Line Strike
//Hammer (done)
//Hanging Man (done)
//Shooting Star (done)

func candlepatternsHammer(candles []*binance.Kline, divergence float64) bool {
	if isDownTrend(candles, 3) {
		if math.Abs(candles[len(candles)-1].Close-candles[len(candles)-2].Close) <= divergence {
			if math.Abs(candles[len(candles)-1].Close-candles[len(candles)-1].High) <= divergence {
				return true
			}
		}
	}
	return false
}

func candlepatternsHangingMan(candles []*binance.Kline, divergence float64) bool {
	if isUpTrend(candles, 3) {
		if math.Abs(candles[len(candles)-1].High-candles[len(candles)-1].Close) <= divergence {
			if candles[len(candles)-1].Low < candles[len(candles)-2].Low {
				return true
			}
		}
	}
	return false
}

func shootingStar(candles []*binance.Kline, divergence float64) bool {
	if isUpTrend(candles, 3) {
		if candles[len(candles)-1].Close > candles[len(candles)-2].Close {
			if candles[len(candles)-1].High > candles[len(candles)-2].High {
				if (candles[len(candles)-1].High - candles[len(candles)-1].Close) > (candles[len(candles)-2].High - candles[len(candles)-2].Close) {
					if (candles[len(candles)-1].Close - candles[len(candles)-1].Open) * 4 < candles[len(candles)-1].High - candles[len(candles)-1].Close {
						return true
					}
				}
			}
		}
	}
	return false
}

func isUpTrend(candles []*binance.Kline, numberOfCandlesToUse int) bool {
	if len(candles) <= numberOfCandlesToUse {
		numberOfCandlesToUse = len(candles) - 1
		println("WARNING WRONG LENGTH AT UPTRENDCHECK")
	}

	for i:= numberOfCandlesToUse - 1; i >= 0; i-- {
		if candles[len(candles) - 1 - numberOfCandlesToUse].High > candles[len(candles) - numberOfCandlesToUse].High {
			return false
		}

		if candles[len(candles) - 1 - numberOfCandlesToUse].Close > candles[len(candles) - numberOfCandlesToUse].Close {
			return false
		}
	}
	return true
}

func isDownTrend(candles []*binance.Kline, numberOfCandlesToUse int) bool {

	if len(candles) <= numberOfCandlesToUse {
		numberOfCandlesToUse = len(candles) - 1
		println("WARNING WRONG LENGTH AT DOWNTRENDCHECK")
	}


	for i:= numberOfCandlesToUse - 1; i >= 0; i-- {
		if candles[len(candles) - 1 - numberOfCandlesToUse].High < candles[len(candles) - numberOfCandlesToUse].High {
			return false
		}

		if candles[len(candles) - 1 - numberOfCandlesToUse].Close < candles[len(candles) - numberOfCandlesToUse].Close {
			return false
		}

		if candles[len(candles) - 1 - numberOfCandlesToUse].Low > candles[len(candles) - numberOfCandlesToUse].Low {
			return false
		}
	}
	return true
}
