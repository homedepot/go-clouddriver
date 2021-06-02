// Code generated by counterfeiter. DO NOT EDIT.
package kubernetesfakes

import (
	"context"
	"sync"

	"github.com/homedepot/go-clouddriver/pkg/kubernetes"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
)

type FakeClient struct {
	ApplyStub        func(*unstructured.Unstructured) (kubernetes.Metadata, error)
	applyMutex       sync.RWMutex
	applyArgsForCall []struct {
		arg1 *unstructured.Unstructured
	}
	applyReturns struct {
		result1 kubernetes.Metadata
		result2 error
	}
	applyReturnsOnCall map[int]struct {
		result1 kubernetes.Metadata
		result2 error
	}
	ApplyWithNamespaceOverrideStub        func(*unstructured.Unstructured, string) (kubernetes.Metadata, error)
	applyWithNamespaceOverrideMutex       sync.RWMutex
	applyWithNamespaceOverrideArgsForCall []struct {
		arg1 *unstructured.Unstructured
		arg2 string
	}
	applyWithNamespaceOverrideReturns struct {
		result1 kubernetes.Metadata
		result2 error
	}
	applyWithNamespaceOverrideReturnsOnCall map[int]struct {
		result1 kubernetes.Metadata
		result2 error
	}
	DeleteResourceByKindAndNameAndNamespaceStub        func(string, string, string, v1.DeleteOptions) error
	deleteResourceByKindAndNameAndNamespaceMutex       sync.RWMutex
	deleteResourceByKindAndNameAndNamespaceArgsForCall []struct {
		arg1 string
		arg2 string
		arg3 string
		arg4 v1.DeleteOptions
	}
	deleteResourceByKindAndNameAndNamespaceReturns struct {
		result1 error
	}
	deleteResourceByKindAndNameAndNamespaceReturnsOnCall map[int]struct {
		result1 error
	}
	DiscoverStub        func() error
	discoverMutex       sync.RWMutex
	discoverArgsForCall []struct {
	}
	discoverReturns struct {
		result1 error
	}
	discoverReturnsOnCall map[int]struct {
		result1 error
	}
	GVRForKindStub        func(string) (schema.GroupVersionResource, error)
	gVRForKindMutex       sync.RWMutex
	gVRForKindArgsForCall []struct {
		arg1 string
	}
	gVRForKindReturns struct {
		result1 schema.GroupVersionResource
		result2 error
	}
	gVRForKindReturnsOnCall map[int]struct {
		result1 schema.GroupVersionResource
		result2 error
	}
	GetStub        func(string, string, string) (*unstructured.Unstructured, error)
	getMutex       sync.RWMutex
	getArgsForCall []struct {
		arg1 string
		arg2 string
		arg3 string
	}
	getReturns struct {
		result1 *unstructured.Unstructured
		result2 error
	}
	getReturnsOnCall map[int]struct {
		result1 *unstructured.Unstructured
		result2 error
	}
	ListByGVRStub        func(schema.GroupVersionResource, v1.ListOptions) (*unstructured.UnstructuredList, error)
	listByGVRMutex       sync.RWMutex
	listByGVRArgsForCall []struct {
		arg1 schema.GroupVersionResource
		arg2 v1.ListOptions
	}
	listByGVRReturns struct {
		result1 *unstructured.UnstructuredList
		result2 error
	}
	listByGVRReturnsOnCall map[int]struct {
		result1 *unstructured.UnstructuredList
		result2 error
	}
	ListByGVRWithContextStub        func(context.Context, schema.GroupVersionResource, v1.ListOptions) (*unstructured.UnstructuredList, error)
	listByGVRWithContextMutex       sync.RWMutex
	listByGVRWithContextArgsForCall []struct {
		arg1 context.Context
		arg2 schema.GroupVersionResource
		arg3 v1.ListOptions
	}
	listByGVRWithContextReturns struct {
		result1 *unstructured.UnstructuredList
		result2 error
	}
	listByGVRWithContextReturnsOnCall map[int]struct {
		result1 *unstructured.UnstructuredList
		result2 error
	}
	ListResourceStub        func(string, v1.ListOptions) (*unstructured.UnstructuredList, error)
	listResourceMutex       sync.RWMutex
	listResourceArgsForCall []struct {
		arg1 string
		arg2 v1.ListOptions
	}
	listResourceReturns struct {
		result1 *unstructured.UnstructuredList
		result2 error
	}
	listResourceReturnsOnCall map[int]struct {
		result1 *unstructured.UnstructuredList
		result2 error
	}
	ListResourceWithContextStub        func(context.Context, string, v1.ListOptions) (*unstructured.UnstructuredList, error)
	listResourceWithContextMutex       sync.RWMutex
	listResourceWithContextArgsForCall []struct {
		arg1 context.Context
		arg2 string
		arg3 v1.ListOptions
	}
	listResourceWithContextReturns struct {
		result1 *unstructured.UnstructuredList
		result2 error
	}
	listResourceWithContextReturnsOnCall map[int]struct {
		result1 *unstructured.UnstructuredList
		result2 error
	}
	ListResourcesByKindAndNamespaceStub        func(string, string, v1.ListOptions) (*unstructured.UnstructuredList, error)
	listResourcesByKindAndNamespaceMutex       sync.RWMutex
	listResourcesByKindAndNamespaceArgsForCall []struct {
		arg1 string
		arg2 string
		arg3 v1.ListOptions
	}
	listResourcesByKindAndNamespaceReturns struct {
		result1 *unstructured.UnstructuredList
		result2 error
	}
	listResourcesByKindAndNamespaceReturnsOnCall map[int]struct {
		result1 *unstructured.UnstructuredList
		result2 error
	}
	PatchStub        func(string, string, string, []byte) (kubernetes.Metadata, *unstructured.Unstructured, error)
	patchMutex       sync.RWMutex
	patchArgsForCall []struct {
		arg1 string
		arg2 string
		arg3 string
		arg4 []byte
	}
	patchReturns struct {
		result1 kubernetes.Metadata
		result2 *unstructured.Unstructured
		result3 error
	}
	patchReturnsOnCall map[int]struct {
		result1 kubernetes.Metadata
		result2 *unstructured.Unstructured
		result3 error
	}
	PatchUsingStrategyStub        func(string, string, string, []byte, types.PatchType) (kubernetes.Metadata, *unstructured.Unstructured, error)
	patchUsingStrategyMutex       sync.RWMutex
	patchUsingStrategyArgsForCall []struct {
		arg1 string
		arg2 string
		arg3 string
		arg4 []byte
		arg5 types.PatchType
	}
	patchUsingStrategyReturns struct {
		result1 kubernetes.Metadata
		result2 *unstructured.Unstructured
		result3 error
	}
	patchUsingStrategyReturnsOnCall map[int]struct {
		result1 kubernetes.Metadata
		result2 *unstructured.Unstructured
		result3 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeClient) Apply(arg1 *unstructured.Unstructured) (kubernetes.Metadata, error) {
	fake.applyMutex.Lock()
	ret, specificReturn := fake.applyReturnsOnCall[len(fake.applyArgsForCall)]
	fake.applyArgsForCall = append(fake.applyArgsForCall, struct {
		arg1 *unstructured.Unstructured
	}{arg1})
	stub := fake.ApplyStub
	fakeReturns := fake.applyReturns
	fake.recordInvocation("Apply", []interface{}{arg1})
	fake.applyMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeClient) ApplyCallCount() int {
	fake.applyMutex.RLock()
	defer fake.applyMutex.RUnlock()
	return len(fake.applyArgsForCall)
}

func (fake *FakeClient) ApplyCalls(stub func(*unstructured.Unstructured) (kubernetes.Metadata, error)) {
	fake.applyMutex.Lock()
	defer fake.applyMutex.Unlock()
	fake.ApplyStub = stub
}

func (fake *FakeClient) ApplyArgsForCall(i int) *unstructured.Unstructured {
	fake.applyMutex.RLock()
	defer fake.applyMutex.RUnlock()
	argsForCall := fake.applyArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeClient) ApplyReturns(result1 kubernetes.Metadata, result2 error) {
	fake.applyMutex.Lock()
	defer fake.applyMutex.Unlock()
	fake.ApplyStub = nil
	fake.applyReturns = struct {
		result1 kubernetes.Metadata
		result2 error
	}{result1, result2}
}

func (fake *FakeClient) ApplyReturnsOnCall(i int, result1 kubernetes.Metadata, result2 error) {
	fake.applyMutex.Lock()
	defer fake.applyMutex.Unlock()
	fake.ApplyStub = nil
	if fake.applyReturnsOnCall == nil {
		fake.applyReturnsOnCall = make(map[int]struct {
			result1 kubernetes.Metadata
			result2 error
		})
	}
	fake.applyReturnsOnCall[i] = struct {
		result1 kubernetes.Metadata
		result2 error
	}{result1, result2}
}

func (fake *FakeClient) ApplyWithNamespaceOverride(arg1 *unstructured.Unstructured, arg2 string) (kubernetes.Metadata, error) {
	fake.applyWithNamespaceOverrideMutex.Lock()
	ret, specificReturn := fake.applyWithNamespaceOverrideReturnsOnCall[len(fake.applyWithNamespaceOverrideArgsForCall)]
	fake.applyWithNamespaceOverrideArgsForCall = append(fake.applyWithNamespaceOverrideArgsForCall, struct {
		arg1 *unstructured.Unstructured
		arg2 string
	}{arg1, arg2})
	stub := fake.ApplyWithNamespaceOverrideStub
	fakeReturns := fake.applyWithNamespaceOverrideReturns
	fake.recordInvocation("ApplyWithNamespaceOverride", []interface{}{arg1, arg2})
	fake.applyWithNamespaceOverrideMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeClient) ApplyWithNamespaceOverrideCallCount() int {
	fake.applyWithNamespaceOverrideMutex.RLock()
	defer fake.applyWithNamespaceOverrideMutex.RUnlock()
	return len(fake.applyWithNamespaceOverrideArgsForCall)
}

func (fake *FakeClient) ApplyWithNamespaceOverrideCalls(stub func(*unstructured.Unstructured, string) (kubernetes.Metadata, error)) {
	fake.applyWithNamespaceOverrideMutex.Lock()
	defer fake.applyWithNamespaceOverrideMutex.Unlock()
	fake.ApplyWithNamespaceOverrideStub = stub
}

func (fake *FakeClient) ApplyWithNamespaceOverrideArgsForCall(i int) (*unstructured.Unstructured, string) {
	fake.applyWithNamespaceOverrideMutex.RLock()
	defer fake.applyWithNamespaceOverrideMutex.RUnlock()
	argsForCall := fake.applyWithNamespaceOverrideArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeClient) ApplyWithNamespaceOverrideReturns(result1 kubernetes.Metadata, result2 error) {
	fake.applyWithNamespaceOverrideMutex.Lock()
	defer fake.applyWithNamespaceOverrideMutex.Unlock()
	fake.ApplyWithNamespaceOverrideStub = nil
	fake.applyWithNamespaceOverrideReturns = struct {
		result1 kubernetes.Metadata
		result2 error
	}{result1, result2}
}

func (fake *FakeClient) ApplyWithNamespaceOverrideReturnsOnCall(i int, result1 kubernetes.Metadata, result2 error) {
	fake.applyWithNamespaceOverrideMutex.Lock()
	defer fake.applyWithNamespaceOverrideMutex.Unlock()
	fake.ApplyWithNamespaceOverrideStub = nil
	if fake.applyWithNamespaceOverrideReturnsOnCall == nil {
		fake.applyWithNamespaceOverrideReturnsOnCall = make(map[int]struct {
			result1 kubernetes.Metadata
			result2 error
		})
	}
	fake.applyWithNamespaceOverrideReturnsOnCall[i] = struct {
		result1 kubernetes.Metadata
		result2 error
	}{result1, result2}
}

func (fake *FakeClient) DeleteResourceByKindAndNameAndNamespace(arg1 string, arg2 string, arg3 string, arg4 v1.DeleteOptions) error {
	fake.deleteResourceByKindAndNameAndNamespaceMutex.Lock()
	ret, specificReturn := fake.deleteResourceByKindAndNameAndNamespaceReturnsOnCall[len(fake.deleteResourceByKindAndNameAndNamespaceArgsForCall)]
	fake.deleteResourceByKindAndNameAndNamespaceArgsForCall = append(fake.deleteResourceByKindAndNameAndNamespaceArgsForCall, struct {
		arg1 string
		arg2 string
		arg3 string
		arg4 v1.DeleteOptions
	}{arg1, arg2, arg3, arg4})
	stub := fake.DeleteResourceByKindAndNameAndNamespaceStub
	fakeReturns := fake.deleteResourceByKindAndNameAndNamespaceReturns
	fake.recordInvocation("DeleteResourceByKindAndNameAndNamespace", []interface{}{arg1, arg2, arg3, arg4})
	fake.deleteResourceByKindAndNameAndNamespaceMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3, arg4)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeClient) DeleteResourceByKindAndNameAndNamespaceCallCount() int {
	fake.deleteResourceByKindAndNameAndNamespaceMutex.RLock()
	defer fake.deleteResourceByKindAndNameAndNamespaceMutex.RUnlock()
	return len(fake.deleteResourceByKindAndNameAndNamespaceArgsForCall)
}

func (fake *FakeClient) DeleteResourceByKindAndNameAndNamespaceCalls(stub func(string, string, string, v1.DeleteOptions) error) {
	fake.deleteResourceByKindAndNameAndNamespaceMutex.Lock()
	defer fake.deleteResourceByKindAndNameAndNamespaceMutex.Unlock()
	fake.DeleteResourceByKindAndNameAndNamespaceStub = stub
}

func (fake *FakeClient) DeleteResourceByKindAndNameAndNamespaceArgsForCall(i int) (string, string, string, v1.DeleteOptions) {
	fake.deleteResourceByKindAndNameAndNamespaceMutex.RLock()
	defer fake.deleteResourceByKindAndNameAndNamespaceMutex.RUnlock()
	argsForCall := fake.deleteResourceByKindAndNameAndNamespaceArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3, argsForCall.arg4
}

func (fake *FakeClient) DeleteResourceByKindAndNameAndNamespaceReturns(result1 error) {
	fake.deleteResourceByKindAndNameAndNamespaceMutex.Lock()
	defer fake.deleteResourceByKindAndNameAndNamespaceMutex.Unlock()
	fake.DeleteResourceByKindAndNameAndNamespaceStub = nil
	fake.deleteResourceByKindAndNameAndNamespaceReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeClient) DeleteResourceByKindAndNameAndNamespaceReturnsOnCall(i int, result1 error) {
	fake.deleteResourceByKindAndNameAndNamespaceMutex.Lock()
	defer fake.deleteResourceByKindAndNameAndNamespaceMutex.Unlock()
	fake.DeleteResourceByKindAndNameAndNamespaceStub = nil
	if fake.deleteResourceByKindAndNameAndNamespaceReturnsOnCall == nil {
		fake.deleteResourceByKindAndNameAndNamespaceReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.deleteResourceByKindAndNameAndNamespaceReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeClient) Discover() error {
	fake.discoverMutex.Lock()
	ret, specificReturn := fake.discoverReturnsOnCall[len(fake.discoverArgsForCall)]
	fake.discoverArgsForCall = append(fake.discoverArgsForCall, struct {
	}{})
	stub := fake.DiscoverStub
	fakeReturns := fake.discoverReturns
	fake.recordInvocation("Discover", []interface{}{})
	fake.discoverMutex.Unlock()
	if stub != nil {
		return stub()
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeClient) DiscoverCallCount() int {
	fake.discoverMutex.RLock()
	defer fake.discoverMutex.RUnlock()
	return len(fake.discoverArgsForCall)
}

func (fake *FakeClient) DiscoverCalls(stub func() error) {
	fake.discoverMutex.Lock()
	defer fake.discoverMutex.Unlock()
	fake.DiscoverStub = stub
}

func (fake *FakeClient) DiscoverReturns(result1 error) {
	fake.discoverMutex.Lock()
	defer fake.discoverMutex.Unlock()
	fake.DiscoverStub = nil
	fake.discoverReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeClient) DiscoverReturnsOnCall(i int, result1 error) {
	fake.discoverMutex.Lock()
	defer fake.discoverMutex.Unlock()
	fake.DiscoverStub = nil
	if fake.discoverReturnsOnCall == nil {
		fake.discoverReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.discoverReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeClient) GVRForKind(arg1 string) (schema.GroupVersionResource, error) {
	fake.gVRForKindMutex.Lock()
	ret, specificReturn := fake.gVRForKindReturnsOnCall[len(fake.gVRForKindArgsForCall)]
	fake.gVRForKindArgsForCall = append(fake.gVRForKindArgsForCall, struct {
		arg1 string
	}{arg1})
	stub := fake.GVRForKindStub
	fakeReturns := fake.gVRForKindReturns
	fake.recordInvocation("GVRForKind", []interface{}{arg1})
	fake.gVRForKindMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeClient) GVRForKindCallCount() int {
	fake.gVRForKindMutex.RLock()
	defer fake.gVRForKindMutex.RUnlock()
	return len(fake.gVRForKindArgsForCall)
}

func (fake *FakeClient) GVRForKindCalls(stub func(string) (schema.GroupVersionResource, error)) {
	fake.gVRForKindMutex.Lock()
	defer fake.gVRForKindMutex.Unlock()
	fake.GVRForKindStub = stub
}

func (fake *FakeClient) GVRForKindArgsForCall(i int) string {
	fake.gVRForKindMutex.RLock()
	defer fake.gVRForKindMutex.RUnlock()
	argsForCall := fake.gVRForKindArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeClient) GVRForKindReturns(result1 schema.GroupVersionResource, result2 error) {
	fake.gVRForKindMutex.Lock()
	defer fake.gVRForKindMutex.Unlock()
	fake.GVRForKindStub = nil
	fake.gVRForKindReturns = struct {
		result1 schema.GroupVersionResource
		result2 error
	}{result1, result2}
}

func (fake *FakeClient) GVRForKindReturnsOnCall(i int, result1 schema.GroupVersionResource, result2 error) {
	fake.gVRForKindMutex.Lock()
	defer fake.gVRForKindMutex.Unlock()
	fake.GVRForKindStub = nil
	if fake.gVRForKindReturnsOnCall == nil {
		fake.gVRForKindReturnsOnCall = make(map[int]struct {
			result1 schema.GroupVersionResource
			result2 error
		})
	}
	fake.gVRForKindReturnsOnCall[i] = struct {
		result1 schema.GroupVersionResource
		result2 error
	}{result1, result2}
}

func (fake *FakeClient) Get(arg1 string, arg2 string, arg3 string) (*unstructured.Unstructured, error) {
	fake.getMutex.Lock()
	ret, specificReturn := fake.getReturnsOnCall[len(fake.getArgsForCall)]
	fake.getArgsForCall = append(fake.getArgsForCall, struct {
		arg1 string
		arg2 string
		arg3 string
	}{arg1, arg2, arg3})
	stub := fake.GetStub
	fakeReturns := fake.getReturns
	fake.recordInvocation("Get", []interface{}{arg1, arg2, arg3})
	fake.getMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeClient) GetCallCount() int {
	fake.getMutex.RLock()
	defer fake.getMutex.RUnlock()
	return len(fake.getArgsForCall)
}

func (fake *FakeClient) GetCalls(stub func(string, string, string) (*unstructured.Unstructured, error)) {
	fake.getMutex.Lock()
	defer fake.getMutex.Unlock()
	fake.GetStub = stub
}

func (fake *FakeClient) GetArgsForCall(i int) (string, string, string) {
	fake.getMutex.RLock()
	defer fake.getMutex.RUnlock()
	argsForCall := fake.getArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *FakeClient) GetReturns(result1 *unstructured.Unstructured, result2 error) {
	fake.getMutex.Lock()
	defer fake.getMutex.Unlock()
	fake.GetStub = nil
	fake.getReturns = struct {
		result1 *unstructured.Unstructured
		result2 error
	}{result1, result2}
}

func (fake *FakeClient) GetReturnsOnCall(i int, result1 *unstructured.Unstructured, result2 error) {
	fake.getMutex.Lock()
	defer fake.getMutex.Unlock()
	fake.GetStub = nil
	if fake.getReturnsOnCall == nil {
		fake.getReturnsOnCall = make(map[int]struct {
			result1 *unstructured.Unstructured
			result2 error
		})
	}
	fake.getReturnsOnCall[i] = struct {
		result1 *unstructured.Unstructured
		result2 error
	}{result1, result2}
}

func (fake *FakeClient) ListByGVR(arg1 schema.GroupVersionResource, arg2 v1.ListOptions) (*unstructured.UnstructuredList, error) {
	fake.listByGVRMutex.Lock()
	ret, specificReturn := fake.listByGVRReturnsOnCall[len(fake.listByGVRArgsForCall)]
	fake.listByGVRArgsForCall = append(fake.listByGVRArgsForCall, struct {
		arg1 schema.GroupVersionResource
		arg2 v1.ListOptions
	}{arg1, arg2})
	stub := fake.ListByGVRStub
	fakeReturns := fake.listByGVRReturns
	fake.recordInvocation("ListByGVR", []interface{}{arg1, arg2})
	fake.listByGVRMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeClient) ListByGVRCallCount() int {
	fake.listByGVRMutex.RLock()
	defer fake.listByGVRMutex.RUnlock()
	return len(fake.listByGVRArgsForCall)
}

func (fake *FakeClient) ListByGVRCalls(stub func(schema.GroupVersionResource, v1.ListOptions) (*unstructured.UnstructuredList, error)) {
	fake.listByGVRMutex.Lock()
	defer fake.listByGVRMutex.Unlock()
	fake.ListByGVRStub = stub
}

func (fake *FakeClient) ListByGVRArgsForCall(i int) (schema.GroupVersionResource, v1.ListOptions) {
	fake.listByGVRMutex.RLock()
	defer fake.listByGVRMutex.RUnlock()
	argsForCall := fake.listByGVRArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeClient) ListByGVRReturns(result1 *unstructured.UnstructuredList, result2 error) {
	fake.listByGVRMutex.Lock()
	defer fake.listByGVRMutex.Unlock()
	fake.ListByGVRStub = nil
	fake.listByGVRReturns = struct {
		result1 *unstructured.UnstructuredList
		result2 error
	}{result1, result2}
}

func (fake *FakeClient) ListByGVRReturnsOnCall(i int, result1 *unstructured.UnstructuredList, result2 error) {
	fake.listByGVRMutex.Lock()
	defer fake.listByGVRMutex.Unlock()
	fake.ListByGVRStub = nil
	if fake.listByGVRReturnsOnCall == nil {
		fake.listByGVRReturnsOnCall = make(map[int]struct {
			result1 *unstructured.UnstructuredList
			result2 error
		})
	}
	fake.listByGVRReturnsOnCall[i] = struct {
		result1 *unstructured.UnstructuredList
		result2 error
	}{result1, result2}
}

func (fake *FakeClient) ListByGVRWithContext(arg1 context.Context, arg2 schema.GroupVersionResource, arg3 v1.ListOptions) (*unstructured.UnstructuredList, error) {
	fake.listByGVRWithContextMutex.Lock()
	ret, specificReturn := fake.listByGVRWithContextReturnsOnCall[len(fake.listByGVRWithContextArgsForCall)]
	fake.listByGVRWithContextArgsForCall = append(fake.listByGVRWithContextArgsForCall, struct {
		arg1 context.Context
		arg2 schema.GroupVersionResource
		arg3 v1.ListOptions
	}{arg1, arg2, arg3})
	stub := fake.ListByGVRWithContextStub
	fakeReturns := fake.listByGVRWithContextReturns
	fake.recordInvocation("ListByGVRWithContext", []interface{}{arg1, arg2, arg3})
	fake.listByGVRWithContextMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeClient) ListByGVRWithContextCallCount() int {
	fake.listByGVRWithContextMutex.RLock()
	defer fake.listByGVRWithContextMutex.RUnlock()
	return len(fake.listByGVRWithContextArgsForCall)
}

func (fake *FakeClient) ListByGVRWithContextCalls(stub func(context.Context, schema.GroupVersionResource, v1.ListOptions) (*unstructured.UnstructuredList, error)) {
	fake.listByGVRWithContextMutex.Lock()
	defer fake.listByGVRWithContextMutex.Unlock()
	fake.ListByGVRWithContextStub = stub
}

func (fake *FakeClient) ListByGVRWithContextArgsForCall(i int) (context.Context, schema.GroupVersionResource, v1.ListOptions) {
	fake.listByGVRWithContextMutex.RLock()
	defer fake.listByGVRWithContextMutex.RUnlock()
	argsForCall := fake.listByGVRWithContextArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *FakeClient) ListByGVRWithContextReturns(result1 *unstructured.UnstructuredList, result2 error) {
	fake.listByGVRWithContextMutex.Lock()
	defer fake.listByGVRWithContextMutex.Unlock()
	fake.ListByGVRWithContextStub = nil
	fake.listByGVRWithContextReturns = struct {
		result1 *unstructured.UnstructuredList
		result2 error
	}{result1, result2}
}

func (fake *FakeClient) ListByGVRWithContextReturnsOnCall(i int, result1 *unstructured.UnstructuredList, result2 error) {
	fake.listByGVRWithContextMutex.Lock()
	defer fake.listByGVRWithContextMutex.Unlock()
	fake.ListByGVRWithContextStub = nil
	if fake.listByGVRWithContextReturnsOnCall == nil {
		fake.listByGVRWithContextReturnsOnCall = make(map[int]struct {
			result1 *unstructured.UnstructuredList
			result2 error
		})
	}
	fake.listByGVRWithContextReturnsOnCall[i] = struct {
		result1 *unstructured.UnstructuredList
		result2 error
	}{result1, result2}
}

func (fake *FakeClient) ListResource(arg1 string, arg2 v1.ListOptions) (*unstructured.UnstructuredList, error) {
	fake.listResourceMutex.Lock()
	ret, specificReturn := fake.listResourceReturnsOnCall[len(fake.listResourceArgsForCall)]
	fake.listResourceArgsForCall = append(fake.listResourceArgsForCall, struct {
		arg1 string
		arg2 v1.ListOptions
	}{arg1, arg2})
	stub := fake.ListResourceStub
	fakeReturns := fake.listResourceReturns
	fake.recordInvocation("ListResource", []interface{}{arg1, arg2})
	fake.listResourceMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeClient) ListResourceCallCount() int {
	fake.listResourceMutex.RLock()
	defer fake.listResourceMutex.RUnlock()
	return len(fake.listResourceArgsForCall)
}

func (fake *FakeClient) ListResourceCalls(stub func(string, v1.ListOptions) (*unstructured.UnstructuredList, error)) {
	fake.listResourceMutex.Lock()
	defer fake.listResourceMutex.Unlock()
	fake.ListResourceStub = stub
}

func (fake *FakeClient) ListResourceArgsForCall(i int) (string, v1.ListOptions) {
	fake.listResourceMutex.RLock()
	defer fake.listResourceMutex.RUnlock()
	argsForCall := fake.listResourceArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeClient) ListResourceReturns(result1 *unstructured.UnstructuredList, result2 error) {
	fake.listResourceMutex.Lock()
	defer fake.listResourceMutex.Unlock()
	fake.ListResourceStub = nil
	fake.listResourceReturns = struct {
		result1 *unstructured.UnstructuredList
		result2 error
	}{result1, result2}
}

func (fake *FakeClient) ListResourceReturnsOnCall(i int, result1 *unstructured.UnstructuredList, result2 error) {
	fake.listResourceMutex.Lock()
	defer fake.listResourceMutex.Unlock()
	fake.ListResourceStub = nil
	if fake.listResourceReturnsOnCall == nil {
		fake.listResourceReturnsOnCall = make(map[int]struct {
			result1 *unstructured.UnstructuredList
			result2 error
		})
	}
	fake.listResourceReturnsOnCall[i] = struct {
		result1 *unstructured.UnstructuredList
		result2 error
	}{result1, result2}
}

func (fake *FakeClient) ListResourceWithContext(arg1 context.Context, arg2 string, arg3 v1.ListOptions) (*unstructured.UnstructuredList, error) {
	fake.listResourceWithContextMutex.Lock()
	ret, specificReturn := fake.listResourceWithContextReturnsOnCall[len(fake.listResourceWithContextArgsForCall)]
	fake.listResourceWithContextArgsForCall = append(fake.listResourceWithContextArgsForCall, struct {
		arg1 context.Context
		arg2 string
		arg3 v1.ListOptions
	}{arg1, arg2, arg3})
	stub := fake.ListResourceWithContextStub
	fakeReturns := fake.listResourceWithContextReturns
	fake.recordInvocation("ListResourceWithContext", []interface{}{arg1, arg2, arg3})
	fake.listResourceWithContextMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeClient) ListResourceWithContextCallCount() int {
	fake.listResourceWithContextMutex.RLock()
	defer fake.listResourceWithContextMutex.RUnlock()
	return len(fake.listResourceWithContextArgsForCall)
}

func (fake *FakeClient) ListResourceWithContextCalls(stub func(context.Context, string, v1.ListOptions) (*unstructured.UnstructuredList, error)) {
	fake.listResourceWithContextMutex.Lock()
	defer fake.listResourceWithContextMutex.Unlock()
	fake.ListResourceWithContextStub = stub
}

func (fake *FakeClient) ListResourceWithContextArgsForCall(i int) (context.Context, string, v1.ListOptions) {
	fake.listResourceWithContextMutex.RLock()
	defer fake.listResourceWithContextMutex.RUnlock()
	argsForCall := fake.listResourceWithContextArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *FakeClient) ListResourceWithContextReturns(result1 *unstructured.UnstructuredList, result2 error) {
	fake.listResourceWithContextMutex.Lock()
	defer fake.listResourceWithContextMutex.Unlock()
	fake.ListResourceWithContextStub = nil
	fake.listResourceWithContextReturns = struct {
		result1 *unstructured.UnstructuredList
		result2 error
	}{result1, result2}
}

func (fake *FakeClient) ListResourceWithContextReturnsOnCall(i int, result1 *unstructured.UnstructuredList, result2 error) {
	fake.listResourceWithContextMutex.Lock()
	defer fake.listResourceWithContextMutex.Unlock()
	fake.ListResourceWithContextStub = nil
	if fake.listResourceWithContextReturnsOnCall == nil {
		fake.listResourceWithContextReturnsOnCall = make(map[int]struct {
			result1 *unstructured.UnstructuredList
			result2 error
		})
	}
	fake.listResourceWithContextReturnsOnCall[i] = struct {
		result1 *unstructured.UnstructuredList
		result2 error
	}{result1, result2}
}

func (fake *FakeClient) ListResourcesByKindAndNamespace(arg1 string, arg2 string, arg3 v1.ListOptions) (*unstructured.UnstructuredList, error) {
	fake.listResourcesByKindAndNamespaceMutex.Lock()
	ret, specificReturn := fake.listResourcesByKindAndNamespaceReturnsOnCall[len(fake.listResourcesByKindAndNamespaceArgsForCall)]
	fake.listResourcesByKindAndNamespaceArgsForCall = append(fake.listResourcesByKindAndNamespaceArgsForCall, struct {
		arg1 string
		arg2 string
		arg3 v1.ListOptions
	}{arg1, arg2, arg3})
	stub := fake.ListResourcesByKindAndNamespaceStub
	fakeReturns := fake.listResourcesByKindAndNamespaceReturns
	fake.recordInvocation("ListResourcesByKindAndNamespace", []interface{}{arg1, arg2, arg3})
	fake.listResourcesByKindAndNamespaceMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeClient) ListResourcesByKindAndNamespaceCallCount() int {
	fake.listResourcesByKindAndNamespaceMutex.RLock()
	defer fake.listResourcesByKindAndNamespaceMutex.RUnlock()
	return len(fake.listResourcesByKindAndNamespaceArgsForCall)
}

func (fake *FakeClient) ListResourcesByKindAndNamespaceCalls(stub func(string, string, v1.ListOptions) (*unstructured.UnstructuredList, error)) {
	fake.listResourcesByKindAndNamespaceMutex.Lock()
	defer fake.listResourcesByKindAndNamespaceMutex.Unlock()
	fake.ListResourcesByKindAndNamespaceStub = stub
}

func (fake *FakeClient) ListResourcesByKindAndNamespaceArgsForCall(i int) (string, string, v1.ListOptions) {
	fake.listResourcesByKindAndNamespaceMutex.RLock()
	defer fake.listResourcesByKindAndNamespaceMutex.RUnlock()
	argsForCall := fake.listResourcesByKindAndNamespaceArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *FakeClient) ListResourcesByKindAndNamespaceReturns(result1 *unstructured.UnstructuredList, result2 error) {
	fake.listResourcesByKindAndNamespaceMutex.Lock()
	defer fake.listResourcesByKindAndNamespaceMutex.Unlock()
	fake.ListResourcesByKindAndNamespaceStub = nil
	fake.listResourcesByKindAndNamespaceReturns = struct {
		result1 *unstructured.UnstructuredList
		result2 error
	}{result1, result2}
}

func (fake *FakeClient) ListResourcesByKindAndNamespaceReturnsOnCall(i int, result1 *unstructured.UnstructuredList, result2 error) {
	fake.listResourcesByKindAndNamespaceMutex.Lock()
	defer fake.listResourcesByKindAndNamespaceMutex.Unlock()
	fake.ListResourcesByKindAndNamespaceStub = nil
	if fake.listResourcesByKindAndNamespaceReturnsOnCall == nil {
		fake.listResourcesByKindAndNamespaceReturnsOnCall = make(map[int]struct {
			result1 *unstructured.UnstructuredList
			result2 error
		})
	}
	fake.listResourcesByKindAndNamespaceReturnsOnCall[i] = struct {
		result1 *unstructured.UnstructuredList
		result2 error
	}{result1, result2}
}

func (fake *FakeClient) Patch(arg1 string, arg2 string, arg3 string, arg4 []byte) (kubernetes.Metadata, *unstructured.Unstructured, error) {
	var arg4Copy []byte
	if arg4 != nil {
		arg4Copy = make([]byte, len(arg4))
		copy(arg4Copy, arg4)
	}
	fake.patchMutex.Lock()
	ret, specificReturn := fake.patchReturnsOnCall[len(fake.patchArgsForCall)]
	fake.patchArgsForCall = append(fake.patchArgsForCall, struct {
		arg1 string
		arg2 string
		arg3 string
		arg4 []byte
	}{arg1, arg2, arg3, arg4Copy})
	stub := fake.PatchStub
	fakeReturns := fake.patchReturns
	fake.recordInvocation("Patch", []interface{}{arg1, arg2, arg3, arg4Copy})
	fake.patchMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3, arg4)
	}
	if specificReturn {
		return ret.result1, ret.result2, ret.result3
	}
	return fakeReturns.result1, fakeReturns.result2, fakeReturns.result3
}

func (fake *FakeClient) PatchCallCount() int {
	fake.patchMutex.RLock()
	defer fake.patchMutex.RUnlock()
	return len(fake.patchArgsForCall)
}

func (fake *FakeClient) PatchCalls(stub func(string, string, string, []byte) (kubernetes.Metadata, *unstructured.Unstructured, error)) {
	fake.patchMutex.Lock()
	defer fake.patchMutex.Unlock()
	fake.PatchStub = stub
}

func (fake *FakeClient) PatchArgsForCall(i int) (string, string, string, []byte) {
	fake.patchMutex.RLock()
	defer fake.patchMutex.RUnlock()
	argsForCall := fake.patchArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3, argsForCall.arg4
}

func (fake *FakeClient) PatchReturns(result1 kubernetes.Metadata, result2 *unstructured.Unstructured, result3 error) {
	fake.patchMutex.Lock()
	defer fake.patchMutex.Unlock()
	fake.PatchStub = nil
	fake.patchReturns = struct {
		result1 kubernetes.Metadata
		result2 *unstructured.Unstructured
		result3 error
	}{result1, result2, result3}
}

func (fake *FakeClient) PatchReturnsOnCall(i int, result1 kubernetes.Metadata, result2 *unstructured.Unstructured, result3 error) {
	fake.patchMutex.Lock()
	defer fake.patchMutex.Unlock()
	fake.PatchStub = nil
	if fake.patchReturnsOnCall == nil {
		fake.patchReturnsOnCall = make(map[int]struct {
			result1 kubernetes.Metadata
			result2 *unstructured.Unstructured
			result3 error
		})
	}
	fake.patchReturnsOnCall[i] = struct {
		result1 kubernetes.Metadata
		result2 *unstructured.Unstructured
		result3 error
	}{result1, result2, result3}
}

func (fake *FakeClient) PatchUsingStrategy(arg1 string, arg2 string, arg3 string, arg4 []byte, arg5 types.PatchType) (kubernetes.Metadata, *unstructured.Unstructured, error) {
	var arg4Copy []byte
	if arg4 != nil {
		arg4Copy = make([]byte, len(arg4))
		copy(arg4Copy, arg4)
	}
	fake.patchUsingStrategyMutex.Lock()
	ret, specificReturn := fake.patchUsingStrategyReturnsOnCall[len(fake.patchUsingStrategyArgsForCall)]
	fake.patchUsingStrategyArgsForCall = append(fake.patchUsingStrategyArgsForCall, struct {
		arg1 string
		arg2 string
		arg3 string
		arg4 []byte
		arg5 types.PatchType
	}{arg1, arg2, arg3, arg4Copy, arg5})
	stub := fake.PatchUsingStrategyStub
	fakeReturns := fake.patchUsingStrategyReturns
	fake.recordInvocation("PatchUsingStrategy", []interface{}{arg1, arg2, arg3, arg4Copy, arg5})
	fake.patchUsingStrategyMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3, arg4, arg5)
	}
	if specificReturn {
		return ret.result1, ret.result2, ret.result3
	}
	return fakeReturns.result1, fakeReturns.result2, fakeReturns.result3
}

func (fake *FakeClient) PatchUsingStrategyCallCount() int {
	fake.patchUsingStrategyMutex.RLock()
	defer fake.patchUsingStrategyMutex.RUnlock()
	return len(fake.patchUsingStrategyArgsForCall)
}

func (fake *FakeClient) PatchUsingStrategyCalls(stub func(string, string, string, []byte, types.PatchType) (kubernetes.Metadata, *unstructured.Unstructured, error)) {
	fake.patchUsingStrategyMutex.Lock()
	defer fake.patchUsingStrategyMutex.Unlock()
	fake.PatchUsingStrategyStub = stub
}

func (fake *FakeClient) PatchUsingStrategyArgsForCall(i int) (string, string, string, []byte, types.PatchType) {
	fake.patchUsingStrategyMutex.RLock()
	defer fake.patchUsingStrategyMutex.RUnlock()
	argsForCall := fake.patchUsingStrategyArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3, argsForCall.arg4, argsForCall.arg5
}

func (fake *FakeClient) PatchUsingStrategyReturns(result1 kubernetes.Metadata, result2 *unstructured.Unstructured, result3 error) {
	fake.patchUsingStrategyMutex.Lock()
	defer fake.patchUsingStrategyMutex.Unlock()
	fake.PatchUsingStrategyStub = nil
	fake.patchUsingStrategyReturns = struct {
		result1 kubernetes.Metadata
		result2 *unstructured.Unstructured
		result3 error
	}{result1, result2, result3}
}

func (fake *FakeClient) PatchUsingStrategyReturnsOnCall(i int, result1 kubernetes.Metadata, result2 *unstructured.Unstructured, result3 error) {
	fake.patchUsingStrategyMutex.Lock()
	defer fake.patchUsingStrategyMutex.Unlock()
	fake.PatchUsingStrategyStub = nil
	if fake.patchUsingStrategyReturnsOnCall == nil {
		fake.patchUsingStrategyReturnsOnCall = make(map[int]struct {
			result1 kubernetes.Metadata
			result2 *unstructured.Unstructured
			result3 error
		})
	}
	fake.patchUsingStrategyReturnsOnCall[i] = struct {
		result1 kubernetes.Metadata
		result2 *unstructured.Unstructured
		result3 error
	}{result1, result2, result3}
}

func (fake *FakeClient) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.applyMutex.RLock()
	defer fake.applyMutex.RUnlock()
	fake.applyWithNamespaceOverrideMutex.RLock()
	defer fake.applyWithNamespaceOverrideMutex.RUnlock()
	fake.deleteResourceByKindAndNameAndNamespaceMutex.RLock()
	defer fake.deleteResourceByKindAndNameAndNamespaceMutex.RUnlock()
	fake.discoverMutex.RLock()
	defer fake.discoverMutex.RUnlock()
	fake.gVRForKindMutex.RLock()
	defer fake.gVRForKindMutex.RUnlock()
	fake.getMutex.RLock()
	defer fake.getMutex.RUnlock()
	fake.listByGVRMutex.RLock()
	defer fake.listByGVRMutex.RUnlock()
	fake.listByGVRWithContextMutex.RLock()
	defer fake.listByGVRWithContextMutex.RUnlock()
	fake.listResourceMutex.RLock()
	defer fake.listResourceMutex.RUnlock()
	fake.listResourceWithContextMutex.RLock()
	defer fake.listResourceWithContextMutex.RUnlock()
	fake.listResourcesByKindAndNamespaceMutex.RLock()
	defer fake.listResourcesByKindAndNamespaceMutex.RUnlock()
	fake.patchMutex.RLock()
	defer fake.patchMutex.RUnlock()
	fake.patchUsingStrategyMutex.RLock()
	defer fake.patchUsingStrategyMutex.RUnlock()
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

var _ kubernetes.Client = new(FakeClient)
