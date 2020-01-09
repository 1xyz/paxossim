package components

import (
	v1 "github.com/1xyz/paxossim/v1"
	"github.com/1xyz/paxossim/v1/messages"
	"github.com/1xyz/paxossim/v1/types"
	log "github.com/sirupsen/logrus"
	"sync/atomic"
)

var commanderCount int32 = 0

type Commander struct {
	v1.Process

	exchange v1.MessageExchange

	leader v1.Addr

	acceptors []v1.Addr

	pvalue types.PValue
}

func NewCommander(exchange v1.MessageExchange, leader v1.Addr, acceptors []v1.Addr, pvalue types.PValue) *Commander {
	// it is possibl for leaders across go-routines to increment this
	processID := v1.ProcessID(atomic.AddInt32(&commanderCount, 1))

	cmdr := &Commander{
		Process:   v1.NewProcess(processID, v1.Commander),
		exchange:  exchange,
		leader:    leader,
		acceptors: acceptors,
		pvalue:    pvalue,
	}

	exchange.Register(cmdr)
	return cmdr
}

func (cmdr *Commander) Run() {
	ctxLog := log.WithFields(log.Fields{"Addr": cmdr.GetAddr(), "Method": "Commander.Run"})
	addrSet := cmdr.broadcastToAcceptors()

	for {
		msg, err := cmdr.Process.Recv()
		if err != nil {
			ctxLog.Panicf("error in inbox recv %v", err)
		}

		phase2bMessage, ok := msg.(messages.Phase2bMessage)
		if !ok {
			ctxLog.Panicf("unknown message type %v", msg)
		}

		if !cmdr.handleMessage(phase2bMessage, &addrSet) {
			break
		}
	}

	err := cmdr.exchange.UnRegister(cmdr)
	if err != nil {
		ctxLog.Panicf("cmdr.exchange.UnRegister %v", err)
	}
}

func (cmdr *Commander) broadcastToAcceptors() v1.AddrSet {
	addrSet := make(v1.AddrSet)
	phase2aMessage := messages.NewPhase2aMessage(cmdr.GetAddr(), cmdr.pvalue)
	for _, acceptor := range cmdr.acceptors {
		err := cmdr.exchange.Send(acceptor, phase2aMessage)
		if err != nil {
			log.Panicf("cmdr.exchange.send failed %v", err)
		}

		addrSet.Add(acceptor)
	}

	return addrSet
}

func (cmdr *Commander) handleMessage(phase2bMessage messages.Phase2bMessage, addrSet *v1.AddrSet) bool {
	majority := float64(len(cmdr.acceptors)) / 2
	if types.Compare(&cmdr.pvalue.BN, &phase2bMessage.BallotNumber) == 0 && addrSet.Contains(phase2bMessage.Src()) {
		addrSet.Remove(phase2bMessage.Src())
		if float64(addrSet.Len()) < majority {
			decisionMessage := messages.NewDecisionMessage(cmdr.GetAddr(), cmdr.pvalue.Slot, cmdr.pvalue.Command)
			err := cmdr.exchange.SendAll(v1.Replica, decisionMessage)
			if err != nil {
				log.Panicf("cmdr.exchange.sendAll failed %v", err)
			}

			return false
		}
	} else {
		premptedMessage := messages.NewPremptedMessage(cmdr.GetAddr(), phase2bMessage.BallotNumber)
		err := cmdr.exchange.Send(cmdr.leader, premptedMessage)
		if err != nil {
			log.Panicf("cmdr.exchange.send failed %v", err)
		}

		return false
	}

	return true
}
