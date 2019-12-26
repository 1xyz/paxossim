package paxossim

import log "github.com/sirupsen/logrus"

type Commander struct {
	Process
	Leader    Entity
	Acceptors []Entity
	Replicas  []Entity
	PV        *PValue
}

func NewCommander(commanderID string, leader Entity, acceptors []Entity, replicas []Entity, pv *PValue) *Commander {
	return &Commander{
		Process:   *NewProcess(commanderID),
		Leader:    leader,
		Acceptors: acceptors,
		Replicas:  replicas,
		PV:        pv,
	}
}

func (cmd *Commander) Run() {
	ctxLog := log.WithFields(log.Fields{"id": cmd.pid})
	addrSet := cmd.BroadcastToAcceptors()

	for {
		msg := cmd.inbox.WaitForItem()

		p2Resp, ok := msg.(*P2bMessage)
		if !ok {
			ctxLog.Panicf("encountered unknown message type %v", msg)
		}

		if !cmd.HandleMessage(p2Resp, &addrSet) {
			break
		}
	}
}

func (cmd *Commander) BroadcastToAcceptors() StringSet {
	addrSet := make(StringSet)
	p2Request := NewP2aMessage(cmd.pid, cmd, cmd.PV)
	for _, acceptor := range cmd.Acceptors {
		acceptor.SendMessage(p2Request)
		addrSet.Add(acceptor.ID())
	}
	return addrSet
}

func (cmd *Commander) HandleMessage(p2Resp *P2bMessage, addrSet *StringSet) bool {
	majority := float64(len(cmd.Acceptors)) / 2
	if cmd.PV.BN.CompareTo(p2Resp.BN) == 0 && addrSet.Contains(p2Resp.Src) {
		addrSet.Remove(p2Resp.Src)
		if float64(addrSet.Len()) < majority {
			decisionMessage := NewDecisionMessage(cmd.ID(), cmd.PV.Slot, cmd.PV.C)
			for _, replica := range cmd.Replicas {
				replica.SendMessage(decisionMessage)
			}
			return false
		}
	} else {
		premptedMessage := NewPremptedMessage(cmd.ID(), p2Resp.BN)
		cmd.Leader.SendMessage(premptedMessage)
		return false
	}

	return true
}

type StringSet map[string]bool

func (ss StringSet) Contains(key string) bool {
	_, ok := ss[key]
	return ok
}

func (ss StringSet) Add(key string) {
	ss[key] = true
}

func (ss StringSet) Remove(key string) {
	delete(ss, key)
}

func (ss StringSet) Len() int {
	return len(ss)
}
