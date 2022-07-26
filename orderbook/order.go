package orderbook

import (
	"github.com/shopspring/decimal"
)

type OrderAction int

const (
	Buy OrderAction = iota
	Sell
)

type Order struct {
	Id             string
	Action         OrderAction
	quantity       int
	TS             int64
	DoneTS         int64
	Price          decimal.Decimal
	TradedAVGPrice decimal.Decimal
}

func NewOrder(
	id string,
	action OrderAction,
	quantity int,
	ts int64,
	doneTS int64,
	price decimal.Decimal,
	tradedAVGPrice decimal.Decimal,
) *Order {
	return &Order{
		Id:             id,
		Action:         action,
		quantity:       quantity,
		TS:             ts,
		DoneTS:         doneTS,
		Price:          price,
		TradedAVGPrice: tradedAVGPrice,
	}
}

func (o *Order) OrderKey() *OrderKey {
	return NewOrderKey(o.TS, o.Price)
}

func (o *Order) Quantity() int {
	return o.quantity
}

func (o *Order) SetQuantity(nQ int) {
	o.quantity = nQ
}

func (o *Order) SetTradedAVGPrice(p decimal.Decimal) {
	o.TradedAVGPrice = p
}

func (o *Order) Completed() bool {
	return o.quantity == 0
}
