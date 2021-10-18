package v1_test

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"

	// . "github.com/homedepot/go-clouddriver/internal/api/v1"

	"github.com/homedepot/go-clouddriver/internal/kubernetes"
	"github.com/jinzhu/gorm"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var _ = Describe("Resource", func() {
	Describe("#LoadKubernetesResources", func() {
		BeforeEach(func() {
			setup()
			fakeSQLClient.DeleteKubernetesResourcesByAccountNameReturns(nil)
			fakeKubeClient.ListResourceWithContextReturns(&unstructured.UnstructuredList{}, nil)
			log.SetOutput(ioutil.Discard)

			uri = svr.URL + "/v1/kubernetes/providers/test-account/resources"
			body.Write([]byte(payloadRequestKubernetesProviders))
			createRequest(http.MethodPut)
		})

		AfterEach(func() {
			teardown()
		})

		JustBeforeEach(func() {
			doRequest()
		})

		When("getting the kubernetes provider for an account errors", func() {
			BeforeEach(func() {
				fakeSQLClient.GetKubernetesProviderReturns(kubernetes.Provider{}, errors.New("error getting provider"))
			})

			It("returns internal server error", func() {
				Expect(res.StatusCode).To(Equal(http.StatusInternalServerError))
			})
		})

		When("discovering API errors", func() {
			BeforeEach(func() {
				fakeKubeClient.DiscoverReturns(errors.New("error discovering API"))
			})

			It("returns internal server error", func() {
				Expect(res.StatusCode).To(Equal(http.StatusInternalServerError))
			})
		})

		When("listing deployments returns an error", func() {
			BeforeEach(func() {
				fakeKubeClient.ListResourceWithContextReturnsOnCall(0, nil, errors.New("error listing deployments"))
			})

			It("continues", func() {
				Expect(res.StatusCode).To(Equal(http.StatusOK))
			})
		})

		When("deleting kubernetes resources errors", func() {
			BeforeEach(func() {
				fakeSQLClient.DeleteKubernetesResourcesByAccountNameReturns(errors.New("error deleting resources"))
			})

			It("returns internal server error", func() {
				Expect(res.StatusCode).To(Equal(http.StatusInternalServerError))
			})
		})

		When("there are no previously deployed resources", func() {
			BeforeEach(func() {
				fakeKubeClient.ListResourceWithContextReturns(&unstructured.UnstructuredList{}, nil)
			})

			It("deletes existing resources and returns status OK", func() {
				Expect(fakeSQLClient.DeleteKubernetesResourcesByAccountNameCallCount()).To(Equal(1))
				Expect(res.StatusCode).To(Equal(http.StatusOK))
			})
		})

		When("using a namespace-scoped provider", func() {
			BeforeEach(func() {
				namespace := "test-namespace"
				fakeSQLClient.GetKubernetesProviderReturns(kubernetes.Provider{
					Name:      "test-account",
					Host:      "http://localhost",
					CAData:    "",
					Namespace: &namespace,
				}, nil)
			})

			It("sets field selector in list options", func() {
				_, _, lo := fakeKubeClient.ListResourceWithContextArgsForCall(0)
				Expect(lo.FieldSelector).To(Equal("metadata.namespace=test-namespace"))
				Expect(res.StatusCode).To(Equal(http.StatusOK))
			})
		})

		When("creating kubernetes resource errors", func() {
			BeforeEach(func() {
				fakeKubeClient.ListResourceWithContextReturnsOnCall(0, fakeDaemonSets, nil)
				fakeKubeClient.ListResourceWithContextReturns(&unstructured.UnstructuredList{}, nil)
				fakeSQLClient.CreateKubernetesResourceReturns(errors.New("error creating resource"))
			})

			It("returns internal server error", func() {
				Expect(res.StatusCode).To(Equal(http.StatusInternalServerError))
			})
		})

		When("list resources returns a daemonset", func() {
			BeforeEach(func() {
				fakeKubeClient.ListResourceWithContextReturnsOnCall(0, fakeDaemonSets, nil)
				fakeKubeClient.ListResourceWithContextReturns(&unstructured.UnstructuredList{}, nil)
				fakeKubeClient.GVRForKindReturns(schema.GroupVersionResource{
					Group:    "apps",
					Version:  "v1",
					Resource: "replicasets",
				}, nil)
			})

			It("returns status OK", func() {
				Expect(res.StatusCode).To(Equal(http.StatusOK))
				Expect(fakeKubeClient.ListResourceWithContextCallCount()).To(Equal(6))
				Expect(fakeSQLClient.DeleteKubernetesResourcesByAccountNameCallCount()).To(Equal(1))
				Expect(fakeSQLClient.CreateKubernetesResourceCallCount()).To(Equal(1))

				kr := fakeSQLClient.CreateKubernetesResourceArgsForCall(0)
				Expect(kr.AccountName).To(Equal("test-account"))
				Expect(kr.APIGroup).To(Equal("apps"))
				Expect(kr.ArtifactName).To(Equal("test-daemonset"))
				Expect(kr.Cluster).To(Equal("daemonSet test-daemonset"))
				Expect(kr.ID).ToNot(BeEmpty())
				Expect(kr.Kind).To(Equal("DaemonSet"))
				Expect(kr.Name).To(Equal("test-daemonset"))
				Expect(kr.Namespace).To(Equal("test-namespace"))
				Expect(kr.Resource).To(Equal("replicasets"))
				Expect(kr.SpinnakerApp).To(Equal("test-application"))
				Expect(kr.TaskID).ToNot(BeEmpty())
				Expect(kr.TaskType).To(BeEmpty())
				Expect(kr.Timestamp).ToNot(BeNil())
				Expect(kr.Version).To(Equal("v1"))
			})
		})

		When("list resources returns multiple deployments", func() {
			BeforeEach(func() {
				fakeKubeClient.ListResourceWithContextReturnsOnCall(1, fakeDeployments, nil)
				fakeKubeClient.ListResourceWithContextReturns(&unstructured.UnstructuredList{}, nil)
				fakeKubeClient.GVRForKindReturns(schema.GroupVersionResource{
					Group:    "apps",
					Version:  "v1",
					Resource: "deployments",
				}, nil)
			})

			It("returns status OK", func() {
				Expect(res.StatusCode).To(Equal(http.StatusOK))
				Expect(fakeKubeClient.ListResourceWithContextCallCount()).To(Equal(6))
				Expect(fakeSQLClient.DeleteKubernetesResourcesByAccountNameCallCount()).To(Equal(1))
				Expect(fakeSQLClient.CreateKubernetesResourceCallCount()).To(Equal(2))

				kr := fakeSQLClient.CreateKubernetesResourceArgsForCall(0)
				Expect(kr.AccountName).To(Equal("test-account"))
				Expect(kr.APIGroup).To(Equal("apps"))
				Expect(kr.ArtifactName).To(Equal("test-deployment1"))
				Expect(kr.Cluster).To(Equal("deployment test-deployment1"))
				Expect(kr.ID).ToNot(BeEmpty())
				Expect(kr.Kind).To(Equal("Deployment"))
				Expect(kr.Name).To(Equal("test-deployment1"))
				Expect(kr.Namespace).To(Equal("test-namespace1"))
				Expect(kr.Resource).To(Equal("deployments"))
				Expect(kr.SpinnakerApp).To(Equal("test-application1"))
				Expect(kr.TaskID).ToNot(BeEmpty())
				Expect(kr.TaskType).To(BeEmpty())
				Expect(kr.Timestamp).ToNot(BeNil())
				Expect(kr.Version).To(Equal("v1"))

				kr = fakeSQLClient.CreateKubernetesResourceArgsForCall(1)
				Expect(kr.AccountName).To(Equal("test-account"))
				Expect(kr.APIGroup).To(Equal("apps"))
				Expect(kr.ArtifactName).To(Equal("test-deployment2"))
				Expect(kr.Cluster).To(Equal("deployment test-deployment2"))
				Expect(kr.ID).ToNot(BeEmpty())
				Expect(kr.Kind).To(Equal("Deployment"))
				Expect(kr.Name).To(Equal("test-deployment2"))
				Expect(kr.Namespace).To(Equal("test-namespace2"))
				Expect(kr.Resource).To(Equal("deployments"))
				Expect(kr.SpinnakerApp).To(Equal("test-application2"))
				Expect(kr.TaskID).ToNot(BeEmpty())
				Expect(kr.TaskType).To(BeEmpty())
				Expect(kr.Timestamp).ToNot(BeNil())
				Expect(kr.Version).To(Equal("v1"))
			})
		})

		When("list resources returns a ingress", func() {
			BeforeEach(func() {
				fakeKubeClient.ListResourceWithContextReturnsOnCall(2, fakeIngresses, nil)
				fakeKubeClient.ListResourceWithContextReturns(&unstructured.UnstructuredList{}, nil)
				fakeKubeClient.GVRForKindReturns(schema.GroupVersionResource{
					Group:    "",
					Version:  "v1beta1",
					Resource: "ingresses",
				}, nil)
			})

			It("returns status OK", func() {
				Expect(res.StatusCode).To(Equal(http.StatusOK))
				Expect(fakeKubeClient.ListResourceWithContextCallCount()).To(Equal(6))
				Expect(fakeSQLClient.DeleteKubernetesResourcesByAccountNameCallCount()).To(Equal(1))
				Expect(fakeSQLClient.CreateKubernetesResourceCallCount()).To(Equal(1))

				kr := fakeSQLClient.CreateKubernetesResourceArgsForCall(0)
				Expect(kr.AccountName).To(Equal("test-account"))
				Expect(kr.APIGroup).To(BeEmpty())
				Expect(kr.ArtifactName).To(Equal("test-ingress"))
				Expect(kr.Cluster).To(Equal("ingress test-ingress"))
				Expect(kr.ID).ToNot(BeEmpty())
				Expect(kr.Kind).To(Equal("Ingress"))
				Expect(kr.Name).To(Equal("test-ingress"))
				Expect(kr.Namespace).To(Equal("test-namespace"))
				Expect(kr.Resource).To(Equal("ingresses"))
				Expect(kr.SpinnakerApp).To(Equal("test-application"))
				Expect(kr.TaskID).ToNot(BeEmpty())
				Expect(kr.TaskType).To(BeEmpty())
				Expect(kr.Timestamp).ToNot(BeNil())
				Expect(kr.Version).To(Equal("v1beta1"))
			})
		})

		When("list resources returns a versioned replicaset", func() {
			BeforeEach(func() {
				fakeKubeClient.ListResourceWithContextReturnsOnCall(3, fakeReplicaSets, nil)
				fakeKubeClient.ListResourceWithContextReturns(&unstructured.UnstructuredList{}, nil)
				fakeKubeClient.GVRForKindReturns(schema.GroupVersionResource{
					Group:    "",
					Version:  "v1",
					Resource: "replicasets",
				}, nil)
			})

			It("returns status OK and NameWithoutVersion is called", func() {
				Expect(res.StatusCode).To(Equal(http.StatusOK))
				Expect(fakeKubeClient.ListResourceWithContextCallCount()).To(Equal(6))
				Expect(fakeSQLClient.DeleteKubernetesResourcesByAccountNameCallCount()).To(Equal(1))
				Expect(fakeSQLClient.CreateKubernetesResourceCallCount()).To(Equal(1))

				kr := fakeSQLClient.CreateKubernetesResourceArgsForCall(0)
				Expect(kr.AccountName).To(Equal("test-account"))
				Expect(kr.APIGroup).To(BeEmpty())
				Expect(kr.ArtifactName).To(Equal("test-replicaset"))
				Expect(kr.Cluster).To(Equal("replicaSet test-replicaset"))
				Expect(kr.ID).ToNot(BeEmpty())
				Expect(kr.Kind).To(Equal("ReplicaSet"))
				Expect(kr.Name).To(Equal("test-replicaset-v001"))
				Expect(kr.Namespace).To(Equal("test-namespace"))
				Expect(kr.Resource).To(Equal("replicasets"))
				Expect(kr.SpinnakerApp).To(Equal("test-application"))
				Expect(kr.TaskID).ToNot(BeEmpty())
				Expect(kr.TaskType).To(BeEmpty())
				Expect(kr.Timestamp).ToNot(BeNil())
				Expect(kr.Version).To(Equal("v1"))
			})
		})

		When("list resources returns a service", func() {
			BeforeEach(func() {
				fakeKubeClient.ListResourceWithContextReturnsOnCall(4, fakeServices, nil)
				fakeKubeClient.ListResourceWithContextReturns(&unstructured.UnstructuredList{}, nil)
				fakeKubeClient.GVRForKindReturns(schema.GroupVersionResource{
					Group:    "",
					Version:  "v1",
					Resource: "services",
				}, nil)
			})

			It("returns status OK", func() {
				Expect(res.StatusCode).To(Equal(http.StatusOK))
				Expect(fakeKubeClient.ListResourceWithContextCallCount()).To(Equal(6))
				Expect(fakeSQLClient.DeleteKubernetesResourcesByAccountNameCallCount()).To(Equal(1))
				Expect(fakeSQLClient.CreateKubernetesResourceCallCount()).To(Equal(1))

				kr := fakeSQLClient.CreateKubernetesResourceArgsForCall(0)
				Expect(kr.AccountName).To(Equal("test-account"))
				Expect(kr.APIGroup).To(BeEmpty())
				Expect(kr.ArtifactName).To(Equal("test-service"))
				Expect(kr.Cluster).To(Equal("service test-service"))
				Expect(kr.ID).ToNot(BeEmpty())
				Expect(kr.Kind).To(Equal("Service"))
				Expect(kr.Name).To(Equal("test-service"))
				Expect(kr.Namespace).To(Equal("test-namespace"))
				Expect(kr.Resource).To(Equal("services"))
				Expect(kr.SpinnakerApp).To(Equal("test-application"))
				Expect(kr.TaskID).ToNot(BeEmpty())
				Expect(kr.TaskType).To(BeEmpty())
				Expect(kr.Timestamp).ToNot(BeNil())
				Expect(kr.Version).To(Equal("v1"))
			})
		})

		When("list resources returns a statefulset", func() {
			BeforeEach(func() {
				fakeKubeClient.ListResourceWithContextReturnsOnCall(5, fakeStatefulSets, nil)
				fakeKubeClient.ListResourceWithContextReturns(&unstructured.UnstructuredList{}, nil)
				fakeKubeClient.GVRForKindReturns(schema.GroupVersionResource{
					Group:    "",
					Version:  "v1",
					Resource: "statefulsets",
				}, nil)
			})

			It("returns status OK", func() {
				Expect(res.StatusCode).To(Equal(http.StatusOK))
				Expect(fakeKubeClient.ListResourceWithContextCallCount()).To(Equal(6))
				Expect(fakeSQLClient.DeleteKubernetesResourcesByAccountNameCallCount()).To(Equal(1))
				Expect(fakeSQLClient.CreateKubernetesResourceCallCount()).To(Equal(1))

				kr := fakeSQLClient.CreateKubernetesResourceArgsForCall(0)
				Expect(kr.AccountName).To(Equal("test-account"))
				Expect(kr.APIGroup).To(BeEmpty())
				Expect(kr.ArtifactName).To(Equal("test-statefulset"))
				Expect(kr.Cluster).To(Equal("statefulSet test-statefulset"))
				Expect(kr.ID).ToNot(BeEmpty())
				Expect(kr.Kind).To(Equal("StatefulSet"))
				Expect(kr.Name).To(Equal("test-statefulset"))
				Expect(kr.Namespace).To(Equal("test-namespace"))
				Expect(kr.Resource).To(Equal("statefulsets"))
				Expect(kr.SpinnakerApp).To(Equal("test-application"))
				Expect(kr.TaskID).ToNot(BeEmpty())
				Expect(kr.TaskType).To(BeEmpty())
				Expect(kr.Timestamp).ToNot(BeNil())
				Expect(kr.Version).To(Equal("v1"))
			})
		})

		When("All list resources call return resources", func() {
			BeforeEach(func() {
				fakeKubeClient.ListResourceWithContextReturnsOnCall(0, fakeDaemonSets, nil)
				fakeKubeClient.ListResourceWithContextReturnsOnCall(1, fakeDeployments, nil)
				fakeKubeClient.ListResourceWithContextReturnsOnCall(2, fakeIngresses, nil)
				fakeKubeClient.ListResourceWithContextReturnsOnCall(3, fakeReplicaSets, nil)
				fakeKubeClient.ListResourceWithContextReturnsOnCall(4, fakeServices, nil)
				fakeKubeClient.ListResourceWithContextReturnsOnCall(5, fakeStatefulSets, nil)
				fakeKubeClient.GVRForKindReturns(schema.GroupVersionResource{
					Group:    "",
					Version:  "v1",
					Resource: "statefulsets",
				}, nil)
			})

			It("returns status OK and creates all resource rows", func() {
				Expect(res.StatusCode).To(Equal(http.StatusOK))
				Expect(fakeKubeClient.ListResourceWithContextCallCount()).To(Equal(6))
				Expect(fakeSQLClient.DeleteKubernetesResourcesByAccountNameCallCount()).To(Equal(1))
				Expect(fakeSQLClient.CreateKubernetesResourceCallCount()).To(Equal(7))
			})
		})
	})

	Describe("#DeleteKubernetesResources", func() {
		BeforeEach(func() {
			setup()
			uri = svr.URL + "/v1/kubernetes/providers/test-name/resources"
			createRequest(http.MethodDelete)
		})

		AfterEach(func() {
			teardown()
		})

		JustBeforeEach(func() {
			doRequest()
		})

		When("the record is not found", func() {
			BeforeEach(func() {
				fakeSQLClient.GetKubernetesProviderReturns(kubernetes.Provider{}, gorm.ErrRecordNotFound)
			})

			It("returns an error", func() {
				Expect(res.StatusCode).To(Equal(http.StatusNotFound))
				validateResponse(payloadKubernetesProviderNotFound)
			})
		})

		When("getting the provider returns a generic error", func() {
			BeforeEach(func() {
				fakeSQLClient.GetKubernetesProviderReturns(kubernetes.Provider{}, errors.New("error getting provider"))
			})

			It("returns an error", func() {
				Expect(res.StatusCode).To(Equal(http.StatusInternalServerError))
				validateResponse(payloadKubernetesProviderGetGenericError)
			})
		})

		When("deleting the resources returns an error", func() {
			BeforeEach(func() {
				fakeSQLClient.DeleteKubernetesResourcesByAccountNameReturns(errors.New("error deleting resources"))
			})

			It("returns an error", func() {
				Expect(res.StatusCode).To(Equal(http.StatusInternalServerError))
				validateResponse(payloadKubernetesResourcesDeleteGenericError)
			})
		})

		When("it succeeds", func() {
			It("returns status no content", func() {
				Expect(res.StatusCode).To(Equal(http.StatusNoContent))
			})
		})
	})
})
