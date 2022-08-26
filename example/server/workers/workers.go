package workers

import (
	"errors"
	"github.com/gojuukaze/YTask/v3/core/server"
)

type User struct {
	Id   int
	Name string
}

func Add(a int, b int) int {

	return a + b
}

func AddSub(ctl *server.TaskCtl, a int, b int) (int, int) {
	// do not retry
	ctl.SetRetryCount(0)

	return a + b, a - b
}

func Retry(ctl *server.TaskCtl, a int, b int) (int, int) {
	if ctl.GetRetryCount()%2 == 0 {
		// use ctl.Retry
		ctl.Retry(errors.New("ctl.Retry"))
	} else {
		// or use panic
		panic("panic retry")
	}

	return a + b, a - b
}

func AppendUser(user User, ids []int, names []string) []User {
	var r = make([]User, 0)
	r = append(r, user)
	for i := range ids {
		r = append(r, User{
			Id:   ids[i],
			Name: names[i],
		})
	}
	return r
}
