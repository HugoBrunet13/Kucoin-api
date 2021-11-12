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

func placeOrderMulti(orders []*kucoin.CreateOrderModel, symbol string) {
	rsp, err := session.CreateMultiOrder(symbol, orders)
	if err != nil {
		log.Printf("Failed to Place Muliple Order %s", err)
	}
	r := &kucoin.CreateMultiOrderResultModel{}
	if err := rsp.ReadData(r); err != nil {
		log.Printf("Failed to Parse Muliple Order reponse%s", err)
	}
	log.Print("multiple resp: %v", r)
	// r[0].OrderId
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

func getPriceFromOrderbookThenSendOrder(symbol string, depth int64, askLevel int64, defaultPrice float64, orderType string, dp int, amount float64) {

	orderbook, err := getOrderbook(symbol, depth) // XXX
	if err != nil {
		log.Printf("ERROR: %s", err)
	}
	log.Printf("orderbook: %v", orderbook)

	buyPrice := getBuyPrice(10, orderbook)

	if buyPrice == "" {
		buyPrice = strconv.FormatFloat(defaultPrice, 'f', dp, 64)
	}

	buyPriceFloat, err := strconv.ParseFloat(buyPrice, 32)

	// Compute Buy Qty
	qtyBuy := int(amount / buyPriceFloat)

	log.Printf("Initial amt avail 2: %s  Buy Price: %s buy qty: %d", amount, buyPrice, qtyBuy)

	// BUY --> Delay
	buyOrder, err := placeOrder(symbol, "buy", orderType, strconv.Itoa(qtyBuy), buyPrice, false)
	if err != nil {
		log.Printf("Failed to place order: %s", err)
	} else {
		log.Printf("Order placed! Order-id: %s", buyOrder)
	}
}

func getBalance(token, balanceType string) string {
	rsp, err := session.Accounts(token, balanceType)
	if err != nil {
		log.Printf("Failed to get balances: %s", err)
	}
	cl := kucoin.AccountsModel{}
	if err := rsp.ReadData(&cl); err != nil {
		log.Printf("Failed to parse balances: %s", err)
	}
	return cl[0].Balance
}

func main() {

	SYMBOL := "SHIB-USDT"
	TYPE := "LIMIT" // order type to send
	TOKEN := "SHIB"

	INITIAL_AMOUNT1 := 0.27 // For Buy order 1
	INITIAL_AMOUNT2 := 0.27 // For Buy order 2

	ORDERBOOK_DEPTH := int64(100)
	ASK_LEVEL := int64(70)

	DP := 8 // Decimal
	DEFAULT_PRICE := 0.00005

	// { % of quantity, x of avg fill price}
	SELL_ORDER_BATCH_1 := [][]float64{
		{0.25, 4},
		{0.25, 8},
		{0.15, 12},
		{0.15, 25},
		{0.03, 27},
	}

	//{ % of quantity, x of avg fill price}
	SELL_ORDER_BATCH_2 := [][]float64{
		{0.1, 1000},
		{0.05, 1800},
		{0.02, 18000},
	}
	// // 1. Wait until market open
	startTime := viper.GetString("targetTime")
	wait(startTime)

	// 1. Get orderbook price and send order
	go getPriceFromOrderbookThenSendOrder(SYMBOL, ORDERBOOK_DEPTH, ASK_LEVEL, DEFAULT_PRICE, TYPE, DP, INITIAL_AMOUNT2)

	// 2. Send buy order with default price
	buyPrice := strconv.FormatFloat(DEFAULT_PRICE, 'f', DP, 64)
	buyPriceFloat, err := strconv.ParseFloat(buyPrice, 32)

	// Compute Buy Qty
	qtyBuy := int(INITIAL_AMOUNT1 / buyPriceFloat)

	log.Printf("Initial amt avail : %f  Buy Price: %s buy qty: %d", INITIAL_AMOUNT1, buyPrice, qtyBuy)

	// BUY --> Delay
	buyOrder, err := placeOrder(SYMBOL, "buy", TYPE, strconv.Itoa(qtyBuy), buyPrice, false)
	if err != nil {
		log.Printf("Failed to place order: %s", err)
	} else {
		log.Printf("Order placed! Order-id: %s", buyOrder)
	}

	// // wait 1000 miliecond ?
	time.Sleep(1000 * time.Millisecond)

	// Get balances token
	balanceToken := getBalance(TOKEN, "trade")
	// Get balances token
	balanceUsdt := getBalance("USDT", "trade")

	if balanceToken == "" {
		balanceToken = "200" //XXX
	}
	if balanceUsdt == "" {
		balanceUsdt = "2" //XXX
	}

	balanceTokenFloat, err := strconv.ParseFloat(balanceToken, 32)
	log.Printf("Balance Token: %s", balanceToken)

	balanceUsdtFloat, err := strconv.ParseFloat(balanceUsdt, 32)
	log.Printf("Balance USDT: %s", balanceUsdt)

	filledAmt := (INITIAL_AMOUNT1 + INITIAL_AMOUNT2) - balanceUsdtFloat
	log.Printf("filled amt: %f", filledAmt)

	//filledAmt/baktoken
	avgFilledPrice := DEFAULT_PRICE // XXX TODO

	log.Printf("Avg filled price: %f", avgFilledPrice)

	// Send first batch

	batchOrders1 := make([]*kucoin.CreateOrderModel, 0, len(SELL_ORDER_BATCH_1))

	for i := 0; i < len(SELL_ORDER_BATCH_1); i++ {
		qtySell := int(balanceTokenFloat * SELL_ORDER_BATCH_1[i][0])
		priceSell := avgFilledPrice * SELL_ORDER_BATCH_1[i][1]
		sellOrder := &kucoin.CreateOrderModel{
			ClientOid: kucoin.IntToString(time.Now().UnixNano() + int64(i)),
			Side:      "sell",
			Price:     strconv.FormatFloat(priceSell, 'f', DP, 64),
			Size:      strconv.Itoa(qtySell),
		}
		log.Printf("OrderId: %s Qty Sell : %d  Price Sell : %f", sellOrder.ClientOid, qtySell, priceSell)
		batchOrders1 = append(batchOrders1, sellOrder)
	}
	placeOrderMulti(batchOrders1, SYMBOL)

	// // wait 10 second ?
	time.Sleep(10 * time.Second)

	// Send batch 2
	// Send first batch

	batchOrders2 := make([]*kucoin.CreateOrderModel, 0, len(SELL_ORDER_BATCH_2))

	for i := 0; i < len(SELL_ORDER_BATCH_2); i++ {
		qtySell := int(balanceTokenFloat * SELL_ORDER_BATCH_2[i][0])
		priceSell := avgFilledPrice * SELL_ORDER_BATCH_2[i][1]
		sellOrder := &kucoin.CreateOrderModel{
			ClientOid: kucoin.IntToString(time.Now().UnixNano() + int64(i)),
			Side:      "sell",
			Price:     strconv.FormatFloat(priceSell, 'f', DP, 64),
			Size:      strconv.Itoa(qtySell),
		}
		log.Printf("OrderId: %s Qty Sell : %d  Price Sell : %f", sellOrder.ClientOid, qtySell, priceSell)
		batchOrders2 = append(batchOrders2, sellOrder)
	}
	placeOrderMulti(batchOrders2, SYMBOL)
}
