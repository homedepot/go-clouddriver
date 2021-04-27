// Code generated by counterfeiter. DO NOT EDIT.
package kubernetesfakes

import (
	"sync"

	"github.com/homedepot/go-clouddriver/pkg/kubernetes"
	"k8s.io/client-go/rest"
)

type FakeController struct {
	NewClientStub        func(*rest.Config) (kubernetes.Client, error)
	newClientMutex       sync.RWMutex
	newClientArgsForCall []struct {
		arg1 *rest.Config
	}
	newClientReturns struct {
		result1 kubernetes.Client
		result2 error
	}
	newClientReturnsOnCall map[int]struct {
		result1 kubernetes.Client
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeController) NewClient(arg1 *rest.Config) (kubernetes.Client, error) {
	fake.newClientMutex.Lock()
	ret, specificReturn := fake.newClientReturnsOnCall[len(fake.newClientArgsForCall)]
	fake.newClientArgsForCall = append(fake.newClientArgsForCall, struct {
		arg1 *rest.Config
	}{arg1})
	fake.recordInvocation("NewClient", []interface{}{arg1})
	fake.newClientMutex.Unlock()
	if fake.NewClientStub != nil {
		return fake.NewClientStub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.newClientReturns.result1, fake.newClientReturns.result2
}

func (fake *FakeController) NewClientCallCount() int {
	fake.newClientMutex.RLock()
	defer fake.newClientMutex.RUnlock()
	return len(fake.newClientArgsForCall)
}

func (fake *FakeController) NewClientArgsForCall(i int) *rest.Config {
	fake.newClientMutex.RLock()
	defer fake.newClientMutex.RUnlock()
	return fake.newClientArgsForCall[i].arg1
}

func (fake *FakeController) NewClientReturns(result1 kubernetes.Client, result2 error) {
	fake.NewClientStub = nil
	fake.newClientReturns = struct {
		result1 kubernetes.Client
		result2 error
	}{result1, result2}
}

func (fake *FakeController) NewClientReturnsOnCall(i int, result1 kubernetes.Client, result2 error) {
	fake.NewClientStub = nil
	if fake.newClientReturnsOnCall == nil {
		fake.newClientReturnsOnCall = make(map[int]struct {
			result1 kubernetes.Client
			result2 error
		})
	}
	fake.newClientReturnsOnCall[i] = struct {
		result1 kubernetes.Client
		result2 error
	}{result1, result2}
}

func (fake *FakeController) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.newClientMutex.RLock()
	defer fake.newClientMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeController) recordInvocation(key string, args []interface{}) {
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

var _ kubernetes.Controller = new(FakeController)
