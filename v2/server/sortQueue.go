package server

import (
	"fmt"
	"sync"
	"time"

	"github.com/vua/YTask/v2/message"
)

const maxLen = 20

type SortQueue struct {
	sync.Mutex

	Queue [maxLen + 1]message.Message
	len   int
}

func (s *SortQueue) IsFull() bool {
	return s.len == maxLen
}

func (s *SortQueue) Insert(msg message.Message) *message.Message {
	s.Lock()
	defer s.Unlock()
	if s.len == 0 {
		s.Queue[0] = msg

	} else {
		t := msg.TaskCtl.GetRunTime()
		if s.Queue[s.len-1].RunTimeBeforeOrEqual(t) {
			s.Queue[s.len] = msg

		} else {
			index := s.binarySearch(0, s.len-1, t)
			copy(s.Queue[index+1:], s.Queue[index:])
			s.Queue[index] = msg
		}

	}

	if s.len == maxLen {
		return &s.Queue[s.len]
	} else {
		s.len++
		return nil
	}

}

func (s *SortQueue) binarySearch(leftIndex int, rightIndex int, t time.Time) int {

	if leftIndex > rightIndex {
		return leftIndex
	}
	middleIndex := (leftIndex + rightIndex) / 2

	if s.Queue[middleIndex].RunTimeAfter(t) {
		return s.binarySearch(leftIndex, middleIndex-1, t)
	} else if s.Queue[middleIndex].RunTimeBefore(t) {
		return s.binarySearch(middleIndex+1, rightIndex, t)
	} else {
		return middleIndex + 1
	}
}

func (s *SortQueue) Pop() *message.Message {
	s.Lock()
	defer s.Unlock()
	if s.len == 0 {
		return nil
	}
	if s.Queue[0].IsRunTime() {
		msg := s.Queue[0]
		copy(s.Queue[:], s.Queue[1:])
		s.len--
		return &msg
	}
	return nil
}

func (s *SortQueue) Get(i int) message.Message {
	return s.Queue[i]
}

func (s *SortQueue) print() {
	for i := 0; i < s.len; i++ {
		fmt.Print(s.Queue[i].WorkerName, ",")
	}
	fmt.Println()

}
