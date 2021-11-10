package main

import (
	// "io"
	"flag"
	"log"
	"strconv"
	"time"

	// "net/http"
	"github.com/Kucoin/kucoin-go-sdk"
	"github.com/spf13/viper"
)

var session *kucoin.ApiService

func getBuyPrice(askLevel int64, orderbook *kucoin.PartOrderBookModel) string {
	if len(orderbook.Asks) == 0 {
		return ""
	}
	if len(orderbook.Asks) > int(askLevel) {
		return orderbook.Bids[len(orderbook.Bids)-1][0]
	}
	return orderbook.Bids[askLevel][0]
}

func getOrderbook(symbol string, depth int64) (*kucoin.PartOrderBookModel, error) {
	rsp, err := session.AggregatedPartOrderBook(symbol, depth)
	if err != nil {
		return nil, err
	}
	c := &kucoin.PartOrderBookModel{}
	if err := rsp.ReadData(c); err != nil {
		return nil, err
	}
	return c, nil
}

func wait(startTime string) {
	for {
		now := time.Now()
		nowS := now.Format("2006-01-02 15:04:05")
		log.Printf("Wait....")
		if nowS >= startTime {
			break
		}
	}
}

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
		wait(targetTime)
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

	SYMBOL := "SHIB-USDT" // XXX Put in config?
	// TYPE := "limmit"

	// 1. Wait until market open
	startTime := viper.GetString("targetTime")
	wait(startTime)

	// 2. Get orderbook
	orderbook, err := getOrderbook(SYMBOL, 20)
	if err != nil {
		log.Printf("ERROR: %s", err)
	}
	log.Printf("orderbook: %v", orderbook)

	// 3. Get buy price from orderbook
	buyPrice := getBuyPrice(10, orderbook)

	log.Printf("buy price: %s", buyPrice)
	// 4. If buy price is null, use default price
	if buyPrice == "" {
		buyPrice = "0.0005" // XXXX
	}

	// Compute Buy Qty
	buyPriceFloat, err := strconv.ParseFloat(buyPrice, 32)
	buyQty := 200.0 / buyPriceFloat

	log.Printf("buy qty: %s", int64(buyQty))

	/**
	// BUY --> Delay
	buyOrder, err := placeOrder(SYMBOL, "buy", TYPE, buyQty.toString(), buyPrice, true)
	if err != nil {
		log.Printf("Failed to place order: %s", err)
	} else {
		log.Printf("Order placed! Order-id: %s", buyOrder)
	}

	//Derive BUY price and BUY quantity to get Sell price and Sell quantity

	//SELL 1
	sellOrder1, err := placeOrder(SYMBOL, "sell", TYPE, "10000", "0.0001", false)
	if err != nil {
		log.Printf("Failed to place order: %s", err)
	} else {
		log.Printf("Order placed! Order-id: %s", sellOrder1)
	}

	// SELL 2
	sellOrder2, err := placeOrder(SYMBOL, "sell", TYPE, "10000", "0.0001", false)
	if err != nil {
		log.Printf("Failed to place order: %s", err)
	} else {
		log.Printf("Order placed! Order-id: %s", sellOrder2)
	}

	// SELL 3
	sellOrder3, err := placeOrder(SYMBOL, "sell", TYPE, "10000", "0.0001", false)
	if err != nil {
		log.Printf("Failed to place order: %s", err)
	} else {
		log.Printf("Order placed! Order-id: %s", sellOrder3)
	}

	// SELL 4
	sellOrder4, err := placeOrder(SYMBOL, "sell", TYPE, "10000", "0.0001", false)
	if err != nil {
		log.Printf("Failed to place order: %s", err)
	} else {
		log.Printf("Order placed! Order-id: %s", sellOrder4)
	}

	**/

}
