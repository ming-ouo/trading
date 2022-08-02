package env

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func Init() {
	initEnvVar()
	initZap()
}

func initEnvVar() {
	v := viper.New()
	v.SetConfigName("config.yaml")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")

	e := v.ReadInConfig()
	if e != nil {
		panic(e)
	}

	RabbitMQHost = v.GetString("rabbitmq.host")
	RabbitMQPort = v.GetInt("rabbitmq.port")
	RabbitMQUser = v.GetString("rabbitmq.user")
	RabbitMQPassword = v.GetString("rabbitmq.password")
	InputStreamName = v.GetString("input_stream.name")
	OutputStreamName = v.GetString("output_stream.name")
}

func initZap() {
	var l *zap.Logger
	var e error

	if Debug {
		l, e = zap.NewDevelopment()
	} else {
		config := zap.NewProductionConfig()
		config.DisableStacktrace = true
		l, e = config.Build()
	}

	if e != nil {
		panic(e)
	}

	zap.ReplaceGlobals(l)
}
