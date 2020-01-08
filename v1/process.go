package v1

import (
	"fmt"
	"github.com/1xyz/paxossim/queue"
)

// ProcessType - The different paxos process types
type ProcessType int

const (
	Acceptor ProcessType = iota
	Commander
	Leader
	Replica
	Scout
	Client
)

var ptStrings = map[ProcessType]string{
	Acceptor:  "Acceptor",
	Commander: "Commander",
	Leader:    "leader",
	Replica:   "Replica",
	Scout:     "Scout",
	Client:    "Client",
}

// stringer implementation for ProcessType
func (pt ProcessType) String() string {
	v, ok := ptStrings[pt]
	if !ok {
		return "Unknown"
	} else {
		return v
	}
}

// ProcessInbox Identifier in this Paxos system
type ProcessID int

type Addr interface {
	// Return a unique identifier (aka. address) for this process
	ID() ProcessID

	// Return the ProcessInbox type for this process
	Type() ProcessType
}

// ProcessInbox - interface allowing a Paxos process to be addressed & sent messages
type ProcessInbox interface {
	// Address of the inbox
	Addr

	// Send a message to this process
	Send(m Message) error
}

// ProcessOutbox - interface allowing a process to recv messages
type ProcessOutbox interface {
	// Recv for the next message
	Recv() (Message, error)
}

// Process - interface allowing a runnable entity with message I/O
type Process interface {
	ProcessInbox
	ProcessOutbox
	GetAddr() Addr
}

func NewProcess(id ProcessID, pt ProcessType) Process {
	return newBasicProcess(id, pt)
}

type basicAddr struct {
	// A globally unique process identifier for this
	id ProcessID

	// A Process type
	pt ProcessType
}

func (b basicAddr) ID() ProcessID {
	return b.id
}

func (b basicAddr) Type() ProcessType {
	return b.Type()
}

func (b basicAddr) String() string {
	return fmt.Sprintf("(%v-%v)", b.pt, b.id)
}

type basicProcess struct {
	basicAddr

	// inbox of incoming messages
	inbox queue.Queue
}

func newBasicProcess(id ProcessID, pt ProcessType) *basicProcess {
	return &basicProcess{
		basicAddr: basicAddr{
			id: id,
			pt: pt,
		},
		inbox: queue.NewQueue(),
	}
}

func (b basicProcess) ID() ProcessID {
	return b.id
}

func (b basicProcess) Type() ProcessType {
	return b.pt
}

func (b basicProcess) GetAddr() Addr {
	return basicAddr{
		id: b.ID(),
		pt: b.Type(),
	}
}

func (b basicProcess) Send(m Message) error {
	b.inbox.Enqueue(m)
	return nil
}

func (b basicProcess) Recv() (Message, error) {
	entry, ok := b.inbox.WaitForItem().(Message)
	if !ok {
		return nil, fmt.Errorf("cast-error message entry not found %T", entry)
	}
	return entry, nil
}
