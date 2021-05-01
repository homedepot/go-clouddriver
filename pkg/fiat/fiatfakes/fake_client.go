// Code generated by counterfeiter. DO NOT EDIT.
package fiatfakes

import (
	"sync"

	"github.com/homedepot/go-clouddriver/pkg/fiat"
)

type FakeClient struct {
	AuthorizeStub        func(string) (fiat.Response, error)
	authorizeMutex       sync.RWMutex
	authorizeArgsForCall []struct {
		arg1 string
	}
	authorizeReturns struct {
		result1 fiat.Response
		result2 error
	}
	authorizeReturnsOnCall map[int]struct {
		result1 fiat.Response
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeClient) Authorize(arg1 string) (fiat.Response, error) {
	fake.authorizeMutex.Lock()
	ret, specificReturn := fake.authorizeReturnsOnCall[len(fake.authorizeArgsForCall)]
	fake.authorizeArgsForCall = append(fake.authorizeArgsForCall, struct {
		arg1 string
	}{arg1})
	stub := fake.AuthorizeStub
	fakeReturns := fake.authorizeReturns
	fake.recordInvocation("Authorize", []interface{}{arg1})
	fake.authorizeMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeClient) AuthorizeCallCount() int {
	fake.authorizeMutex.RLock()
	defer fake.authorizeMutex.RUnlock()
	return len(fake.authorizeArgsForCall)
}

func (fake *FakeClient) AuthorizeCalls(stub func(string) (fiat.Response, error)) {
	fake.authorizeMutex.Lock()
	defer fake.authorizeMutex.Unlock()
	fake.AuthorizeStub = stub
}

func (fake *FakeClient) AuthorizeArgsForCall(i int) string {
	fake.authorizeMutex.RLock()
	defer fake.authorizeMutex.RUnlock()
	argsForCall := fake.authorizeArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeClient) AuthorizeReturns(result1 fiat.Response, result2 error) {
	fake.authorizeMutex.Lock()
	defer fake.authorizeMutex.Unlock()
	fake.AuthorizeStub = nil
	fake.authorizeReturns = struct {
		result1 fiat.Response
		result2 error
	}{result1, result2}
}

func (fake *FakeClient) AuthorizeReturnsOnCall(i int, result1 fiat.Response, result2 error) {
	fake.authorizeMutex.Lock()
	defer fake.authorizeMutex.Unlock()
	fake.AuthorizeStub = nil
	if fake.authorizeReturnsOnCall == nil {
		fake.authorizeReturnsOnCall = make(map[int]struct {
			result1 fiat.Response
			result2 error
		})
	}
	fake.authorizeReturnsOnCall[i] = struct {
		result1 fiat.Response
		result2 error
	}{result1, result2}
}

func (fake *FakeClient) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.authorizeMutex.RLock()
	defer fake.authorizeMutex.RUnlock()
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

var _ fiat.Client = new(FakeClient)
