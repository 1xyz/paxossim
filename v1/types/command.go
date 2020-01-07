package types

import (
	"fmt"
	v1 "github.com/1xyz/paxossim/v1"
)

// Represents a slot which is assigned to a Command in Paxos
type Slot int

// The initial slot ID
const InitialSlotID Slot = 1

type Command interface {
	// The client unique id
	ClientID() string

	// The command unique id
	CommandID() string

	// Operation
	Op() string
}

// Indicates a unique Command issued by the client to a Replica
// For a unique <ClientID, CommandID> will always map to a unique Op.
type basicCommand struct {
	// Client unique Id
	clientID string

	// Unique command ID
	commandID string

	// Operation associated with this command
	op string
}

func (b basicCommand) ClientID() string {
	return b.clientID
}

func (b basicCommand) Op() string {
	return b.op
}

func (b basicCommand) CommandID() string {
	return b.commandID
}

func (b basicCommand) String() string {
	return fmt.Sprintf("ClientID: %v CommandID: %v Op: %v",
		b.clientID, b.commandID, b.op)
}

// A Reconfiguration Command issued
type reConfigCommand struct {
	basicCommand

	// New Leader configuration
	NewLeaders []v1.Addr
}
