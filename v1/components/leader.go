package components

import (
	v1 "github.com/1xyz/paxossim/v1"
	"github.com/1xyz/paxossim/v1/messages"
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

func (leader *Leader) Run() {
	ctxLog := log.WithFields(log.Fields{"Addr": leader.GetAddr()})
	ctxLog.Debugf("Running Leader")
	leader.spawnNewScout()
	for {
		msg, err := leader.Process.Recv()
		if err != nil {
			ctxLog.Panicf("error in inbox recv %v", err)
		}

		leader.handleMessage(msg)
	}
}

func (leader *Leader) spawnNewScout() {
	ctxLog := log.WithFields(log.Fields{"Addr": leader.GetAddr()})
	s := NewScout(leader.exchange, leader.GetAddr(), leader.acceptors, leader.ballotNumber)
	go s.Run()
	ctxLog.Debugf("Spawned a new Scout")
}

func (leader *Leader) spawnNewCommander(slot types.Slot) {
	ctxLog := log.WithFields(log.Fields{"Addr": leader.GetAddr()})
	command, found := leader.proposals.Get(slot)
	if !found {
		log.Panicf("no command found for slot %v", slot)
	}

	pValue := types.PValue{
		BN:      leader.ballotNumber,
		Slot:    slot,
		Command: command,
	}
	c := NewCommander(leader.exchange, leader.GetAddr(), leader.acceptors, pValue)
	go c.Run()
	ctxLog.Debugf("Spawned a new Commander")
}

func (leader *Leader) handleMessage(message v1.Message) {
	ctxLog := log.WithFields(log.Fields{
		"Addr": leader.GetAddr(),
		"Method": "leader.handleMessage",
	})
	ctxLog.Debugf("Recd a message of type %T", message)

	switch v := message.(type) {
	case messages.ProposeMessage:
		pm := message.(messages.ProposeMessage)
		// Check if this slot has already been assigned here
		if leader.proposals.Contains(pm.Slot) {
			ctxLog.Debugf("the corresponding slot %v has been assigned", pm.Slot)
			return
		}

		// Assign the slot to this command in this leader and spawn a new commander
		leader.proposals.Assign(pm.Slot, pm.Command)
		if !leader.active {
			return
		}

		leader.spawnNewCommander(pm.Slot)

	case messages.AdoptedMessage:
		am := message.(messages.AdoptedMessage)
		if types.Compare(&leader.ballotNumber, &am.BallotNumber) != 0 {
			return
		}

		pMax := make(map[types.Slot]types.BallotNumber)
		for pv, _ := range am.Accepted {
			e, ok := pMax[pv.Slot]
			if !ok || (types.Compare(&e, &pv.BN) < 0) {
				pMax[pv.Slot] = pv.BN
				leader.proposals.Assign(pv.Slot, pv.Command)
			}
		}

		for slot, _ := range leader.proposals {
			leader.spawnNewCommander(slot)
		}

		// Activate the leader
		leader.active = true

	case messages.PreemptMessage:
		pm := message.(messages.PreemptMessage)
		res := types.Compare(&pm.BallotNumber, &leader.ballotNumber)
		if res <= 0 {
			ctxLog.Debugf("expected remote ballot-number to be greater could be a delayed message")
			return
		}

		leader.active = false
		leader.ballotNumber.Round++
		leader.spawnNewScout()

	default:
		log.Panicf("Unknown message type %v", v)
	}
}
