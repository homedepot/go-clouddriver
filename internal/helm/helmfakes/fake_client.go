// Code generated by counterfeiter. DO NOT EDIT.
package helmfakes

import (
	"sync"

	"github.com/homedepot/go-clouddriver/internal/helm"
)

type FakeClient struct {
	GetChartStub        func(string, string) ([]byte, error)
	getChartMutex       sync.RWMutex
	getChartArgsForCall []struct {
		arg1 string
		arg2 string
	}
	getChartReturns struct {
		result1 []byte
		result2 error
	}
	getChartReturnsOnCall map[int]struct {
		result1 []byte
		result2 error
	}
	GetIndexStub        func() (helm.Index, error)
	getIndexMutex       sync.RWMutex
	getIndexArgsForCall []struct {
	}
	getIndexReturns struct {
		result1 helm.Index
		result2 error
	}
	getIndexReturnsOnCall map[int]struct {
		result1 helm.Index
		result2 error
	}
	WithUsernameAndPasswordStub        func(string, string)
	withUsernameAndPasswordMutex       sync.RWMutex
	withUsernameAndPasswordArgsForCall []struct {
		arg1 string
		arg2 string
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeClient) GetChart(arg1 string, arg2 string) ([]byte, error) {
	fake.getChartMutex.Lock()
	ret, specificReturn := fake.getChartReturnsOnCall[len(fake.getChartArgsForCall)]
	fake.getChartArgsForCall = append(fake.getChartArgsForCall, struct {
		arg1 string
		arg2 string
	}{arg1, arg2})
	stub := fake.GetChartStub
	fakeReturns := fake.getChartReturns
	fake.recordInvocation("GetChart", []interface{}{arg1, arg2})
	fake.getChartMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeClient) GetChartCallCount() int {
	fake.getChartMutex.RLock()
	defer fake.getChartMutex.RUnlock()
	return len(fake.getChartArgsForCall)
}

func (fake *FakeClient) GetChartCalls(stub func(string, string) ([]byte, error)) {
	fake.getChartMutex.Lock()
	defer fake.getChartMutex.Unlock()
	fake.GetChartStub = stub
}

func (fake *FakeClient) GetChartArgsForCall(i int) (string, string) {
	fake.getChartMutex.RLock()
	defer fake.getChartMutex.RUnlock()
	argsForCall := fake.getChartArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeClient) GetChartReturns(result1 []byte, result2 error) {
	fake.getChartMutex.Lock()
	defer fake.getChartMutex.Unlock()
	fake.GetChartStub = nil
	fake.getChartReturns = struct {
		result1 []byte
		result2 error
	}{result1, result2}
}

func (fake *FakeClient) GetChartReturnsOnCall(i int, result1 []byte, result2 error) {
	fake.getChartMutex.Lock()
	defer fake.getChartMutex.Unlock()
	fake.GetChartStub = nil
	if fake.getChartReturnsOnCall == nil {
		fake.getChartReturnsOnCall = make(map[int]struct {
			result1 []byte
			result2 error
		})
	}
	fake.getChartReturnsOnCall[i] = struct {
		result1 []byte
		result2 error
	}{result1, result2}
}

func (fake *FakeClient) GetIndex() (helm.Index, error) {
	fake.getIndexMutex.Lock()
	ret, specificReturn := fake.getIndexReturnsOnCall[len(fake.getIndexArgsForCall)]
	fake.getIndexArgsForCall = append(fake.getIndexArgsForCall, struct {
	}{})
	stub := fake.GetIndexStub
	fakeReturns := fake.getIndexReturns
	fake.recordInvocation("GetIndex", []interface{}{})
	fake.getIndexMutex.Unlock()
	if stub != nil {
		return stub()
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeClient) GetIndexCallCount() int {
	fake.getIndexMutex.RLock()
	defer fake.getIndexMutex.RUnlock()
	return len(fake.getIndexArgsForCall)
}

func (fake *FakeClient) GetIndexCalls(stub func() (helm.Index, error)) {
	fake.getIndexMutex.Lock()
	defer fake.getIndexMutex.Unlock()
	fake.GetIndexStub = stub
}

func (fake *FakeClient) GetIndexReturns(result1 helm.Index, result2 error) {
	fake.getIndexMutex.Lock()
	defer fake.getIndexMutex.Unlock()
	fake.GetIndexStub = nil
	fake.getIndexReturns = struct {
		result1 helm.Index
		result2 error
	}{result1, result2}
}

func (fake *FakeClient) GetIndexReturnsOnCall(i int, result1 helm.Index, result2 error) {
	fake.getIndexMutex.Lock()
	defer fake.getIndexMutex.Unlock()
	fake.GetIndexStub = nil
	if fake.getIndexReturnsOnCall == nil {
		fake.getIndexReturnsOnCall = make(map[int]struct {
			result1 helm.Index
			result2 error
		})
	}
	fake.getIndexReturnsOnCall[i] = struct {
		result1 helm.Index
		result2 error
	}{result1, result2}
}

func (fake *FakeClient) WithUsernameAndPassword(arg1 string, arg2 string) {
	fake.withUsernameAndPasswordMutex.Lock()
	fake.withUsernameAndPasswordArgsForCall = append(fake.withUsernameAndPasswordArgsForCall, struct {
		arg1 string
		arg2 string
	}{arg1, arg2})
	stub := fake.WithUsernameAndPasswordStub
	fake.recordInvocation("WithUsernameAndPassword", []interface{}{arg1, arg2})
	fake.withUsernameAndPasswordMutex.Unlock()
	if stub != nil {
		fake.WithUsernameAndPasswordStub(arg1, arg2)
	}
}

func (fake *FakeClient) WithUsernameAndPasswordCallCount() int {
	fake.withUsernameAndPasswordMutex.RLock()
	defer fake.withUsernameAndPasswordMutex.RUnlock()
	return len(fake.withUsernameAndPasswordArgsForCall)
}

func (fake *FakeClient) WithUsernameAndPasswordCalls(stub func(string, string)) {
	fake.withUsernameAndPasswordMutex.Lock()
	defer fake.withUsernameAndPasswordMutex.Unlock()
	fake.WithUsernameAndPasswordStub = stub
}

func (fake *FakeClient) WithUsernameAndPasswordArgsForCall(i int) (string, string) {
	fake.withUsernameAndPasswordMutex.RLock()
	defer fake.withUsernameAndPasswordMutex.RUnlock()
	argsForCall := fake.withUsernameAndPasswordArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeClient) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.getChartMutex.RLock()
	defer fake.getChartMutex.RUnlock()
	fake.getIndexMutex.RLock()
	defer fake.getIndexMutex.RUnlock()
	fake.withUsernameAndPasswordMutex.RLock()
	defer fake.withUsernameAndPasswordMutex.RUnlock()
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

var _ helm.Client = new(FakeClient)
