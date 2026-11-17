package patcher_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	. "github.com/homedepot/go-clouddriver/internal/kubernetes/patcher"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/cli-runtime/pkg/resource"
	"k8s.io/client-go/rest/fake"
	"k8s.io/kubectl/pkg/util"
)

const (
	testNamespace = "test-ns"
	testName      = "test-deployment"
)

type recordedRequest struct {
	method      string
	path        string
	query       url.Values
	contentType string
	body        []byte
}

var _ = Describe("Patcher", func() {
	var (
		p                *Patcher
		requests         []recordedRequest
		respond          func(req *http.Request) *http.Response
		current          *unstructured.Unstructured
		modified         []byte
		serverSideApply  bool
		patch            []byte
		patchedObject    runtime.Object
		err              error
		deploymentsPath  = fmt.Sprintf("/namespaces/%s/deployments", testNamespace)
		deploymentPath   = deploymentsPath + "/" + testName
		deploymentGV     = schema.GroupVersion{Group: "apps", Version: "v1"}
		deploymentGVR    = deploymentGV.WithResource("deployments")
		deploymentGVK    = deploymentGV.WithKind("Deployment")
		deploymentMapper = &meta.RESTMapping{
			Resource:         deploymentGVR,
			GroupVersionKind: deploymentGVK,
			Scope:            meta.RESTScopeNamespace,
		}
	)

	BeforeEach(func() {
		requests = nil
		serverSideApply = false

		client := &fake.RESTClient{
			NegotiatedSerializer: resource.UnstructuredPlusDefaultContentConfig().NegotiatedSerializer,
			GroupVersion:         deploymentGV,
			Client: fake.CreateHTTPClient(func(req *http.Request) (*http.Response, error) {
				var body []byte
				if req.Body != nil {
					body, _ = io.ReadAll(req.Body)
				}

				requests = append(requests, recordedRequest{
					method:      req.Method,
					path:        req.URL.Path,
					query:       req.URL.Query(),
					contentType: req.Header.Get("Content-Type"),
					body:        body,
				})

				return respond(req), nil
			}),
		}

		p, err = New(&resource.Info{Mapping: deploymentMapper}, resource.NewHelper(client, deploymentMapper))
		Expect(err).ToNot(HaveOccurred())

		current = newDeployment(1, "nginx:1.0")
		Expect(util.CreateApplyAnnotation(current, unstructured.UnstructuredJSONScheme)).To(Succeed())

		respond = func(*http.Request) *http.Response {
			return objectResponse(http.StatusOK, current)
		}
	})

	JustBeforeEach(func() {
		patch, patchedObject, err = p.Patch(current, modified, testNamespace, testName, serverSideApply)
	})

	Context("client-side apply", func() {
		When("the manifest changes a field", func() {
			BeforeEach(func() {
				modified = modifiedConfiguration(newDeployment(2, "nginx:1.1"))
			})

			It("sends one strategic merge patch with the change", func() {
				Expect(err).ToNot(HaveOccurred())
				Expect(patchedObject).ToNot(BeNil())
				Expect(requests).To(HaveLen(1))

				req := requests[0]
				Expect(req.method).To(Equal(http.MethodPatch))
				Expect(req.path).To(Equal(deploymentPath))
				Expect(req.contentType).To(Equal(string(types.StrategicMergePatchType)))
				Expect(req.body).To(Equal(patch))

				var sent map[string]interface{}
				Expect(json.Unmarshal(req.body, &sent)).To(Succeed())
				Expect(sent).To(HaveKeyWithValue("spec", HaveKeyWithValue("replicas", BeNumerically("==", 2))))
				Expect(string(req.body)).To(ContainSubstring(`"$setElementOrder/containers":[{"name":"app"}]`))
				Expect(string(req.body)).To(ContainSubstring(`"image":"nginx:1.1"`))
			})
		})

		When("the manifest changes lifecycle.preStop.sleep", func() {
			BeforeEach(func() {
				current = withPreStopSleep(newDeployment(1, "nginx:1.0"), 20)
				Expect(util.CreateApplyAnnotation(current, unstructured.UnstructuredJSONScheme)).To(Succeed())

				modified = modifiedConfiguration(withPreStopSleep(newDeployment(1, "nginx:1.0"), 30))
			})

			It("sends one strategic merge patch with the new sleep duration", func() {
				Expect(err).ToNot(HaveOccurred())
				Expect(requests).To(HaveLen(1))

				req := requests[0]
				Expect(req.method).To(Equal(http.MethodPatch))
				Expect(req.contentType).To(Equal(string(types.StrategicMergePatchType)))
				Expect(string(req.body)).To(ContainSubstring(`"lifecycle":{"preStop":{"sleep":{"seconds":30}}}`))
			})
		})

		When("the manifest is unchanged", func() {
			BeforeEach(func() {
				modified = modifiedConfiguration(newDeployment(1, "nginx:1.0"))
			})

			It("returns an empty patch without calling the API server", func() {
				Expect(err).ToNot(HaveOccurred())
				Expect(string(patch)).To(Equal("{}"))
				Expect(requests).To(BeEmpty())
			})
		})
	})

	Context("server-side apply", func() {
		BeforeEach(func() {
			serverSideApply = true
			modified = modifiedConfiguration(current)
		})

		When("force is not set", func() {
			It("sends an apply patch as the spinnaker field manager without forcing", func() {
				Expect(err).ToNot(HaveOccurred())
				Expect(requests).To(HaveLen(1))

				req := requests[0]
				Expect(req.method).To(Equal(http.MethodPatch))
				Expect(req.path).To(Equal(deploymentPath))
				Expect(req.contentType).To(Equal(string(types.ApplyPatchType)))
				Expect(req.query.Get("fieldManager")).To(Equal("spinnaker"))
				Expect(req.query.Get("force")).To(Equal("false"))

				var sent map[string]interface{}
				Expect(json.Unmarshal(req.body, &sent)).To(Succeed())
				Expect(sent).To(HaveKeyWithValue("kind", "Deployment"))
				Expect(sent).To(HaveKeyWithValue("spec", HaveKeyWithValue("replicas", BeNumerically("==", 1))))
			})
		})

		When("force is set", func() {
			BeforeEach(func() {
				p.Force = true
			})

			It("sends an apply patch that forces conflicts", func() {
				Expect(err).ToNot(HaveOccurred())
				Expect(requests).To(HaveLen(1))
				Expect(requests[0].contentType).To(Equal(string(types.ApplyPatchType)))
				Expect(requests[0].query.Get("fieldManager")).To(Equal("spinnaker"))
				Expect(requests[0].query.Get("force")).To(Equal("true"))
			})
		})
	})
})

func newDeployment(replicas int64, image string) *unstructured.Unstructured {
	return &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "apps/v1",
			"kind":       "Deployment",
			"metadata": map[string]interface{}{
				"name":      testName,
				"namespace": testNamespace,
			},
			"spec": map[string]interface{}{
				"replicas": replicas,
				"selector": map[string]interface{}{
					"matchLabels": map[string]interface{}{"app": "test"},
				},
				"template": map[string]interface{}{
					"metadata": map[string]interface{}{
						"labels": map[string]interface{}{"app": "test"},
					},
					"spec": map[string]interface{}{
						"containers": []interface{}{
							map[string]interface{}{
								"name":  "app",
								"image": image,
							},
						},
					},
				},
			},
		},
	}
}

func withPreStopSleep(u *unstructured.Unstructured, seconds int64) *unstructured.Unstructured {
	containers, found, err := unstructured.NestedSlice(u.Object, "spec", "template", "spec", "containers")
	Expect(err).ToNot(HaveOccurred())
	Expect(found).To(BeTrue())

	containers[0].(map[string]interface{})["lifecycle"] = map[string]interface{}{
		"preStop": map[string]interface{}{
			"sleep": map[string]interface{}{"seconds": seconds},
		},
	}
	Expect(unstructured.SetNestedSlice(u.Object, containers, "spec", "template", "spec", "containers")).To(Succeed())

	return u
}

func modifiedConfiguration(u *unstructured.Unstructured) []byte {
	modified, err := util.GetModifiedConfiguration(u, true, unstructured.UnstructuredJSONScheme)
	Expect(err).ToNot(HaveOccurred())

	return modified
}

func objectResponse(code int, u *unstructured.Unstructured) *http.Response {
	body, err := json.Marshal(u.Object)
	Expect(err).ToNot(HaveOccurred())

	return jsonResponse(code, body)
}

func jsonResponse(code int, body []byte) *http.Response {
	return &http.Response{
		StatusCode: code,
		Header:     http.Header{"Content-Type": []string{"application/json"}},
		Body:       io.NopCloser(bytes.NewReader(body)),
	}
}
