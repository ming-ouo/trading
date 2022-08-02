package orderbook

import (
	"github.com/shopspring/decimal"
	"testing"
)

func TestNewOrderQueueAsc(t *testing.T) {
	oq := NewOrderQueueAsc()

	d10 := decimal.NewFromInt(10)
	d20 := decimal.NewFromInt(20)
	d30 := decimal.NewFromInt(30)

	oq.Put(NewOrderKey(30, d10), nil)
	oq.Put(NewOrderKey(20, d20), nil)
	oq.Put(NewOrderKey(10, d30), nil)

	headOrder, _ := oq.Head()
	if !headOrder.Price.Equal(d10) {
		t.Errorf("got %s, want %s", d10.String(), headOrder.Price.String())
	}

	tailOrder, _ := oq.Tail()
	if !tailOrder.Price.Equal(d30) {
		t.Errorf("got %s, want %s", d30.String(), tailOrder.Price.String())
	}

	oq.Put(NewOrderKey(5, d10), nil)
	headOrder, _ = oq.Head()
	if headOrder.TS != 5 {
		t.Errorf("got %d, want %d", headOrder.TS, 5)
	}
}

func TestNewOrderQueueDesc(t *testing.T) {
	oq := NewOrderQueueDesc()
	d10 := decimal.NewFromInt(10)
	d20 := decimal.NewFromInt(20)
	d30 := decimal.NewFromInt(30)

	oq.Put(NewOrderKey(30, d10), nil)
	oq.Put(NewOrderKey(20, d20), nil)
	oq.Put(NewOrderKey(10, d30), nil)

	headOrder, _ := oq.Head()
	if !headOrder.Price.Equal(d30) {
		t.Errorf("got %s, want %s", d30.String(), headOrder.Price.String())
	}

	tailOrder, _ := oq.Tail()
	if !tailOrder.Price.Equal(d10) {
		t.Errorf("got %s, want %s", d10.String(), tailOrder.Price.String())
	}

	oq.Put(NewOrderKey(5, d30), nil)
	headOrder, _ = oq.Head()
	if headOrder.TS != 5 {
		t.Errorf("got %d, want %d", headOrder.TS, 5)
	}
}
