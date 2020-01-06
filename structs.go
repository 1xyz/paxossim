package paxossim

import (
	"fmt"
	"github.com/1xyz/paxossim/queue"
	"strings"
)

type (
	SlotID int

	Configuration struct {
		Leaders []Entity
	}

	Entity interface {
		ID() string
		Run()
		SendMessage(m Message)
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

	BallotNumber struct {
		Round    int
		LeaderID string
	}

	PValue struct {
		BN   *BallotNumber
		Slot SlotID
		C    Command
	}

	PValues map[PValue]bool

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

func NewConfiguration(nLeaders int) *Configuration {
	leaders := make([]Entity, 0, nLeaders)
	return &Configuration{Leaders: leaders}
}

func (c *Configuration) AppendLeader(leader Entity) {
	c.Leaders = append(c.Leaders, leader)
}

func NewProcess(pid string) *Process {
	return &Process{
		pid:   pid,
		inbox: queue.NewQueue(),
	}
}

func (p Process) ID() string {
	return p.pid
}

func (p Process) String() string {
	return p.ID()
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

// Compare returns an integer comparing two BallotNumbers lexicographically.
// The result will be:
//   0 if bn == otherBn,
//   -1 if bn < otherBn, and
//   +1 if bn > otherBn.
func (bn *BallotNumber) CompareTo(otherBn *BallotNumber) int {
	if bn == otherBn {
		return 0
	}
	c1 := bn.Round - otherBn.Round
	if c1 == 0 {
		return strings.Compare(bn.LeaderID, otherBn.LeaderID)
	} else if c1 > 0 {
		return 1
	} else {
		return -1
	}
}

func (bn *BallotNumber) String() string {
	return fmt.Sprintf("(%d, %s)", bn.Round, bn.LeaderID)
}

func (pValue *PValue) String() string {
	return fmt.Sprintf("(%v, %v, %v)",
		pValue.BN, pValue.Slot, pValue.C)
}

func (pvalues *PValues) Set(value *PValue) {
	_, ok := (*pvalues)[*value]
	if !ok {
		(*pvalues)[*value] = true
	}
}

func (pvalues *PValues) Contains(value *PValue) bool {
	_, ok := (*pvalues)[*value]
	return ok
}

func (pvalues *PValues) Update(values *PValues) {
	for v, _ := range *values {
		pvalues.Set(&v)
	}
}

type StringSet map[string]bool

func (ss StringSet) Contains(key string) bool {
	_, ok := ss[key]
	return ok
}

func (ss StringSet) Add(key string) {
	ss[key] = true
}

func (ss StringSet) Remove(key string) {
	delete(ss, key)
}

func (ss StringSet) Len() int {
	return len(ss)
}

type SlotCommandMap map[SlotID]Command

func (s SlotCommandMap) Get(slot SlotID) (Command, bool) {
	result, ok := s[slot]
	return result, ok
}

func (s SlotCommandMap) Contains(slot SlotID) bool {
	_, ok := s[slot]
	return ok
}

func (s SlotCommandMap) Remove(slot SlotID) {
	delete(s, slot)
}

func (s SlotCommandMap) Assign(slot SlotID, c Command) {
	s[slot] = c
}
