package main

import (
	"context"
	"github.com/ming-ouo/trading/internal/env"
	"github.com/ming-ouo/trading/services/trading"
	"github.com/rabbitmq/rabbitmq-stream-go-client/pkg/stream"
)

func main() {
	env.Init()

	envOptions := stream.NewEnvironmentOptions().
		SetHost(env.RabbitMQHost).
		SetPort(env.RabbitMQPort).
		SetUser(env.RabbitMQUser).
		SetPassword(env.RabbitMQPassword)

	inputOptions := stream.NewConsumerOptions().
		SetConsumerName("my_consumer").                  // set a consumer name
		SetOffset(stream.OffsetSpecification{}.First()). // start consuming from the beginning
		SetCRCCheck(false)                               // Disable crc control, increase the performances

	outputOptions := stream.NewProducerOptions().
		SetSubEntrySize(500).
		SetCompression(stream.Compression{}.None())

	streamOptions := &stream.StreamOptions{
		MaxLengthBytes: stream.ByteCapacity{}.GB(4),
	}

	t := trading.NewTrading(
		env.InputStreamName,
		env.OutputStreamName,
		envOptions,
		streamOptions,
		streamOptions,
		inputOptions,
		outputOptions)

	err := t.Start(context.TODO())
	if err != nil {
		panic(err) // fail fast
	}
}
