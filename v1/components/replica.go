package components

import (
	v1 "github.com/1xyz/paxossim/v1"
	"github.com/1xyz/paxossim/v1/messages"
	"github.com/1xyz/paxossim/v1/types"
	log "github.com/sirupsen/logrus"
)

const (
	Window             types.Slot = 5
	InitialSlotID      types.Slot = 1
	InitialRequestSize            = 100
)

var replicaCount = 0

type Replica struct {
	v1.Process

	exchange v1.MessageExchange

	// Index of the next slot which can be proposed
	slotIn types.Slot

	// Index of the slot for which a decision needs to be made
	slotOut types.Slot

	// Requests which have not been proposed or decided
	requests []types.Command

	// Requests which have been proposed but not decided, indexed by slot
	proposals types.SlotCommandMap

	// Requests which have been decided, indexed by the slot
	decisions types.SlotCommandMap

	// Configuration; primarily the leader configuration
	leaders []v1.Addr
}

func NewReplica(exchange v1.MessageExchange, leaders []v1.Addr) *Replica {
	processID := replicaCount
	replicaCount++

	r := &Replica{
		Process:   v1.NewProcess(v1.ProcessID(processID), v1.Replica),
		slotIn:    InitialSlotID,
		slotOut:   InitialSlotID,
		requests:  make([]types.Command, 0, InitialRequestSize),
		proposals: make(types.SlotCommandMap),
		decisions: make(types.SlotCommandMap),
		exchange:  exchange,
		leaders:   leaders,
	}

	exchange.Register(r)
	return r
}

func (r *Replica) Run() {
	ctxLog := log.WithFields(log.Fields{
		"id": r.ID(), "type": r.Type(),
	})

	for {
		msg, err := r.Process.Recv()
		if err != nil {
			ctxLog.Panicf("error in inbox recv %v", err)
		}

		r.handleMessage(msg)
		r.propose()
	}
}

func (r *Replica) handleMessage(message v1.Message) {
	ctxLog := log.WithFields(log.Fields{
		"id": r.ID(), "type": r.Type(),
	})

	switch v := message.(type) {
	case messages.RequestMessage:
		rm := message.(messages.RequestMessage)
		ctxLog.Debugf("ReqMessage %v", rm)
		r.requests = append(r.requests, rm.Command)

	default:
		log.Panicf("Unknown message type %v", v)
	}
}

func (r *Replica) propose() {

}
