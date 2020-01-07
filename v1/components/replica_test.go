package components

import (
	"fmt"
	"github.com/1xyz/paxossim"
	v1 "github.com/1xyz/paxossim/v1"
	"github.com/1xyz/paxossim/v1/messages"
	"github.com/1xyz/paxossim/v1/types"
	"github.com/1xyz/paxossim/v1/v1fakes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewReplica(t *testing.T) {
	r := NewReplica(&v1fakes.FakeMessageExchange{}, make([]v1.Addr, 0, 3))
	assert.NotNil(t, r)
}

func TestReplica_HandlesRequestMessage(t *testing.T) {
	fakeExchange := v1fakes.FakeMessageExchange{}
	leaders := make([]v1.Addr, 0, 3)

	r := NewReplica(&fakeExchange, leaders)
	r.handleMessage(newTestRequestMessage("1"))

	// Assert this request has been queued
	assert.Equal(t, 1, len(r.requests))
}

func TestReplica_Propose_NewRequest(t *testing.T) {
	fakeExchange := v1fakes.FakeMessageExchange{}
	leaders := newLeaders()

	r := NewReplica(&fakeExchange, leaders)

	r.handleMessage(newTestRequestMessage("1"))
	r.propose()

	assert.Equal(t, types.Slot(2), r.slotIn)
	assert.Equal(t, 1, len(r.proposals))
	assert.Equal(t, 0, len(r.requests))
	assert.Equal(t, 0, len(r.decisions))

	assert.Equal(t, 2, fakeExchange.SendCallCount())
	addr, message := fakeExchange.SendArgsForCall(0)
	assert.Equal(t, v1.Leader, addr.Type())
	assert.Equal(t, v1.ProcessID(100), addr.ID())
	assert.Equal(t, v1.Replica, message.Src().Type())
	assert.Equal(t, r.ID(), message.Src().ID())

	_, ok := message.(messages.ProposeMessage)
	assert.True(t, ok)
}

func TestReplica_Propose_NewRequestWithWindowLimit(t *testing.T) {
	fakeExchange := v1fakes.FakeMessageExchange{}
	leaders := newLeaders()

	r := NewReplica(&fakeExchange, leaders)

	for i := 0; i < int(paxossim.Window)+1; i++ {
		r.handleMessage(newTestRequestMessage(fmt.Sprintf("%d", i)))
		r.propose()
	}

	assert.Equal(t, 1, len(r.requests))
}

func newTestRequestMessage(commandID string) messages.RequestMessage {
	return messages.RequestMessage{
		Command: types.BasicCommand{
			ClientID:  "client:1",
			CommandID: commandID,
			Op:        "ADD",
		},
	}
}

func newLeaders() []v1.Addr {
	nLeaders := 2
	leaders := make([]v1.Addr, nLeaders, nLeaders)
	for i := 0; i < nLeaders; i++ {
		leaders[i] = &v1fakes.FakeAddr{
			IDStub: func() v1.ProcessID {
				return v1.ProcessID(100)
			},
			TypeStub: func() v1.ProcessType {
				return v1.Leader
			},
		}
	}
	return leaders
}
