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

var order_1 = "O-00032"

// var order_2 = "O-0021"
// var order_3 = "O-0022"
// var order_4 = "O-0023"
// var order_5 = "O-0024"

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

func placeOrder(orderId string, symbol string, side string, orderType string, size string, price string, delay bool) (string, error) {
	params := &kucoin.CreateOrderModel{
		ClientOid: orderId,
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
			log.Printf("Now: %s - Target: %s", nowS, targetTime)
			if nowS >= targetTime {
				break
			}
		}
	}

	rsp, err := session.CreateOrder(params)
	if err != nil {
		log.Printf("Error1: %s", err)
		return "", err
	}
	os := kucoin.OrderModel{}
	err2 := rsp.ReadData(&os)
	if err2 != nil {
		log.Printf("Error2: %s", err2)
		return "", err2
	}
	log.Printf("Res %s", rsp)
	return rsp.Code, nil

}

var session *kucoin.ApiService

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

	order1, err := placeOrder(order_1, "SHIB-USDT", "sell", "limit", "10000", "0.0001", false)
	if err != nil {
		log.Printf("Failed to place order: %s", err)
	} else {
		log.Printf("Order placed! Order-id: %s", order1)
	}
	// //Sell 1
	// paramsOrder2 := &kucoin.CreateOrderModel{
	// 	ClientOid: order_2,
	// 	Side:      "sell",
	// 	Symbol:    "CERE-USDT",
	// 	Type:      "limit",
	// 	Size:      "500",
	// 	Price:     "0.4",
	// }
	// _, errOrder2 := session.CreateOrder(paramsOrder2)
	// if errOrder2 != nil {
	// 	log.Printf("Error: %s", err)
	// 	// return
	// }
	// // osOrder2 := kucoin.OrdersModel{}
	// // _errOrder2 := rspOrder2.ReadData(&osOrder2)
	// // if _errOrder2 != nil {
	// // 	log.Printf("Error: %s", _errOrder2)
	// // 	return
	// // }
	// // log.Printf("Res %s", rspOrder2)

	// //Sell 2
	// paramsOrder3 := &kucoin.CreateOrderModel{
	// 	ClientOid: order_3,
	// 	Side:      "sell",
	// 	Symbol:    "CERE-USDT",
	// 	Type:      "limit",
	// 	Size:      "500",
	// 	Price:     "0.8",
	// }
	// _, errOrder3 := session.CreateOrder(paramsOrder3)
	// if errOrder3 != nil {
	// 	log.Printf("Error: %s", errOrder3)
	// 	// return
	// }
	// // osOrder3 := kucoin.OrdersModel{}
	// // _errOrder3 := rspOrder3.ReadData(&osOrder3)
	// // if _errOrder3 != nil {
	// // 	log.Printf("Error: %s", _errOrder3)
	// // 	return
	// // }
	// // log.Printf("Res %s", rspOrder3)

	// //Sell 3
	// paramsOrder4 := &kucoin.CreateOrderModel{
	// 	ClientOid: order_4,
	// 	Side:      "sell",
	// 	Symbol:    "CERE-USDT",
	// 	Type:      "limit",
	// 	Size:      "500",
	// 	Price:     "1.5",
	// }
	// _, errOrder4 := session.CreateOrder(paramsOrder4)
	// if errOrder4 != nil {
	// 	log.Printf("Error: %s", errOrder4)
	// 	// return
	// }
	// // osOrder4 := kucoin.OrdersModel{}
	// // _errOrder4 := rspOrder4.ReadData(&osOrder4)
	// // if _errOrder4 != nil {
	// // 	log.Printf("Error: %s", _errOrder4)
	// // 	return
	// // }
	// // log.Printf("Res %s", rspOrder4)

	// //Sell 4
	// paramsOrder5 := &kucoin.CreateOrderModel{
	// 	ClientOid: order_5,
	// 	Side:      "sell",
	// 	Symbol:    "CERE-USDT",
	// 	Type:      "limit",
	// 	Size:      "500",
	// 	Price:     "3",
	// }
	// _, errOrder5 := session.CreateOrder(paramsOrder5)
	// if errOrder5 != nil {
	// 	log.Printf("Error: %s", errOrder5)
	// 	// return
	// }
	// // osOrder5 := kucoin.OrdersModel{}
	// // _errOrder5 := rspOrder5.ReadData(&osOrder5)
	// // if _errOrder5 != nil {
	// // 	log.Printf("Error: %s", _errOrder5)
	// // 	return
	// // }
	// // log.Printf("Res %s", rspOrder5)

}
