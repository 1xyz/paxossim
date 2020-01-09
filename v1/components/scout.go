package components

import (
	v1 "github.com/1xyz/paxossim/v1"
	"github.com/1xyz/paxossim/v1/messages"
	"github.com/1xyz/paxossim/v1/types"
	log "github.com/sirupsen/logrus"
	"sync/atomic"
)

var scoutCount int32 = 0

type Scout struct {
	v1.Process

	exchange v1.MessageExchange

	leader v1.Addr

	acceptors []v1.Addr

	bn types.BallotNumber

	pvalues types.PValues
}

func NewScout(exchange v1.MessageExchange, leader v1.Addr, acceptors []v1.Addr, number types.BallotNumber) *Scout {
	id := atomic.AddInt32(&scoutCount, 1)
	processID := v1.ProcessID(id)
	s := &Scout{
		exchange:  exchange,
		Process:   v1.NewProcess(processID, v1.Scout),
		leader:    leader,
		acceptors: acceptors,
		bn:        number,
		pvalues:   make(types.PValues),
	}

	exchange.Register(s)
	return s
}

func (scout *Scout) Run() {
	ctxLog := log.WithFields(log.Fields{"Addr": scout.GetAddr(), "Method": "Scout.Run"})
	addrSet := scout.broadcastToAcceptors()

	for {
		msg, err := scout.Process.Recv()
		if err != nil {
			ctxLog.Panicf("error in inbox recv %v", err)
		}

		phase1bMessage, ok := msg.(messages.Phase1bMessage)
		if !ok {
			ctxLog.Panicf("unknown message type %v", msg)
		}

		if !scout.handleMessage(phase1bMessage, &addrSet) {
			break
		}
	}

	err := scout.exchange.UnRegister(scout)
	if err != nil {
		ctxLog.Panicf("scout.exchange.UnRegister %v", err)
	}
}
func (scout *Scout) broadcastToAcceptors() v1.AddrSet {
	addrSet := make(v1.AddrSet)
	phase1aMessage := messages.NewPhase1aMessage(scout.GetAddr(), scout.bn)
	for _, acceptor := range scout.acceptors {
		err := scout.exchange.Send(acceptor, phase1aMessage)
		if err != nil {
			log.Panicf("scout.exchange.send failed %v", err)
		}

		addrSet.Add(acceptor)
	}

	return addrSet
}

func (scout *Scout) handleMessage(phase1bMessage messages.Phase1bMessage, addrSet *v1.AddrSet) bool {
	majority := float64(len(scout.acceptors)) / 2

	if types.Compare(&scout.bn, &phase1bMessage.BallotNumber) == 0 && addrSet.Contains(phase1bMessage.Src()) {
		addrSet.Remove(phase1bMessage.Src())
		scout.pvalues.Update(phase1bMessage.PValues)
		if float64(addrSet.Len()) < majority {
			adoptedMessage := messages.NewAdoptedMessage(scout.GetAddr(), scout.bn, scout.pvalues)
			err := scout.exchange.Send(scout.leader, adoptedMessage)
			if err != nil {
				log.Panicf("scout.exchange.send failed %v", err)
			}

			return false
		}
	} else {
		premptedMessage := messages.NewPremptedMessage(scout.GetAddr(), phase1bMessage.BallotNumber)
		err := scout.exchange.Send(scout.leader, premptedMessage)
		if err != nil {
			log.Panicf("scout.exchange.send failed %v", err)
		}

		return false
	}

	return true
}
