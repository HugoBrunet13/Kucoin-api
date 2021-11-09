package main

import (
	// "io"
	"flag"
	"log"
	"time"

	// "net/http"
	"github.com/Kucoin/kucoin-go-sdk"
	"github.com/spf13/viper"
)

var session *kucoin.ApiService

func placeOrder(symbol string, side string, orderType string, size string, price string, delay bool) (string, error) {
	log.Printf("PlaceOrder: Sym: %s Side: %s Type: %s Size %s Price: %s Delay: %t", symbol, side, orderType, size, price, delay)
	params := &kucoin.CreateOrderModel{
		ClientOid: kucoin.IntToString(time.Now().UnixNano()),
		Side:      side,
		Symbol:    symbol,
		Type:      orderType,
		Size:      size,
	}
	if orderType != "market" {
		params.Price = price
	}

	if delay {
		targetTime := viper.GetString("targetTime")
		for {
			now := time.Now()
			nowS := now.Format("2006-01-02 15:04:05")
			log.Printf("Wait....")
			if nowS >= targetTime {
				break
			}
		}
	}

	rsp, err := session.CreateOrder(params)
	if err != nil {
		return "", err
	}

	orderResult := &kucoin.CreateOrderResultModel{}
	errParse := rsp.ReadData(orderResult)
	if errParse != nil {
		return "", errParse
	}
	return orderResult.OrderId, nil

}

func loadConfig(config_path string) {
	log.Printf("load_config config_path=%v\n", config_path)
	viper.SetConfigName(config_path)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		log.Printf("Fatal error config file: %s\n", err)
	}
}

func init() {
	var config_path string
	flag.StringVar(&config_path, "config", "config.yaml", "Config file")
	flag.Parse()
	loadConfig(config_path)

	apiKey := viper.GetString("key")
	apiSecret := viper.GetString("secret")
	passphrase := viper.GetString("passphrase")

	session = kucoin.NewApiService(
		kucoin.ApiKeyOption(apiKey),
		kucoin.ApiSecretOption(apiSecret),
		kucoin.ApiPassPhraseOption(passphrase),
	)
}

func main() {

	// BUY --> Delay
	buyOrder, err := placeOrder("SHIB-USDT", "buy", "limit", "10000", "0.0001", true)
	if err != nil {
		log.Printf("Failed to place order: %s", err)
	} else {
		log.Printf("Order placed! Order-id: %s", buyOrder)
	}

	//SELL 1
	sellOrder1, err := placeOrder("SHIB-USDT", "sell", "limit", "10000", "0.0001", false)
	if err != nil {
		log.Printf("Failed to place order: %s", err)
	} else {
		log.Printf("Order placed! Order-id: %s", sellOrder1)
	}

	// SELL 2
	sellOrder2, err := placeOrder("SHIB-USDT", "sell", "limit", "10000", "0.0001", false)
	if err != nil {
		log.Printf("Failed to place order: %s", err)
	} else {
		log.Printf("Order placed! Order-id: %s", sellOrder2)
	}

	// SELL 3
	sellOrder3, err := placeOrder("SHIB-USDT", "sell", "limit", "10000", "0.0001", false)
	if err != nil {
		log.Printf("Failed to place order: %s", err)
	} else {
		log.Printf("Order placed! Order-id: %s", sellOrder3)
	}

	// SELL 4
	sellOrder4, err := placeOrder("SHIB-USDT", "sell", "limit", "10000", "0.0001", false)
	if err != nil {
		log.Printf("Failed to place order: %s", err)
	} else {
		log.Printf("Order placed! Order-id: %s", sellOrder4)
	}

}
