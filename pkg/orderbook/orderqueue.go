package orderbook

import (
	"fmt"
	"github.com/emirpasic/gods/maps/treemap"
)

type OrderQueue struct {
	m *treemap.Map
}

func orderKeyAscComparator(a, b interface{}) int {
	aAss := a.(*OrderKey)
	bAss := b.(*OrderKey)

	return aAss.CompareAsc(bAss)
}

func orderKeyDescComparator(a, b interface{}) int {
	aAss := a.(*OrderKey)
	bAss := b.(*OrderKey)

	return aAss.CompareDesc(bAss)
}

func NewOrderQueueAsc() *OrderQueue {
	return &OrderQueue{
		m: treemap.NewWith(orderKeyAscComparator),
	}
}

func NewOrderQueueDesc() *OrderQueue {
	return &OrderQueue{
		m: treemap.NewWith(orderKeyDescComparator),
	}
}

func (oq *OrderQueue) Put(k *OrderKey, v *Order) {
	oq.m.Put(k, v)
}

func (oq *OrderQueue) Remove(k *OrderKey) {
	oq.m.Remove(k)
}

func (oq *OrderQueue) Head() (*OrderKey, *Order) {
	if oq.m.Size() == 0 {
		return nil, nil
	}

	k, v := oq.m.Min()

	return k.(*OrderKey), v.(*Order)
}

func (oq *OrderQueue) Tail() (*OrderKey, *Order) {
	if oq.m.Size() == 0 {
		return nil, nil
	}

	k, v := oq.m.Max()

	return k.(*OrderKey), v.(*Order)
}

func (oq *OrderQueue) Size() int {
	return oq.m.Size()
}

func (oq *OrderQueue) Empty() bool {
	return oq.m.Empty()
}

func (oq *OrderQueue) DebugKeys() {
	keys := oq.m.Keys()
	for _, v := range keys {
		ok := v.(*OrderKey)
		fmt.Printf("%+v ", ok)
	}

	fmt.Println()
}
func (oq *OrderQueue) DebugValues() {
	values := oq.m.Values()
	for _, v := range values {
		ok := v.(*Order)
		fmt.Printf("%+v ", ok)
	}

	fmt.Println()
}
