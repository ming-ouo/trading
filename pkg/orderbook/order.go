package orderbook

import (
	"github.com/shopspring/decimal"
)

type OrderAction int

const (
	Buy OrderAction = iota
	Sell
)

type OrderType int

const (
	LimitOrder OrderType = iota
)

type Order struct {
	Id             string          `json:"id"`
	Type           OrderType       `json:"type"`
	Action         OrderAction     `json:"action"`
	Quantity       int             `json:"quantity"`
	TS             int64           `json:"ts"`
	DoneTS         int64           `json:"doneTS"`
	Price          decimal.Decimal `json:"price"`
	TradedAVGPrice decimal.Decimal `json:"tradedAvgPrice"`
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
		Quantity:       quantity,
		TS:             ts,
		DoneTS:         doneTS,
		Price:          price,
		TradedAVGPrice: tradedAVGPrice,
	}
}

func (o *Order) OrderKey() *OrderKey {
	return NewOrderKey(o.TS, o.Price)
}

func (o *Order) GetQuantity() int {
	return o.Quantity
}

func (o *Order) SetQuantity(nQ int) {
	o.Quantity = nQ
}

func (o *Order) SetTradedAVGPrice(p decimal.Decimal) {
	o.TradedAVGPrice = p
}

func (o *Order) Completed() bool {
	return o.Quantity == 0
}
