package workers

import (
	"errors"
	"fmt"
	"github.com/gojuukaze/YTask/v2/controller"
	"time"
)

func Add(ctl *controller.TaskCtl, a int, b int) int {
	fmt.Printf("add %+v\n", ctl)

	if ctl.RetryCount >=2 {
		ctl.Retry(errors.New("233"))
		return 1
	}
	time.Sleep(400*time.Millisecond)
	fmt.Printf("%d+%d=%d\n", a, b, a+b)
	//panic("pppp")
	return a + b
}
