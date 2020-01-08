package components

import (
	v1 "github.com/1xyz/paxossim/v1"
	"github.com/1xyz/paxossim/v1/messages"
	"github.com/1xyz/paxossim/v1/types"
	"github.com/1xyz/paxossim/v1/v1fakes"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestNewScout(t *testing.T) {
	Convey("when a new scout is initialized", t, func() {
		exchange := &v1fakes.FakeMessageExchange{}
		leader := newFakeAddr(fakeLeaderID, v1.Leader)
		acceptors := newFakeAddrs(3, fakeAcceptorID, v1.Acceptor)
		scout := NewScout(exchange, leader, acceptors, newFakeBallot(0, leader))

		Convey("the resuting ptr should not be nil", func() {
			So(scout, ShouldNotBeNil)
		})
	})
}

func TestScout_PreEmptsOnNewerBallot(t *testing.T) {
	Convey("Given a commander configured for a ballot number", t, func() {
		exchange := &v1fakes.FakeMessageExchange{}
		leader := newFakeAddr(fakeLeaderID, v1.Leader)
		acceptors := newFakeAddrs(3, fakeAcceptorID, v1.Acceptor)
		scout := NewScout(exchange, leader, acceptors, newFakeBallot(0, leader))

		Convey("when it receives a phase1 response with a newer Ballot number", func() {
			newBN := newFakeBallot(3, newFakeAddr(fakeLeaderID+10, v1.Leader))
			responders := makeSet(acceptors)
			bContinue := scout.handleMessage(messages.NewPhase1bMessage(acceptors[0], newBN, nil), &responders)

			Convey("the scout signals an exit", func() {
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

func TestScout_AdoptsOnMatchingBallot(t *testing.T) {
	Convey("Given a scout configured for a ballot number", t, func() {
		exchange := &v1fakes.FakeMessageExchange{}
		leader := newFakeAddr(fakeLeaderID, v1.Leader)
		acceptors := newFakeAddrs(3, fakeAcceptorID, v1.Acceptor)
		bn := newFakeBallot(10, leader)
		scout := NewScout(exchange, leader, acceptors, bn)

		Convey("when it receives a phase1response for the same ballotNumber", func() {
			responders := makeSet(acceptors)
			pValues := make(types.PValues)
			pValues.Set(newFakePValue(8, leader))

			Convey("from one acceptor", func() {
				bContinue := scout.handleMessage(messages.NewPhase1bMessage(acceptors[0], bn, pValues), &responders)

				Convey("it continues to wait for more responses", func() {
					So(bContinue, ShouldBeTrue)
					So(2, ShouldEqual, responders.Len())
				})

				Convey("it adds the response to its pvalues", func() {
					So(len(scout.pvalues), ShouldEqual, 1)
				})

				Convey("and from a majority of acceptors", func() {
					pValues.Set(newFakePValue(7, leader))
					bContinue := scout.handleMessage(messages.NewPhase1bMessage(acceptors[1], bn, pValues), &responders)

					Convey("it signals an exit", func() {
						So(bContinue, ShouldBeFalse)
						So(1, ShouldEqual, responders.Len())
					})

					Convey("it adds the response to its pvalues", func() {
						So(len(scout.pvalues), ShouldEqual, 2)
					})

					Convey("a message is sent", func() {
						So(exchange.SendCallCount(), ShouldEqual, 1)
						addr, msg := exchange.SendArgsForCall(0)

						Convey("to its leader", func() {
							So(addr, ShouldEqual, leader)
						})

						Convey("of type AdoptedMessage", func() {
							adoptedMessage, ok := msg.(messages.AdoptedMessage)
							So(ok, ShouldBeTrue)

							Convey("containing the newer ballot", func() {
								So(adoptedMessage.BallotNumber, ShouldResemble, bn)
							})

							Convey("and the pvalues", func() {
								So(len(adoptedMessage.Accepted), ShouldEqual, 2)
							})
						})
					})
				})
			})
		})
	})
}

//
