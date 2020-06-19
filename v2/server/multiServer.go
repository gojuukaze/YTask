package server

import (
	"github.com/gojuukaze/YTask/v2/config"
	"github.com/gojuukaze/YTask/v2/log"
	"github.com/sirupsen/logrus"
)

type MultiServer struct {
	ServerMap map[string]*Server // groupName:server

	config config.Config
}

func NewMultiServer(c config.Config) MultiServer {

	s := make(map[string]*Server)
	if c.Debug {
		log.YTaskLog.SetLevel(logrus.DebugLevel)
	}
	return MultiServer{
		ServerMap: s,
		config:    c,
	}
}

func (t MultiServer) CloneConfig() config.Config {
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
func (t *MultiServer) Add(groupName string, workerName string, w interface{}) {
	server := t.GetOrCreateServer(groupName)
	server.Add(groupName, workerName, w)

}

func (t *MultiServer) GetOrCreateServer(groupName string) *Server {
	server, ok := t.ServerMap[groupName]
	if ok {
		return server
	} else {
		newServer := NewServer(t.config.Clone())
		t.ServerMap[groupName] = &newServer
		return t.ServerMap[groupName]
	}

}

func (t *MultiServer) Run(groupName string, numWorkers int) {
	server, ok := t.ServerMap[groupName]
	if !ok {
		panic("YTask: not found group: " + groupName)
	}
	server.Run(groupName, numWorkers)

}

func (t *MultiServer) GetClient() Client {
	server := NewServer(t.config.Clone())

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
