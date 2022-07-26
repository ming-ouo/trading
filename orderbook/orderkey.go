package orderbook

import "github.com/shopspring/decimal"

type OrderKey struct {
	TS    int64
	Price decimal.Decimal
}

func NewOrderKey(ts int64, price decimal.Decimal) *OrderKey {
	return &OrderKey{
		TS:    ts,
		Price: price,
	}
}

func (ok *OrderKey) CompareAsc(b *OrderKey) int {
	priceCmp := ok.Price.Cmp(b.Price)
	if priceCmp != 0 {
		return priceCmp
	}

	return ok.compareTS(b)
}

func (ok *OrderKey) CompareDesc(b *OrderKey) int {
	priceCmp := b.Price.Cmp(ok.Price)
	if priceCmp != 0 {
		return priceCmp
	}

	return ok.compareTS(b)
}

func (ok *OrderKey) compareTS(b *OrderKey) int {
	if ok.TS > b.TS {
		return 1
	}

	if ok.TS < b.TS {
		return -1
	}

	return 0
}
