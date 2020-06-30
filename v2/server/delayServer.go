package server

import (
	"github.com/gojuukaze/YTask/v2/config"
	"github.com/gojuukaze/YTask/v2/message"
	"sync"
)

type DelayServer struct {
	sync.Mutex
	groupName string

	// 延时任务获取到本地队列，用于在本地排序
	queue [20]message.DelayMessage

	// 到处理时间的队列
	readyChan chan message.Message

}

func NewDelayServer(groupName string, c config.Config) DelayServer {
	c.Backend = nil
	return DelayServer{
		groupName: groupName,
	}
}
