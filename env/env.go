package env

import (
	"fmt"
	"github.com/1xyz/paxossim"
	log "github.com/sirupsen/logrus"
	"time"
)

const (
	ClientReqInterval = 1 * time.Second
)

type Env struct {
	config   *paxossim.Configuration
	replicas []*paxossim.Replica
	clients  []*paxossim.Client
}

func NewEnv(nReplicas int, nClients int, nLeaders int) *Env {
	leaders := make([]paxossim.Entity, 0, nLeaders)
	for i := 0; i < nLeaders; i++ {
		id := fmt.Sprintf("Leader %d", i)
		leaders = append(leaders, paxossim.NewLeader(id))
	}

	config := &paxossim.Configuration{
		Leaders: leaders,
	}

	// construct the replicas
	replicas := make([]*paxossim.Replica, 0, nReplicas)
	for i := 0; i < nReplicas; i++ {
		id := fmt.Sprintf("replica: %d", i)
		replicas = append(replicas, paxossim.NewReplica(id, config))
	}
	log.WithFields(log.Fields{
		"nReplicas": nReplicas,
	}).Debug("Constructed replicas")

	// construct the clients
	clients := make([]*paxossim.Client, 0, nClients)
	for i := 0; i < nClients; i++ {
		id := fmt.Sprintf("client %d", i)
		clients = append(clients, paxossim.NewClient(id, replicas, ClientReqInterval))
	}

	return &Env{
		config:   config,
		replicas: replicas,
		clients:  clients,
	}
}

func (e *Env) Run() {
	for _, l := range e.config.Leaders {
		go l.Run()
	}
	for _, r := range e.replicas {
		go r.Run()
	}
	for _, c := range e.clients {
		go c.Run()
	}
}

func (e *Env) Stop() {
	for _, c := range e.clients {
		log.Debugf("Stopping client %v", c)
		c.Stop()
	}
}
