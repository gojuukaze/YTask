package yworker

import "github.com/vua/YTask/v1/ymsg"

type WorkerInterface interface {
	Run(msg ymsg.Message) error
	Name() string
}
