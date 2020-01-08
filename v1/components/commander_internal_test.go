package components

import (
	v1 "github.com/1xyz/paxossim/v1"
	"github.com/1xyz/paxossim/v1/messages"
	"github.com/1xyz/paxossim/v1/types"
	"github.com/1xyz/paxossim/v1/v1fakes"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestNewCommander(t *testing.T) {
	Convey("when a new commander is initialized", t, func() {
		exchange := &v1fakes.FakeMessageExchange{}
		leader := newFakeAddr(fakeLeaderID, v1.Leader)
		acceptors := newFakeAddrs(3, fakeAcceptorID, v1.Acceptor)
		cmdr := NewCommander(exchange, leader, acceptors, newFakePValue(0, leader))

		Convey("the resuting ptr should not be nil", func() {
			So(cmdr, ShouldNotBeNil)
		})
	})
}

func TestCommander_BroadcastToAcceptors(t *testing.T) {
	Convey("when a new commander is initialized", t, func() {
		exchange := &v1fakes.FakeMessageExchange{}
		leader := newFakeAddr(fakeLeaderID, v1.Leader)
		nAcceptors := 3
		acceptors := newFakeAddrs(nAcceptors, fakeAcceptorID, v1.Acceptor)
		pValue := newFakePValue(0, leader)
		cmdr := NewCommander(exchange, leader, acceptors, pValue)

		Convey("ensure a message is sent to all acceptors", func() {
			result := cmdr.broadcastToAcceptors()

			So(result.Contains(acceptors[0]), ShouldBeTrue)
			So(result.Contains(acceptors[1]), ShouldBeTrue)
			So(result.Contains(acceptors[2]), ShouldBeTrue)

			Convey("and is received by the acceptors", func() {
				So(exchange.SendCallCount(), ShouldEqual, nAcceptors)
				addr, msg := exchange.SendArgsForCall(0)
				So(acceptors[0], ShouldEqual, addr)

				Convey("as a Phase2aMessage with the correct PVale", func() {
					phase2aMessage, ok := msg.(messages.Phase2aMessage)
					So(ok, ShouldBeTrue)
					So(phase2aMessage.PValue, ShouldResemble, pValue)
				})
			})
		})
	})
}

func TestCommander_PreEmptsOnNewerBallot(t *testing.T) {
	Convey("Given a commander configured for a ballot number", t, func() {
		exchange := &v1fakes.FakeMessageExchange{}
		leader := newFakeAddr(fakeLeaderID, v1.Leader)
		acceptors := newFakeAddrs(3, fakeAcceptorID, v1.Acceptor)
		cmdr := NewCommander(exchange, leader, acceptors, newFakePValue(0, leader))

		Convey("when it receives a newer Ballot number", func() {
			newBN := newFakeBallot(3, newFakeAddr(fakeLeaderID+10, v1.Leader))
			responders := makeSet(acceptors)
			bContinue := cmdr.handleMessage(messages.NewPhase2bMessage(acceptors[0], newBN), &responders)

			Convey("the commander signals an exit", func() {
				So(bContinue, ShouldBeFalse)
			})

			Convey("a message is sent", func() {
				So(exchange.SendCallCount(), ShouldEqual, 1)
				addr, msg := exchange.SendArgsForCall(0)

				Convey("to its leader", func() {
					So(addr, ShouldEqual, leader)
				})

				Convey("of type PreEmptMessage", func() {
					preEmptMessage, ok := msg.(messages.PreemptMessage)
					So(ok, ShouldBeTrue)

					Convey("containing the newer ballot", func() {
						So(preEmptMessage.BallotNumber, ShouldResemble, newBN)
					})
				})
			})
		})
	})
}

func TestCommander_SendsDecisionToReplicas(t *testing.T) {
	Convey("Given a commander configured for a ballot number", t, func() {
		exchange := &v1fakes.FakeMessageExchange{}
		leader := newFakeAddr(fakeLeaderID, v1.Leader)
		acceptors := newFakeAddrs(3, fakeAcceptorID, v1.Acceptor)
		pValue := newFakePValue(0, leader)
		cmdr := NewCommander(exchange, leader, acceptors, pValue)

		Convey("when it receives the same Ballot ", func() {
			responders := makeSet(acceptors)

			Convey("from one acceptor", func() {
				bContinue := cmdr.handleMessage(messages.NewPhase2bMessage(acceptors[0], pValue.BN), &responders)

				Convey("it continues to wait for more responses", func() {
					So(bContinue, ShouldBeTrue)
					So(2, ShouldEqual, responders.Len())
				})

				Convey("from a majority of acceptors", func() {
					bContinue := cmdr.handleMessage(messages.NewPhase2bMessage(acceptors[2], pValue.BN), &responders)

					Convey("it signals an exit", func() {
						So(bContinue, ShouldBeFalse)
						So(1, ShouldEqual, responders.Len())
					})

					Convey("and notification is sent", func() {
						So(exchange.SendAllCallCount(), ShouldEqual, 1)
						pt, msg := exchange.SendAllArgsForCall(0)

						Convey("all replicas", func() {
							So(pt, ShouldEqual, v1.Replica)
						})

						Convey("with a DecisionMessage", func() {
							decisionMessage, ok := msg.(messages.DecisionMessage)
							So(ok, ShouldBeTrue)

							Convey("and the assigned command and slot", func() {
								So(pValue.Command, ShouldResemble, decisionMessage.Command)
								So(pValue.Slot, ShouldResemble, decisionMessage.Slot)
							})
						})
					})
				})
			})
		})
	})
}

func makeSet(acceptors []v1.Addr) v1.AddrSet {
	result := make(v1.AddrSet)
	for _, e := range acceptors {
		result.Add(e)
	}
	return result
}

const fakeLeaderID = v1.ProcessID(100)
const fakeAcceptorID = v1.ProcessID(200)

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
