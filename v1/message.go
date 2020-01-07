package v1

import (
	"fmt"
)

// Message - interface intended to be implemented by
// every message exchanged between Paxos Processes
type Message interface {
	// The source Address of this message
	Src() Addr
}

// MessageExchange - Facilitates message exchanges between Paxos processes.
type MessageExchange interface {
	// Send a message to this specific Paxos process identified by its ProcessID
	Send(dest Addr, m Message) error

	// Broadcast a message all Paxos process of a specified type
	SendAll(pt ProcessType, m Message) error

	// Register a process with this exchange
	Register(p ProcessInbox) error

	// UnRegister a process with this exchange
	UnRegister(p ProcessInbox) error
}

// basicMessageExchange - bare-bones implementation of MessageExchange
// here the sendAll method sequentially sends message to each of recipients
type basicMessageExchange struct {
	// Lookup ProcessInbox by the process identifier
	addrToProcessInbox map[Addr]ProcessInbox

	// Lookup processes by the process type
	typeToProcessInbox *typeToProcessMap
}

func (bme basicMessageExchange) Send(dest Addr, m Message) error {
	v, ok := bme.addrToProcessInbox[dest]
	if !ok {
		return fmt.Errorf("not-found: process with id %v not-found", dest)
	}
	return v.Send(m)
}

func (bme basicMessageExchange) SendAll(pt ProcessType, m Message) error {
	entries, ok := bme.typeToProcessInbox.get(pt)
	if !ok || entries.Len() == 0 {
		return fmt.Errorf("not-found: No process(es) with type:%v found", pt)
	}
	for e := entries.Front(); e != nil; e = e.Next() {
		p := e.Value.(ProcessInbox)
		err := p.Send(m)
		if err != nil {
			return fmt.Errorf("send failed: to process=%v %v", p, err)
		}
	}
	return nil
}

func (bme basicMessageExchange) Register(p ProcessInbox) error {
	_, ok := bme.addrToProcessInbox[p.(Addr)]
	if ok {
		return fmt.Errorf("duplicate: process with id %v", p.ID())
	}
	bme.addrToProcessInbox[p.(Addr)] = p
	bme.typeToProcessInbox.put(p)
	return nil
}

func (bme basicMessageExchange) UnRegister(p ProcessInbox) error {
	_, ok := bme.addrToProcessInbox[p.(Addr)]
	if !ok {
		return fmt.Errorf("not-found: process with id %v", p.ID())
	}
	delete(bme.addrToProcessInbox, p.(Addr))
	bme.typeToProcessInbox.remove(p)
	return nil
}
