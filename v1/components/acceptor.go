package components

import (
	v1 "github.com/1xyz/paxossim/v1"
	"github.com/1xyz/paxossim/v1/messages"
	"github.com/1xyz/paxossim/v1/types"
	log "github.com/sirupsen/logrus"
)

var acceptorCount = 0

type Acceptor struct {
	v1.Process

	exchange v1.MessageExchange

	// Set of PValues accepted so far
	Accepted types.PValues

	// Last Adopted ballot number
	BN *types.BallotNumber
}

func NewAcceptor(exchange v1.MessageExchange) *Acceptor {
	processId := v1.ProcessID(acceptorCount)
	acceptorCount++

	a := &Acceptor{
		Process:  v1.NewProcess(processId, v1.Acceptor),
		Accepted: make(types.PValues),
		BN:       nil,
		exchange: exchange,
	}

	err := exchange.Register(a)
	if err != nil {
		log.Panicf("exchange.Register error %v", err)
	}

	return a
}

func (accp *Acceptor) Run() {
	ctxLog := log.WithFields(log.Fields{"Addr": accp.GetAddr()})

	for {
		msg, err := accp.Process.Recv()
		if err != nil {
			ctxLog.Panicf("error in inbox recv %v", err)
		}

		accp.handleMessage(msg)
	}
}

func (accp *Acceptor) handleMessage(message v1.Message) {
	ctxLog := log.WithFields(log.Fields{"Addr": accp.GetAddr(), "Method": "Acceptor.handleMessage"})
	switch v := message.(type) {
	case messages.Phase1aMessage:
		phase1aMessage := message.(messages.Phase1aMessage)
		ctxLog.Debugf("ReceivedMessage %T", phase1aMessage)
		if accp.BN == nil || types.Compare(&phase1aMessage.BallotNumber, accp.BN) > 0 {
			ctxLog.Debugf("Adopting ballot %v", phase1aMessage.BallotNumber)
			accp.BN = &phase1aMessage.BallotNumber
		}

		phase1bMessage := messages.NewPhase1bMessage(accp.GetAddr(), *accp.BN, accp.Accepted)
		err := accp.exchange.Send(phase1aMessage.Src(), phase1bMessage)
		if err != nil {
			log.Warnf("accp.exchange.send failed %v", err)
		}

		return

	case messages.Phase2aMessage:
		if accp.BN == nil {
			log.Panicf("unexpected: acceptor.bn is nil")
		}

		phase2aMessage := message.(messages.Phase2aMessage)
		if types.Compare(accp.BN, &phase2aMessage.PValue.BN) == 0 {
			ctxLog.Debugf("Accepted pvalue %v", phase2aMessage.PValue)
			accp.Accepted.Set(phase2aMessage.PValue)
		}

		phase2bMessage := messages.NewPhase2bMessage(accp.GetAddr(), *accp.BN)
		err := accp.exchange.Send(phase2aMessage.Src(), phase2bMessage)
		if err != nil {
			log.Warnf("accp.exchange.send failed %v", err)
		}

		return

	default:
		ctxLog.Panicf("Unknown message type %v", v)
	}
}
