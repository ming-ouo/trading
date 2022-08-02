package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"sync/atomic"
	"time"

	"github.com/bytedance/sonic"
	"github.com/ming-ouo/trading/pkg/orderbook"
	"github.com/rabbitmq/rabbitmq-stream-go-client/pkg/amqp"
	"github.com/rabbitmq/rabbitmq-stream-go-client/pkg/message"
	"github.com/rabbitmq/rabbitmq-stream-go-client/pkg/stream"
	"github.com/shopspring/decimal"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func CheckErr(err error) {
	if err != nil {
		fmt.Printf("%s ", err)
		os.Exit(1)
	}
}

var messagesConfirmed int32

func handlePublishConfirm(confirms stream.ChannelPublishConfirm) {
	go func() {
		for confirmed := range confirms {
			for _, msg := range confirmed {
				if msg.IsConfirmed() {
					atomic.AddInt32(&messagesConfirmed, 1)
				}
			}
		}
	}()
}

func main() {
	log.Println("Start loading test data")
	// Connect to the broker ( or brokers )
	sEnv, err := stream.NewEnvironment(
		stream.NewEnvironmentOptions().
			SetHost("localhost").
			SetPort(5552).
			SetUser("guest").
			SetPassword("guest"))
	CheckErr(err)

	numOfOrders := 10000000
	numOfOrdersInt64 := int64(numOfOrders)

	streamName := "symbol_input_streams"

	err = sEnv.DeclareStream(streamName,
		&stream.StreamOptions{
			MaxLengthBytes: stream.ByteCapacity{}.GB(4),
		},
	)
	CheckErr(err)

	producer, err := sEnv.NewProducer(streamName, stream.NewProducerOptions().
		SetSubEntrySize(500).
		SetCompression(stream.Compression{}.None()))
	CheckErr(err)

	//optional publish confirmation channel
	chPublishConfirm := producer.NotifyPublishConfirmation()
	handlePublishConfirm(chPublishConfirm)

	arr := make([]message.StreamMessage, 0, numOfOrdersInt64)

	for i := int64(0); i < numOfOrdersInt64/2; i++ {
		newOrder := orderbook.NewOrder(RandStringRunes(5), orderbook.Sell, 10, time.Now().UnixNano(), 0, decimal.NewFromInt(i+1), decimal.NewFromInt(0))
		b, err := sonic.Marshal(newOrder)
		CheckErr(err)
		arr = append(arr, amqp.NewMessage(b))
	}

	for i := int64(0); i < numOfOrdersInt64/2; i++ {
		newOrder := orderbook.NewOrder(RandStringRunes(5), orderbook.Buy, 10, time.Now().UnixNano(), 0, decimal.NewFromInt(i+1), decimal.NewFromInt(0))
		b, err := sonic.Marshal(newOrder)
		CheckErr(err)
		arr = append(arr, amqp.NewMessage(b))
	}

	log.Println(numOfOrdersInt64)
	log.Println(len(arr))

	for i := int64(0); i < numOfOrdersInt64; i++ {
		err = producer.Send(arr[i])
		CheckErr(err)
	}

	log.Println("Finish loading test data")

	err = sEnv.Close()
	CheckErr(err)
}
