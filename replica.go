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
		SlotIn    SlotID         // Index of the next slot which can be proposed
		SlotOut   SlotID         // Index of the slot for which a decision needs to be Made
		Requests  []Command      // Requests which have not been proposed or decided
		Proposals SlotCommandMap // Requests which have been proposed but not decided indexed by slot number
		Decisions SlotCommandMap // Requests which have been decided indexed by the slot number
		Config    *Configuration // Configuration; primarily the leader configuration
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
		SlotIn:    InitialSlotID,
		SlotOut:   InitialSlotID,
		Requests:  make([]Command, 0, InitialRequestSize),
		Proposals: make(SlotCommandMap),
		Decisions: make(SlotCommandMap),
		Config:    initialConfig,
	}
}

func (r *Replica) Run() {
	for {
		msg := r.inbox.WaitForItem()
		r.HandleMsg(msg.(Message))
		r.Propose()
	}
}

func (r *Replica) HandleMsg(msg Message) {
	ctxLog := log.WithFields(log.Fields{
		"id": r.pid,
	})
	switch v := msg.(type) {
	case *RequestMessage:
		rm := msg.(*RequestMessage)
		ctxLog.Debugf("ReqMessage %v", rm)
		r.Requests = append(r.Requests, rm.C)

	case *DecisionMessage:
		dm := msg.(*DecisionMessage)
		ctxLog.Debugf("DecisionMessage %v", dm)
		// ToDo: check to see if a command already exists
		r.Decisions[dm.SlotID] = dm.C
		for r.Decisions.contains(r.SlotOut) {
			decidedCmd := r.Decisions[r.SlotOut]
			if r.Proposals.contains(r.SlotOut) {
				proposedCmd := r.Proposals[r.SlotOut]
				if proposedCmd != decidedCmd {
					r.Requests = append(r.Requests, proposedCmd)
				}
				r.Proposals.remove(r.SlotOut)
			}
			r.Perform(decidedCmd)
			r.SlotOut++
		}

	default:
		log.Panicf("Unknown message type %v", v)
	}
}

func (r *Replica) Perform(c Command) {
	// Different replicas might have proposed the same command for
	// different slots. In this case we don't really want to apply
	// the command at this replica more than once
	for slot := InitialSlotID; slot < r.SlotOut; slot++ {
		if r.Decisions[slot] == c {
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

func (r *Replica) Propose() {
	for {
		if len(r.Requests) == 0 || r.SlotIn >= (r.SlotOut+Window) {
			break
		}

		req := r.Requests[0]
		r.Requests = r.Requests[1:]
		if r.SlotIn > Window && r.Decisions.contains(r.SlotIn-Window) {
			cmd, ok := r.Proposals[r.SlotIn-Window].(*ReconfigCommand)
			if ok {
				log.Debugf("Updating configuration %v", cmd.Config)
				r.Config = cmd.Config
			}
		}

		r.Proposals[r.SlotIn] = req
		for _, leader := range r.Config.Leaders {
			leader.SendMessage(NewProposeMessage(r.pid, r.SlotIn, req))
		}
		r.SlotIn++
	}
}
