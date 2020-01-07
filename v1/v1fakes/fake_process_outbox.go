// Code generated by counterfeiter. DO NOT EDIT.
package v1fakes

import (
	"sync"

	v1 "github.com/1xyz/paxossim/v1"
)

type FakeProcessOutbox struct {
	RecvStub        func() (v1.Message, error)
	recvMutex       sync.RWMutex
	recvArgsForCall []struct {
	}
	recvReturns struct {
		result1 v1.Message
		result2 error
	}
	recvReturnsOnCall map[int]struct {
		result1 v1.Message
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeProcessOutbox) Recv() (v1.Message, error) {
	fake.recvMutex.Lock()
	ret, specificReturn := fake.recvReturnsOnCall[len(fake.recvArgsForCall)]
	fake.recvArgsForCall = append(fake.recvArgsForCall, struct {
	}{})
	fake.recordInvocation("Recv", []interface{}{})
	fake.recvMutex.Unlock()
	if fake.RecvStub != nil {
		return fake.RecvStub()
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.recvReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeProcessOutbox) RecvCallCount() int {
	fake.recvMutex.RLock()
	defer fake.recvMutex.RUnlock()
	return len(fake.recvArgsForCall)
}

func (fake *FakeProcessOutbox) RecvCalls(stub func() (v1.Message, error)) {
	fake.recvMutex.Lock()
	defer fake.recvMutex.Unlock()
	fake.RecvStub = stub
}

func (fake *FakeProcessOutbox) RecvReturns(result1 v1.Message, result2 error) {
	fake.recvMutex.Lock()
	defer fake.recvMutex.Unlock()
	fake.RecvStub = nil
	fake.recvReturns = struct {
		result1 v1.Message
		result2 error
	}{result1, result2}
}

func (fake *FakeProcessOutbox) RecvReturnsOnCall(i int, result1 v1.Message, result2 error) {
	fake.recvMutex.Lock()
	defer fake.recvMutex.Unlock()
	fake.RecvStub = nil
	if fake.recvReturnsOnCall == nil {
		fake.recvReturnsOnCall = make(map[int]struct {
			result1 v1.Message
			result2 error
		})
	}
	fake.recvReturnsOnCall[i] = struct {
		result1 v1.Message
		result2 error
	}{result1, result2}
}

func (fake *FakeProcessOutbox) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.recvMutex.RLock()
	defer fake.recvMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeProcessOutbox) recordInvocation(key string, args []interface{}) {
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

var _ v1.ProcessOutbox = new(FakeProcessOutbox)
