package server

import (
	"fmt"
	"github.com/gojuukaze/YTask/v2/message"
	"github.com/gojuukaze/YTask/v2/yerrors"
	"time"
)

// 获取延时任务到本地队列
func (s *DelayServer) GetDelayMsgGoroutine() {

	//log.YTaskLog.WithField("server", s.delayGroupName).WithField("goroutine", "get_delay_message").Info("start")
	s.logger.InfoWithField("goroutine get_delay_message start", "server", s.delayGroupName)

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
				//log.YTaskLog.WithField("server", s.delayGroupName).WithField("goroutine", "get_delay_message").Error("get msg error, ", err)
				s.logger.ErrorWithField(fmt.Sprint("goroutine get_delay_message get msg error, ", err), "server", s.delayGroupName)
			}
			continue
		}

		//log.YTaskLog.WithField("server", s.delayGroupName).WithField("goroutine", "get_delay_message").Debug("get delay msg, ", msg)
		s.logger.DebugWithField(fmt.Sprint("goroutine get_delay_message get delay msg, ", msg), "server", s.delayGroupName)

		s.GetDelayMsgGoroutine_UpdateQueue(msg)

	}
	s.getDelayMsgStopChan <- struct{}{}

	//log.YTaskLog.WithField("server", s.delayGroupName).WithField("goroutine", "get_delay_message").Info("stop")
	s.logger.InfoWithField("goroutine get_delay_message stop", "server", s.delayGroupName)
}

func (s *DelayServer) GetDelayMsgGoroutine_UpdateQueue(msg message.Message) {

	popMsg := s.queue.Insert(msg)
	if popMsg != nil {
		//log.YTaskLog.WithField("server", s.delayGroupName).WithField("goroutine", "get_delay_message").Debug("pop msg, ", *popMsg)
		s.logger.DebugWithField(fmt.Sprint("goroutine get_delay_message pop msg, ", *popMsg), "server", s.delayGroupName)

		err := s.SendMsg(s.delayGroupName, *popMsg)
		if err != nil {
			//log.YTaskLog.WithField("server", s.delayGroupName).WithField("goroutine", "get_delay_message").Error("Send msg error: ", err, " [msg=", *popMsg, "]")
			s.logger.ErrorWithField(fmt.Sprint("goroutine get_delay_message Send msg error: ", err, " [msg=", *popMsg, "]"), "server", s.delayGroupName)
		}

	}

}

// 从本地队列中获取到处理时间的任务，发送到readyMsgChan
func (s *DelayServer) GetReadyMsgGoroutine() {
	//log.YTaskLog.WithField("server", s.delayGroupName).WithField("goroutine", "get_ready_message").Info("start")
	s.logger.InfoWithField("goroutine get_ready_message start", "server", s.delayGroupName)

	for true {
		if s.IsStop() {
			break
		}
		msg := s.queue.Pop()
		if msg == nil {
			time.Sleep(300 * time.Millisecond)
			continue
		}
		//log.YTaskLog.WithField("server", s.delayGroupName).WithField("goroutine", "get_ready_message").Debug("get ready msg: ", msg)
		s.logger.DebugWithField(fmt.Sprint("goroutine get_ready_message get ready msg: ", msg), "server", s.delayGroupName)

		err := s.GetReadyMsgGoroutine_Send(*msg)
		// 这里只有停止服务时才会报错
		if err != nil {
			err = s.broker.LSend(s.GetQueueName(s.delayGroupName), *msg)
			if err != nil {
				//log.YTaskLog.WithField("server", s.delayGroupName).WithField("goroutine", "get_ready_message").Error("LSend msg error: ", err, " [msg=", msg, "]")
				s.logger.ErrorWithField(fmt.Sprint("goroutine get_ready_message LSend msg error: ", err, " [msg=", msg, "]"), "server", s.delayGroupName)
			}
		}
	}
	s.getReadyMsgStopChan <- struct{}{}
	//log.YTaskLog.WithField("server", s.delayGroupName).WithField("goroutine", "get_ready_message").Info("stop")
	s.logger.InfoWithField("goroutine get_ready_message stop", "server", s.delayGroupName)
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
	//log.YTaskLog.WithField("server", s.delayGroupName).WithField("goroutine", "send_ready_message").Info("start")
	s.logger.InfoWithField("goroutine send_ready_message start", "server", s.delayGroupName)

	for msg := range s.readyMsgChan {

		//log.YTaskLog.WithField("server", s.delayGroupName).WithField("goroutine", "send_ready_message").Info("send ready msg: ", msg)
		s.logger.InfoWithField(fmt.Sprint("goroutine send_ready_message send ready msg: ", msg), "server", s.delayGroupName)

		err := s.SendReadyMsgGoroutine_Send(msg)
		// 这里只有停止服务时才会报错
		if err != nil {
			err = s.broker.LSend(s.GetQueueName(s.delayGroupName), msg)
			if err != nil {
				//log.YTaskLog.WithField("server", s.delayGroupName).WithField("goroutine", "send_ready_message").Error("LSend msg error: ", err, " [msg=", msg, "]")
				s.logger.ErrorWithField(fmt.Sprint("goroutine send_ready_message LSend msg error: ", err, " [msg=", msg, "]"), "server", s.delayGroupName)
			}
		}
	}
	//log.YTaskLog.WithField("server", s.delayGroupName).WithField("goroutine", "send_ready_message").Info("stop")
	s.logger.InfoWithField("goroutine send_ready_message stop", "server", s.delayGroupName)
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
