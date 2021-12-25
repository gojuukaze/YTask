package server

import (
	"time"

	"github.com/vua/YTask/v2/log"
	"github.com/vua/YTask/v2/message"
	"github.com/vua/YTask/v2/yerrors"
)

// 获取延时任务到本地队列
func (s *DelayServer) GetDelayMsgGoroutine() {

	log.YTaskLog.WithField("server", s.delayGroupName).WithField("goroutine", "get_delay_message").Info("start")

	for true {
		if s.IsStop() {
			break
		}
		if s.queue.IsFull() {
			time.Sleep(300 * time.Millisecond)
		}
		msg, err := s.Next(s.delayGroupName)

		if err != nil {
			if !yerrors.IsEqual(err, yerrors.ErrTypeEmptyQuery) {
				log.YTaskLog.WithField("server", s.delayGroupName).WithField("goroutine", "get_delay_message").Error("get msg error, ", err)
			}
			continue
		}
		log.YTaskLog.WithField("server", s.delayGroupName).WithField("goroutine", "get_delay_message").Debug("get delay msg", msg)

		s.GetDelayMsgGoroutine_UpdateQueue(msg)

	}
	s.getDelayMsgStopChan <- struct{}{}
	log.YTaskLog.WithField("server", s.delayGroupName).WithField("goroutine", "get_delay_message").Info("stop")

}

func (s *DelayServer) GetDelayMsgGoroutine_UpdateQueue(msg message.Message) {

	popMsg := s.queue.Insert(msg)
	if popMsg != nil {
		log.YTaskLog.WithField("server", s.delayGroupName).WithField("goroutine", "get_delay_message").Debug("pop msg", *popMsg)

		err := s.SendMsg(s.delayGroupName, *popMsg)
		if err != nil {
			log.YTaskLog.WithField("server", s.delayGroupName).WithField("goroutine", "get_delay_message").Error("Send msg error: ", err, " [msg=", *popMsg, "]")
		}

	}

}

// 从本地队列中获取到处理时间的任务，发送到readyMsgChan
func (s *DelayServer) GetReadyMsgGoroutine() {
	log.YTaskLog.WithField("server", s.delayGroupName).WithField("goroutine", "get_ready_message").Info("start")

	for true {
		if s.IsStop() {
			break
		}
		msg := s.queue.Pop()
		if msg == nil {
			time.Sleep(300 * time.Millisecond)
			continue
		}
		log.YTaskLog.WithField("server", s.delayGroupName).WithField("goroutine", "get_ready_message").Debug("get ready msg", msg)

		err := s.GetReadyMsgGoroutine_Send(*msg)
		// 这里只有停止服务时才会报错
		if err != nil {
			err = s.broker.LSend(s.GetQueueName(s.delayGroupName), *msg)
			if err != nil {
				log.YTaskLog.WithField("server", s.delayGroupName).WithField("goroutine", "get_ready_message").Error("LSend msg error: ", err, " [msg=", msg, "]")
			}
		}
	}
	s.getReadyMsgStopChan <- struct{}{}
	log.YTaskLog.WithField("server", s.delayGroupName).WithField("goroutine", "get_ready_message").Info("stop")

}

func (s *DelayServer) GetReadyMsgGoroutine_Send(msg message.Message) (err error) {

	defer func() {
		e := recover()
		if e != nil {
			err = e.(error)
		}
	}()
	if s.IsStop() {
		err = yerrors.ErrServerStop{}
		return
	}
	s.readyMsgChan <- msg

	return

}

// 从readyMsgChan中读取任务，传给inlineServer的Chan处理
func (s *DelayServer) SendReadyMsgGoroutine() {
	log.YTaskLog.WithField("server", s.delayGroupName).WithField("goroutine", "send_ready_message").Info("start")

	for msg := range s.readyMsgChan {

		log.YTaskLog.WithField("server", s.delayGroupName).WithField("goroutine", "send_ready_message").Info("send ready msg", msg)

		err := s.SendReadyMsgGoroutine_Send(msg)
		// 这里只有停止服务时才会报错
		if err != nil {
			err = s.broker.LSend(s.GetQueueName(s.delayGroupName), msg)
			if err != nil {
				log.YTaskLog.WithField("server", s.delayGroupName).WithField("goroutine", "send_ready_message").Error("LSend msg error: ", err, " [msg=", msg, "]")
			}
		}
	}
	log.YTaskLog.WithField("server", s.delayGroupName).WithField("goroutine", "send_ready_message").Info("stop")
	s.sendReadyMsgStopChan <- struct{}{}

}

func (s *DelayServer) SendReadyMsgGoroutine_Send(msg message.Message) (err error) {
	defer func() {
		e := recover()
		if e != nil {
			err = e.(error)
		}
	}()
	if s.IsStop() {
		err = yerrors.ErrServerStop{}
		return
	}
	s.inlineServerMsgChan <- msg
	return

}
