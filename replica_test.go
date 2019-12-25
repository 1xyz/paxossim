package paxossim_test

import (
	"fmt"
	"github.com/1xyz/paxossim"
	"github.com/1xyz/paxossim/paxossimfakes"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	nLeaders = 2
)

func newFakeConfigWithLeaders(leaders []*paxossimfakes.FakeEntity) *paxossim.Configuration {
	entities := toEntities(leaders)
	return &paxossim.Configuration{Leaders: entities}
}

func newFakeConfig() *paxossim.Configuration {
	leaders := newLeaders()
	return newFakeConfigWithLeaders(leaders)
}

func toEntities(leaders []*paxossimfakes.FakeEntity) []paxossim.Entity {
	nLen := len(leaders)
	entities := make([]paxossim.Entity, nLen, nLen)
	for i := range leaders {
		entities[i] = leaders[i]
	}
	return entities
}

func newLeaders() []*paxossimfakes.FakeEntity {
	leaders := make([]*paxossimfakes.FakeEntity, nLeaders, nLeaders)
	for i := 0; i < nLeaders; i++ {
		fe := paxossimfakes.FakeEntity{}
		leaders[i] = &fe
	}
	return leaders
}

func nextCommand(commandId int) *paxossim.BasicCommand {
	cmd := &paxossim.BasicCommand{
		ClientID:  "client:1",
		CommandID: fmt.Sprintf("%d", commandId),
		Op:        "ADD",
	}
	return cmd
}

func TestNewReplica(t *testing.T) {
	r := paxossim.NewReplica("r:1", newFakeConfig())
	assert.Equal(t, 0, len(r.Requests))
	assert.Equal(t, 0, len(r.Decisions))
	assert.Equal(t, 0, len(r.Proposals))
	assert.Equal(t, paxossim.InitialSlotID, r.SlotIn)
	assert.Equal(t, paxossim.InitialSlotID, r.SlotOut)
}

func TestReplica_HandleMsg_NewRequest(t *testing.T) {
	r := paxossim.NewReplica("r:1", newFakeConfig())
	req := paxossim.NewRequestMessage("client:1", &paxossim.BasicCommand{
		ClientID:  "client:1",
		CommandID: "1",
		Op:        "ADD",
	})
	r.HandleMsg(req)
	// Assert this request has been queued
	assert.Equal(t, 1, len(r.Requests))
}

func TestReplica_Propose_NewRequest(t *testing.T) {
	leaders := newLeaders()
	cmd := nextCommand(1)
	r := paxossim.NewReplica("r:1", newFakeConfigWithLeaders(leaders))
	req := paxossim.NewRequestMessage("client:1", cmd)
	r.HandleMsg(req)
	r.Propose()

	assert.Equal(t, paxossim.SlotID(2), r.SlotIn)
	for _, leader := range leaders {
		// All the leaders have been sent messages
		assert.Equal(t, 1, leader.SendMessageCallCount())
		pm, ok := leader.SendMessageArgsForCall(0).(*paxossim.ProposeMessage)
		assert.True(t, ok)
		assert.Equal(t, paxossim.SlotID(1), pm.SlotID)
		assert.Equal(t, cmd, pm.C)
	}

	assert.Equal(t, 1, len(r.Proposals))
	assert.Equal(t, 0, len(r.Requests))
	assert.Equal(t, 0, len(r.Decisions))
	assert.Equal(t, paxossim.InitialSlotID, r.SlotOut)
}

func TestReplica_Propose_NewRequestWithWindowLimit(t *testing.T) {
	leaders := newLeaders()
	r := paxossim.NewReplica("r:1", newFakeConfigWithLeaders(leaders))

	for i := 0; i < int(paxossim.Window)+1; i++ {
		req := paxossim.NewRequestMessage("client:1", nextCommand(i))
		r.HandleMsg(req)
		r.Propose()
	}

	assert.Equal(t, 1, len(r.Requests))
}

func TestReplica_HandleMsg_NewDecisionMatchingProposal(t *testing.T) {
	cmd := nextCommand(1)
	req := paxossim.NewRequestMessage("client:1", cmd)
	r := paxossim.NewReplica("r:1", newFakeConfigWithLeaders(newLeaders()))
	s := r.SlotIn

	r.HandleMsg(req)
	r.Propose()

	decision := paxossim.NewDecisionMessage("leader:1", s, cmd)
	r.HandleMsg(decision)

	assert.Equal(t, 0, len(r.Requests))
	assert.Equal(t, 0, len(r.Proposals))
	assert.Equal(t, 1, len(r.Decisions))
	assert.Equal(t, paxossim.SlotID(2), r.SlotOut)
}

func TestReplica_HandleMsg_NewDecisionNotMatchingProposal(t *testing.T) {
	cmd := nextCommand(1)
	req := paxossim.NewRequestMessage("client:1", cmd)
	r := paxossim.NewReplica("r:1", newFakeConfigWithLeaders(newLeaders()))
	s := r.SlotIn

	r.HandleMsg(req)
	r.Propose()

	cmd2 := &paxossim.BasicCommand{
		ClientID:  "client:2",
		CommandID: "1",
		Op:        "SUBTRACT",
	}
	decision := paxossim.NewDecisionMessage("leader:1", s, cmd2)
	r.HandleMsg(decision)

	assert.Equal(t, 1, len(r.Requests))
	assert.Equal(t, 0, len(r.Proposals))
	assert.Equal(t, 1, len(r.Decisions))
	assert.Equal(t, paxossim.SlotID(2), r.SlotOut)
}
