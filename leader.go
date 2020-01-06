package paxossim

import (
	"fmt"
	log "github.com/sirupsen/logrus"
)

type Leader struct {
	*Process
	BN        *BallotNumber
	Active    bool
	Proposals *SlotCommandMap
	Acceptors []Entity
	Replicas  []Entity
}

func NewLeader(leaderID string, acceptors []Entity) *Leader {
	proposals := make(SlotCommandMap)
	replicas := make([]Entity, 0, 3)
	return &Leader{
		Process: NewProcess(leaderID),
		BN: &BallotNumber{
			Round:    0,
			LeaderID: leaderID,
		},
		Proposals: &proposals,
		Acceptors: acceptors,
		Replicas:  replicas,
		Active:    false,
	}
}

// Append a replica to this leader
func (leader *Leader) AppendReplica(replica Entity) {
	leader.Replicas = append(leader.Replicas, replica)
}

// spawnNewScout - creates a new scout with relevant information
// and schedules it to be run in a separate go-routine
func (leader *Leader) spawnNewScout() {
	// We need to make a copy of this ballot number
	// instead of sending this ballot as a pointer.
	bn := *leader.BN
	scoutID := fmt.Sprintf("scout-%v", bn)
	s := NewScout(scoutID, leader, leader.Acceptors, &bn)
	go s.Run()
}

func (leader *Leader) spawnNewCommander(slot SlotID) {
	// we need to make a copy of this ballot number
	// instead of sending this ballot as a pointer.
	bn := *leader.BN
	cmd, _ := leader.Proposals.Get(slot)
	pv := &PValue{
		BN:   &bn,
		Slot: slot,
		C:    cmd,
	}
	commanderID := fmt.Sprintf("commander-%v", pv)
	c := NewCommander(commanderID, leader, leader.Acceptors, leader.Replicas, pv)
	go c.Run()
	log.WithFields(log.Fields{
		"id":   leader.pid,
		"bn":   bn,
		"slot": slot,
		"C":    cmd,
	}).Debugf("new commander")
}

func (leader *Leader) Run() {
	ctxLog := log.WithFields(log.Fields{"id": leader.pid})
	leader.spawnNewScout()
	for {
		entry := leader.inbox.WaitForItem()
		msg, ok := entry.(Message)
		if !ok {
			ctxLog.Panicf("Un-Handled message type recv %v", msg)
		}

		ctxLog.Debugf("Got a new message %v", msg)
		leader.HandleMessage(msg)
	}
}

func (leader *Leader) HandleMessage(msg Message) {
	ctxLog := log.WithFields(log.Fields{
		"id": leader.pid,
	})

	switch v := msg.(type) {
	case *ProposeMessage:
		pm := msg.(*ProposeMessage)
		ctxLog.Debugf("ProposedMessage %v", pm)
		// Check if this slot has already been assigned here
		if leader.Proposals.Contains(pm.SlotID) {
			ctxLog.Debugf("the corresponding slot %v has been assigned", pm.SlotID)
			return
		}
		// Assign the slot to this command in this leader
		// and spawn a new commander
		leader.Proposals.Assign(pm.SlotID, pm.C)
		if !leader.Active {
			return
		}
		leader.spawnNewCommander(pm.SlotID)

	case *AdoptedMessage:
		am := msg.(*AdoptedMessage)
		ctxLog.Debugf("AdoptedMessage %v", am)
		if leader.BN.CompareTo(am.BN) != 0 {
			return
		}
		pmax := make(map[SlotID]*BallotNumber)
		for pv, _ := range *am.Accepted {
			e, ok := pmax[pv.Slot]
			if !ok || (e.CompareTo(pv.BN) < 0) {
				pmax[pv.Slot] = pv.BN
				leader.Proposals.Assign(pv.Slot, pv.C)
			}
		}
		for slot, _ := range *leader.Proposals {
			leader.spawnNewCommander(slot)
		}
		// Activate the leader
		leader.Active = true

	case *PreemptMessage:
		pm := msg.(*PreemptMessage)
		ctxLog.Debugf("Pre-empt message %v", pm)
		res := pm.BN.CompareTo(leader.BN)
		if res <= 0 {
			ctxLog.Debugf("expected remote ballot-number to be greater could be a delayed message")
			return
		}
		leader.Active = false
		leader.BN.Round++
		leader.spawnNewScout()

	default:
		log.Panicf("Unknown message type %v", v)
	}
}
