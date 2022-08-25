package ytask

import (
	"github.com/gojuukaze/YTask/v3/backends"
	"github.com/gojuukaze/YTask/v3/brokers"
	"github.com/gojuukaze/YTask/v3/config"
	"github.com/gojuukaze/YTask/v3/consts"
	"github.com/gojuukaze/YTask/v3/log"
	"github.com/gojuukaze/YTask/v3/server"
	"github.com/gojuukaze/YTask/v3/util/yjson"
	jsoniter "github.com/json-iterator/go"
)

var (
	Server = iServer{}
	Logger = iLogger{}
	Config = iConfig{}
)

type iServer struct {
}

func (i iServer) NewServer(setConfigFunc ...config.SetConfigFunc) server.Server {
	c := config.NewConfig(setConfigFunc...)
	return server.NewServer(c)
}

type iConfig struct {
}

func (i iConfig) Broker(b brokers.BrokerInterface) config.SetConfigFunc {
	return config.Broker(b)
}

func (i iConfig) Backend(b backends.BackendInterface) config.SetConfigFunc {
	return config.Backend(b)
}
func (i iConfig) Debug(debug bool) config.SetConfigFunc {
	return config.Debug(debug)
}

func (i iConfig) EnableDelayServer(enable bool) config.SetConfigFunc {
	return config.EnableDelayServer(enable)
}

func (i iConfig) DelayServerQueueSize(size int) config.SetConfigFunc {
	return config.DelayServerQueueSize(size)
}

// default: 1 day
// task status expires in ex seconds, -1:forever,
func (i iConfig) StatusExpires(ex int) config.SetConfigFunc {
	return config.StatusExpires(ex)
}

// default: 1day
// task result expires in ex seconds, -1:forever,
func (i iConfig) ResultExpires(ex int) config.SetConfigFunc {
	return config.ResultExpires(ex)
}

type iLogger struct {
}

func (i iLogger) NewYTaskLogger() log.LoggerInterface {
	return log.NewYTaskLogger(log.YTaskLog)
}

func UseV2Name() {
	consts.UserV2Name = true
	yjson.YJson = jsoniter.Config{
		EscapeHTML:             true,
		ValidateJsonRawMessage: true,
		TagKey:                 "v2JsonName",
	}.Froze()
}
