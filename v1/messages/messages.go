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

func NewRequestMessage(source v1.Addr, command types.Command) RequestMessage {
	return RequestMessage{
		basicMessage: basicMessage{src: source},
		Command:      command,
	}
}

func (rm RequestMessage) String() string {
	return fmt.Sprintf("RequestMessage: %v command: %v", rm.basicMessage, rm.Command)
}

// DecisionMessage - Decision from the leader to the Replica with assigned slot for a command
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

// Message sent by the leader(Scout) to the  acceptors containing the BallotNumber during the Phase1 of Paxos
type Phase1aMessage struct {
	basicMessage
	BallotNumber types.BallotNumber
}

func NewPhase1aMessage(source v1.Addr, number types.BallotNumber) Phase1aMessage {
	return Phase1aMessage{
		basicMessage: basicMessage{src: source},
		BallotNumber: number,
	}
}

// Message sent by the Acceptor in response to the Phase1aMessage containing the current BallotNumber and
// the list of accepted PValues
type Phase1bMessage struct {
	basicMessage
	BallotNumber types.BallotNumber
	PValues      types.PValues
}

func NewPhase1bMessage(source v1.Addr, number types.BallotNumber, values types.PValues) Phase1bMessage {
	return Phase1bMessage{
		basicMessage: basicMessage{src: source},
		BallotNumber: number,
		PValues:      values,
	}
}

// Message sent by the leader(Commander) to the acceptors containing the PValue (BallotNum, Slot, Command)
type Phase2aMessage struct {
	basicMessage
	PValue types.PValue
}

func NewPhase2aMessage(source v1.Addr, value types.PValue) Phase2aMessage {
	return Phase2aMessage{
		basicMessage: basicMessage{src: source},
		PValue:       value,
	}
}

// Message returned by the Acceptor back to Commander as a response to the Phase2aMessage
type Phase2bMessage struct {
	basicMessage
	BallotNumber types.BallotNumber
}

func NewPhase2bMessage(addr v1.Addr, number types.BallotNumber) Phase2bMessage {
	return Phase2bMessage{
		basicMessage: basicMessage{src: addr},
		BallotNumber: number,
	}
}

// Message sent by the Scout or a Commander indicating that a ballot-number is pre-empted by a new ballot number
type PreemptMessage struct {
	basicMessage
	BallotNumber types.BallotNumber
}

func NewPremptedMessage(addr v1.Addr, number types.BallotNumber) PreemptMessage {
	return PreemptMessage{
		basicMessage: basicMessage{src: addr},
		BallotNumber: number,
	}
}

// Message sent by the Scout to the leader on a successful adoption of ballot by majority of the acceptors.
type AdoptedMessage struct {
	basicMessage
	BallotNumber types.BallotNumber
	Accepted     types.PValues
}

func NewAdoptedMessage(addr v1.Addr, number types.BallotNumber, values types.PValues) AdoptedMessage {
	return AdoptedMessage{
		basicMessage: basicMessage{src: addr},
		BallotNumber: number,
		Accepted:     values,
	}
}
