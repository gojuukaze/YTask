package server

import (
	"context"
	"fmt"
	"github.com/gojuukaze/YTask/v2/config"
	"github.com/gojuukaze/YTask/v2/message"
	"sync"
)


type DelayServer struct {
	sync.Map
	serverUtils
	delayGroupName string

	// 延时任务的本地队列，用于在本地排序
	queue SortQueue

	// 到处理时间的任务先放入readyMsgChan暂存，然后在放入inlineServerMsgChan
	readyMsgChan chan message.Message
	// inlineServer中的chan
	inlineServerMsgChan chan message.Message

	// stop chan
	safeStopChan         chan struct{}
	getDelayMsgStopChan  chan struct{}
	getReadyMsgStopChan  chan struct{}
	sendReadyMsgStopChan chan struct{}
}

func NewDelayServer(groupName string, c config.Config, msgChan chan message.Message) DelayServer {
	ds := DelayServer{
		serverUtils:          newServerUtils(c.Broker, nil, c.Logger, 0, 0),
		queue:                NewSortQueue(c.DelayServerQueueSize),
		readyMsgChan:         make(chan message.Message, 5),
		inlineServerMsgChan:  msgChan,
		safeStopChan:         make(chan struct{}),
		getDelayMsgStopChan:  make(chan struct{}),
		getReadyMsgStopChan:  make(chan struct{}),
		sendReadyMsgStopChan: make(chan struct{}),
	}
	ds.delayGroupName = ds.GetDelayGroupName(groupName)
	return ds
}

func (s *DelayServer) IsStop() bool {
	_, ok := s.Load("isStop")
	return ok
}

func (s *DelayServer) SetStop() {
	s.Store("isStop", struct{}{})

}

func (s *DelayServer) IsRunning() bool {
	_, ok := s.Load("isRunning")
	return ok
}

func (s *DelayServer) SetRunning() {
	s.Store("isRunning", struct{}{})

}

func (s *DelayServer) Run() {
	if s.IsRunning() {
		panic("DelayServer " + s.delayGroupName + " is running")
	}

	s.Store("isRunning", struct{}{})
	s.SetBrokerPoolSize(11)
	s.BrokerActivate()

	//log.YTaskLog.WithField("server", s.delayGroupName).Infof("Start delayServer[%s] ", s.delayGroupName)
	s.logger.InfoWithField(fmt.Sprintf("Start delayServer[%s] ", s.delayGroupName), "server", s.delayGroupName)

	go s.GetDelayMsgGoroutine()
	go s.GetReadyMsgGoroutine()
	go s.SendReadyMsgGoroutine()

}

func (s *DelayServer) safeStop() {
	//log.YTaskLog.WithField("server", s.delayGroupName).Info("waiting for incomplete goroutine ")
	s.logger.InfoWithField("waiting for incomplete goroutine ", "server", s.delayGroupName)

	s.SetStop()
	close(s.readyMsgChan)

	<-s.getDelayMsgStopChan
	<-s.getReadyMsgStopChan
	// 必须要等前两个结束才能执行这个
	s.LSendQueue()
	<-s.sendReadyMsgStopChan

}

func (s *DelayServer) Shutdown(ctx context.Context) error {

	go func() {
		s.safeStop()
		s.safeStopChan <- struct{}{}
	}()

	select {
	case <-s.safeStopChan:
	case <-ctx.Done():
		return ctx.Err()
	}

	//log.YTaskLog.WithField("server", s.delayGroupName).Info("Shutdown!")
	s.logger.InfoWithField("Shutdown!", "server", s.delayGroupName)
	return nil
}

func (s *DelayServer) LSendQueue() {
	for i := s.queue.len-1; i >=0; i-- {
		msg := s.queue.Get(i)
		err := s.LSendMsg(s.delayGroupName, msg)
		if err != nil {
			//log.YTaskLog.WithField("server", s.delayGroupName).Error("SendQueue msg error: ", err, " [msg=", msg, "]")
			s.logger.ErrorWithField(fmt.Sprint("SendQueue msg error: ", err, " [msg=", msg, "]"), "server", s.delayGroupName)
		}
	}
}
