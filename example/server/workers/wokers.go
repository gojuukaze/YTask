package workers

import (
	"github.com/gojuukaze/YTask/v1/ymsg"
	"encoding/json"
	"fmt"
)

type NumArgs struct {
	A int `json:"a"`
	B int `json:"b"`
}

type AddWorker struct {
}

func (a AddWorker) Name() string {
	return "add"
}

func (a AddWorker) Run(msg ymsg.Message) error {
	var args NumArgs
	_ = json.Unmarshal([]byte(msg.JsonArgs), &args)

	fmt.Println(args.A + args.B)
	return nil
}

type SubWorker struct {
}

func (s SubWorker) Name() string {
	return "sub"
}

func (s SubWorker) Run(msg ymsg.Message) error {
	var args NumArgs
	_ = json.Unmarshal([]byte(msg.JsonArgs), &args)

	fmt.Println(args.A - args.B)
	return nil
}

type MulWorker struct {
}

func (m MulWorker) Name() string {
	return "mul"
}

func (m MulWorker) Run(msg ymsg.Message) error {
	var args NumArgs
	_ = json.Unmarshal([]byte(msg.JsonArgs), &args)

	fmt.Println(args.A * args.B)
	return nil
}
