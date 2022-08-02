package orderbook

import (
	"github.com/shopspring/decimal"
	"testing"
	"time"
)

func TestOrderBook_NewLimitPriceOrder(t *testing.T) {
	ob := NewOrderBook()

	arr := make([]*Order, 0)
	for i := int64(0); i < 3; i++ {
		newOrder := NewOrder("", Sell, 10, time.Now().UnixNano(), 0, decimal.NewFromInt(i+1), decimal.NewFromInt(0))
		arr = append(arr, newOrder)
		ob.NewLimitPriceOrder(newOrder)
	}

	if ob.SellQueue.Size() != 3 {
		t.Errorf("got %d, want %d", ob.SellQueue.Size(), 3)
	}

	newOrder := NewOrder("", Buy, 10, time.Now().UnixNano(), 0, decimal.NewFromInt(1), decimal.NewFromInt(0))
	arr = append(arr, newOrder)
	ob.NewLimitPriceOrder(newOrder)

	if ob.SellQueue.Size() != 2 {
		t.Errorf("got %d, want %d", ob.SellQueue.Size(), 2)
	}
}
