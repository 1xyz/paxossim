package paxossim

import (
	log "github.com/sirupsen/logrus"
)

const (
	Window             SlotID = 5
	InitialSlotID      SlotID = 1
	InitialRequestSize        = 100
)

type (
	SlotCommandMap map[SlotID]Command
	Replica        struct {
		*Process                 // Incoming process and mailbox
		slotIn    SlotID         // Index of the next slot which can be proposed
		slotOut   SlotID         // Index of the slot for which a decision needs to be Made
		requests  []Command      // Requests which have not been proposed or decided
		proposals SlotCommandMap // Requests which have been proposed but not decided indexed by slot number
		decisions SlotCommandMap // Requests which have been decided indexed by the slot number
		config    *Configuration // Configuration; primarily the leader configuration
	}
)

func (s SlotCommandMap) contains(slot SlotID) bool {
	_, ok := s[slot]
	return ok
}

func (s SlotCommandMap) remove(slot SlotID) {
	delete(s, slot)
}

func NewReplica(replicaID string, initialConfig *Configuration) *Replica {
	return &Replica{
		Process:   NewProcess(replicaID),
		slotIn:    InitialSlotID,
		slotOut:   InitialSlotID,
		requests:  make([]Command, 0, InitialRequestSize),
		proposals: make(SlotCommandMap),
		decisions: make(SlotCommandMap),
		config:    initialConfig,
	}
}

func (r *Replica) Run() {
	ctxLog := log.WithFields(log.Fields{
		"id": r.pid,
	})
	ctxLog.Debug("Run")
	for {
		msg := r.inbox.WaitForItem()
		switch v := msg.(type) {
		case *RequestMessage:
			rm := msg.(*RequestMessage)
			ctxLog.Debugf("ReqMessage %v", rm)
			r.requests = append(r.requests, rm.C)

		case *DecisionMessage:
			dm := msg.(*DecisionMessage)
			ctxLog.Debugf("DecisionMessage %v", dm)
			// ToDo: check to see if a command already exists
			r.decisions[dm.SlotID] = dm.C
			for r.decisions.contains(r.slotOut) {
				decidedCmd := r.decisions[r.slotOut]
				if r.proposals.contains(r.slotOut) {
					proposedCmd := r.proposals[r.slotOut]
					if proposedCmd != decidedCmd {
						r.requests = append(r.requests, proposedCmd)
					}
					r.proposals.remove(r.slotOut)
				}
				r.perform(decidedCmd)
				r.slotOut++
			}

		default:
			log.Panicf("Unknown message type %v", v)
		}
		r.propose()
	}
}

func (r *Replica) perform(c Command) {
	for slot := InitialSlotID; slot < r.slotOut; slot++ {
		if r.decisions[slot] == c {
			log.Debugf("Command %v detected in decision history", c)
			return
		}
	}

	reconfigCmd, ok := c.(*ReconfigCommand)
	if ok {
		log.Debugf("Reconfig command %v", reconfigCmd)
		return
	}

	// ToDo: Apply state here!!
}

func (r *Replica) propose() {
	for  {
		if len(r.requests) == 0 ||  r.slotIn >= (r.slotOut + Window) {
			break
		}

		req := r.requests[0]
		r.requests = r.requests[1:]
		if r.slotIn > Window && r.decisions.contains(r.slotIn-Window) {
			cmd, ok := r.proposals[r.slotIn-Window].(*ReconfigCommand)
			if ok {
				log.Debugf("Updating configuration %v", cmd.Config)
				r.config = cmd.Config
			}
		}

		r.proposals[r.slotIn] = req
		for _, leader := range r.config.Leaders {
			leader.SendMessage(NewProposeMessage(r.pid, r.slotIn, req))
		}
		r.slotIn++
	}
}
