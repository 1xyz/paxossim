package env

import (
	v1 "github.com/1xyz/paxossim/v1"
	"github.com/1xyz/paxossim/v1/components"
	log "github.com/sirupsen/logrus"
	"time"
)

const (
	ClientReqInterval = 1 * time.Second
)

type Env struct {
	exchange v1.MessageExchange

	replicas []*components.Replica

	leaders []*components.Leader

	clients []*components.Client

	acceptors []*components.Acceptor
}

func NewEnv(nFailures int, nClients int) *Env {
	nReplicas := nFailures + 1
	nLeaders := nFailures + 1
	nAcceptors := (2 * nFailures) + 1
	exchange := v1.NewMessageExchange()

	acceptorAddr := make([]v1.Addr, nAcceptors, nAcceptors)
	acceptors := make([]*components.Acceptor, nAcceptors, nAcceptors)
	for i := 0; i < nAcceptors; i++ {
		acceptors[i] = components.NewAcceptor(exchange)
		acceptorAddr[i] = acceptors[i].GetAddr()
	}

	leaderAddr := make([]v1.Addr, nLeaders, nLeaders)
	leaders := make([]*components.Leader, nLeaders, nLeaders)
	for i := 0; i < nLeaders; i++ {
		leaders[i] = components.NewLeader(exchange, acceptorAddr)
		leaderAddr[i] = leaders[i].GetAddr()
	}

	replicas := make([]*components.Replica, nReplicas, nReplicas)
	for i := 0; i < nReplicas; i++ {
		replicas[i] = components.NewReplica(exchange, leaderAddr)
	}

	log.WithFields(log.Fields{
		"nFailures":  nFailures,
		"nReplicas":  nReplicas,
		"nClients":   nClients,
		"nLeaders":   nLeaders,
		"nAcceptors": nAcceptors,
	}).Debug("Components constructed")

	// construct the clients
	clients := make([]*components.Client, nClients, nClients)
	for i := 0; i < nClients; i++ {
		clients[i] = components.NewClient(exchange, ClientReqInterval)
	}

	return &Env{
		exchange:  exchange,
		leaders:   leaders,
		replicas:  replicas,
		clients:   clients,
		acceptors: acceptors,
	}
}

func (e *Env) Run() {
	for _, a := range e.acceptors {
		go a.Run()
	}
	for _, l := range e.leaders {
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
		log.Infof("Stopping client %v", c)
		c.Stop()
	}
}
