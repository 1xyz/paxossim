package v1

import (
	"fmt"
	log "github.com/sirupsen/logrus"
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

func NewMessageExchange() MessageExchange {
	return &basicMessageExchange{
		addrToProcessInbox: make(map[Addr]ProcessInbox),
		typeToProcessInbox: make(typeToProcessMap),
	}
}

// basicMessageExchange - bare-bones implementation of MessageExchange
// here the sendAll method sequentially sends message to each of recipients
type basicMessageExchange struct {
	// ToDo: in retrospect, it is not a good idea to use an interface
	// ToDo: as a key, we go around it by creating a default instance
	// Lookup ProcessInbox by the process identifier
	addrToProcessInbox map[Addr]ProcessInbox

	// Lookup processes by the process type
	typeToProcessInbox typeToProcessMap
}

func (bme basicMessageExchange) Send(dest Addr, m Message) error {
	addr := NewAddress(dest.ID(), dest.Type())
	v, ok := bme.addrToProcessInbox[addr]
	if !ok {
		return fmt.Errorf("not-found: process with id %v not-found", dest)
	}
	log.WithFields(log.Fields{
		"Method":      "exchange.send",
		"MessageType": fmt.Sprintf("%T", m),
		"Dest":        addr,
		"Source":      m.Src()}).Debugf("SendMessage")
	return v.Send(m)
}

func (bme basicMessageExchange) SendAll(pt ProcessType, m Message) error {
	entries, ok := bme.typeToProcessInbox.get(pt)
	if !ok || entries.Len() == 0 {
		return fmt.Errorf("not-found: No process(es) with type:%v found", pt)
	}
	for e := entries.Front(); e != nil; e = e.Next() {
		p := e.Value.(ProcessInbox)
		ctxLog := log.WithFields(log.Fields{
			"MessageType": fmt.Sprintf("%T", m),
			"Dest":        fmt.Sprintf("(%v-%v)", p.ID(), p.Type()),
			"Source":      m.Src()})
		err := p.Send(m)
		ctxLog.Debugf("SendMessage ")
		if err != nil {
			return fmt.Errorf("send failed: to process=%v %v", p, err)
		}
	}
	return nil
}

func (bme basicMessageExchange) Register(p ProcessInbox) error {
	addr := NewAddress(p.ID(), p.Type())
	_, ok := bme.addrToProcessInbox[addr]
	if ok {
		return fmt.Errorf("duplicate: process with id %v", p.ID())
	}
	bme.addrToProcessInbox[addr] = p
	bme.typeToProcessInbox.put(p)
	return nil
}

func (bme basicMessageExchange) UnRegister(p ProcessInbox) error {
	addr := NewAddress(p.ID(), p.Type())
	_, ok := bme.addrToProcessInbox[addr]
	if !ok {
		return fmt.Errorf("not-found: process with id %v", p.ID())
	}
	delete(bme.addrToProcessInbox, addr)
	bme.typeToProcessInbox.remove(p)
	return nil
}
