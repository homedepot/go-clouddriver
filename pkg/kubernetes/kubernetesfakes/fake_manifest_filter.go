// Code generated by counterfeiter. DO NOT EDIT.
package kubernetesfakes

import (
	"sync"

	"github.com/homedepot/go-clouddriver/pkg/kubernetes"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type FakeManifestFilter struct {
	FilterOnClusterAnnotationStub        func([]unstructured.Unstructured, string) []unstructured.Unstructured
	filterOnClusterAnnotationMutex       sync.RWMutex
	filterOnClusterAnnotationArgsForCall []struct {
		arg1 []unstructured.Unstructured
		arg2 string
	}
	filterOnClusterAnnotationReturns struct {
		result1 []unstructured.Unstructured
	}
	filterOnClusterAnnotationReturnsOnCall map[int]struct {
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

func (fake *FakeManifestFilter) FilterOnClusterAnnotation(arg1 []unstructured.Unstructured, arg2 string) []unstructured.Unstructured {
	var arg1Copy []unstructured.Unstructured
	if arg1 != nil {
		arg1Copy = make([]unstructured.Unstructured, len(arg1))
		copy(arg1Copy, arg1)
	}
	fake.filterOnClusterAnnotationMutex.Lock()
	ret, specificReturn := fake.filterOnClusterAnnotationReturnsOnCall[len(fake.filterOnClusterAnnotationArgsForCall)]
	fake.filterOnClusterAnnotationArgsForCall = append(fake.filterOnClusterAnnotationArgsForCall, struct {
		arg1 []unstructured.Unstructured
		arg2 string
	}{arg1Copy, arg2})
	stub := fake.FilterOnClusterAnnotationStub
	fakeReturns := fake.filterOnClusterAnnotationReturns
	fake.recordInvocation("FilterOnClusterAnnotation", []interface{}{arg1Copy, arg2})
	fake.filterOnClusterAnnotationMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeManifestFilter) FilterOnClusterAnnotationCallCount() int {
	fake.filterOnClusterAnnotationMutex.RLock()
	defer fake.filterOnClusterAnnotationMutex.RUnlock()
	return len(fake.filterOnClusterAnnotationArgsForCall)
}

func (fake *FakeManifestFilter) FilterOnClusterAnnotationCalls(stub func([]unstructured.Unstructured, string) []unstructured.Unstructured) {
	fake.filterOnClusterAnnotationMutex.Lock()
	defer fake.filterOnClusterAnnotationMutex.Unlock()
	fake.FilterOnClusterAnnotationStub = stub
}

func (fake *FakeManifestFilter) FilterOnClusterAnnotationArgsForCall(i int) ([]unstructured.Unstructured, string) {
	fake.filterOnClusterAnnotationMutex.RLock()
	defer fake.filterOnClusterAnnotationMutex.RUnlock()
	argsForCall := fake.filterOnClusterAnnotationArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeManifestFilter) FilterOnClusterAnnotationReturns(result1 []unstructured.Unstructured) {
	fake.filterOnClusterAnnotationMutex.Lock()
	defer fake.filterOnClusterAnnotationMutex.Unlock()
	fake.FilterOnClusterAnnotationStub = nil
	fake.filterOnClusterAnnotationReturns = struct {
		result1 []unstructured.Unstructured
	}{result1}
}

func (fake *FakeManifestFilter) FilterOnClusterAnnotationReturnsOnCall(i int, result1 []unstructured.Unstructured) {
	fake.filterOnClusterAnnotationMutex.Lock()
	defer fake.filterOnClusterAnnotationMutex.Unlock()
	fake.FilterOnClusterAnnotationStub = nil
	if fake.filterOnClusterAnnotationReturnsOnCall == nil {
		fake.filterOnClusterAnnotationReturnsOnCall = make(map[int]struct {
			result1 []unstructured.Unstructured
		})
	}
	fake.filterOnClusterAnnotationReturnsOnCall[i] = struct {
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
	stub := fake.FilterOnLabelStub
	fakeReturns := fake.filterOnLabelReturns
	fake.recordInvocation("FilterOnLabel", []interface{}{arg1Copy, arg2})
	fake.filterOnLabelMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1
	}
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
	fake.filterOnClusterAnnotationMutex.RLock()
	defer fake.filterOnClusterAnnotationMutex.RUnlock()
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
