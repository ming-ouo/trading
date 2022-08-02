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

```
2022/08/02 13:34:51 now:  1021181 # 目前總共處理了多少訂單
2022/08/02 13:34:51 diff: 335894  # 過去一秒內處理了多少訂單
2022/08/02 13:34:51 volume:  0    # 總成交量
2022/08/02 13:34:51 sellQueueSize 1021192 # Sell Queue 的資料數量
2022/08/02 13:34:51 buyQueueSize 0 # Buy Queue 的資料數量
```

## 簡易設計說明

### OrderBook

`OrderBook` 會有兩個 Queue 分別為 `SellQueue/BuyQueue(封裝為 OrderQueue)`

`OrderQueue` 採用 `TreeMap` 實現，底層的資料結構是 `Red-Black Tree`(紅黑樹)

`OrderQueue` 可以根據 `Comparator` 對資料的儲存進行排序

這裡分別會對`交易價格Price`、`時間TS` 進行 Asc/Desc 排序

- `SellQueue`: order by price asc, ts
- `BuyQueue`: order by price desc, ts

這樣可以使用同一套程式碼去處理交易 eg. `NewLimitPriceOrder` 可以同時針對 `Sell/Buy`

