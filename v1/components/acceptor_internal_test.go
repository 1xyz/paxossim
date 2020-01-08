package components

import (
	v1 "github.com/1xyz/paxossim/v1"
	"github.com/1xyz/paxossim/v1/messages"
	"github.com/1xyz/paxossim/v1/v1fakes"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestNewAcceptor(t *testing.T) {
	Convey("Initializing a new Acceptor", t, func() {
		acceptor := NewAcceptor(&v1fakes.FakeMessageExchange{})

		Convey("returns a valid ptr", func() {
			So(acceptor, ShouldNotBeNil)
		})
	})
}

func TestAcceptor_Run_Phase1(t *testing.T) {
	Convey("Given an acceptor", t, func() {
		exchange := &v1fakes.FakeMessageExchange{}
		acceptor := NewAcceptor(exchange)
		leader := newFakeAddr(fakeLeaderID, v1.Leader)
		scout := newFakeAddr(fakeScoutID, v1.Scout)

		Convey("Encounters a ballot for the first time from a scout", func() {
			bn := newFakeBallot(10, leader)
			acceptor.handleMessage(messages.NewPhase1aMessage(scout, bn))

			Convey("the ballot is adopted", func() {
				So(*acceptor.BN, ShouldResemble, bn)

				Convey("And a Phase1bMessage is responded to the scout", func() {
					So(exchange.SendCallCount(), ShouldEqual, 1)

					addr, msg := exchange.SendArgsForCall(0)
					So(addr, ShouldEqual, scout)

					phase1bMessage, ok := msg.(messages.Phase1bMessage)
					So(ok, ShouldBeTrue)
					So(phase1bMessage.BallotNumber, ShouldResemble, bn)
				})
			})

			Convey("Encounters an older ballot from a scout", func() {
				olderBN := newFakeBallot(2, leader)
				acceptor.handleMessage(messages.NewPhase1aMessage(scout, olderBN))

				Convey("the ballot is not adopted", func() {
					So(*acceptor.BN, ShouldNotResemble, olderBN)
					So(*acceptor.BN, ShouldResemble, bn)

					Convey("And a Phase1bMessage is responded to the scout", func() {
						So(exchange.SendCallCount(), ShouldEqual, 2)

						addr, msg := exchange.SendArgsForCall(1)
						So(addr, ShouldEqual, scout)

						phase1bMessage, ok := msg.(messages.Phase1bMessage)
						So(ok, ShouldBeTrue)
						So(phase1bMessage.BallotNumber, ShouldResemble, bn)
					})
				})
			})

			Convey("Encounters a newer ballot from a scout", func() {
				newerBN := newFakeBallot(20, leader)
				acceptor.handleMessage(messages.NewPhase1aMessage(scout, newerBN))

				Convey("the ballot is adopted", func() {
					So(*acceptor.BN, ShouldResemble, newerBN)
					So(*acceptor.BN, ShouldNotResemble, bn)

					Convey("And a Phase1bMessage is responded to the scout", func() {
						So(exchange.SendCallCount(), ShouldEqual, 2)

						addr, msg := exchange.SendArgsForCall(1)
						So(addr, ShouldEqual, scout)

						phase1bMessage, ok := msg.(messages.Phase1bMessage)
						So(ok, ShouldBeTrue)
						So(phase1bMessage.BallotNumber, ShouldResemble, newerBN)
					})
				})
			})
		})
	})
}

func TestAcceptor_Run_Phase2(t *testing.T) {
	Convey("Given an acceptor", t, func() {
		exchange := &v1fakes.FakeMessageExchange{}
		acceptor := NewAcceptor(exchange)
		leader := newFakeAddr(fakeLeaderID, v1.Leader)
		scout := newFakeAddr(fakeScoutID, v1.Scout)
		cmdr := newFakeAddr(fakeCommanderID, v1.Commander)

		bn := newFakeBallot(10, leader)
		acceptor.handleMessage(messages.NewPhase1aMessage(scout, bn))

		Convey("Encounters a phase2 request matching its adopted ballot", func() {
			pValue := newFakePValue(10, leader)
			phase2aMessage := messages.NewPhase2aMessage(cmdr, pValue)
			acceptor.handleMessage(phase2aMessage)

			Convey("the PValue is accepted", func() {
				So(len(acceptor.Accepted), ShouldEqual, 1)

				Convey("And a Phase2bMessage is responded to the commander", func() {
					So(exchange.SendCallCount(), ShouldEqual, 2)

					addr, msg := exchange.SendArgsForCall(1)
					So(addr, ShouldEqual, cmdr)

					phase2bMessage, ok := msg.(messages.Phase2bMessage)
					So(ok, ShouldBeTrue)
					So(phase2bMessage.BallotNumber, ShouldResemble, pValue.BN)
				})
			})
		})

		Convey("Encounters a phase2 request with not matching its adopted ballot", func() {
			pValue := newFakePValue(20, leader)
			phase2aMessage := messages.NewPhase2aMessage(cmdr, pValue)
			acceptor.handleMessage(phase2aMessage)

			Convey("the PValue is rejected", func() {
				So(len(acceptor.Accepted), ShouldEqual, 0)

				Convey("And a Phase2bMessage is responded to the commander", func() {
					So(exchange.SendCallCount(), ShouldEqual, 2)

					addr, msg := exchange.SendArgsForCall(1)
					So(addr, ShouldEqual, cmdr)

					Convey("with the acceptors adopted ballot number", func() {
						phase2bMessage, ok := msg.(messages.Phase2bMessage)
						So(ok, ShouldBeTrue)
						So(phase2bMessage.BallotNumber, ShouldResemble, bn)
					})
				})
			})
		})
	})
}