package messages

import (
	"fmt"
	v1 "github.com/1xyz/paxossim/v1"
	"github.com/1xyz/paxossim/v1/types"
)

type basicMessage struct {
	src v1.Addr
}

func (bm basicMessage) Src() v1.Addr {
	return bm.src
}

func (bm basicMessage) String() string {
	return fmt.Sprintf("source %v", bm.src)
}

// RequestMessage - Request from Client to all Replicas encapsulating a Command
type RequestMessage struct {
	basicMessage
	Command types.Command
}

func (rm RequestMessage) String() string {
	return fmt.Sprintf("RequestMessage: %v command: %v", rm.basicMessage, rm.Command)
}

// DecisionMessage - Decision from the Leader to the Replica with assigned slot for a command
type DecisionMessage struct {
	basicMessage
	Slot    types.Slot
	Command types.Command
}

func NewDecisionMessage(source v1.Addr, slot types.Slot, command types.Command) DecisionMessage {
	return DecisionMessage{
		basicMessage: basicMessage{src: source},
		Slot:         slot,
		Command:      command,
	}
}
func (dm DecisionMessage) String() string {
	return fmt.Sprintf("DecisionMessage: %v slot: %v command: %v",
		dm.basicMessage, dm.Slot, dm.Command)
}

// ProposeMessage - Proposal from the replica to the leader containing a proposed (slot, command)
type ProposeMessage struct {
	basicMessage
	Slot    types.Slot
	Command types.Command
}

func NewProposedMessage(source v1.Addr, slot types.Slot, command types.Command) ProposeMessage {
	return ProposeMessage{
		basicMessage: basicMessage{src: source},
		Slot:         slot,
		Command:      command,
	}
}

func (pm ProposeMessage) String() string {
	return fmt.Sprintf("ProposeMessage: %v slot: %v command: %v",
		pm.basicMessage, pm.Slot, pm.Command)
}
