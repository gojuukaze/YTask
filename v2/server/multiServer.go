package server

import (
	"context"
	"github.com/gojuukaze/YTask/v2/config"
	"golang.org/x/sync/errgroup"
)

type Server struct {
	ServerMap      map[string]*InlineServer // groupName:server
	DelayServerMap map[string]*DelayServer  // groupName:server

	config config.Config
}

func NewServer(c config.Config) Server {

	if c.Debug {
		//log.YTaskLog.SetLevel(logrus.DebugLevel)
		c.Logger.SetLevel("debug")
	}
	return Server{
		ServerMap:      make(map[string]*InlineServer),
		DelayServerMap: make(map[string]*DelayServer),
		config:         c,
	}
}

// Add worker to group
// w : worker func
// callbackFunc:callbackFunc
func (t *Server) Add(groupName string, workerName string, w interface{}, callbackFunc ...interface{}) {
	server := t.getOrCreateInlineServer(groupName)
	server.Add(workerName, w, callbackFunc...)

}

func (t *Server) getOrCreateInlineServer(groupName string) *InlineServer {
	server, ok := t.ServerMap[groupName]
	if ok {
		return server
	} else {
		newServer := NewInlineServer(groupName, t.config.Clone())
		t.ServerMap[groupName] = &newServer
		return t.ServerMap[groupName]
	}

}

func (t *Server) getOrCreateDelayServer(groupName string) *DelayServer {
	ds, ok := t.DelayServerMap[groupName]
	if ok {
		return ds
	} else {
		is := t.ServerMap[groupName]
		ds := NewDelayServer(groupName, t.config.Clone(), is.msgChan)
		t.DelayServerMap[groupName] = &ds
		return t.DelayServerMap[groupName]
	}

}

func (t *Server) Run(groupName string, numWorkers int, enableDelayServer ...bool) {
	server, ok := t.ServerMap[groupName]
	if !ok {
		panic("YTask: not found group: " + groupName)
	}
	server.Run(numWorkers)
	if (t.config.EnableDelayServer && t.config.DelayServerQueueSize > 0) ||
		(len(enableDelayServer) > 0 && enableDelayServer[0]) {
		ds := t.getOrCreateDelayServer(groupName)
		ds.Run()
	}
}

func (t *Server) GetClient() Client {

	return NewClient(t.config.Clone())
}

func (t *Server) Shutdown(ctx context.Context) error {

	var eg = errgroup.Group{}
	for _, s := range t.ServerMap {
		s := s
		if s.IsRunning() {
			eg.Go(func() error {
				return s.Shutdown(ctx)
			})
		}
	}

	for _, s := range t.DelayServerMap {
		s := s
		if s.IsRunning() {
			eg.Go(func() error {
				return s.Shutdown(ctx)
			})
		}
	}

	return eg.Wait()

}
