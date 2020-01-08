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
		ctxLog.Debugf("%v", rm)
		r.requests = append(r.requests, rm.Command)

	case messages.DecisionMessage:
		dm := message.(messages.DecisionMessage)
		ctxLog.Debugf("%v", dm)

		// ToDo: check to see if a command already exists
		// record the slot for the decided command
		r.decisions[dm.Slot] = dm.Command

		// run through all decisions starting from slotOut
		// and attempt to apply them until we find an undecided slot
		for r.decisions.Contains(r.slotOut) {
			decidedCmd := r.decisions[r.slotOut]
			// check to see if this replica made a proposal for this slotOut
			if r.proposals.Contains(r.slotOut) {
				proposedCmd := r.proposals[r.slotOut]
				if proposedCmd != decidedCmd {
					// looks like the leader decided another slot for slotOut
					// ReQueue this command back to the request queue
					r.requests = append(r.requests, proposedCmd)
				}
				// this command is either re-queued or decided so remove
				// from the proposal queue
				r.proposals.Remove(r.slotOut)
			}
			r.perform(decidedCmd)
			r.slotOut++
		}

	default:
		log.Panicf("Unknown message type %v", v)
	}
}

// propose - if there are any pending requests, create proposals by assigning slots to
// the request's command and send to leaders
func (r *Replica) propose() {
	for {
		// check to see if the requests queue is empty or if we have reached the window limit
		if len(r.requests) == 0 || r.slotIn >= (r.slotOut+Window) {
			break
		}

		// Dequeue this request from the requests queue
		req := r.requests[0]
		r.requests = r.requests[1:]

		// check to see if a reconfiguration command needs to be applied
		if r.slotIn > Window && r.decisions.Contains(r.slotIn-Window) {
			cmd, ok := r.proposals[r.slotIn-Window].(types.ReConfigCommand)
			if ok {
				log.Debugf("Updating configuration %v", cmd.NewLeaders)
				r.leaders = cmd.NewLeaders
			}
		}

		// enqueue this proposal and sent it to all leaders
		r.proposals[r.slotIn] = req
		pm := messages.NewProposedMessage(r, r.slotIn, req)
		for _, addr := range r.leaders {
			r.exchange.Send(addr, pm)
		}
		r.slotIn++
	}
}

func (r *Replica) perform(command types.Command) {
	// Different replicas might have proposed the same command for
	// different slots. In this case we don't really want to apply
	// the command at this replica more than once
	for slot := InitialSlotID; slot < r.slotOut; slot++ {
		// log.Infof("slot %v command %v decisions[slot]: %v", slot, command, r.Decisions[slot])
		if r.decisions[slot] == command {
			// log.Infof("Command %v detected in decision history", command)
			return
		}
	}

	recfgCommand, ok := command.(types.ReConfigCommand)
	if ok {
		log.Debugf("Reconfig command %v", recfgCommand)
		// ToDo: we needd to apply the updated configuration
		return
	}

	// ToDo: Apply state here!!
	log.Infof("(%v, %v, %v) r=%v-%v",
		command.GetClientID(), command.GetCommandID(), command.GetOp(), r.Type(), r.ID())
}
