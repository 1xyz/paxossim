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
	config    *paxossim.Configuration
	replicas  []*paxossim.Replica
	clients   []*paxossim.Client
	acceptors []paxossim.Entity
}

func NewEnv(nFailures int, nClients int) *Env {
	nReplicas := nFailures + 1
	nLeaders := nFailures + 1
	nAcceptors := (2 * nFailures) + 1

	// create the acceptors
	acceptors := make([]paxossim.Entity, nAcceptors, nAcceptors)
	for i := 0; i < nAcceptors; i++ {
		acceptors[i] = paxossim.NewAcceptor(fmt.Sprintf("Acceptor %d", i))
	}

	config := paxossim.NewConfiguration(nLeaders)
	leaders := make([]*paxossim.Leader, nLeaders, nLeaders)
	for i := 0; i < nLeaders; i++ {
		leaders[i] = paxossim.NewLeader(fmt.Sprintf("leader %d", i), acceptors)
		config.AppendLeader(leaders[i])
	}

	// construct the replicas
	replicas := make([]*paxossim.Replica, nReplicas, nReplicas)
	for i := 0; i < nReplicas; i++ {
		replicas[i] = paxossim.NewReplica(fmt.Sprintf("replica: %d", i), config)
		// register this replica with the leader
		for _, leader := range leaders {
			leader.AppendReplica(replicas[i])
		}
	}
	log.WithFields(log.Fields{"nReplicas": nReplicas}).Debug("Constructed replicas")

	// construct the clients
	clients := make([]*paxossim.Client, 0, nClients)
	for i := 0; i < nClients; i++ {
		id := fmt.Sprintf("client %d", i)
		clients[i] = append(clients, paxossim.NewClient(id, replicas, ClientReqInterval))
	}

	return &Env{
		config:    config,
		replicas:  replicas,
		clients:   clients,
		acceptors: acceptors,
	}
}

func (e *Env) Run() {
	for _, a := range e.acceptors {
		go a.Run()
	}
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
