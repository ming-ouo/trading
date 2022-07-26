package orderbook

import "time"

func getTS() int64 {
	return time.Now().UnixNano()
}
