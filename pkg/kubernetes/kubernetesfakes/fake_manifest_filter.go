// Code generated by counterfeiter. DO NOT EDIT.
package kubernetesfakes

import (
	"sync"

	"github.com/billiford/go-clouddriver/pkg/kubernetes"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type FakeManifestFilter struct {
	FilterOnClusterStub        func([]unstructured.Unstructured, string) []unstructured.Unstructured
	filterOnClusterMutex       sync.RWMutex
	filterOnClusterArgsForCall []struct {
		arg1 []unstructured.Unstructured
		arg2 string
	}
	filterOnClusterReturns struct {
		result1 []unstructured.Unstructured
	}
	filterOnClusterReturnsOnCall map[int]struct {
		result1 []unstructured.Unstructured
	}
	FilterOnLabelStub        func([]unstructured.Unstructured, string) []unstructured.Unstructured
	filterOnLabelMutex       sync.RWMutex
	filterOnLabelArgsForCall []struct {
		arg1 []unstructured.Unstructured
		arg2 string
	}
	filterOnLabelReturns struct {
		result1 []unstructured.Unstructured
	}
	filterOnLabelReturnsOnCall map[int]struct {
		result1 []unstructured.Unstructured
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeManifestFilter) FilterOnCluster(arg1 []unstructured.Unstructured, arg2 string) []unstructured.Unstructured {
	var arg1Copy []unstructured.Unstructured
	if arg1 != nil {
		arg1Copy = make([]unstructured.Unstructured, len(arg1))
		copy(arg1Copy, arg1)
	}
	fake.filterOnClusterMutex.Lock()
	ret, specificReturn := fake.filterOnClusterReturnsOnCall[len(fake.filterOnClusterArgsForCall)]
	fake.filterOnClusterArgsForCall = append(fake.filterOnClusterArgsForCall, struct {
		arg1 []unstructured.Unstructured
		arg2 string
	}{arg1Copy, arg2})
	fake.recordInvocation("FilterOnCluster", []interface{}{arg1Copy, arg2})
	fake.filterOnClusterMutex.Unlock()
	if fake.FilterOnClusterStub != nil {
		return fake.FilterOnClusterStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.filterOnClusterReturns
	return fakeReturns.result1
}

func (fake *FakeManifestFilter) FilterOnClusterCallCount() int {
	fake.filterOnClusterMutex.RLock()
	defer fake.filterOnClusterMutex.RUnlock()
	return len(fake.filterOnClusterArgsForCall)
}

func (fake *FakeManifestFilter) FilterOnClusterCalls(stub func([]unstructured.Unstructured, string) []unstructured.Unstructured) {
	fake.filterOnClusterMutex.Lock()
	defer fake.filterOnClusterMutex.Unlock()
	fake.FilterOnClusterStub = stub
}

func (fake *FakeManifestFilter) FilterOnClusterArgsForCall(i int) ([]unstructured.Unstructured, string) {
	fake.filterOnClusterMutex.RLock()
	defer fake.filterOnClusterMutex.RUnlock()
	argsForCall := fake.filterOnClusterArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeManifestFilter) FilterOnClusterReturns(result1 []unstructured.Unstructured) {
	fake.filterOnClusterMutex.Lock()
	defer fake.filterOnClusterMutex.Unlock()
	fake.FilterOnClusterStub = nil
	fake.filterOnClusterReturns = struct {
		result1 []unstructured.Unstructured
	}{result1}
}

func (fake *FakeManifestFilter) FilterOnClusterReturnsOnCall(i int, result1 []unstructured.Unstructured) {
	fake.filterOnClusterMutex.Lock()
	defer fake.filterOnClusterMutex.Unlock()
	fake.FilterOnClusterStub = nil
	if fake.filterOnClusterReturnsOnCall == nil {
		fake.filterOnClusterReturnsOnCall = make(map[int]struct {
			result1 []unstructured.Unstructured
		})
	}
	fake.filterOnClusterReturnsOnCall[i] = struct {
		result1 []unstructured.Unstructured
	}{result1}
}

func (fake *FakeManifestFilter) FilterOnLabel(arg1 []unstructured.Unstructured, arg2 string) []unstructured.Unstructured {
	var arg1Copy []unstructured.Unstructured
	if arg1 != nil {
		arg1Copy = make([]unstructured.Unstructured, len(arg1))
		copy(arg1Copy, arg1)
	}
	fake.filterOnLabelMutex.Lock()
	ret, specificReturn := fake.filterOnLabelReturnsOnCall[len(fake.filterOnLabelArgsForCall)]
	fake.filterOnLabelArgsForCall = append(fake.filterOnLabelArgsForCall, struct {
		arg1 []unstructured.Unstructured
		arg2 string
	}{arg1Copy, arg2})
	fake.recordInvocation("FilterOnLabel", []interface{}{arg1Copy, arg2})
	fake.filterOnLabelMutex.Unlock()
	if fake.FilterOnLabelStub != nil {
		return fake.FilterOnLabelStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.filterOnLabelReturns
	return fakeReturns.result1
}

func (fake *FakeManifestFilter) FilterOnLabelCallCount() int {
	fake.filterOnLabelMutex.RLock()
	defer fake.filterOnLabelMutex.RUnlock()
	return len(fake.filterOnLabelArgsForCall)
}

func (fake *FakeManifestFilter) FilterOnLabelCalls(stub func([]unstructured.Unstructured, string) []unstructured.Unstructured) {
	fake.filterOnLabelMutex.Lock()
	defer fake.filterOnLabelMutex.Unlock()
	fake.FilterOnLabelStub = stub
}

func (fake *FakeManifestFilter) FilterOnLabelArgsForCall(i int) ([]unstructured.Unstructured, string) {
	fake.filterOnLabelMutex.RLock()
	defer fake.filterOnLabelMutex.RUnlock()
	argsForCall := fake.filterOnLabelArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeManifestFilter) FilterOnLabelReturns(result1 []unstructured.Unstructured) {
	fake.filterOnLabelMutex.Lock()
	defer fake.filterOnLabelMutex.Unlock()
	fake.FilterOnLabelStub = nil
	fake.filterOnLabelReturns = struct {
		result1 []unstructured.Unstructured
	}{result1}
}

func (fake *FakeManifestFilter) FilterOnLabelReturnsOnCall(i int, result1 []unstructured.Unstructured) {
	fake.filterOnLabelMutex.Lock()
	defer fake.filterOnLabelMutex.Unlock()
	fake.FilterOnLabelStub = nil
	if fake.filterOnLabelReturnsOnCall == nil {
		fake.filterOnLabelReturnsOnCall = make(map[int]struct {
			result1 []unstructured.Unstructured
		})
	}
	fake.filterOnLabelReturnsOnCall[i] = struct {
		result1 []unstructured.Unstructured
	}{result1}
}

func (fake *FakeManifestFilter) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.filterOnClusterMutex.RLock()
	defer fake.filterOnClusterMutex.RUnlock()
	fake.filterOnLabelMutex.RLock()
	defer fake.filterOnLabelMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeManifestFilter) recordInvocation(key string, args []interface{}) {
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

var _ kubernetes.ManifestFilter = new(FakeManifestFilter)
