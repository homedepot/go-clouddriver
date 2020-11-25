// Code generated by counterfeiter. DO NOT EDIT.
package arcadefakes

import (
	"sync"

	"github.com/homedepot/go-clouddriver/pkg/arcade"
)

type FakeClient struct {
	TokenStub        func() (string, error)
	tokenMutex       sync.RWMutex
	tokenArgsForCall []struct {
	}
	tokenReturns struct {
		result1 string
		result2 error
	}
	tokenReturnsOnCall map[int]struct {
		result1 string
		result2 error
	}
	WithAPIKeyStub        func(string)
	withAPIKeyMutex       sync.RWMutex
	withAPIKeyArgsForCall []struct {
		arg1 string
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeClient) Token() (string, error) {
	fake.tokenMutex.Lock()
	ret, specificReturn := fake.tokenReturnsOnCall[len(fake.tokenArgsForCall)]
	fake.tokenArgsForCall = append(fake.tokenArgsForCall, struct {
	}{})
	stub := fake.TokenStub
	fakeReturns := fake.tokenReturns
	fake.recordInvocation("Token", []interface{}{})
	fake.tokenMutex.Unlock()
	if stub != nil {
		return stub()
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeClient) TokenCallCount() int {
	fake.tokenMutex.RLock()
	defer fake.tokenMutex.RUnlock()
	return len(fake.tokenArgsForCall)
}

func (fake *FakeClient) TokenCalls(stub func() (string, error)) {
	fake.tokenMutex.Lock()
	defer fake.tokenMutex.Unlock()
	fake.TokenStub = stub
}

func (fake *FakeClient) TokenReturns(result1 string, result2 error) {
	fake.tokenMutex.Lock()
	defer fake.tokenMutex.Unlock()
	fake.TokenStub = nil
	fake.tokenReturns = struct {
		result1 string
		result2 error
	}{result1, result2}
}

func (fake *FakeClient) TokenReturnsOnCall(i int, result1 string, result2 error) {
	fake.tokenMutex.Lock()
	defer fake.tokenMutex.Unlock()
	fake.TokenStub = nil
	if fake.tokenReturnsOnCall == nil {
		fake.tokenReturnsOnCall = make(map[int]struct {
			result1 string
			result2 error
		})
	}
	fake.tokenReturnsOnCall[i] = struct {
		result1 string
		result2 error
	}{result1, result2}
}

func (fake *FakeClient) WithAPIKey(arg1 string) {
	fake.withAPIKeyMutex.Lock()
	fake.withAPIKeyArgsForCall = append(fake.withAPIKeyArgsForCall, struct {
		arg1 string
	}{arg1})
	stub := fake.WithAPIKeyStub
	fake.recordInvocation("WithAPIKey", []interface{}{arg1})
	fake.withAPIKeyMutex.Unlock()
	if stub != nil {
		fake.WithAPIKeyStub(arg1)
	}
}

func (fake *FakeClient) WithAPIKeyCallCount() int {
	fake.withAPIKeyMutex.RLock()
	defer fake.withAPIKeyMutex.RUnlock()
	return len(fake.withAPIKeyArgsForCall)
}

func (fake *FakeClient) WithAPIKeyCalls(stub func(string)) {
	fake.withAPIKeyMutex.Lock()
	defer fake.withAPIKeyMutex.Unlock()
	fake.WithAPIKeyStub = stub
}

func (fake *FakeClient) WithAPIKeyArgsForCall(i int) string {
	fake.withAPIKeyMutex.RLock()
	defer fake.withAPIKeyMutex.RUnlock()
	argsForCall := fake.withAPIKeyArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeClient) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.tokenMutex.RLock()
	defer fake.tokenMutex.RUnlock()
	fake.withAPIKeyMutex.RLock()
	defer fake.withAPIKeyMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeClient) recordInvocation(key string, args []interface{}) {
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

var _ arcade.Client = new(FakeClient)
