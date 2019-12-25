package paxossim

import (
	"fmt"
)

type (
	// An Message in this Paxos system
	Message interface {
		Source() string
	}

	BasicMessage struct {
		Src string
	}

	// Request from Client to all Replicas encapsulating a Command
	RequestMessage struct {
		*BasicMessage
		C Command
	}

	// Decision from the Leader to the Replica containing
	// the slot and its assigned command
	DecisionMessage struct {
		*BasicMessage
		SlotID SlotID
		C      Command
	}

	// Proposal from the replica to the leader containing
	// a proposed slot and an associated command
	ProposeMessage struct {
		*BasicMessage
		SlotID SlotID
		C      Command
	}

	PhaseMessage struct {
		*BasicMessage
		S Entity
	}

	// Message sent by the Leader(Scout) to the
	// Acceptors containing the BallotNumber during the
	// Phase1 of Paxos
	P1aMessage struct {
		*PhaseMessage
		BN *BallotNumber
	}

	// Message sent by the Acceptor in response to the
	// P1aMessage containing the current BallotNumber and
	// the list of accepted PValues
	P1bMessage struct {
		*BasicMessage
		BN       *BallotNumber
		Accepted *PValues
	}

	// Message sent by the Leader(Commander) to the
	// Acceptors containing the PValue (BallotNum, Slot, Command)
	P2aMessage struct {
		*PhaseMessage
		PV *PValue
	}

	// Message returned by the Acceptor back to Commander
	// as a response to the P2aMessage
	P2bMessage struct {
		*BasicMessage
		BN *BallotNumber
	}
)

func (b *BasicMessage) Source() string {
	return b.Src
}

func (b *BasicMessage) String() string {
	return fmt.Sprintf("Source: %v", b.Src)
}

func NewRequestMessage(source string, command Command) *RequestMessage {
	return &RequestMessage{
		BasicMessage: &BasicMessage{source},
		C:            command,
	}
}

func (r *RequestMessage) String() string {
	return fmt.Sprintf("%v Command: %v",
		r.BasicMessage, r.C)
}

func NewProposeMessage(source string, id SlotID, command Command) *ProposeMessage {
	return &ProposeMessage{
		BasicMessage: &BasicMessage{source},
		SlotID:       id,
		C:            command,
	}
}

func (p *ProposeMessage) String() string {
	return fmt.Sprintf("%v Slot: %v Command: %v",
		p.BasicMessage, p.SlotID, p.C)
}

func NewDecisionMessage(source string, id SlotID, command Command) *DecisionMessage {
	return &DecisionMessage{
		BasicMessage: &BasicMessage{source},
		SlotID:       id,
		C:            command,
	}
}

func (d *DecisionMessage) String() string {
	return fmt.Sprintf("%v Slot: %v Command: %v",
		d.BasicMessage, d.SlotID, d.C)
}

func NewP1aMessage(source string, entity Entity, bn *BallotNumber) *P1aMessage {
	return &P1aMessage{
		PhaseMessage: &PhaseMessage{
			BasicMessage: &BasicMessage{Src: source},
			S:            entity,
		},
		BN: bn,
	}
}

func NewP1bMessage(source string, bn *BallotNumber, pv *PValues) *P1bMessage {
	return &P1bMessage{
		BasicMessage: &BasicMessage{Src: source},
		BN:           bn,
		Accepted:     pv,
	}
}

func NewP2aMessage(source string, entity Entity, pv *PValue) *P2aMessage {
	return &P2aMessage{
		PhaseMessage: &PhaseMessage{
			BasicMessage: &BasicMessage{Src: source},
			S:            entity,
		},
		PV: pv,
	}
}

func NewP2bMessage(source string, bn *BallotNumber) *P2bMessage {
	return &P2bMessage{
		BasicMessage: &BasicMessage{Src: source},
		BN:           bn,
	}
}
