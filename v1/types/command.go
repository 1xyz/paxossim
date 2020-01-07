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
	GetClientID() string

	// The command unique id
	GetCommandID() string

	// Operation
	GetOp() string
}

// Indicates a unique Command issued by the client to a Replica
// For a unique <ClientID, CommandID> will always map to a unique Op.
type BasicCommand struct {
	// Client unique Id
	ClientID string

	// Unique command ID
	CommandID string

	// Operation associated with this command
	Op string
}

func (b BasicCommand) GetClientID() string {
	return b.ClientID
}

func (b BasicCommand) GetOp() string {
	return b.Op
}

func (b BasicCommand) GetCommandID() string {
	return b.CommandID
}

func (b BasicCommand) String() string {
	return fmt.Sprintf("ClientID: %v CommandID: %v Op: %v",
		b.ClientID, b.CommandID, b.Op)
}

// A Reconfiguration Command issued
type ReConfigCommand struct {
	BasicCommand

	// New Leader configuration
	NewLeaders []v1.Addr
}

type SlotCommandMap map[Slot]Command

func (s SlotCommandMap) Get(slot Slot) (Command, bool) {
	result, ok := s[slot]
	return result, ok
}

func (s SlotCommandMap) Contains(slot Slot) bool {
	_, ok := s[slot]
	return ok
}

func (s SlotCommandMap) Remove(slot Slot) {
	delete(s, slot)
}

func (s SlotCommandMap) Assign(slot Slot, c Command) {
	s[slot] = c
}
