package workers

import (
	"encoding/json"
	"fmt"
	"github.com/gojuukaze/YTask/v2/message"
	"time"
)

type NumArgs struct {
	A int `json:"a"`
	B int `json:"b"`
}

type AddStruct struct {
}


func (a AddStruct) Run(msg message.Message) error {
	var args NumArgs
	_ = json.Unmarshal([]byte(msg.JsonArgs), &args)
	fmt.Println("run ",args.A)
	time.Sleep(2*time.Second)
	fmt.Println(args.A + args.B)
	return nil
}

type SubStruct struct {
}


func (s SubStruct) Run(msg message.Message) error {
	var args NumArgs
	_ = json.Unmarshal([]byte(msg.JsonArgs), &args)

	fmt.Println(args.A - args.B)
	return nil
}

type MulStruct struct {
}


func (m MulStruct) Run(msg message.Message) error {
	var args NumArgs
	_ = json.Unmarshal([]byte(msg.JsonArgs), &args)

	fmt.Println(args.A * args.B)
	return nil
}
