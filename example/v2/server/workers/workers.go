package workers

import (
	"errors"
	"github.com/gojuukaze/YTask/v2/controller"
	"time"
)

func Add(ctl *controller.TaskCtl, a int, b int) int {

	if ctl.RetryCount >=2 {
		ctl.Retry(errors.New("233"))
		return 1
	}
	time.Sleep(400*time.Millisecond)
	//panic("pppp")
	return a + b
}
