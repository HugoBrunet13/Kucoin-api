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

	SYMBOL := "FTG-USDT" // XXX Put in config?
	TYPE := "limit"

	INITIAL_AMOUNT := 200.0
	DP := 3 // XXX
	DEFAULT_PRICE := 0.055

	// 1. Wait until market open
	startTime := viper.GetString("targetTime")
	wait(startTime)

	// 2. Get orderbook
	orderbook, err := getOrderbook(SYMBOL, 100)
	if err != nil {
		log.Printf("ERROR: %s", err)
	}
	log.Printf("orderbook: %v", orderbook)

	// 3. Get buy price from orderbook
	buyPrice := getBuyPrice(50, orderbook)

	// 4. If buy price is null, use default price

	if buyPrice == "" {
		buyPrice = strconv.FormatFloat(DEFAULT_PRICE, 'f', DP, 64)
	}

	buyPriceFloat, err := strconv.ParseFloat(buyPrice, 32)
	if (buyPriceFloat - DEFAULT_PRICE) > (DEFAULT_PRICE / 2) {
		buyPrice = strconv.FormatFloat(DEFAULT_PRICE, 'f', DP, 64)
		buyPriceFloat, err = strconv.ParseFloat(buyPrice, 32)
	}

	log.Printf("buy price: %s", buyPrice)

	// Compute Buy Qty
	qtyBuy := int(INITIAL_AMOUNT / buyPriceFloat)

	log.Printf("INITAL AMOUNT: %s  buyPriceFloat: %s", INITIAL_AMOUNT, buyPriceFloat)
	log.Printf("buy qty: %s", qtyBuy)

	// BUY --> Delay
	buyOrder, err := placeOrder(SYMBOL, "buy", TYPE, strconv.Itoa(qtyBuy), buyPrice, false)
	if err != nil {
		log.Printf("Failed to place order: %s", err)
	} else {
		log.Printf("Order placed! Order-id: %s", buyOrder)
	}

	// // wait 500 miliecond ?
	time.Sleep(500 * time.Millisecond)

	//SELL 1
	qtySell1 := int(float64(qtyBuy) * 0.25)
	priceSell1 := buyPriceFloat * 4

	log.Printf("PriceSell1: %s", priceSell1)

	sellOrder1, err := placeOrder(SYMBOL, "sell", TYPE, strconv.Itoa(qtySell1), strconv.FormatFloat(priceSell1, 'f', DP, 64), false)
	if err != nil {
		log.Printf("Failed to place order: %s", err)
	} else {
		log.Printf("Order placed! Order-id: %s", sellOrder1)
	}

	// SELL 2
	qtySell2 := int(float64(qtyBuy) * 0.20)
	priceSell2 := buyPriceFloat * 10

	sellOrder2, err := placeOrder(SYMBOL, "sell", TYPE, strconv.Itoa(qtySell2), strconv.FormatFloat(priceSell2, 'f', DP, 64), false)
	if err != nil {
		log.Printf("Failed to place order: %s", err)
	} else {
		log.Printf("Order placed! Order-id: %s", sellOrder2)
	}

	// SELL 3
	qtySell3 := int(float64(qtyBuy) * 0.15)
	priceSell3 := buyPriceFloat * 50

	sellOrder3, err := placeOrder(SYMBOL, "sell", TYPE, strconv.Itoa(qtySell3), strconv.FormatFloat(priceSell3, 'f', DP, 64), false)
	if err != nil {
		log.Printf("Failed to place order: %s", err)
	} else {
		log.Printf("Order placed! Order-id: %s", sellOrder3)
	}

	// SELL 4
	qtySell4 := int(float64(qtyBuy) * 0.10)
	priceSell4 := buyPriceFloat * 100

	sellOrder4, err := placeOrder(SYMBOL, "sell", TYPE, strconv.Itoa(qtySell4), strconv.FormatFloat(priceSell4, 'f', DP, 64), false)
	if err != nil {
		log.Printf("Failed to place order: %s", err)
	} else {
		log.Printf("Order placed! Order-id: %s", sellOrder4)
	}

	// SELL 5
	qtySell5 := int(float64(qtyBuy) * 0.13)
	priceSell5 := buyPriceFloat * 1000

	sellOrder5, err := placeOrder(SYMBOL, "sell", TYPE, strconv.Itoa(qtySell5), strconv.FormatFloat(priceSell5, 'f', DP, 64), false)
	if err != nil {
		log.Printf("Failed to place order: %s", err)
	} else {
		log.Printf("Order placed! Order-id: %s", sellOrder5)
	}

	// SELL 6
	qtySell6 := int(float64(qtyBuy) * 0.12)
	priceSell6 := buyPriceFloat * 1800

	sellOrder6, err := placeOrder(SYMBOL, "sell", TYPE, strconv.Itoa(qtySell6), strconv.FormatFloat(priceSell6, 'f', DP, 64), false)
	if err != nil {
		log.Printf("Failed to place order: %s", err)
	} else {
		log.Printf("Order placed! Order-id: %s", sellOrder6)
	}

	// SELL 7
	qtySell7 := int(float64(qtyBuy) * 0.05)
	priceSell7 := buyPriceFloat * 18000

	sellOrder7, err := placeOrder(SYMBOL, "sell", TYPE, strconv.Itoa(qtySell7), strconv.FormatFloat(priceSell7, 'f', DP, 64), false)
	if err != nil {
		log.Printf("Failed to place order: %s", err)
	} else {
		log.Printf("Order placed! Order-id: %s", sellOrder7)
	}
}
