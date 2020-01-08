package components

import (
	v1 "github.com/1xyz/paxossim/v1"
	"github.com/1xyz/paxossim/v1/types"
	log "github.com/sirupsen/logrus"
)

var leaderCount = 0

type Leader struct {
	v1.Process

	exchange v1.MessageExchange

	ballotNumber types.BallotNumber

	active bool

	proposals types.SlotCommandMap

	acceptors []v1.Addr
}

func NewLeader(exchange v1.MessageExchange, acceptors []v1.Addr) *Leader {
	processID := leaderCount
	leaderCount++

	p := v1.NewProcess(v1.ProcessID(processID), v1.Leader)
	l := &Leader{
		Process:   p,
		exchange:  exchange,
		proposals: make(types.SlotCommandMap),
		active:    false,
		acceptors: acceptors,
		ballotNumber: types.BallotNumber{
			Round:    0,
			LeaderID: p.GetAddr(),
		},
	}

	ctxLog := log.WithFields(log.Fields{"Addr": l.GetAddr()})
	ctxLog.Debugf("Created leader")

	err := exchange.Register(l)
	if err != nil {
		log.Panicf("exchange.Register error %v", err)
	}
	return l
}
