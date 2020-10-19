// Code generated by counterfeiter. DO NOT EDIT.
package kubernetesfakes

import (
	"sync"

	"github.com/billiford/go-clouddriver/pkg/kubernetes"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/rest"
)

type FakeController struct {
	AddSpinnakerAnnotationsStub        func(*unstructured.Unstructured, string) error
	addSpinnakerAnnotationsMutex       sync.RWMutex
	addSpinnakerAnnotationsArgsForCall []struct {
		arg1 *unstructured.Unstructured
		arg2 string
	}
	addSpinnakerAnnotationsReturns struct {
		result1 error
	}
	addSpinnakerAnnotationsReturnsOnCall map[int]struct {
		result1 error
	}
	AddSpinnakerLabelsStub        func(*unstructured.Unstructured, string) error
	addSpinnakerLabelsMutex       sync.RWMutex
	addSpinnakerLabelsArgsForCall []struct {
		arg1 *unstructured.Unstructured
		arg2 string
	}
	addSpinnakerLabelsReturns struct {
		result1 error
	}
	addSpinnakerLabelsReturnsOnCall map[int]struct {
		result1 error
	}
	AddSpinnakerVersionAnnotationsStub        func(*unstructured.Unstructured, string, kubernetes.SpinnakerVersion) error
	addSpinnakerVersionAnnotationsMutex       sync.RWMutex
	addSpinnakerVersionAnnotationsArgsForCall []struct {
		arg1 *unstructured.Unstructured
		arg2 string
		arg3 kubernetes.SpinnakerVersion
	}
	addSpinnakerVersionAnnotationsReturns struct {
		result1 error
	}
	addSpinnakerVersionAnnotationsReturnsOnCall map[int]struct {
		result1 error
	}
	AddSpinnakerVersionLabelsStub        func(*unstructured.Unstructured, string, kubernetes.SpinnakerVersion) error
	addSpinnakerVersionLabelsMutex       sync.RWMutex
	addSpinnakerVersionLabelsArgsForCall []struct {
		arg1 *unstructured.Unstructured
		arg2 string
		arg3 kubernetes.SpinnakerVersion
	}
	addSpinnakerVersionLabelsReturns struct {
		result1 error
	}
	addSpinnakerVersionLabelsReturnsOnCall map[int]struct {
		result1 error
	}
	GetCurrentVersionStub        func(*unstructured.UnstructuredList, string, string) string
	getCurrentVersionMutex       sync.RWMutex
	getCurrentVersionArgsForCall []struct {
		arg1 *unstructured.UnstructuredList
		arg2 string
		arg3 string
	}
	getCurrentVersionReturns struct {
		result1 string
	}
	getCurrentVersionReturnsOnCall map[int]struct {
		result1 string
	}
	IncrementVersionStub        func(string) kubernetes.SpinnakerVersion
	incrementVersionMutex       sync.RWMutex
	incrementVersionArgsForCall []struct {
		arg1 string
	}
	incrementVersionReturns struct {
		result1 kubernetes.SpinnakerVersion
	}
	incrementVersionReturnsOnCall map[int]struct {
		result1 kubernetes.SpinnakerVersion
	}
	IsVersionedStub        func(*unstructured.Unstructured) bool
	isVersionedMutex       sync.RWMutex
	isVersionedArgsForCall []struct {
		arg1 *unstructured.Unstructured
	}
	isVersionedReturns struct {
		result1 bool
	}
	isVersionedReturnsOnCall map[int]struct {
		result1 bool
	}
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
	ToUnstructuredStub        func(map[string]interface{}) (*unstructured.Unstructured, error)
	toUnstructuredMutex       sync.RWMutex
	toUnstructuredArgsForCall []struct {
		arg1 map[string]interface{}
	}
	toUnstructuredReturns struct {
		result1 *unstructured.Unstructured
		result2 error
	}
	toUnstructuredReturnsOnCall map[int]struct {
		result1 *unstructured.Unstructured
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeController) AddSpinnakerAnnotations(arg1 *unstructured.Unstructured, arg2 string) error {
	fake.addSpinnakerAnnotationsMutex.Lock()
	ret, specificReturn := fake.addSpinnakerAnnotationsReturnsOnCall[len(fake.addSpinnakerAnnotationsArgsForCall)]
	fake.addSpinnakerAnnotationsArgsForCall = append(fake.addSpinnakerAnnotationsArgsForCall, struct {
		arg1 *unstructured.Unstructured
		arg2 string
	}{arg1, arg2})
	fake.recordInvocation("AddSpinnakerAnnotations", []interface{}{arg1, arg2})
	fake.addSpinnakerAnnotationsMutex.Unlock()
	if fake.AddSpinnakerAnnotationsStub != nil {
		return fake.AddSpinnakerAnnotationsStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.addSpinnakerAnnotationsReturns
	return fakeReturns.result1
}

func (fake *FakeController) AddSpinnakerAnnotationsCallCount() int {
	fake.addSpinnakerAnnotationsMutex.RLock()
	defer fake.addSpinnakerAnnotationsMutex.RUnlock()
	return len(fake.addSpinnakerAnnotationsArgsForCall)
}

func (fake *FakeController) AddSpinnakerAnnotationsCalls(stub func(*unstructured.Unstructured, string) error) {
	fake.addSpinnakerAnnotationsMutex.Lock()
	defer fake.addSpinnakerAnnotationsMutex.Unlock()
	fake.AddSpinnakerAnnotationsStub = stub
}

func (fake *FakeController) AddSpinnakerAnnotationsArgsForCall(i int) (*unstructured.Unstructured, string) {
	fake.addSpinnakerAnnotationsMutex.RLock()
	defer fake.addSpinnakerAnnotationsMutex.RUnlock()
	argsForCall := fake.addSpinnakerAnnotationsArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeController) AddSpinnakerAnnotationsReturns(result1 error) {
	fake.addSpinnakerAnnotationsMutex.Lock()
	defer fake.addSpinnakerAnnotationsMutex.Unlock()
	fake.AddSpinnakerAnnotationsStub = nil
	fake.addSpinnakerAnnotationsReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeController) AddSpinnakerAnnotationsReturnsOnCall(i int, result1 error) {
	fake.addSpinnakerAnnotationsMutex.Lock()
	defer fake.addSpinnakerAnnotationsMutex.Unlock()
	fake.AddSpinnakerAnnotationsStub = nil
	if fake.addSpinnakerAnnotationsReturnsOnCall == nil {
		fake.addSpinnakerAnnotationsReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.addSpinnakerAnnotationsReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeController) AddSpinnakerLabels(arg1 *unstructured.Unstructured, arg2 string) error {
	fake.addSpinnakerLabelsMutex.Lock()
	ret, specificReturn := fake.addSpinnakerLabelsReturnsOnCall[len(fake.addSpinnakerLabelsArgsForCall)]
	fake.addSpinnakerLabelsArgsForCall = append(fake.addSpinnakerLabelsArgsForCall, struct {
		arg1 *unstructured.Unstructured
		arg2 string
	}{arg1, arg2})
	fake.recordInvocation("AddSpinnakerLabels", []interface{}{arg1, arg2})
	fake.addSpinnakerLabelsMutex.Unlock()
	if fake.AddSpinnakerLabelsStub != nil {
		return fake.AddSpinnakerLabelsStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.addSpinnakerLabelsReturns
	return fakeReturns.result1
}

func (fake *FakeController) AddSpinnakerLabelsCallCount() int {
	fake.addSpinnakerLabelsMutex.RLock()
	defer fake.addSpinnakerLabelsMutex.RUnlock()
	return len(fake.addSpinnakerLabelsArgsForCall)
}

func (fake *FakeController) AddSpinnakerLabelsCalls(stub func(*unstructured.Unstructured, string) error) {
	fake.addSpinnakerLabelsMutex.Lock()
	defer fake.addSpinnakerLabelsMutex.Unlock()
	fake.AddSpinnakerLabelsStub = stub
}

func (fake *FakeController) AddSpinnakerLabelsArgsForCall(i int) (*unstructured.Unstructured, string) {
	fake.addSpinnakerLabelsMutex.RLock()
	defer fake.addSpinnakerLabelsMutex.RUnlock()
	argsForCall := fake.addSpinnakerLabelsArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeController) AddSpinnakerLabelsReturns(result1 error) {
	fake.addSpinnakerLabelsMutex.Lock()
	defer fake.addSpinnakerLabelsMutex.Unlock()
	fake.AddSpinnakerLabelsStub = nil
	fake.addSpinnakerLabelsReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeController) AddSpinnakerLabelsReturnsOnCall(i int, result1 error) {
	fake.addSpinnakerLabelsMutex.Lock()
	defer fake.addSpinnakerLabelsMutex.Unlock()
	fake.AddSpinnakerLabelsStub = nil
	if fake.addSpinnakerLabelsReturnsOnCall == nil {
		fake.addSpinnakerLabelsReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.addSpinnakerLabelsReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeController) AddSpinnakerVersionAnnotations(arg1 *unstructured.Unstructured, arg2 string, arg3 kubernetes.SpinnakerVersion) error {
	fake.addSpinnakerVersionAnnotationsMutex.Lock()
	ret, specificReturn := fake.addSpinnakerVersionAnnotationsReturnsOnCall[len(fake.addSpinnakerVersionAnnotationsArgsForCall)]
	fake.addSpinnakerVersionAnnotationsArgsForCall = append(fake.addSpinnakerVersionAnnotationsArgsForCall, struct {
		arg1 *unstructured.Unstructured
		arg2 string
		arg3 kubernetes.SpinnakerVersion
	}{arg1, arg2, arg3})
	fake.recordInvocation("AddSpinnakerVersionAnnotations", []interface{}{arg1, arg2, arg3})
	fake.addSpinnakerVersionAnnotationsMutex.Unlock()
	if fake.AddSpinnakerVersionAnnotationsStub != nil {
		return fake.AddSpinnakerVersionAnnotationsStub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.addSpinnakerVersionAnnotationsReturns
	return fakeReturns.result1
}

func (fake *FakeController) AddSpinnakerVersionAnnotationsCallCount() int {
	fake.addSpinnakerVersionAnnotationsMutex.RLock()
	defer fake.addSpinnakerVersionAnnotationsMutex.RUnlock()
	return len(fake.addSpinnakerVersionAnnotationsArgsForCall)
}

func (fake *FakeController) AddSpinnakerVersionAnnotationsCalls(stub func(*unstructured.Unstructured, string, kubernetes.SpinnakerVersion) error) {
	fake.addSpinnakerVersionAnnotationsMutex.Lock()
	defer fake.addSpinnakerVersionAnnotationsMutex.Unlock()
	fake.AddSpinnakerVersionAnnotationsStub = stub
}

func (fake *FakeController) AddSpinnakerVersionAnnotationsArgsForCall(i int) (*unstructured.Unstructured, string, kubernetes.SpinnakerVersion) {
	fake.addSpinnakerVersionAnnotationsMutex.RLock()
	defer fake.addSpinnakerVersionAnnotationsMutex.RUnlock()
	argsForCall := fake.addSpinnakerVersionAnnotationsArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *FakeController) AddSpinnakerVersionAnnotationsReturns(result1 error) {
	fake.addSpinnakerVersionAnnotationsMutex.Lock()
	defer fake.addSpinnakerVersionAnnotationsMutex.Unlock()
	fake.AddSpinnakerVersionAnnotationsStub = nil
	fake.addSpinnakerVersionAnnotationsReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeController) AddSpinnakerVersionAnnotationsReturnsOnCall(i int, result1 error) {
	fake.addSpinnakerVersionAnnotationsMutex.Lock()
	defer fake.addSpinnakerVersionAnnotationsMutex.Unlock()
	fake.AddSpinnakerVersionAnnotationsStub = nil
	if fake.addSpinnakerVersionAnnotationsReturnsOnCall == nil {
		fake.addSpinnakerVersionAnnotationsReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.addSpinnakerVersionAnnotationsReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeController) AddSpinnakerVersionLabels(arg1 *unstructured.Unstructured, arg2 string, arg3 kubernetes.SpinnakerVersion) error {
	fake.addSpinnakerVersionLabelsMutex.Lock()
	ret, specificReturn := fake.addSpinnakerVersionLabelsReturnsOnCall[len(fake.addSpinnakerVersionLabelsArgsForCall)]
	fake.addSpinnakerVersionLabelsArgsForCall = append(fake.addSpinnakerVersionLabelsArgsForCall, struct {
		arg1 *unstructured.Unstructured
		arg2 string
		arg3 kubernetes.SpinnakerVersion
	}{arg1, arg2, arg3})
	fake.recordInvocation("AddSpinnakerVersionLabels", []interface{}{arg1, arg2, arg3})
	fake.addSpinnakerVersionLabelsMutex.Unlock()
	if fake.AddSpinnakerVersionLabelsStub != nil {
		return fake.AddSpinnakerVersionLabelsStub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.addSpinnakerVersionLabelsReturns
	return fakeReturns.result1
}

func (fake *FakeController) AddSpinnakerVersionLabelsCallCount() int {
	fake.addSpinnakerVersionLabelsMutex.RLock()
	defer fake.addSpinnakerVersionLabelsMutex.RUnlock()
	return len(fake.addSpinnakerVersionLabelsArgsForCall)
}

func (fake *FakeController) AddSpinnakerVersionLabelsCalls(stub func(*unstructured.Unstructured, string, kubernetes.SpinnakerVersion) error) {
	fake.addSpinnakerVersionLabelsMutex.Lock()
	defer fake.addSpinnakerVersionLabelsMutex.Unlock()
	fake.AddSpinnakerVersionLabelsStub = stub
}

func (fake *FakeController) AddSpinnakerVersionLabelsArgsForCall(i int) (*unstructured.Unstructured, string, kubernetes.SpinnakerVersion) {
	fake.addSpinnakerVersionLabelsMutex.RLock()
	defer fake.addSpinnakerVersionLabelsMutex.RUnlock()
	argsForCall := fake.addSpinnakerVersionLabelsArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *FakeController) AddSpinnakerVersionLabelsReturns(result1 error) {
	fake.addSpinnakerVersionLabelsMutex.Lock()
	defer fake.addSpinnakerVersionLabelsMutex.Unlock()
	fake.AddSpinnakerVersionLabelsStub = nil
	fake.addSpinnakerVersionLabelsReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeController) AddSpinnakerVersionLabelsReturnsOnCall(i int, result1 error) {
	fake.addSpinnakerVersionLabelsMutex.Lock()
	defer fake.addSpinnakerVersionLabelsMutex.Unlock()
	fake.AddSpinnakerVersionLabelsStub = nil
	if fake.addSpinnakerVersionLabelsReturnsOnCall == nil {
		fake.addSpinnakerVersionLabelsReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.addSpinnakerVersionLabelsReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeController) GetCurrentVersion(arg1 *unstructured.UnstructuredList, arg2 string, arg3 string) string {
	fake.getCurrentVersionMutex.Lock()
	ret, specificReturn := fake.getCurrentVersionReturnsOnCall[len(fake.getCurrentVersionArgsForCall)]
	fake.getCurrentVersionArgsForCall = append(fake.getCurrentVersionArgsForCall, struct {
		arg1 *unstructured.UnstructuredList
		arg2 string
		arg3 string
	}{arg1, arg2, arg3})
	fake.recordInvocation("GetCurrentVersion", []interface{}{arg1, arg2, arg3})
	fake.getCurrentVersionMutex.Unlock()
	if fake.GetCurrentVersionStub != nil {
		return fake.GetCurrentVersionStub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.getCurrentVersionReturns
	return fakeReturns.result1
}

func (fake *FakeController) GetCurrentVersionCallCount() int {
	fake.getCurrentVersionMutex.RLock()
	defer fake.getCurrentVersionMutex.RUnlock()
	return len(fake.getCurrentVersionArgsForCall)
}

func (fake *FakeController) GetCurrentVersionCalls(stub func(*unstructured.UnstructuredList, string, string) string) {
	fake.getCurrentVersionMutex.Lock()
	defer fake.getCurrentVersionMutex.Unlock()
	fake.GetCurrentVersionStub = stub
}

func (fake *FakeController) GetCurrentVersionArgsForCall(i int) (*unstructured.UnstructuredList, string, string) {
	fake.getCurrentVersionMutex.RLock()
	defer fake.getCurrentVersionMutex.RUnlock()
	argsForCall := fake.getCurrentVersionArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *FakeController) GetCurrentVersionReturns(result1 string) {
	fake.getCurrentVersionMutex.Lock()
	defer fake.getCurrentVersionMutex.Unlock()
	fake.GetCurrentVersionStub = nil
	fake.getCurrentVersionReturns = struct {
		result1 string
	}{result1}
}

func (fake *FakeController) GetCurrentVersionReturnsOnCall(i int, result1 string) {
	fake.getCurrentVersionMutex.Lock()
	defer fake.getCurrentVersionMutex.Unlock()
	fake.GetCurrentVersionStub = nil
	if fake.getCurrentVersionReturnsOnCall == nil {
		fake.getCurrentVersionReturnsOnCall = make(map[int]struct {
			result1 string
		})
	}
	fake.getCurrentVersionReturnsOnCall[i] = struct {
		result1 string
	}{result1}
}

func (fake *FakeController) IncrementVersion(arg1 string) kubernetes.SpinnakerVersion {
	fake.incrementVersionMutex.Lock()
	ret, specificReturn := fake.incrementVersionReturnsOnCall[len(fake.incrementVersionArgsForCall)]
	fake.incrementVersionArgsForCall = append(fake.incrementVersionArgsForCall, struct {
		arg1 string
	}{arg1})
	fake.recordInvocation("IncrementVersion", []interface{}{arg1})
	fake.incrementVersionMutex.Unlock()
	if fake.IncrementVersionStub != nil {
		return fake.IncrementVersionStub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.incrementVersionReturns
	return fakeReturns.result1
}

func (fake *FakeController) IncrementVersionCallCount() int {
	fake.incrementVersionMutex.RLock()
	defer fake.incrementVersionMutex.RUnlock()
	return len(fake.incrementVersionArgsForCall)
}

func (fake *FakeController) IncrementVersionCalls(stub func(string) kubernetes.SpinnakerVersion) {
	fake.incrementVersionMutex.Lock()
	defer fake.incrementVersionMutex.Unlock()
	fake.IncrementVersionStub = stub
}

func (fake *FakeController) IncrementVersionArgsForCall(i int) string {
	fake.incrementVersionMutex.RLock()
	defer fake.incrementVersionMutex.RUnlock()
	argsForCall := fake.incrementVersionArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeController) IncrementVersionReturns(result1 kubernetes.SpinnakerVersion) {
	fake.incrementVersionMutex.Lock()
	defer fake.incrementVersionMutex.Unlock()
	fake.IncrementVersionStub = nil
	fake.incrementVersionReturns = struct {
		result1 kubernetes.SpinnakerVersion
	}{result1}
}

func (fake *FakeController) IncrementVersionReturnsOnCall(i int, result1 kubernetes.SpinnakerVersion) {
	fake.incrementVersionMutex.Lock()
	defer fake.incrementVersionMutex.Unlock()
	fake.IncrementVersionStub = nil
	if fake.incrementVersionReturnsOnCall == nil {
		fake.incrementVersionReturnsOnCall = make(map[int]struct {
			result1 kubernetes.SpinnakerVersion
		})
	}
	fake.incrementVersionReturnsOnCall[i] = struct {
		result1 kubernetes.SpinnakerVersion
	}{result1}
}

func (fake *FakeController) IsVersioned(arg1 *unstructured.Unstructured) bool {
	fake.isVersionedMutex.Lock()
	ret, specificReturn := fake.isVersionedReturnsOnCall[len(fake.isVersionedArgsForCall)]
	fake.isVersionedArgsForCall = append(fake.isVersionedArgsForCall, struct {
		arg1 *unstructured.Unstructured
	}{arg1})
	fake.recordInvocation("IsVersioned", []interface{}{arg1})
	fake.isVersionedMutex.Unlock()
	if fake.IsVersionedStub != nil {
		return fake.IsVersionedStub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.isVersionedReturns
	return fakeReturns.result1
}

func (fake *FakeController) IsVersionedCallCount() int {
	fake.isVersionedMutex.RLock()
	defer fake.isVersionedMutex.RUnlock()
	return len(fake.isVersionedArgsForCall)
}

func (fake *FakeController) IsVersionedCalls(stub func(*unstructured.Unstructured) bool) {
	fake.isVersionedMutex.Lock()
	defer fake.isVersionedMutex.Unlock()
	fake.IsVersionedStub = stub
}

func (fake *FakeController) IsVersionedArgsForCall(i int) *unstructured.Unstructured {
	fake.isVersionedMutex.RLock()
	defer fake.isVersionedMutex.RUnlock()
	argsForCall := fake.isVersionedArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeController) IsVersionedReturns(result1 bool) {
	fake.isVersionedMutex.Lock()
	defer fake.isVersionedMutex.Unlock()
	fake.IsVersionedStub = nil
	fake.isVersionedReturns = struct {
		result1 bool
	}{result1}
}

func (fake *FakeController) IsVersionedReturnsOnCall(i int, result1 bool) {
	fake.isVersionedMutex.Lock()
	defer fake.isVersionedMutex.Unlock()
	fake.IsVersionedStub = nil
	if fake.isVersionedReturnsOnCall == nil {
		fake.isVersionedReturnsOnCall = make(map[int]struct {
			result1 bool
		})
	}
	fake.isVersionedReturnsOnCall[i] = struct {
		result1 bool
	}{result1}
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
	fakeReturns := fake.newClientReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeController) NewClientCallCount() int {
	fake.newClientMutex.RLock()
	defer fake.newClientMutex.RUnlock()
	return len(fake.newClientArgsForCall)
}

func (fake *FakeController) NewClientCalls(stub func(*rest.Config) (kubernetes.Client, error)) {
	fake.newClientMutex.Lock()
	defer fake.newClientMutex.Unlock()
	fake.NewClientStub = stub
}

func (fake *FakeController) NewClientArgsForCall(i int) *rest.Config {
	fake.newClientMutex.RLock()
	defer fake.newClientMutex.RUnlock()
	argsForCall := fake.newClientArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeController) NewClientReturns(result1 kubernetes.Client, result2 error) {
	fake.newClientMutex.Lock()
	defer fake.newClientMutex.Unlock()
	fake.NewClientStub = nil
	fake.newClientReturns = struct {
		result1 kubernetes.Client
		result2 error
	}{result1, result2}
}

func (fake *FakeController) NewClientReturnsOnCall(i int, result1 kubernetes.Client, result2 error) {
	fake.newClientMutex.Lock()
	defer fake.newClientMutex.Unlock()
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

func (fake *FakeController) ToUnstructured(arg1 map[string]interface{}) (*unstructured.Unstructured, error) {
	fake.toUnstructuredMutex.Lock()
	ret, specificReturn := fake.toUnstructuredReturnsOnCall[len(fake.toUnstructuredArgsForCall)]
	fake.toUnstructuredArgsForCall = append(fake.toUnstructuredArgsForCall, struct {
		arg1 map[string]interface{}
	}{arg1})
	fake.recordInvocation("ToUnstructured", []interface{}{arg1})
	fake.toUnstructuredMutex.Unlock()
	if fake.ToUnstructuredStub != nil {
		return fake.ToUnstructuredStub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.toUnstructuredReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeController) ToUnstructuredCallCount() int {
	fake.toUnstructuredMutex.RLock()
	defer fake.toUnstructuredMutex.RUnlock()
	return len(fake.toUnstructuredArgsForCall)
}

func (fake *FakeController) ToUnstructuredCalls(stub func(map[string]interface{}) (*unstructured.Unstructured, error)) {
	fake.toUnstructuredMutex.Lock()
	defer fake.toUnstructuredMutex.Unlock()
	fake.ToUnstructuredStub = stub
}

func (fake *FakeController) ToUnstructuredArgsForCall(i int) map[string]interface{} {
	fake.toUnstructuredMutex.RLock()
	defer fake.toUnstructuredMutex.RUnlock()
	argsForCall := fake.toUnstructuredArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeController) ToUnstructuredReturns(result1 *unstructured.Unstructured, result2 error) {
	fake.toUnstructuredMutex.Lock()
	defer fake.toUnstructuredMutex.Unlock()
	fake.ToUnstructuredStub = nil
	fake.toUnstructuredReturns = struct {
		result1 *unstructured.Unstructured
		result2 error
	}{result1, result2}
}

func (fake *FakeController) ToUnstructuredReturnsOnCall(i int, result1 *unstructured.Unstructured, result2 error) {
	fake.toUnstructuredMutex.Lock()
	defer fake.toUnstructuredMutex.Unlock()
	fake.ToUnstructuredStub = nil
	if fake.toUnstructuredReturnsOnCall == nil {
		fake.toUnstructuredReturnsOnCall = make(map[int]struct {
			result1 *unstructured.Unstructured
			result2 error
		})
	}
	fake.toUnstructuredReturnsOnCall[i] = struct {
		result1 *unstructured.Unstructured
		result2 error
	}{result1, result2}
}

func (fake *FakeController) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.addSpinnakerAnnotationsMutex.RLock()
	defer fake.addSpinnakerAnnotationsMutex.RUnlock()
	fake.addSpinnakerLabelsMutex.RLock()
	defer fake.addSpinnakerLabelsMutex.RUnlock()
	fake.addSpinnakerVersionAnnotationsMutex.RLock()
	defer fake.addSpinnakerVersionAnnotationsMutex.RUnlock()
	fake.addSpinnakerVersionLabelsMutex.RLock()
	defer fake.addSpinnakerVersionLabelsMutex.RUnlock()
	fake.getCurrentVersionMutex.RLock()
	defer fake.getCurrentVersionMutex.RUnlock()
	fake.incrementVersionMutex.RLock()
	defer fake.incrementVersionMutex.RUnlock()
	fake.isVersionedMutex.RLock()
	defer fake.isVersionedMutex.RUnlock()
	fake.newClientMutex.RLock()
	defer fake.newClientMutex.RUnlock()
	fake.toUnstructuredMutex.RLock()
	defer fake.toUnstructuredMutex.RUnlock()
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
