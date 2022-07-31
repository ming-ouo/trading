package main

import (
	"github.com/ming-ouo/trading/pkg/orderbook"
	"github.com/shopspring/decimal"
	"log"
	"math/rand"
	"time"
)

func track(msg string) (string, time.Time) {
	return msg, time.Now()
}

func duration(msg string, start time.Time) {
	log.Printf("%v: %v\n", msg, time.Since(start))
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func main() {
	numOfOrders := 500000
	numOfOrdersInt64 := int64(numOfOrders)

	ob := orderbook.NewOrderBook()

	// pre-create orders
	orders := make([]*orderbook.Order, 0, numOfOrders*2)

	for i := int64(0); i < numOfOrdersInt64; i++ {
		newOrder := orderbook.NewOrder(RandStringRunes(5), orderbook.Sell, 10, time.Now().UnixNano(), 0, decimal.NewFromInt(1), decimal.NewFromInt(0))
		orders = append(orders, newOrder)
	}

	for i := int64(0); i < numOfOrdersInt64; i++ {
		newOrder := orderbook.NewOrder(RandStringRunes(5), orderbook.Buy, 10, time.Now().UnixNano(), 0, decimal.NewFromInt(1), decimal.NewFromInt(0))
		orders = append(orders, newOrder)
	}

	// Measure execution time
	_, v := track("trading")

	for i := int64(0); i < numOfOrdersInt64; i++ {
		ob.NewLimitPriceOrder(orders[i])
	}

	for i := numOfOrdersInt64; i < numOfOrdersInt64*2; i++ {
		ob.NewLimitPriceOrder(orders[i])
	}

	log.Printf("number of trade orders: %d", numOfOrders*2)
	duration("trading execution time", v)

	//ob.BuyQueue.DebugKeys()
	//ob.BuyQueue.DebugValues()
	//ob.SellQueue.DebugKeys()
	//ob.SellQueue.DebugValues()
}
