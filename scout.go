package paxossim

import log "github.com/sirupsen/logrus"

type Scout struct {
	*Process
	Leader    Entity
	Acceptors []Entity
	BN        *BallotNumber
	pvalues   PValues
}

func NewScout(scoutID string, leader Entity, acceptors []Entity, BN *BallotNumber) *Scout {
	return &Scout{
		Process:   NewProcess(scoutID),
		Leader:    leader,
		Acceptors: acceptors,
		BN:        BN,
		pvalues:   make(PValues),
	}
}

func (scout *Scout) Run() {
	ctxLog := log.WithFields(log.Fields{"id": scout.pid})
	addrSet := scout.BroadcastToAcceptors()

	for {
		msg := scout.inbox.WaitForItem()

		p1Response, ok := msg.(*P1bMessage)
		if !ok {
			ctxLog.Panicf("encountered unknown message type %v", msg)
		}

		if !scout.HandleMessage(p1Response, &addrSet) {
			break
		}
	}
}

// Broadcast the Phase1 request to the Set of all acceptors
func (scout *Scout) BroadcastToAcceptors() StringSet {
	addrSet := make(StringSet)
	p1Request := NewP1aMessage(scout.pid, scout, scout.BN)
	for _, acceptor := range scout.Acceptors {
		acceptor.SendMessage(p1Request)
		addrSet.Add(acceptor.ID())
	}
	return addrSet
}

func (scout *Scout) HandleMessage(p1Response *P1bMessage, addrSet *StringSet) bool {
	majority := float64(len(scout.Acceptors)) / 2
	// ToDo: remove the addrSet compare case, we should just ignore cases when addresses don't match
	if scout.BN.CompareTo(p1Response.BN) == 0 && addrSet.Contains(p1Response.Src) {
		addrSet.Remove(p1Response.Src)
		scout.pvalues.Update(p1Response.Accepted)
		if float64(addrSet.Len()) < majority {
			adoptedMessage := NewAdoptedMessage(scout.ID(), scout.BN, &scout.pvalues)
			scout.Leader.SendMessage(adoptedMessage)
			return false
		}
	} else {
		premptedMessage := NewPremptedMessage(scout.ID(), p1Response.BN)
		scout.Leader.SendMessage(premptedMessage)
		return false
	}

	return true
}
