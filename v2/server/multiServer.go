package server

import (
	"context"
	"github.com/gojuukaze/YTask/v2/config"
	"github.com/gojuukaze/YTask/v2/log"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

type Server struct {
	ServerMap map[string]*InlineServer // groupName:server

	config config.Config
}

func NewServer(c config.Config) Server {

	s := make(map[string]*InlineServer)
	if c.Debug {
		log.YTaskLog.SetLevel(logrus.DebugLevel)
	}
	return Server{
		ServerMap: s,
		config:    c,
	}
}

func (t Server) CloneConfig() config.Config {
	return config.Config{
		Broker:        t.config.Broker,
		Backend:       t.config.Backend,
		Debug:         t.config.Debug,
		StatusExpires: t.config.StatusExpires,
		ResultExpires: t.config.ResultExpires,
	}

}

// add worker to group
// w : worker func
func (t *Server) Add(groupName string, workerName string, w interface{}) {
	server := t.GetOrCreateInlineServer(groupName)
	server.Add(groupName, workerName, w)

}

func (t *Server) GetOrCreateInlineServer(groupName string) *InlineServer {
	server, ok := t.ServerMap[groupName]
	if ok {
		return server
	} else {
		newServer := NewInlineServer(t.config.Clone())
		t.ServerMap[groupName] = &newServer
		return t.ServerMap[groupName]
	}

}

func (t *Server) Run(groupName string, numWorkers int) {
	server, ok := t.ServerMap[groupName]
	if !ok {
		panic("YTask: not found group: " + groupName)
	}
	server.Run(groupName, numWorkers)

}

func (t *Server) GetClient() Client {
	server := NewInlineServer(t.config.Clone())

	if server.broker != nil {
		if server.broker.GetPoolSize() <= 0 {
			server.broker.SetPoolSize(10)
		}
		server.broker.Activate()
	}
	if server.backend != nil {
		if server.backend.GetPoolSize() <= 0 {
			server.backend.SetPoolSize(10)
		}
		server.backend.Activate()
	}
	return NewClient(&server)
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

	return eg.Wait()

}
