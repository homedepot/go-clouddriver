package core_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/homedepot/go-clouddriver/internal/kubernetes"
	clouddriver "github.com/homedepot/go-clouddriver/pkg"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var _ = Describe("Task", func() {
	Describe("#GetTask", func() {
		var logBuf *bytes.Buffer

		BeforeEach(func() {
			setup()
			logBuf = &bytes.Buffer{}
			log.SetOutput(logBuf)
			uri = svr.URL + "/task/task-id"
			createRequest(http.MethodGet)
		})

		AfterEach(func() {
			log.SetOutput(os.Stderr)
			teardown()
		})

		JustBeforeEach(func() {
			doRequest()
		})

		When("listing the resources returns an error", func() {
			BeforeEach(func() {
				fakeSQLClient.ListKubernetesResourcesByTaskIDReturns(nil, errors.New("error listing resources"))
			})

			It("returns a failed task", func() {
				Expect(res.StatusCode).To(Equal(http.StatusInternalServerError))
				t := clouddriver.Task{}
				b, _ := io.ReadAll(res.Body)
				err := json.Unmarshal(b, &t)
				Expect(err).To(BeNil())
				Expect(t.Status.Failed).To(BeTrue())
				Expect(t.Status.Retryable).To(BeTrue())
				Expect(t.Status.Status).To(Equal("Error listing resources for task (id: task-id): error listing resources"))
				Expect(logBuf.String()).To(ContainSubstring("GetTask: error listing resources for task (id: task-id): error listing resources"))
			})
		})

		When("no resources are returned", func() {
			BeforeEach(func() {
				fakeSQLClient.ListKubernetesResourcesByTaskIDReturns([]kubernetes.Resource{}, nil)
			})

			It("returns not found error", func() {
				Expect(res.StatusCode).To(Equal(http.StatusNotFound))
				Expect(logBuf.String()).To(BeEmpty())
			})
		})

		When("getting the provider returns an error", func() {
			BeforeEach(func() {
				fakeSQLClient.GetKubernetesProviderReturns(kubernetes.Provider{}, errors.New("error getting provider"))
			})

			It("returns a failed task", func() {
				Expect(res.StatusCode).To(Equal(http.StatusInternalServerError))
				t := clouddriver.Task{}
				b, _ := io.ReadAll(res.Body)
				err := json.Unmarshal(b, &t)
				Expect(err).To(BeNil())
				Expect(t.Status.Failed).To(BeTrue())
				Expect(t.Status.Retryable).To(BeTrue())
				Expect(t.Status.Status).To(Equal("Error getting kubernetes provider test-account-name for task (id: task-id): internal: error getting kubernetes provider test-account-name: error getting provider"))
				Expect(logBuf.String()).To(ContainSubstring("GetTask: error getting kubernetes provider for task (id: task-id, account: test-account-name)"))
			})
		})

		When("the task type is cleanup", func() {
			BeforeEach(func() {
				fakeSQLClient.ListKubernetesResourcesByTaskIDReturns([]kubernetes.Resource{
					{
						AccountName: "test-account-name",
						TaskType:    clouddriver.TaskTypeCleanup,
					},
				}, nil)
			})

			It("does not call make calls to the server", func() {
				Expect(res.StatusCode).To(Equal(http.StatusOK))
				Expect(fakeKubeClient.GetCallCount()).To(Equal(0))
			})
		})

		When("the task type is noop", func() {
			BeforeEach(func() {
				fakeSQLClient.ListKubernetesResourcesByTaskIDReturns([]kubernetes.Resource{
					{
						AccountName: "test-account-name",
						TaskType:    clouddriver.TaskTypeNoOp,
					},
				}, nil)
			})

			It("does not call make calls to the server", func() {
				Expect(res.StatusCode).To(Equal(http.StatusOK))
				Expect(fakeKubeClient.GetCallCount()).To(Equal(0))
			})
		})

		Context("when the task type is delete", func() {
			BeforeEach(func() {
				fakeSQLClient.ListKubernetesResourcesByTaskIDReturns([]kubernetes.Resource{
					{
						AccountName: "test-account-name",
						// ID:           "",
						Timestamp: time.Time{},
						// TaskID:       "",
						TaskType:     clouddriver.TaskTypeDelete,
						APIGroup:     "apps",
						Name:         "test-deployment",
						ArtifactName: "",
						Namespace:    "test-namespace",
						// Resource:     "",
						Version:      "v1",
						Kind:         "Deployment",
						SpinnakerApp: "test-app",
						Cluster:      "test-cluster",
					},
				}, nil)
			})

			When("the server returns a not found error", func() {
				BeforeEach(func() {
					fakeKubeClient.GetReturns(nil, errors.New(`horizontalpodautoscalers.autoscaling "php-apache1-v008" not found`))
				})

				It("ignores the not found error and returns a complete task", func() {
					Expect(res.StatusCode).To(Equal(http.StatusOK))
					Expect(fakeKubeClient.GetCallCount()).To(Equal(1))
					validateResponse(payloadTaskComplete)
				})
			})

			When("the server returns a generic error", func() {
				BeforeEach(func() {
					fakeKubeClient.GetReturns(nil, errors.New(`generic error`))
				})

				It("ignores the not found error", func() {
					Expect(res.StatusCode).To(Equal(http.StatusInternalServerError))
					Expect(fakeKubeClient.GetCallCount()).To(Equal(1))
				})
			})

			When("the server returns the resource", func() {
				BeforeEach(func() {
					fakeKubeClient.GetReturns(&unstructured.Unstructured{
						Object: map[string]interface{}{
							"kind":       "Deployment",
							"apiVersion": "apps/v1",
							"metadata": map[string]interface{}{
								"name":      "test-deployment",
								"namespace": "test-namespace",
							},
						},
					}, nil)
				})

				It("returns an incomplete task", func() {
					Expect(res.StatusCode).To(Equal(http.StatusOK))
					Expect(fakeKubeClient.GetCallCount()).To(Equal(1))
					validateResponse(payloadTaskIncomplete)
				})
			})
		})

		When("getting the manifest returns an error", func() {
			BeforeEach(func() {
				fakeKubeClient.GetReturns(nil, errors.New("error getting resource"))
			})

			It("returns a failed task", func() {
				Expect(res.StatusCode).To(Equal(http.StatusInternalServerError))
				t := clouddriver.Task{}
				b, _ := io.ReadAll(res.Body)
				err := json.Unmarshal(b, &t)
				Expect(err).To(BeNil())
				Expect(t.Status.Failed).To(BeTrue())
				Expect(t.Status.Retryable).To(BeTrue())
				Expect(t.Status.Status).To(Equal("Error getting resource for task (task ID: task-id, kind: test-kind, name: test-name, namespace: test-namespace): error getting resource"))
				Expect(logBuf.String()).To(ContainSubstring("GetTask: error getting resource for task (id: task-id, account: test-account-name, task type: , group: , version: , resource: test-kind, name: test-name, namespace: test-namespace): error getting resource"))
			})
		})

		Context("when a resource is gone at read-back after a successful create", func() {
			var (
				rows    []kubernetes.Resource
				missing []string
				group   string
				gvkErr  error
			)

			mapping := func(name, namespace string) kubernetes.Resource {
				return kubernetes.Resource{
					AccountName:  "test-account-name",
					APIGroup:     "getambassador.io",
					Version:      "v2",
					Kind:         "Mapping",
					Resource:     "mappings",
					Name:         name,
					Namespace:    namespace,
					ArtifactName: name,
				}
			}

			resultObject := func() clouddriver.TaskResultObject {
				t := clouddriver.Task{}
				b, _ := io.ReadAll(res.Body)
				Expect(json.Unmarshal(b, &t)).To(Succeed())
				Expect(t.ResultObjects).To(HaveLen(1))

				return t.ResultObjects[0]
			}

			BeforeEach(func() {
				rows = nil
				missing = nil
				group = ""
				gvkErr = nil
			})

			// configure sets up the fakes, then repeats the request: the outer
			// JustBeforeEach already sent one with the default fakes.
			configure := func() {
				fakeSQLClient.ListKubernetesResourcesByTaskIDReturns(rows, nil)
				fakeKubeClient.GVKForKindReturns(schema.GroupVersionKind{Group: group}, gvkErr)
				fakeKubeClient.GetStub = func(kind, name, namespace string) (*unstructured.Unstructured, error) {
					for _, m := range missing {
						if m == name {
							return nil, apierrors.NewNotFound(schema.GroupResource{Group: group, Resource: kind}, name)
						}
					}

					return &unstructured.Unstructured{Object: map[string]interface{}{
						"kind":     "Any",
						"metadata": map[string]interface{}{"name": name, "namespace": namespace},
					}}, nil
				}
				logBuf.Reset()
				createRequest(http.MethodGet)
				doRequest()
			}

			When("a paired canary Mapping vanished", func() {
				It("tolerates it, leaves it out of manifests and names, and logs it", func() {
					rows = []kubernetes.Resource{mapping("app", "ns"), mapping("app-canary", "ns")}
					missing = []string{"app-canary"}
					group = "getambassador.io"
					configure()

					Expect(res.StatusCode).To(Equal(http.StatusOK))
					ro := resultObject()
					Expect(ro.Manifests).To(HaveLen(1))
					Expect(ro.Manifests[0]["metadata"].(map[string]interface{})["name"]).To(Equal("app"))
					Expect(ro.ManifestNamesByNamespace).To(Equal(map[string][]string{"ns": {"Mapping app"}}))
					Expect(ro.DeployedNamesByLocation).To(Equal(ro.ManifestNamesByNamespace))
					Expect(ro.ManifestNamesByNamespaceToRefresh).To(Equal(ro.ManifestNamesByNamespace))
					// The artifact lists are unchanged: both rows are still listed.
					Expect(ro.CreatedArtifacts).To(HaveLen(2))
					Expect(ro.BoundArtifacts).To(HaveLen(2))
					Expect(logBuf.String()).To(ContainSubstring("GetTask: tolerated vanished resource (id: task-id, account: test-account-name, group: getambassador.io, resource: mappings, name: app-canary, namespace: ns)"))
					Expect(logBuf.String()).ToNot(ContainSubstring("GetTask: skipped canary Mapping"))
					Expect(logBuf.String()).ToNot(ContainSubstring("GetTask: error getting resource"))
				})
			})

			When("a paired canary Mapping still exists", func() {
				It("keeps it in manifests but not in the names, and logs the skip", func() {
					rows = []kubernetes.Resource{mapping("app", "ns"), mapping("app-canary", "ns")}
					group = "getambassador.io"
					configure()

					Expect(res.StatusCode).To(Equal(http.StatusOK))
					ro := resultObject()
					Expect(ro.Manifests).To(HaveLen(2))
					Expect(ro.ManifestNamesByNamespace).To(Equal(map[string][]string{"ns": {"Mapping app"}}))
					Expect(ro.CreatedArtifacts).To(HaveLen(2))
					Expect(logBuf.String()).To(ContainSubstring("GetTask: skipped canary Mapping (id: task-id, account: test-account-name, group: getambassador.io, resource: mappings, name: app-canary, namespace: ns)"))
					Expect(logBuf.String()).ToNot(ContainSubstring("tolerated vanished"))
				})
			})

			When("a canary Mapping has no base Mapping in the task", func() {
				BeforeEach(func() {
					rows = []kubernetes.Resource{mapping("eligible-canary", "ns")}
					group = "getambassador.io"
				})

				It("keeps it in the names when it exists", func() {
					configure()

					Expect(res.StatusCode).To(Equal(http.StatusOK))
					ro := resultObject()
					Expect(ro.ManifestNamesByNamespace).To(Equal(map[string][]string{"ns": {"Mapping eligible-canary"}}))
					Expect(logBuf.String()).To(BeEmpty())
				})

				It("keeps the 500 when it is gone", func() {
					missing = []string{"eligible-canary"}
					configure()

					Expect(res.StatusCode).To(Equal(http.StatusInternalServerError))
					Expect(logBuf.String()).To(ContainSubstring("GetTask: error getting resource"))
					Expect(logBuf.String()).ToNot(ContainSubstring("tolerated vanished"))
				})
			})

			When("the base Mapping is in another namespace", func() {
				It("does not pair them and keeps the 500", func() {
					rows = []kubernetes.Resource{mapping("app", "other-ns"), mapping("app-canary", "ns")}
					missing = []string{"app-canary"}
					group = "getambassador.io"
					configure()

					Expect(res.StatusCode).To(Equal(http.StatusInternalServerError))
				})
			})

			When("a Job vanished (removed by its TTL)", func() {
				It("tolerates it and returns empty manifests and names", func() {
					rows = []kubernetes.Resource{{
						AccountName: "test-account-name",
						APIGroup:    "batch",
						Version:     "v1",
						Kind:        "Job",
						Resource:    "jobs",
						Name:        "test-job",
						Namespace:   "ns",
					}}
					missing = []string{"test-job"}
					group = "batch"
					configure()

					Expect(res.StatusCode).To(Equal(http.StatusOK))
					ro := resultObject()
					Expect(ro.Manifests).To(BeEmpty())
					Expect(ro.ManifestNamesByNamespace).To(BeEmpty())
					Expect(ro.CreatedArtifacts).To(HaveLen(1))
					Expect(logBuf.String()).To(ContainSubstring("GetTask: tolerated vanished resource (id: task-id, account: test-account-name, group: batch, resource: jobs, name: test-job, namespace: ns)"))
				})
			})

			When("a kind outside the allowlist vanished", func() {
				DescribeTable("keeps the 500 and does not log a tolerated resource",
					func(apiGroup, kind, resource string) {
						rows = []kubernetes.Resource{{
							AccountName: "test-account-name",
							APIGroup:    apiGroup,
							Version:     "v1",
							Kind:        kind,
							Resource:    resource,
							Name:        "gone",
							Namespace:   "ns",
						}}
						missing = []string{"gone"}
						group = apiGroup
						configure()

						Expect(res.StatusCode).To(Equal(http.StatusInternalServerError))
						Expect(logBuf.String()).To(ContainSubstring("GetTask: error getting resource"))
						Expect(logBuf.String()).ToNot(ContainSubstring("tolerated vanished"))
					},
					Entry("Deployment", "apps", "Deployment", "deployments"),
					Entry("StatefulSet", "apps", "StatefulSet", "statefulsets"),
					Entry("Rollout", "argoproj.io", "Rollout", "rollouts"),
					Entry("TinyHome custom kind", "tinyhome.example.com", "TinyHomeDeployment", "tinyhomedeployments"),
					Entry("Mapping that is not a canary", "getambassador.io", "Mapping", "mappings"),
					Entry("Job in the wrong group", "example.com", "Job", "jobs"),
				)
			})

			When("the group the lookup resolves to differs from the stored group", func() {
				It("keeps the 500 (the lookup may have picked the wrong group)", func() {
					rows = []kubernetes.Resource{{
						AccountName: "test-account-name",
						APIGroup:    "batch",
						Version:     "v1",
						Kind:        "Job",
						Resource:    "jobs",
						Name:        "test-job",
						Namespace:   "ns",
					}}
					missing = []string{"test-job"}
					group = "other.example.com"
					configure()

					Expect(res.StatusCode).To(Equal(http.StatusInternalServerError))
					Expect(logBuf.String()).ToNot(ContainSubstring("tolerated vanished"))
				})
			})

			When("resolving the group returns an error", func() {
				It("keeps the 500", func() {
					rows = []kubernetes.Resource{mapping("app", "ns"), mapping("app-canary", "ns")}
					missing = []string{"app-canary"}
					group = "getambassador.io"
					gvkErr = errors.New("discovery error")
					configure()

					Expect(res.StatusCode).To(Equal(http.StatusInternalServerError))
				})
			})

			When("the error is not a not found error", func() {
				It("keeps the 500", func() {
					rows = []kubernetes.Resource{mapping("app", "ns"), mapping("app-canary", "ns")}
					group = "getambassador.io"
					configure()
					fakeKubeClient.GetStub = nil
					fakeKubeClient.GetReturns(nil, apierrors.NewServiceUnavailable("api server unavailable"))
					logBuf.Reset()
					createRequest(http.MethodGet)
					doRequest()

					Expect(res.StatusCode).To(Equal(http.StatusInternalServerError))
					Expect(logBuf.String()).ToNot(ContainSubstring("tolerated vanished"))
				})
			})

			When("a vanished canary Mapping is mixed with a resource that still exists", func() {
				It("keeps the other resource in manifests and names", func() {
					rows = []kubernetes.Resource{
						{
							AccountName: "test-account-name", APIGroup: "apps", Version: "v1",
							Kind: "Deployment", Resource: "deployments", Name: "app", Namespace: "ns",
						},
						mapping("app", "ns"),
						mapping("app-canary", "ns"),
					}
					missing = []string{"app-canary"}
					group = "getambassador.io"
					configure()

					Expect(res.StatusCode).To(Equal(http.StatusOK))
					ro := resultObject()
					Expect(ro.Manifests).To(HaveLen(2))
					Expect(ro.ManifestNamesByNamespace).To(Equal(map[string][]string{"ns": {"Deployment app", "Mapping app"}}))
				})
			})
		})

		When("it succeeds", func() {
			It("succeeds", func() {
				Expect(res.StatusCode).To(Equal(http.StatusOK))
				Expect(logBuf.String()).To(BeEmpty())
			})
		})
	})
})
