package paxossim_test

import (
	"github.com/1xyz/paxossim"
	"github.com/1xyz/paxossim/paxossimfakes"
	"github.com/stretchr/testify/assert"
	"testing"
)

// Assert the state of a newly created acceptor
func TestNewAcceptor(t *testing.T) {
	acceptor := paxossim.NewAcceptor("acceptor:1")
	assert.Equal(t,
		paxossim.BallotNumber{
			Round:    0,
			LeaderID: "",
		}, *acceptor.BN)
	assert.NotNil(t, acceptor.Accepted)
	assert.Equal(t, 0, len(acceptor.Accepted))
}

// Assert a new ballot number is adopted
func TestAcceptor_HandleMsg_NewBN(t *testing.T) {
	acceptor := paxossim.NewAcceptor("acceptor:1")
	bn := &paxossim.BallotNumber{
		Round:    1,
		LeaderID: "Leader:0",
	}
	fakeLeader := &paxossimfakes.FakeEntity{}
	acceptor.HandleMsg(paxossim.NewP1aMessage("Scout:0", fakeLeader, bn))

	// assert the newly adopted ballot number
	assert.Equal(t, *bn, *acceptor.BN)

	// assert that the fake leader was called with a phase1 response
	assert.Equal(t, 1, fakeLeader.SendMessageCallCount())
	p1resp, ok := fakeLeader.SendMessageArgsForCall(0).(*paxossim.P1bMessage)
	assert.True(t, ok)
	assert.Equal(t, *bn, *(p1resp.BN))
}

// Assert a Ballot number that is stale (aka. older) is un-adopted
func TestAcceptor_HandleMsg_OldBN(t *testing.T) {
	acceptor := paxossim.NewAcceptor("acceptor:1")
	fakeLeader := &paxossimfakes.FakeEntity{}
	acceptor.HandleMsg(paxossim.NewP1aMessage("Scout:0", fakeLeader, &paxossim.BallotNumber{
		Round:    100,
		LeaderID: "Leader:0",
	}))

	oldBn := &paxossim.BallotNumber{
		Round:    2,
		LeaderID: "Leader:0",
	}

	// assert no new adoption tool place
	assert.NotEqual(t, *oldBn, *acceptor.BN)
}

// Assert a PValue is accepted if BN matches
func TestAcceptor_HandleMsg_PVAccepted(t *testing.T) {
	acceptor := paxossim.NewAcceptor("acceptor:1")
	fakeLeader := &paxossimfakes.FakeEntity{}
	bn := &paxossim.BallotNumber{
		Round:    100,
		LeaderID: "Leader:0",
	}
	acceptor.HandleMsg(paxossim.NewP1aMessage("Scout:0", fakeLeader, bn))

	// create a new p2 request with a match ballot, leader
	pv := &paxossim.PValue{
		BN:   bn,
		Slot: paxossim.SlotID(10234),
		C: &paxossim.BasicCommand{
			ClientID:  "client:0",
			CommandID: "12",
			Op:        "ADD",
		},
	}
	acceptor.HandleMsg(paxossim.NewP2aMessage("Commander:0", fakeLeader, pv))

	// assert that pvalue has been accepted
	assert.True(t, acceptor.Accepted.Contains(pv))

	// Assert that a Phase2 response was delivered to the Commander
	// Note there must be two calls recorded call 0 is phase1 response
	assert.Equal(t, 2, fakeLeader.SendMessageCallCount())
	p2resp, ok := fakeLeader.SendMessageArgsForCall(1).(*paxossim.P2bMessage)
	assert.True(t, ok)
	assert.Equal(t, *bn, *(p2resp.BN))
}

// Assert a PValue is rejected if BN does not match the adopted BN
func TestAcceptor_HandleMsg_PVRejected(t *testing.T) {
	fakeLeader := &paxossimfakes.FakeEntity{}
	acceptor := paxossim.NewAcceptor("acceptor:1")
	bn := &paxossim.BallotNumber{
		Round:    100,
		LeaderID: "Leader:0",
	}

	// create a new p2 request with a match ballot, leader
	pv := &paxossim.PValue{
		BN:   bn,
		Slot: paxossim.SlotID(10234),
		C: &paxossim.BasicCommand{
			ClientID:  "client:0",
			CommandID: "12",
			Op:        "ADD",
		},
	}
	acceptor.HandleMsg(paxossim.NewP2aMessage("Commander:0", fakeLeader, pv))

	// assert that pvalue has not been accepted
	assert.False(t, acceptor.Accepted.Contains(pv))

	// assert that the Phase2 response contains the older BN
	assert.Equal(t, 1, fakeLeader.SendMessageCallCount())
	p2resp, ok := fakeLeader.SendMessageArgsForCall(0).(*paxossim.P2bMessage)
	assert.True(t, ok)
	assert.Equal(t, paxossim.BallotNumber{
		Round:    0,
		LeaderID: "",
	}, *(p2resp.BN))
}
