package orderbook

import (
	"github.com/shopspring/decimal"
)

type OrderBook struct {
	SellQueue *OrderQueue // order by price asc, ts
	BuyQueue  *OrderQueue // order by price desc, ts
}

func NewOrderBook() *OrderBook {
	return &OrderBook{
		SellQueue: NewOrderQueueAsc(),
		BuyQueue:  NewOrderQueueDesc(),
	}
}

func (ob *OrderBook) NewLimitPriceOrder(newOrder *Order) (decimal.Decimal, []*Order, []*Order) {
	var popQueue *OrderQueue
	var pushQueue *OrderQueue

	comparator := newOrder.Price.LessThanOrEqual

	if newOrder.Action == Buy {
		popQueue, pushQueue = ob.SellQueue, ob.BuyQueue
	}

	if newOrder.Action == Sell {
		popQueue, pushQueue = ob.BuyQueue, ob.SellQueue
	}

	var popOrderKey *OrderKey
	var popOrder *Order

	newOrderQTYLeft := newOrder.GetQuantity()

	doneOrders := make([]*Order, 0)
	partialOrders := make([]*Order, 0)

	popOrderKey, popOrder = popQueue.Head()

	for newOrderQTYLeft > 0 && !popQueue.Empty() && comparator(popOrder.Price) {
		popOrderQTY := popOrder.GetQuantity()
		newOrderQTY := newOrder.GetQuantity()

		if popOrderQTY > newOrderQTY {
			partialOrders = append(partialOrders, NewOrder(popOrder.Id, popOrder.Action, newOrderQTY, popOrder.TS, getTS(), popOrder.Price, popOrder.Price))
			popOrder.SetQuantity(popOrderQTY - newOrderQTY)
			// Update the left quantity of popOrder
			popQueue.Put(popOrderKey, popOrder)
			newOrderQTYLeft = 0
			break
		}

		// popOrder has completed
		ob.removeOrder(popOrderKey, popOrder)
		popOrder.SetTradedAVGPrice(popOrder.Price)
		doneOrders = append(doneOrders, popOrder)
		newOrderQTYLeft -= popOrder.GetQuantity()
		popOrderKey, popOrder = popQueue.Head()
	}

	newOrder.SetQuantity(newOrderQTYLeft)
	completedNewOrder := newOrder.GetQuantity() - newOrderQTYLeft

	totalPrice := decimal.NewFromInt(0)

	// Calculate the traded average price
	totalPrice = totalPrice.Add(ob.sumTotalPrice(doneOrders))
	totalPrice = totalPrice.Add(ob.sumTotalPrice(partialOrders))

	// The new order has not been completed
	// Add to pushQueue
	if newOrderQTYLeft > 0 {
		// Partial of the new order has been completed
		if completedNewOrder != 0 {
			partialOrders = append(partialOrders, NewOrder(newOrder.Id, newOrder.Action, completedNewOrder, newOrder.TS, getTS(), newOrder.Price, newOrder.Price))
		}

		pushQueue.Put(newOrder.OrderKey(), newOrder)
		return totalPrice, partialOrders, doneOrders
	}

	// newOrderQTYLeft == 0
	newOrder.SetTradedAVGPrice(totalPrice)
	doneOrders = append(doneOrders, newOrder)

	return totalPrice, partialOrders, doneOrders
}

func (ob *OrderBook) sumTotalPrice(oq []*Order) decimal.Decimal {
	totalPrice := decimal.NewFromInt(0)

	for _, v := range oq {
		subTotalPrice := decimal.NewFromInt(int64(v.GetQuantity())).Mul(v.TradedAVGPrice)
		totalPrice = totalPrice.Add(subTotalPrice)
	}

	return totalPrice
}

func (ob *OrderBook) removeOrder(orderKey *OrderKey, order *Order) {
	if order.Action == Sell {
		ob.SellQueue.Remove(orderKey)
		return
	}

	ob.BuyQueue.Remove(orderKey)
}
