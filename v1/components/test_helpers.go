package components

import (
	v1 "github.com/1xyz/paxossim/v1"
	"github.com/1xyz/paxossim/v1/messages"
	"github.com/1xyz/paxossim/v1/types"
	"github.com/1xyz/paxossim/v1/v1fakes"
)

func makeSet(acceptors []v1.Addr) v1.AddrSet {
	result := make(v1.AddrSet)
	for _, e := range acceptors {
		result.Add(e)
	}
	return result
}

const (
	fakeLeaderID    = v1.ProcessID(100)
	fakeAcceptorID  = v1.ProcessID(200)
	fakeCommanderID = v1.ProcessID(300)
	fakeScoutID     = v1.ProcessID(400)
)

func newFakePValue(round int, leaderID v1.Addr) types.PValue {
	return types.PValue{
		BN:   newFakeBallot(round, leaderID),
		Slot: 0,
		Command: types.BasicCommand{
			ClientID:  "client:0",
			CommandID: "0",
			Op:        "POP",
		},
	}
}

func newFakeBallot(round int, leaderID v1.Addr) types.BallotNumber {
	return types.BallotNumber{
		Round:    round,
		LeaderID: leaderID,
	}
}

func newFakeAddr(id v1.ProcessID, pt v1.ProcessType) v1.Addr {
	return &v1fakes.FakeAddr{
		IDStub: func() v1.ProcessID {
			return id
		},
		TypeStub: func() v1.ProcessType {
			return pt
		},
	}
}

func newFakeAddrs(count int, idStart v1.ProcessID, pt v1.ProcessType) []v1.Addr {
	entries := make([]v1.Addr, count, count)
	for i := 0; i < count; i++ {
		entries[i] = newFakeAddr(idStart, pt)
		idStart++
	}

	return entries
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
	return newFakeAddrs(2, fakeLeaderID, v1.Leader)
}
