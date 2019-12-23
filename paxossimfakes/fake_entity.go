// Code generated by counterfeiter. DO NOT EDIT.
package paxossimfakes

import (
	"sync"

	"github.com/1xyz/paxossim"
)

type FakeEntity struct {
	RunStub        func()
	runMutex       sync.RWMutex
	runArgsForCall []struct {
	}
	SendMessageStub        func(paxossim.Message)
	sendMessageMutex       sync.RWMutex
	sendMessageArgsForCall []struct {
		arg1 paxossim.Message
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeEntity) Run() {
	fake.runMutex.Lock()
	fake.runArgsForCall = append(fake.runArgsForCall, struct {
	}{})
	fake.recordInvocation("Run", []interface{}{})
	fake.runMutex.Unlock()
	if fake.RunStub != nil {
		fake.RunStub()
	}
}

func (fake *FakeEntity) RunCallCount() int {
	fake.runMutex.RLock()
	defer fake.runMutex.RUnlock()
	return len(fake.runArgsForCall)
}

func (fake *FakeEntity) RunCalls(stub func()) {
	fake.runMutex.Lock()
	defer fake.runMutex.Unlock()
	fake.RunStub = stub
}

func (fake *FakeEntity) SendMessage(arg1 paxossim.Message) {
	fake.sendMessageMutex.Lock()
	fake.sendMessageArgsForCall = append(fake.sendMessageArgsForCall, struct {
		arg1 paxossim.Message
	}{arg1})
	fake.recordInvocation("SendMessage", []interface{}{arg1})
	fake.sendMessageMutex.Unlock()
	if fake.SendMessageStub != nil {
		fake.SendMessageStub(arg1)
	}
}

func (fake *FakeEntity) SendMessageCallCount() int {
	fake.sendMessageMutex.RLock()
	defer fake.sendMessageMutex.RUnlock()
	return len(fake.sendMessageArgsForCall)
}

func (fake *FakeEntity) SendMessageCalls(stub func(paxossim.Message)) {
	fake.sendMessageMutex.Lock()
	defer fake.sendMessageMutex.Unlock()
	fake.SendMessageStub = stub
}

func (fake *FakeEntity) SendMessageArgsForCall(i int) paxossim.Message {
	fake.sendMessageMutex.RLock()
	defer fake.sendMessageMutex.RUnlock()
	argsForCall := fake.sendMessageArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeEntity) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.runMutex.RLock()
	defer fake.runMutex.RUnlock()
	fake.sendMessageMutex.RLock()
	defer fake.sendMessageMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeEntity) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ paxossim.Entity = new(FakeEntity)