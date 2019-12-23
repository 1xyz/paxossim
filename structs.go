package paxossim

import (
	"fmt"
	"github.com/1xyz/paxossim/queue"
)

type (
	SlotID int

	Configuration struct {
		Leaders []*Leader
	}

	Process struct {
		pid   string      // A unique process id for this process
		inbox queue.Queue // the incoming message queue for this process
	}

	// Indicates a unique Command issued by the client
	Command interface {
		GetClientID() string
		GetCommandID() string
		GetOp() string
	}

	// BasicCommand -  a unique command issued by the client
	// to a Replica such that: For a unique <ClientID, CommandID>
	// will always map to a unique Op.
	BasicCommand struct {
		ClientID  string // Client unique Id
		CommandID string // Client specific unique Id
		Op        string // Operation associated with this command
	}

	// ReconfigCommand - Represents the Command to re-configure
	// the Configuration
	ReconfigCommand struct {
		BasicCommand
		Config *Configuration // Represents the new updated configuration
	}
)

func NewProcess(pid string) *Process {
	return &Process{
		pid:   pid,
		inbox: queue.NewQueue(),
	}
}

func (p Process) String() string {
	return p.pid
}

func (p Process) SendMessage(m Message) {
	p.inbox.Enqueue(m)
}

func (b *BasicCommand) GetClientID() string {
	return b.ClientID
}

func (b *BasicCommand) GetOp() string {
	return b.Op
}

func (b *BasicCommand) GetCommandID() string {
	return b.CommandID
}

func (b *BasicCommand) String() string {
	return fmt.Sprintf("ClientID: %v CommandID: %v Op: %v",
		b.ClientID, b.CommandID, b.Op)
}
