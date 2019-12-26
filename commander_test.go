package paxossim_test

import (
	"fmt"
	"github.com/1xyz/paxossim"
	"github.com/1xyz/paxossim/paxossimfakes"
	"github.com/stretchr/testify/assert"
	"testing"
)

// Construct fake entities
func createFakeEntities(nFailures int) (paxossim.Entity, []paxossim.Entity, []paxossim.Entity) {
	leader := &paxossimfakes.FakeEntity{}

	nReplicas := nFailures + 1
	replicas := make([]paxossim.Entity, nReplicas, nReplicas)
	for i := range replicas {
		replicas[i] = &paxossimfakes.FakeEntity{}
	}

	nAcceptors := (2 * nFailures) + 1
	acceptors := make([]paxossim.Entity, nAcceptors, nAcceptors)
	for i := range acceptors {
		id := fmt.Sprintf("acceptor:%d", i)
		acceptors[i] = &paxossimfakes.FakeEntity{
			IDStub: func() string {
				return id
			},
		}
	}

	return leader, acceptors, replicas
}

// create a default ballot
func createDefaultPV() *paxossim.PValue {
	return &paxossim.PValue{
		BN: &paxossim.BallotNumber{
			Round:    0,
			LeaderID: "leader:1",
		},
		Slot: 0,
		C: &paxossim.BasicCommand{
			ClientID:  "client:0",
			CommandID: "0",
			Op:        "POP",
		},
	}
}

// Assert the creation a new Commander
func TestNewCommander(t *testing.T) {
	leader, acceptors, replicas := createFakeEntities(1)
	pv := createDefaultPV()

	cmd := paxossim.NewCommander("commander:1", leader, acceptors, replicas, pv)
	assert.NotNil(t, cmd.Leader)
	assert.Equal(t, len(acceptors), len(cmd.Acceptors))
	assert.Equal(t, len(replicas), len(cmd.Replicas))
	assert.Equal(t, *pv, *cmd.PV)
}

func TestCommander_HandleMessage_PreEmptsOnNewerBallot(t *testing.T) {
	newBN := &paxossim.BallotNumber{
		Round:    3,
		LeaderID: "leader:1",
	}
	p2Response := paxossim.NewP2bMessage("acceptor:0", newBN)
	leader, acceptors, replicas := createFakeEntities(1)
	pv := createDefaultPV()
	cmd := paxossim.NewCommander("commander:1", leader, acceptors, replicas, pv)

	ss := make(paxossim.StringSet)
	ss.Add("acceptor:0")
	bContinue := cmd.HandleMessage(p2Response, &ss)

	// Assert HandleMessage signals an exit
	assert.False(t, bContinue)

	// Assert that pre-empted message was sent to leadder
	fakeLeader := leader.(*paxossimfakes.FakeEntity)
	assert.Equal(t, 1, fakeLeader.SendMessageCallCount())

	// Assert the message & ballot number returned
	premptMessage := fakeLeader.SendMessageArgsForCall(0).(*paxossim.PreemptMessage)
	assert.Equal(t, *newBN, *premptMessage.BN)
}

func TestCommander_HandleMessage_DecisionResponse(t *testing.T) {
	leader, acceptors, replicas := createFakeEntities(1)
	pv := createDefaultPV()
	cmd := paxossim.NewCommander("commander:1", leader, acceptors, replicas, pv)

	addrSet := cmd.BroadcastToAcceptors()

	bn0 := *pv.BN
	// assert we continue for messages
	assert.True(t, cmd.HandleMessage(paxossim.NewP2bMessage("acceptor:0", &bn0), &addrSet))
	// asser the address list has shrunk
	assert.Equal(t, 2, addrSet.Len())

	bn1 := *pv.BN
	// assert we don't continue for messages, we have reached a majority
	assert.False(t, cmd.HandleMessage(paxossim.NewP2bMessage("acceptor:2", &bn1), &addrSet))
	// asser the address list has shrunk
	assert.Equal(t, 1, addrSet.Len())

	// assert all replicas have been notified
	for _, replica := range replicas {
		// every replica gets a message
		fr := replica.(*paxossimfakes.FakeEntity)
		assert.Equal(t, 1, fr.SendMessageCallCount())

		// validate this message
		decisionMessage := fr.SendMessageArgsForCall(0).(*paxossim.DecisionMessage)
		assert.Equal(t, pv.Slot, decisionMessage.SlotID)
		assert.Equal(t, pv.C.GetClientID(), decisionMessage.C.GetClientID())
		assert.Equal(t, pv.C.GetCommandID(), decisionMessage.C.GetCommandID())
		assert.Equal(t, pv.C.GetOp(), decisionMessage.C.GetOp())
	}
}
