package paxossim

import "fmt"

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
