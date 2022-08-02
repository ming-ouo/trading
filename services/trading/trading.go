package trading

import (
	"context"
	json "github.com/bytedance/sonic"
	"github.com/ming-ouo/trading/pkg/orderbook"
	"github.com/rabbitmq/rabbitmq-stream-go-client/pkg/amqp"
	"github.com/rabbitmq/rabbitmq-stream-go-client/pkg/stream"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
	"log"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"
)

type Trading struct {
	inputStreamName     string
	outputStreamName    string
	inputOptions        *stream.ConsumerOptions
	outputOptions       *stream.ProducerOptions
	streamEnv           *stream.Environment
	streamEnvOptions    *stream.EnvironmentOptions
	outputProducer      *stream.Producer
	inputConsumer       *stream.Consumer
	inputStreamOptions  *stream.StreamOptions
	outputStreamOptions *stream.StreamOptions
	chWaitTrading       chan *orderbook.Order
	chWaitOutput        chan *orderbook.Order
	orderBook           *orderbook.OrderBook
	inputMsgCount       int64
	volume              decimal.Decimal
}

func NewTrading(
	inputStreamName string,
	outputStreamName string,
	streamEnvOptions *stream.EnvironmentOptions,
	inputStreamOptions *stream.StreamOptions,
	outputStreamOptions *stream.StreamOptions,
	inputOptions *stream.ConsumerOptions,
	outputOptions *stream.ProducerOptions,
) *Trading {
	return &Trading{
		inputStreamName:     inputStreamName,
		inputOptions:        inputOptions,
		inputStreamOptions:  inputStreamOptions,
		outputStreamName:    outputStreamName,
		outputOptions:       outputOptions,
		outputStreamOptions: outputStreamOptions,
		streamEnvOptions:    streamEnvOptions,
		chWaitTrading:       make(chan *orderbook.Order, 1),
		chWaitOutput:        make(chan *orderbook.Order, 1),
		orderBook:           orderbook.NewOrderBook(),
		volume:              decimal.NewFromFloat(0.0),
	}
}

func (t *Trading) handleInputMessageFunc(consumerContext stream.ConsumerContext, message *amqp.Message) {
	atomic.AddInt64(&t.inputMsgCount, 1)

	messageDataBytes := message.GetData()
	if messageDataBytes == nil {
		zap.L().Error("message.GetData() is nil", zap.Any("message.Header", message.Header))
		return
	}

	var data orderbook.Order

	err := json.Unmarshal(messageDataBytes, &data)
	if err != nil {
		zap.L().Error("json.Unmarshal()", zap.Error(err))
		return
	}

	t.chWaitTrading <- &data
}

func (t *Trading) processTrading() {
	var o *orderbook.Order
	for o = range t.chWaitTrading {
		volume, _, _ := t.orderBook.NewLimitPriceOrder(o)
		t.volume = t.volume.Add(volume)
	}
}

func (t *Trading) processOutput() {
	// TODO
}

func (t *Trading) DebugCount() {
	var prev int64
	var now int64

	for {
		now = atomic.LoadInt64(&t.inputMsgCount)
		log.Println("now: ", now)
		log.Printf("diff: %d \n", now-prev)
		log.Println("volume: ", t.volume.String())
		log.Println("sellQueueSize", t.orderBook.SellQueue.Size())
		log.Println("buyQueueSize", t.orderBook.BuyQueue.Size())
		log.Println("------------------")
		prev = now

		<-time.After(1 * time.Second)
	}
}

func (t *Trading) Start(ctx context.Context) error {
	var err error
	t.streamEnv, err = stream.NewEnvironment(t.streamEnvOptions)
	if err != nil {
		zap.L().Error("stream.NewEnvironment", zap.Error(err))
		return err
	}

	err = t.streamEnv.DeclareStream(t.inputStreamName, t.inputStreamOptions)
	if err != nil {
		zap.L().Error("t.streamEnv.DeclareStream", zap.Error(err))
		return err
	}

	err = t.streamEnv.DeclareStream(t.outputStreamName, t.outputStreamOptions)
	if err != nil {
		zap.L().Error("t.streamEnv.DeclareStream", zap.Error(err))
		return err
	}

	t.inputConsumer, err = t.streamEnv.NewConsumer(
		t.inputStreamName,
		t.handleInputMessageFunc,
		t.inputOptions,
	)
	if err != nil {
		zap.L().Error("t.streamEnv.NewConsumer", zap.Error(err))
		return err
	}

	t.outputProducer, err = t.streamEnv.NewProducer(
		t.outputStreamName,
		t.outputOptions,
	)
	if err != nil {
		zap.L().Error("t.streamEnv.NewProducer", zap.Error(err))
		return err
	}

	go t.processTrading()
	go t.DebugCount()
	//go t.processOutput()

	zap.L().Info("Trading starts running ...")

	// `signal.Notify` registers the given channel to
	// receive notifications of the specified signals.
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case <-ctx.Done():
			t.shutdown()
			return nil
		case <-sigs:
			t.shutdown()
			return nil
		default:
			<-time.After(100 * time.Millisecond)
		}
	}
}

func (t *Trading) shutdown() {
	//TODO
	zap.L().Error("shutdown...")
}
