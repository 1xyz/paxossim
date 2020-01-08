package components

import (
	"fmt"
	v1 "github.com/1xyz/paxossim/v1"
	"github.com/1xyz/paxossim/v1/messages"
	"github.com/1xyz/paxossim/v1/types"
	"github.com/1xyz/paxossim/v1/v1fakes"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestNewReplica(t *testing.T) {
	Convey("When a new replica is created", t, func() {
		r := NewReplica(&v1fakes.FakeMessageExchange{}, make([]v1.Addr, 0, 3))

		Convey("the resulting ptr should not be nil", func() {
			So(r, ShouldNotBeNil)
		})

		Convey("The Slots should be set to InitialSlotID", func() {
			So(r.slotIn, ShouldEqual, InitialSlotID)
			So(r.slotOut, ShouldEqual, InitialSlotID)
		})

		Convey("The request queue is empty", func() {
			So(len(r.requests), ShouldEqual, 0)
		})

		Convey("The proposals & decisions dictionary are empty", func() {
			So(len(r.proposals), ShouldEqual, 0)
			So(len(r.decisions), ShouldEqual, 0)
		})
	})
}

func TestReplica_Proposal(t *testing.T) {
	Convey("Given a replica", t, func() {

		fakeExchange := v1fakes.FakeMessageExchange{}
		leaders := newLeaders()
		r := NewReplica(&fakeExchange, leaders)

		Convey("When a new request is sent to it", func() {
			r.handleMessage(newTestRequestMessage("1"))

			Convey("the request is queued", func() {
				So(1, ShouldEqual, len(r.requests))
			})

			Convey("When a command is proposed for that request", func() {
				r.propose()

				Convey("The slotIn should increment", func() {
					So(r.slotIn, ShouldEqual, types.Slot(2))
				})

				Convey("The request is de-queued from requests queue", func() {
					So(len(r.requests), ShouldEqual, 0)
				})

				Convey("The request is queued to proposals dictionary", func() {
					So(len(r.proposals), ShouldEqual, 1)
				})

				Convey("Messages are sent", func() {
					So(fakeExchange.SendCallCount(), ShouldEqual, len(leaders))

					addr, message := fakeExchange.SendArgsForCall(0)
					Convey("to leaders", func() {
						So(addr.Type(), ShouldEqual, v1.Leader)
						So(addr.ID(), ShouldEqual, v1.ProcessID(100))
					})

					Convey("from the replicas", func() {
						So(message.Src().Type(), ShouldEqual, v1.Replica)
						So(message.Src().ID(), ShouldEqual, r.ID())
					})

					Convey("of type ProposeMessage", func() {
						_, ok := message.(messages.ProposeMessage)
						So(ok, ShouldBeTrue)
					})
				})
			})
		})
	})
}

func TestReplica_ProposalReachesWindowLimit(t *testing.T) {
	Convey("Given a new replica", t, func() {
		fakeExchange := v1fakes.FakeMessageExchange{}
		leaders := newLeaders()

		r := NewReplica(&fakeExchange, leaders)

		Convey("When the window limit is reached", func() {
			for i := 0; i < int(Window)+1; i++ {
				r.handleMessage(newTestRequestMessage(fmt.Sprintf("%d", i)))
				r.propose()
			}

			Convey("new requests are queued to the request queue", func() {
				So(1, ShouldEqual, len(r.requests))
			})
		})
	})
}

func TestReplica_NewDecisionMatchingProposal(t *testing.T) {
	Convey("Given a new replica", t, func() {
		fakeExchange := v1fakes.FakeMessageExchange{}
		leaders := newLeaders()

		r := NewReplica(&fakeExchange, leaders)
		requestMessage := newTestRequestMessage("1")
		slot := r.slotIn

		Convey("And new proposal is sent to the leades", func() {
			r.handleMessage(requestMessage)
			r.propose()

			Convey("And a decision is received matching the proposal", func() {
				decisionMessage := messages.NewDecisionMessage(leaders[0], slot, requestMessage.Command)
				r.handleMessage(decisionMessage)

				Convey("Decision is added to the decision map", func() {
					So(1, ShouldEqual, len(r.decisions))
				})

				Convey("SlotOut is incremented", func() {
					So(r.slotOut, ShouldEqual, types.Slot(2))
				})

				Convey("Proposal and request queues are empty", func() {
					So(len(r.requests), ShouldEqual, 0)
					So(len(r.proposals), ShouldEqual, 0)
				})
			})
		})
	})
}

func TestReplica_NewDecisionNotMatchingProposal(t *testing.T) {
	Convey("Given a new replica", t, func() {

		fakeExchange := v1fakes.FakeMessageExchange{}
		leaders := newLeaders()

		r := NewReplica(&fakeExchange, leaders)
		requestMessage := newTestRequestMessage("1")
		slot := r.slotIn

		Convey("And new proposal is sent to the leades", func() {
			r.handleMessage(requestMessage)
			r.propose()

			Convey("And a decision is received not-matching the proposal for that slot", func() {
				unmatchedCommand := types.BasicCommand{
					ClientID:  "client:100",
					CommandID: "2",
					Op:        "SUBTRACT",
				}

				decisionMessage := messages.NewDecisionMessage(leaders[0], slot, unmatchedCommand)
				r.handleMessage(decisionMessage)

				Convey("Decision is added to the decision map", func() {
					So(1, ShouldEqual, len(r.decisions))
				})

				Convey("SlotOut is incremented", func() {
					So(r.slotOut, ShouldEqual, types.Slot(2))
				})

				Convey("Proposal map is emptied", func() {
					So(len(r.proposals), ShouldEqual, 0)
				})

				Convey("conflicting request/command is enqueued to request q", func() {
					So(len(r.requests), ShouldEqual, 1)
				})
			})
		})
	})
}

