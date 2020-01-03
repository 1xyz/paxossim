package paxossim_test

import (
	"github.com/1xyz/paxossim"
	"github.com/1xyz/paxossim/paxossimfakes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func newFakeBN(round int) *paxossim.BallotNumber {
	return &paxossim.BallotNumber{
		Round:    round,
		LeaderID: "leader:1",
	}
}

func TestNewScout(t *testing.T) {
	leader, _, acceptors := createFakeEntities(1)
	bn := newFakeBN(0)

	scout := paxossim.NewScout("scout:1", leader, acceptors, bn)
	assert.NotNil(t, scout)
}

func TestScout_HandleMessage_SendsPremptedMessage(t *testing.T) {
	l, _, a := createFakeEntities(1)
	bn1 := newFakeBN(0)
	scout := paxossim.NewScout("scout:1", l, a, bn1)

	ss := make(paxossim.StringSet)
	ss.Add("acceptor:0")

	bn2 := newFakeBN(1)
	p1bMessage := paxossim.NewP1bMessage("acceptor:1", bn2, nil)
	bContinue := scout.HandleMessage(p1bMessage, &ss)

	assert.False(t, bContinue)
	assert.Equal(t, 1, l.(*paxossimfakes.FakeEntity).SendMessageCallCount())

	msg := l.(*paxossimfakes.FakeEntity).SendMessageArgsForCall(0).(*paxossim.PreemptMessage)
	assert.Equal(t, 0, msg.BN.CompareTo(bn2))
}

func TestScout_HandleMessage_HandlesAdoptedMessage(t *testing.T) {
	leader, acceptors, _ := createFakeEntities(1)
	pv := createDefaultPV()
	scout := paxossim.NewScout("scout:1", leader, acceptors, pv.BN)

	addrSet := scout.BroadcastToAcceptors()
	pValues := make(paxossim.PValues)
	pValues.Set(pv)

	bContinue := scout.HandleMessage(paxossim.NewP1bMessage("acceptor:1",
		newFakeBN(0), &pValues), &addrSet)
	assert.True(t, bContinue)

	bContinue = scout.HandleMessage(paxossim.NewP1bMessage("acceptor:0",
		newFakeBN(0), &pValues), &addrSet)
	assert.False(t, bContinue)

	assert.Equal(t, 1, leader.(*paxossimfakes.FakeEntity).SendMessageCallCount())
	msg := leader.(*paxossimfakes.FakeEntity).SendMessageArgsForCall(0).(*paxossim.AdoptedMessage)
	assert.Equal(t, 0, msg.BN.CompareTo(newFakeBN(0)))
}
