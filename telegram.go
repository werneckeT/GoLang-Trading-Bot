package main

import (
	"github.com/binance-exchange/go-binance"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"strconv"
	"strings"
	"time"
)

var newMessages []string

//var groupChatId = -1001443686647
var levitationGroup = -507027995
var ownChatId = 756435376

var spotCoins []string

func startTelegramBot(b binance.Binance, telegramToken string) {
	bot, tgBotErr := tgbotapi.NewBotAPI(telegramToken)

	if tgBotErr != nil {
		panic(tgBotErr)
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, _ := bot.GetUpdatesChan(u)

	go alertLoop(b)
	go sendAlerts(bot)
	go fetchSpotCoin(b)
	go updatePinned(b, bot)

	for update := range updates {
		if update.Message != nil {
			if update.Message.Chat.ID == int64(levitationGroup) || update.Message.Chat.ID == int64(ownChatId) {
				if strings.HasPrefix(update.Message.Text, "/price ") || strings.HasPrefix(update.Message.Text, "/price@") {
					if strings.Contains(update.Message.Text, " ") {
						if len(strings.Split(update.Message.Text, " ")) > 1 {
							priceFloat := getCurrentPrice(b, strings.Split(update.Message.Text, " ")[1])
							if priceFloat < 0 {
								errorMessage := "Error! Coin does not exist! $" + strings.Split(update.Message.Text, " ")[1]
								bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, errorMessage))
							} else {
								currentPrice := "$" + strconv.FormatFloat(priceFloat, 'f', 4, 64)
								bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, currentPrice))
							}
						} else {
							errorMessage := "Error! Usage: /price <Coin>"
							bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, errorMessage))
						}
					}
					continue
				}

				if strings.HasPrefix(update.Message.Text, "/add ") || strings.HasPrefix(update.Message.Text, "/add@") {
					if update.Message.From.UserName == "whoisazer" {
						var coin = strings.Split(update.Message.Text, " ")[1]
						if coinExists(b, coin) {
							addCoin(coin)
						} else {
							errorMessage := "Error! Coin does not exist! $" + coin
							bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, errorMessage))
						}
					} else {
						errorMessage := "no. :))))))))"
						bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, errorMessage))
					}
					continue
				}

				if strings.HasPrefix(update.Message.Text, "/watch ") || strings.HasPrefix(update.Message.Text, "/watch@") {
					if strings.Contains(update.Message.Text, " ") {
						coin := strings.Split(update.Message.Text, " ")[1]
						if coinExists(b, coin) {
							if !contains(spotCoins, coin) {
								spotCoins = append(spotCoins, coin)
							} else {
								errorMessage := "Error! This Coin is already displayed!"
								bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, errorMessage))
							}
						}else{
							errorMessage := "Error! Coin does not exist! $" + coin
							bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, errorMessage))
						}
					}
				}

				if strings.HasPrefix(update.Message.Text, "/setalert ") || strings.HasPrefix(update.Message.Text, "/setalert@") {
					str := strings.Split(update.Message.Text, " ")
					if len(str) != 3 {
						errorMessage := "Error! Usage: /setalert <Coin> <Price>"
						bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, errorMessage))
					} else {
						if len(str[1]) < 17 {
							priceFloat, _ := strconv.ParseFloat(str[2], 4)
							if coinExists(b, str[1]) {
								if priceFloat > getCurrentPrice(b, str[1]) {
									newAlert := priceAlert{
										coin:       str[1],
										price:      priceFloat,
										candleType: "unused",
										username:   update.Message.From.UserName,
									}
									if alertContains(alerts, newAlert) {
										errorMessage := "Error! This Alert does already exist!"
										bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, errorMessage))
									} else {
										alerts = append(alerts, newAlert)
										successMessage := "Alert for " + str[1] + " at $" + str[2] + " has been added!"
										bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, successMessage))
									}
								} else {
									errorMessage := "Error! The Coin is already above the Alert!"
									bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, errorMessage))
								}
							} else {
								errorMessage := "Error! This Coin does not exist!"
								bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, errorMessage))
							}
						}
					}
					continue
				}
			}
		}
	}
}

func fetchSpotCoin(b binance.Binance) {
	for true {
		spotCoins = getSpotCoins(b)
		time.Sleep(time.Duration(600000000000))
	}
}

func sendAlerts(bot *tgbotapi.BotAPI) {
	for true {
		for len(newMessages) > 0 {
			msg := newMessages[0]
			newMessages = newMessages[1:]

			sendMessage := tgbotapi.NewMessage(int64(ownChatId), msg)
			bot.Send(sendMessage)
		}
		time.Sleep(time.Duration(1000000000))
	}
}

func updatePinned(b binance.Binance, bot *tgbotapi.BotAPI) {
	/*
	if !pathExists("messages.lev") {
		os.Create("messages.lev")
		pinnedMessage := tgbotapi.NewMessage(int64(levitationGroup), "TO UPADTE")
		msg, err := bot.Send(pinnedMessage)
		if err != nil {
			println(err)
			return
		}
		fileWriteErr := ioutil.WriteFile("messages.lev", []byte(strconv.FormatInt(int64(msg.MessageID), 10)), 0644)
		if fileWriteErr != nil {
			panic(fileWriteErr)
		}
	}
	coinIndex := 0
	for true {
		currentCoins := getCoins("config.lev")
		coins := currentCoins

		for _, c := range spotCoins {
			if !contains(coins, c) {
				coins = append(coins, c)
			}
		}

		newPinned := ""
		for i := 0; i < 3; i++ {
			if coinIndex+i == len(coins) {
				coinIndex = 0
			}
			newPinned += getCoinString(b, coins[coinIndex+i])
			if i != 2 {
				newPinned += " | "
			}
		}

		msgID, _ := strconv.ParseInt(getPinnedMessageID("messages.lev"), 10, 0)
		messageConfig := tgbotapi.NewEditMessageText(int64(levitationGroup), int(msgID), newPinned)
		bot.Send(messageConfig)
		coinIndex += 3
		time.Sleep(time.Duration(5000000000))
	}*/
}

func getCoinString(b binance.Binance, coin string) string {
	return coin[0:len(coin)-4] + ": $" + strconv.FormatFloat(getCurrentPrice(b, coin), 'f', 2, 64)
}
