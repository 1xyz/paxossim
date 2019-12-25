package paxossim

import log "github.com/sirupsen/logrus"

type (
	Acceptor struct {
		*Process               // Incoming process and mailbox
		Accepted PValues       // Set of PValues accepted so far
		BN       *BallotNumber // Last Adopted ballot number
	}
)

func NewAcceptor(acceptorID string) *Acceptor {
	return &Acceptor{
		Process:  NewProcess(acceptorID),
		Accepted: make(PValues),
		BN: &BallotNumber{
			Round:    0,
			LeaderID: "",
		},
	}
}

func (a *Acceptor) Run() {
	for {
		msg := a.inbox.WaitForItem()
		a.HandleMsg(msg.(Message))
	}
}

func (a *Acceptor) HandleMsg(msg Message) {
	ctxLog := log.WithFields(log.Fields{
		"id": a.pid,
	})
	switch v := msg.(type) {
	case *P1aMessage:
		p1Req := msg.(*P1aMessage)
		if p1Req.BN.CompareTo(a.BN) > 0 {
			a.BN = p1Req.BN
		}
		p1Resp := NewP1bMessage(a.pid, a.BN, &a.Accepted)
		p1Req.S.SendMessage(p1Resp)
		ctxLog.Debugf("P1bMessage to %v", p1Req.Src)
		break

	case *P2aMessage:
		p2Req := msg.(*P2aMessage)
		if p2Req.PV.BN.CompareTo(a.BN) == 0 {
			a.Accepted.set(p2Req.PV)
		}
		p2Resp := NewP2bMessage(a.pid, a.BN)
		p2Req.S.SendMessage(p2Resp)
		break

	default:
		ctxLog.Panicf("Unknown message type %v", v)
	}
}
