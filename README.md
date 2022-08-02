# trading

## 如何使用

### 1. 安裝暫時的 RabbitMQ

```shell 
docker run -it --rm --name rabbitmq -p 5552:5552 -p 15672:15672\
    -e RABBITMQ_SERVER_ADDITIONAL_ERL_ARGS='-rabbitmq_stream advertised_host localhost -rabbit loopback_users "none"' \
    rabbitmq:3.9-management
```

```shell
docker exec rabbitmq rabbitmq-plugins enable rabbitmq_stream_management
```

### 2. 載入測試資料

```bash
~$ go run cmd/loadtest/loadtest.go
```

### 3. 執行 Trading Server

```bash
~$ go run main.go
{"level":"info","ts":1659418488.8361995,"caller":"trading/trading.go:154","msg":"Trading starts running ..."}
```

## Output 說明

Trading Server 會每秒持續印出 server 的簡易資訊，資訊如下：

```
2022/08/02 13:34:51 now:  1021181 # 目前總共處理了多少訂單
2022/08/02 13:34:51 diff: 335894  # 過去一秒內處理了多少訂單
2022/08/02 13:34:51 volume:  0    # 總成交量
2022/08/02 13:34:51 sellQueueSize 1021192 # Sell Queue 的資料數量
2022/08/02 13:34:51 buyQueueSize 0 # Buy Queue 的資料數量
```

## 簡易設計說明

### OrderBook

`pkg/orderbook/*`

`OrderBook` 會有兩個 Queue 分別為 `SellQueue/BuyQueue(封裝為 OrderQueue)`

`OrderQueue` 採用 `TreeMap` 實現，底層的資料結構是 `Red-Black Tree`(紅黑樹)

`OrderQueue` 可以根據 `Comparator` 對資料的儲存進行排序

這裡分別會對`交易價格Price`、`時間TS` 進行 Asc/Desc 排序

- `SellQueue`: order by price asc, ts
- `BuyQueue`: order by price desc, ts

這樣可以使用同一套程式碼去處理交易 eg. `NewLimitPriceOrder` 可以同時針對 `Sell/Buy`

### Trading

`services/trading/trading.go`

主要使用 RabbitMQ Streams 作為下單輸入

`handleInputMessageFunc`: 從 Queue 讀出資料後會 `Unmarshal` 然後放到 channel `chWaitTrading` 等待後續處理
`processTrading`: 理論上要針對 order 的 `Type/Action` 做出相對應的訂單處理

## 簡易 Benchmark

在筆電上簡單進行測試

```
CPU: Intel(R) Core(TM) i5-1035G4 CPU @ 1.10GHz 4C8T
Memory: 32 GB
SSD: XPG SX8200 Pro PCIe Gen3x4 M.2 2280 1TB
```

- 平均每秒可以處理的 Order 數量：~500,000/s
