package workers

import (
	"encoding/json"
	"fmt"
	"github.com/gojuukaze/YTask/v1.1/message"
)

func AddFunc(msg message.Message) error {
	var num NumArgs
	json.Unmarshal([]byte(msg.JsonArgs), &num)
	fmt.Println(num.A + num.B)
	return nil
}