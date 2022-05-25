package kubernetes_test

import (
	"encoding/json"

	"github.com/homedepot/go-clouddriver/internal/artifact"
	. "github.com/homedepot/go-clouddriver/internal/kubernetes"
	clouddriver "github.com/homedepot/go-clouddriver/pkg"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

var _ = Describe("Artifact", func() {
	Describe("#BindArtifacts", func() {
		var (
			resource  *unstructured.Unstructured
			artifacts []clouddriver.Artifact
		)

		BeforeEach(func() {
			artifacts = []clouddriver.Artifact{
				{
					Name:      "gcr.io/test-project/test-container-image",
					Type:      artifact.TypeDockerImage,
					Reference: "gcr.io/test-project/test-container-image:v1.0.0",
				},
				{
					Name:      "my-config-map",
					Type:      artifact.TypeKubernetesConfigMap,
					Reference: "my-config-map-v000",
				},
				{
					Name:      "my-config-map2",
					Type:      artifact.TypeKubernetesConfigMap,
					Reference: "my-config-map2-v000",
				},
				{
					Name:      "my-secret",
					Type:      artifact.TypeKubernetesSecret,
					Reference: "my-secret-v000",
				},
				{
					Name:      "my-secret2",
					Type:      artifact.TypeKubernetesSecret,
					Reference: "my-secret2-v000",
				},
				{
					Name:      "my-deployment",
					Type:      artifact.TypeKubernetesDeployment,
					Reference: "my-deployment-v000",
				},
				{
					Name:      "my-replicaSet",
					Type:      artifact.TypeKubernetesReplicaSet,
					Reference: "my-replicaSet-v000",
				},
			}
		})

		JustBeforeEach(func() {
			BindArtifacts(resource, artifacts)
		})

		When("the iterable path is not of type []interface{}", func() {
			BeforeEach(func() {
				resource = &unstructured.Unstructured{
					Object: map[string]interface{}{
						"kind":       "Pod",
						"apiVersion": "v1",
						"spec": map[string]interface{}{
							"containers": map[string]interface{}{},
						},
					},
				}
			})

			It("leaves the resource as is", func() {
				b, err := json.Marshal(resource)
				Expect(err).To(BeNil())
				Expect(string(b)).To(MatchJSON(`{
          "apiVersion": "v1",
          "kind": "Pod",
          "spec": {
            "containers": {}
          }
        }`))
			})
		})

		When("the value of an iterable path is not of type map[string]interface{}", func() {
			BeforeEach(func() {
				resource = &unstructured.Unstructured{
					Object: map[string]interface{}{
						"kind":       "Deployment",
						"apiVersion": "apps/v1",
						"spec": map[string]interface{}{
							"template": map[string]interface{}{
								"spec": map[string]interface{}{
									"containers": []interface{}{
										"string",
										"slice",
									},
								},
							},
						},
					},
				}
			})

			It("leaves the resource as is", func() {
				b, err := json.Marshal(resource)
				Expect(err).To(BeNil())
				Expect(string(b)).To(MatchJSON(`{
          "apiVersion": "apps/v1",
          "kind": "Deployment",
          "spec": {
            "template": {
              "spec": {
                "containers": [
                  "string",
                  "slice"
                ]
              }
            }
          }
        }`))
			})
		})

		When("the expected final nested string value is not a string", func() {
			BeforeEach(func() {
				resource = &unstructured.Unstructured{
					Object: map[string]interface{}{
						"kind":       "Pod",
						"apiVersion": "v1",
						"spec": map[string]interface{}{
							"containers": []interface{}{
								map[string]interface{}{
									"name":  "another-test-container-name",
									"image": []interface{}{"gcr.io/test-project/another-test-container-image"},
								},
								map[string]interface{}{
									"name":  "test-container-name",
									"image": "gcr.io/test-project/fake-test-container-image",
								},
							},
						},
					},
				}
			})

			It("leaves the resource as is", func() {
				b, err := json.Marshal(resource)
				Expect(err).To(BeNil())
				Expect(string(b)).To(MatchJSON(`{
          "apiVersion": "v1",
          "kind": "Pod",
          "spec": {
            "containers": [
              {
                "image": [
                  "gcr.io/test-project/another-test-container-image"
                ],
                "name": "another-test-container-name"
              },
              {
                "image": "gcr.io/test-project/fake-test-container-image",
                "name": "test-container-name"
              }
            ]
          }
        }`))
			})
		})

		Describe("docker/image", func() {
			When(".spec.containers.*.image", func() {
				BeforeEach(func() {
					resource = &unstructured.Unstructured{
						Object: map[string]interface{}{
							"kind":       "Pod",
							"apiVersion": "v1",
							"spec": map[string]interface{}{
								"containers": []interface{}{
									map[string]interface{}{
										"name":  "another-test-container-name",
										"image": "gcr.io/test-project/another-test-container-image",
									},
									map[string]interface{}{
										"name":  "test-container-name",
										"image": "gcr.io/test-project/test-container-image",
									},
								},
							},
						},
					}
				})

				It("updates the container image to contain the reference", func() {
					o := NewPod(resource.Object)
					containers := o.Object().Spec.Containers
					Expect(containers).To(HaveLen(2))
					Expect(containers[0].Image).To(Equal("gcr.io/test-project/another-test-container-image"))
					Expect(containers[1].Image).To(Equal("gcr.io/test-project/test-container-image:v1.0.0"))
				})
			})

			When(".spec.template.spec.containers.*.image", func() {
				BeforeEach(func() {
					resource = &unstructured.Unstructured{
						Object: map[string]interface{}{
							"kind":       "Deployment",
							"apiVersion": "apps/v1",
							"spec": map[string]interface{}{
								"template": map[string]interface{}{
									"spec": map[string]interface{}{
										"containers": []interface{}{
											map[string]interface{}{
												"name":  "test-container-name",
												"image": "gcr.io/test-project/test-container-image",
											},
											map[string]interface{}{
												"name":  "another-test-container-name",
												"image": "gcr.io/test-project/another-test-container-image",
											},
										},
									},
								},
							},
						},
					}
				})

				It("only updates the container in the artifact", func() {
					o := NewDeployment(resource.Object)
					containers := o.Object().Spec.Template.Spec.Containers
					Expect(containers).To(HaveLen(2))
					Expect(containers[0].Image).To(Equal("gcr.io/test-project/test-container-image:v1.0.0"))
					Expect(containers[1].Image).To(Equal("gcr.io/test-project/another-test-container-image"))
				})
			})

			When(".spec.initContainers.*.image", func() {
				BeforeEach(func() {
					resource = &unstructured.Unstructured{
						Object: map[string]interface{}{
							"kind":       "Pod",
							"apiVersion": "v1",
							"spec": map[string]interface{}{
								"initContainers": []interface{}{
									map[string]interface{}{
										"name":  "test-init-container-name",
										"image": "gcr.io/test-project/test-container-image",
									},
									map[string]interface{}{
										"name":  "another-test-init-container-name",
										"image": "gcr.io/test-project/another-test-container-image",
									},
								},
							},
						},
					}
				})

				It("updates the init container image to contain the reference", func() {
					o := NewPod(resource.Object)
					initContainers := o.Object().Spec.InitContainers
					Expect(initContainers).To(HaveLen(2))
					Expect(initContainers[0].Image).To(Equal("gcr.io/test-project/test-container-image:v1.0.0"))
					Expect(initContainers[1].Image).To(Equal("gcr.io/test-project/another-test-container-image"))
				})
			})

			When(".spec.template.spec.initContainers.*.image", func() {
				BeforeEach(func() {
					resource = &unstructured.Unstructured{
						Object: map[string]interface{}{
							"kind":       "Deployment",
							"apiVersion": "apps/v1",
							"spec": map[string]interface{}{
								"template": map[string]interface{}{
									"spec": map[string]interface{}{
										"initContainers": []interface{}{
											map[string]interface{}{
												"name":  "another-test-init-container-name",
												"image": "gcr.io/test-project/another-test-container-image",
											},
											map[string]interface{}{
												"name":  "test-init-container-name",
												"image": "gcr.io/test-project/test-container-image",
											},
										},
									},
								},
							},
						},
					}
				})

				It("only updates the init container in the artifact", func() {
					o := NewDeployment(resource.Object)
					initContainers := o.Object().Spec.Template.Spec.InitContainers
					Expect(initContainers).To(HaveLen(2))
					Expect(initContainers[0].Image).To(Equal("gcr.io/test-project/another-test-container-image"))
					Expect(initContainers[1].Image).To(Equal("gcr.io/test-project/test-container-image:v1.0.0"))
				})
			})
		})

		Describe("kubernetes/configMap", func() {
			Context(".spec.template.spec.volumes.*.configMap.name", func() {
				BeforeEach(func() {
					resource = &unstructured.Unstructured{
						Object: map[string]interface{}{
							"kind":       "Deployment",
							"apiVersion": "apps/v1",
							"spec": map[string]interface{}{
								"template": map[string]interface{}{
									"spec": map[string]interface{}{
										"volumes": []interface{}{
											map[string]interface{}{
												"configMap": map[string]interface{}{
													"name": "not-my-config-map",
												},
											},
											map[string]interface{}{
												"configMap": map[string]interface{}{
													"name": "my-config-map",
												},
											},
										},
									},
								},
							},
						},
					}
				})

				It("replaces the configMap", func() {
					o := NewDeployment(resource.Object)
					volumes := o.Object().Spec.Template.Spec.Volumes
					Expect(volumes).To(HaveLen(2))
					Expect(volumes[0].VolumeSource.ConfigMap.Name).To(Equal("not-my-config-map"))
					Expect(volumes[1].VolumeSource.ConfigMap.Name).To(Equal("my-config-map-v000"))
				})
			})

			Context(".spec.volumes.*.configMap.name", func() {
				BeforeEach(func() {
					resource = &unstructured.Unstructured{
						Object: map[string]interface{}{
							"kind":       "Pod",
							"apiVersion": "v1",
							"spec": map[string]interface{}{
								"volumes": []interface{}{
									map[string]interface{}{
										"configMap": map[string]interface{}{
											"name": "not-my-config-map",
										},
									},
									map[string]interface{}{
										"configMap": map[string]interface{}{
											"name": "my-config-map",
										},
									},
								},
							},
						},
					}
				})

				It("replaces the configMap", func() {
					o := NewPod(resource.Object)
					volumes := o.Object().Spec.Volumes
					Expect(volumes).To(HaveLen(2))
					Expect(volumes[0].VolumeSource.ConfigMap.Name).To(Equal("not-my-config-map"))
					Expect(volumes[1].VolumeSource.ConfigMap.Name).To(Equal("my-config-map-v000"))
				})
			})

			Context(".spec.template.spec.volumes.*.projected.sources.*.configMap.name", func() {
				BeforeEach(func() {
					resource = &unstructured.Unstructured{
						Object: map[string]interface{}{
							"kind":       "Deployment",
							"apiVersion": "apps/v1",
							"spec": map[string]interface{}{
								"template": map[string]interface{}{
									"spec": map[string]interface{}{
										"volumes": []interface{}{
											map[string]interface{}{
												"projected": map[string]interface{}{
													"sources": []interface{}{
														map[string]interface{}{
															"configMap": map[string]interface{}{
																"name": "not-my-config-map",
															},
														},
														map[string]interface{}{
															"secret": map[string]interface{}{
																"name": "not-my-secret",
															},
														},
														map[string]interface{}{
															"configMap": map[string]interface{}{
																"name": "my-config-map",
															},
														},
														map[string]interface{}{
															"secret": map[string]interface{}{
																"name": "not-my-secret2",
															},
														},
													},
												},
											},
											map[string]interface{}{
												"projected": map[string]interface{}{
													"sources": []interface{}{
														map[string]interface{}{
															"configMap": map[string]interface{}{
																"name": "not-my-config-map",
															},
														},
														map[string]interface{}{
															"secret": map[string]interface{}{
																"name": "not-my-secret",
															},
														},
														map[string]interface{}{
															"configMap": map[string]interface{}{
																"name": "my-config-map2",
															},
														},
														map[string]interface{}{
															"secret": map[string]interface{}{
																"name": "not-my-secret2",
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					}
				})

				It("replaces the configMap", func() {
					o := NewDeployment(resource.Object)
					volumes := o.Object().Spec.Template.Spec.Volumes
					Expect(volumes).To(HaveLen(2))
					Expect(volumes[0].VolumeSource.Projected.Sources).To(HaveLen(4))
					Expect(volumes[0].VolumeSource.Projected.Sources[0].ConfigMap.Name).To(Equal("not-my-config-map"))
					Expect(volumes[0].VolumeSource.Projected.Sources[1].Secret.Name).To(Equal("not-my-secret"))
					Expect(volumes[0].VolumeSource.Projected.Sources[2].ConfigMap.Name).To(Equal("my-config-map-v000"))
					Expect(volumes[0].VolumeSource.Projected.Sources[3].Secret.Name).To(Equal("not-my-secret2"))
					Expect(volumes[1].VolumeSource.Projected.Sources).To(HaveLen(4))
					Expect(volumes[1].VolumeSource.Projected.Sources[0].ConfigMap.Name).To(Equal("not-my-config-map"))
					Expect(volumes[1].VolumeSource.Projected.Sources[1].Secret.Name).To(Equal("not-my-secret"))
					Expect(volumes[1].VolumeSource.Projected.Sources[2].ConfigMap.Name).To(Equal("my-config-map2-v000"))
					Expect(volumes[1].VolumeSource.Projected.Sources[3].Secret.Name).To(Equal("not-my-secret2"))
				})
			})

			Context(".spec.volumes.*.projected.sources.*configMap.name", func() {
				BeforeEach(func() {
					resource = &unstructured.Unstructured{
						Object: map[string]interface{}{
							"kind":       "Pod",
							"apiVersion": "v1",
							"spec": map[string]interface{}{
								"volumes": []interface{}{
									map[string]interface{}{
										"projected": map[string]interface{}{
											"sources": []interface{}{
												map[string]interface{}{
													"configMap": map[string]interface{}{
														"name": "not-my-config-map",
													},
												},
												map[string]interface{}{
													"secret": map[string]interface{}{
														"name": "not-my-secret",
													},
												},
												map[string]interface{}{
													"configMap": map[string]interface{}{
														"name": "my-config-map",
													},
												},
												map[string]interface{}{
													"secret": map[string]interface{}{
														"name": "not-my-secret2",
													},
												},
											},
										},
									},
									map[string]interface{}{
										"projected": map[string]interface{}{
											"sources": []interface{}{
												map[string]interface{}{
													"configMap": map[string]interface{}{
														"name": "not-my-config-map",
													},
												},
												map[string]interface{}{
													"secret": map[string]interface{}{
														"name": "not-my-secret",
													},
												},
												map[string]interface{}{
													"configMap": map[string]interface{}{
														"name": "my-config-map2",
													},
												},
												map[string]interface{}{
													"secret": map[string]interface{}{
														"name": "not-my-secret2",
													},
												},
											},
										},
									},
								},
							},
						},
					}
				})

				It("replaces the configMap", func() {
					o := NewPod(resource.Object)
					volumes := o.Object().Spec.Volumes
					Expect(volumes).To(HaveLen(2))
					Expect(volumes[0].VolumeSource.Projected.Sources).To(HaveLen(4))
					Expect(volumes[0].VolumeSource.Projected.Sources[0].ConfigMap.Name).To(Equal("not-my-config-map"))
					Expect(volumes[0].VolumeSource.Projected.Sources[1].Secret.Name).To(Equal("not-my-secret"))
					Expect(volumes[0].VolumeSource.Projected.Sources[2].ConfigMap.Name).To(Equal("my-config-map-v000"))
					Expect(volumes[0].VolumeSource.Projected.Sources[3].Secret.Name).To(Equal("not-my-secret2"))
					Expect(volumes[1].VolumeSource.Projected.Sources).To(HaveLen(4))
					Expect(volumes[1].VolumeSource.Projected.Sources[0].ConfigMap.Name).To(Equal("not-my-config-map"))
					Expect(volumes[1].VolumeSource.Projected.Sources[1].Secret.Name).To(Equal("not-my-secret"))
					Expect(volumes[1].VolumeSource.Projected.Sources[2].ConfigMap.Name).To(Equal("my-config-map2-v000"))
					Expect(volumes[1].VolumeSource.Projected.Sources[3].Secret.Name).To(Equal("not-my-secret2"))
				})
			})

			Context(".spec.template.spec.containers.*.env.*.valueFrom.configMapKeyRef.name", func() {
				BeforeEach(func() {
					resource = &unstructured.Unstructured{
						Object: map[string]interface{}{
							"kind":       "Deployment",
							"apiVersion": "apps/v1",
							"spec": map[string]interface{}{
								"template": map[string]interface{}{
									"spec": map[string]interface{}{
										"containers": []interface{}{
											map[string]interface{}{
												"env": []interface{}{
													map[string]interface{}{
														"valueFrom": map[string]interface{}{
															"configMapKeyRef": map[string]interface{}{
																"name": "not-my-config-map",
															},
														},
													},
													map[string]interface{}{
														"valueFrom": map[string]interface{}{
															"secretKeyRef": map[string]interface{}{
																"name": "not-my-secret",
															},
														},
													},
													map[string]interface{}{
														"valueFrom": map[string]interface{}{
															"configMapKeyRef": map[string]interface{}{
																"name": "my-config-map",
															},
														},
													},
													map[string]interface{}{
														"valueFrom": map[string]interface{}{
															"secretKeyRef": map[string]interface{}{
																"name": "not-my-secret2",
															},
														},
													},
												},
											},
											map[string]interface{}{
												"env": []interface{}{
													map[string]interface{}{
														"valueFrom": map[string]interface{}{
															"configMapKeyRef": map[string]interface{}{
																"name": "not-my-config-map",
															},
														},
													},
													map[string]interface{}{
														"valueFrom": map[string]interface{}{
															"secretKeyRef": map[string]interface{}{
																"name": "not-my-secret",
															},
														},
													},
													map[string]interface{}{
														"valueFrom": map[string]interface{}{
															"configMapKeyRef": map[string]interface{}{
																"name": "my-config-map2",
															},
														},
													},
													map[string]interface{}{
														"valueFrom": map[string]interface{}{
															"secretKeyRef": map[string]interface{}{
																"name": "not-my-secret2",
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					}
				})

				It("replaces the configMap", func() {
					o := NewDeployment(resource.Object)
					containers := o.Object().Spec.Template.Spec.Containers
					Expect(containers).To(HaveLen(2))
					Expect(containers[0].Env).To(HaveLen(4))
					Expect(containers[0].Env[0].ValueFrom.ConfigMapKeyRef.Name).To(Equal("not-my-config-map"))
					Expect(containers[0].Env[1].ValueFrom.SecretKeyRef.Name).To(Equal("not-my-secret"))
					Expect(containers[0].Env[2].ValueFrom.ConfigMapKeyRef.Name).To(Equal("my-config-map-v000"))
					Expect(containers[0].Env[3].ValueFrom.SecretKeyRef.Name).To(Equal("not-my-secret2"))
					Expect(containers[1].Env).To(HaveLen(4))
					Expect(containers[1].Env[0].ValueFrom.ConfigMapKeyRef.Name).To(Equal("not-my-config-map"))
					Expect(containers[1].Env[1].ValueFrom.SecretKeyRef.Name).To(Equal("not-my-secret"))
					Expect(containers[1].Env[2].ValueFrom.ConfigMapKeyRef.Name).To(Equal("my-config-map2-v000"))
					Expect(containers[1].Env[3].ValueFrom.SecretKeyRef.Name).To(Equal("not-my-secret2"))
				})
			})

			Context(".spec.containers.*.env.*.valueFrom.configMapKeyRef.name", func() {
				BeforeEach(func() {
					resource = &unstructured.Unstructured{
						Object: map[string]interface{}{
							"kind":       "Pod",
							"apiVersion": "v1",
							"spec": map[string]interface{}{
								"containers": []interface{}{
									map[string]interface{}{
										"env": []interface{}{
											map[string]interface{}{
												"valueFrom": map[string]interface{}{
													"configMapKeyRef": map[string]interface{}{
														"name": "not-my-config-map",
													},
												},
											},
											map[string]interface{}{
												"valueFrom": map[string]interface{}{
													"secretKeyRef": map[string]interface{}{
														"name": "not-my-secret",
													},
												},
											},
											map[string]interface{}{
												"valueFrom": map[string]interface{}{
													"configMapKeyRef": map[string]interface{}{
														"name": "my-config-map",
													},
												},
											},
											map[string]interface{}{
												"valueFrom": map[string]interface{}{
													"secretKeyRef": map[string]interface{}{
														"name": "not-my-secret2",
													},
												},
											},
										},
									},
									map[string]interface{}{
										"env": []interface{}{
											map[string]interface{}{
												"valueFrom": map[string]interface{}{
													"configMapKeyRef": map[string]interface{}{
														"name": "not-my-config-map",
													},
												},
											},
											map[string]interface{}{
												"valueFrom": map[string]interface{}{
													"secretKeyRef": map[string]interface{}{
														"name": "not-my-secret",
													},
												},
											},
											map[string]interface{}{
												"valueFrom": map[string]interface{}{
													"configMapKeyRef": map[string]interface{}{
														"name": "my-config-map2",
													},
												},
											},
											map[string]interface{}{
												"valueFrom": map[string]interface{}{
													"secretKeyRef": map[string]interface{}{
														"name": "not-my-secret2",
													},
												},
											},
										},
									},
								},
							},
						},
					}
				})

				It("replaces the configMap", func() {
					o := NewPod(resource.Object)
					containers := o.Object().Spec.Containers
					Expect(containers).To(HaveLen(2))
					Expect(containers[0].Env).To(HaveLen(4))
					Expect(containers[0].Env[0].ValueFrom.ConfigMapKeyRef.Name).To(Equal("not-my-config-map"))
					Expect(containers[0].Env[1].ValueFrom.SecretKeyRef.Name).To(Equal("not-my-secret"))
					Expect(containers[0].Env[2].ValueFrom.ConfigMapKeyRef.Name).To(Equal("my-config-map-v000"))
					Expect(containers[0].Env[3].ValueFrom.SecretKeyRef.Name).To(Equal("not-my-secret2"))
					Expect(containers[1].Env).To(HaveLen(4))
					Expect(containers[1].Env[0].ValueFrom.ConfigMapKeyRef.Name).To(Equal("not-my-config-map"))
					Expect(containers[1].Env[1].ValueFrom.SecretKeyRef.Name).To(Equal("not-my-secret"))
					Expect(containers[1].Env[2].ValueFrom.ConfigMapKeyRef.Name).To(Equal("my-config-map2-v000"))
					Expect(containers[1].Env[3].ValueFrom.SecretKeyRef.Name).To(Equal("not-my-secret2"))
				})
			})

			Context(".spec.template.spec.initContainers.*.env.*.valueFrom.configMapKeyRef.name", func() {
				BeforeEach(func() {
					resource = &unstructured.Unstructured{
						Object: map[string]interface{}{
							"kind":       "Deployment",
							"apiVersion": "apps/v1",
							"spec": map[string]interface{}{
								"template": map[string]interface{}{
									"spec": map[string]interface{}{
										"initContainers": []interface{}{
											map[string]interface{}{
												"env": []interface{}{
													map[string]interface{}{
														"valueFrom": map[string]interface{}{
															"configMapKeyRef": map[string]interface{}{
																"name": "not-my-config-map",
															},
														},
													},
													map[string]interface{}{
														"valueFrom": map[string]interface{}{
															"secretKeyRef": map[string]interface{}{
																"name": "not-my-secret",
															},
														},
													},
													map[string]interface{}{
														"valueFrom": map[string]interface{}{
															"configMapKeyRef": map[string]interface{}{
																"name": "my-config-map",
															},
														},
													},
													map[string]interface{}{
														"valueFrom": map[string]interface{}{
															"secretKeyRef": map[string]interface{}{
																"name": "not-my-secret2",
															},
														},
													},
												},
											},
											map[string]interface{}{
												"env": []interface{}{
													map[string]interface{}{
														"valueFrom": map[string]interface{}{
															"configMapKeyRef": map[string]interface{}{
																"name": "not-my-config-map",
															},
														},
													},
													map[string]interface{}{
														"valueFrom": map[string]interface{}{
															"secretKeyRef": map[string]interface{}{
																"name": "not-my-secret",
															},
														},
													},
													map[string]interface{}{
														"valueFrom": map[string]interface{}{
															"configMapKeyRef": map[string]interface{}{
																"name": "my-config-map2",
															},
														},
													},
													map[string]interface{}{
														"valueFrom": map[string]interface{}{
															"secretKeyRef": map[string]interface{}{
																"name": "not-my-secret2",
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					}
				})

				It("replaces the configMap", func() {
					o := NewDeployment(resource.Object)
					containers := o.Object().Spec.Template.Spec.InitContainers
					Expect(containers).To(HaveLen(2))
					Expect(containers[0].Env).To(HaveLen(4))
					Expect(containers[0].Env[0].ValueFrom.ConfigMapKeyRef.Name).To(Equal("not-my-config-map"))
					Expect(containers[0].Env[1].ValueFrom.SecretKeyRef.Name).To(Equal("not-my-secret"))
					Expect(containers[0].Env[2].ValueFrom.ConfigMapKeyRef.Name).To(Equal("my-config-map-v000"))
					Expect(containers[0].Env[3].ValueFrom.SecretKeyRef.Name).To(Equal("not-my-secret2"))
					Expect(containers[1].Env).To(HaveLen(4))
					Expect(containers[1].Env[0].ValueFrom.ConfigMapKeyRef.Name).To(Equal("not-my-config-map"))
					Expect(containers[1].Env[1].ValueFrom.SecretKeyRef.Name).To(Equal("not-my-secret"))
					Expect(containers[1].Env[2].ValueFrom.ConfigMapKeyRef.Name).To(Equal("my-config-map2-v000"))
					Expect(containers[1].Env[3].ValueFrom.SecretKeyRef.Name).To(Equal("not-my-secret2"))
				})
			})

			Context(".spec.initContainers.*.env.*.valueFrom.configMapKeyRef.name", func() {
				BeforeEach(func() {
					resource = &unstructured.Unstructured{
						Object: map[string]interface{}{
							"kind":       "Pod",
							"apiVersion": "v1",
							"spec": map[string]interface{}{
								"initContainers": []interface{}{
									map[string]interface{}{
										"env": []interface{}{
											map[string]interface{}{
												"valueFrom": map[string]interface{}{
													"configMapKeyRef": map[string]interface{}{
														"name": "not-my-config-map",
													},
												},
											},
											map[string]interface{}{
												"valueFrom": map[string]interface{}{
													"secretKeyRef": map[string]interface{}{
														"name": "not-my-secret",
													},
												},
											},
											map[string]interface{}{
												"valueFrom": map[string]interface{}{
													"configMapKeyRef": map[string]interface{}{
														"name": "my-config-map",
													},
												},
											},
											map[string]interface{}{
												"valueFrom": map[string]interface{}{
													"secretKeyRef": map[string]interface{}{
														"name": "not-my-secret2",
													},
												},
											},
										},
									},
									map[string]interface{}{
										"env": []interface{}{
											map[string]interface{}{
												"valueFrom": map[string]interface{}{
													"configMapKeyRef": map[string]interface{}{
														"name": "not-my-config-map",
													},
												},
											},
											map[string]interface{}{
												"valueFrom": map[string]interface{}{
													"secretKeyRef": map[string]interface{}{
														"name": "not-my-secret",
													},
												},
											},
											map[string]interface{}{
												"valueFrom": map[string]interface{}{
													"configMapKeyRef": map[string]interface{}{
														"name": "my-config-map2",
													},
												},
											},
											map[string]interface{}{
												"valueFrom": map[string]interface{}{
													"secretKeyRef": map[string]interface{}{
														"name": "not-my-secret2",
													},
												},
											},
										},
									},
								},
							},
						},
					}
				})

				It("replaces the configMap", func() {
					o := NewPod(resource.Object)
					containers := o.Object().Spec.InitContainers
					Expect(containers).To(HaveLen(2))
					Expect(containers[0].Env).To(HaveLen(4))
					Expect(containers[0].Env[0].ValueFrom.ConfigMapKeyRef.Name).To(Equal("not-my-config-map"))
					Expect(containers[0].Env[1].ValueFrom.SecretKeyRef.Name).To(Equal("not-my-secret"))
					Expect(containers[0].Env[2].ValueFrom.ConfigMapKeyRef.Name).To(Equal("my-config-map-v000"))
					Expect(containers[0].Env[3].ValueFrom.SecretKeyRef.Name).To(Equal("not-my-secret2"))
					Expect(containers[1].Env).To(HaveLen(4))
					Expect(containers[1].Env[0].ValueFrom.ConfigMapKeyRef.Name).To(Equal("not-my-config-map"))
					Expect(containers[1].Env[1].ValueFrom.SecretKeyRef.Name).To(Equal("not-my-secret"))
					Expect(containers[1].Env[2].ValueFrom.ConfigMapKeyRef.Name).To(Equal("my-config-map2-v000"))
					Expect(containers[1].Env[3].ValueFrom.SecretKeyRef.Name).To(Equal("not-my-secret2"))
				})
			})

			Context(".spec.template.spec.containers.*.envFrom.*.configMapRef.name", func() {
				BeforeEach(func() {
					resource = &unstructured.Unstructured{
						Object: map[string]interface{}{
							"kind":       "Deployment",
							"apiVersion": "apps/v1",
							"spec": map[string]interface{}{
								"template": map[string]interface{}{
									"spec": map[string]interface{}{
										"containers": []interface{}{
											map[string]interface{}{
												"envFrom": []interface{}{
													map[string]interface{}{
														"configMapRef": map[string]interface{}{
															"name": "not-my-config-map",
														},
													},
													map[string]interface{}{
														"secretRef": map[string]interface{}{
															"name": "not-my-secret",
														},
													},
													map[string]interface{}{
														"configMapRef": map[string]interface{}{
															"name": "my-config-map",
														},
													},
													map[string]interface{}{
														"secretRef": map[string]interface{}{
															"name": "not-my-secret2",
														},
													},
												},
											},
											map[string]interface{}{
												"envFrom": []interface{}{
													map[string]interface{}{
														"configMapRef": map[string]interface{}{
															"name": "not-my-config-map",
														},
													},
													map[string]interface{}{
														"secretRef": map[string]interface{}{
															"name": "not-my-secret",
														},
													},
													map[string]interface{}{
														"configMapRef": map[string]interface{}{
															"name": "my-config-map2",
														},
													},
													map[string]interface{}{
														"secretRef": map[string]interface{}{
															"name": "not-my-secret2",
														},
													},
												},
											},
										},
									},
								},
							},
						},
					}
				})

				It("replaces the configMap", func() {
					o := NewDeployment(resource.Object)
					containers := o.Object().Spec.Template.Spec.Containers
					Expect(containers).To(HaveLen(2))
					Expect(containers[0].EnvFrom).To(HaveLen(4))
					Expect(containers[0].EnvFrom[0].ConfigMapRef.Name).To(Equal("not-my-config-map"))
					Expect(containers[0].EnvFrom[1].SecretRef.Name).To(Equal("not-my-secret"))
					Expect(containers[0].EnvFrom[2].ConfigMapRef.Name).To(Equal("my-config-map-v000"))
					Expect(containers[0].EnvFrom[3].SecretRef.Name).To(Equal("not-my-secret2"))
					Expect(containers[1].EnvFrom).To(HaveLen(4))
					Expect(containers[1].EnvFrom[0].ConfigMapRef.Name).To(Equal("not-my-config-map"))
					Expect(containers[1].EnvFrom[1].SecretRef.Name).To(Equal("not-my-secret"))
					Expect(containers[1].EnvFrom[2].ConfigMapRef.Name).To(Equal("my-config-map2-v000"))
					Expect(containers[1].EnvFrom[3].SecretRef.Name).To(Equal("not-my-secret2"))
				})
			})

			Context(".spec.containers.*.envFrom.*.configMapRef.name", func() {
				BeforeEach(func() {
					resource = &unstructured.Unstructured{
						Object: map[string]interface{}{
							"kind":       "Pod",
							"apiVersion": "v1",
							"spec": map[string]interface{}{
								"containers": []interface{}{
									map[string]interface{}{
										"envFrom": []interface{}{
											map[string]interface{}{
												"configMapRef": map[string]interface{}{
													"name": "not-my-config-map",
												},
											},
											map[string]interface{}{
												"secretRef": map[string]interface{}{
													"name": "not-my-secret",
												},
											},
											map[string]interface{}{
												"configMapRef": map[string]interface{}{
													"name": "my-config-map",
												},
											},
											map[string]interface{}{
												"secretRef": map[string]interface{}{
													"name": "not-my-secret2",
												},
											},
										},
									},
									map[string]interface{}{
										"envFrom": []interface{}{
											map[string]interface{}{
												"configMapRef": map[string]interface{}{
													"name": "not-my-config-map",
												},
											},
											map[string]interface{}{
												"secretRef": map[string]interface{}{
													"name": "not-my-secret",
												},
											},
											map[string]interface{}{
												"configMapRef": map[string]interface{}{
													"name": "my-config-map2",
												},
											},
											map[string]interface{}{
												"secretRef": map[string]interface{}{
													"name": "not-my-secret2",
												},
											},
										},
									},
								},
							},
						},
					}
				})

				It("replaces the configMap", func() {
					o := NewPod(resource.Object)
					containers := o.Object().Spec.Containers
					Expect(containers).To(HaveLen(2))
					Expect(containers[0].EnvFrom).To(HaveLen(4))
					Expect(containers[0].EnvFrom[0].ConfigMapRef.Name).To(Equal("not-my-config-map"))
					Expect(containers[0].EnvFrom[1].SecretRef.Name).To(Equal("not-my-secret"))
					Expect(containers[0].EnvFrom[2].ConfigMapRef.Name).To(Equal("my-config-map-v000"))
					Expect(containers[0].EnvFrom[3].SecretRef.Name).To(Equal("not-my-secret2"))
					Expect(containers[1].EnvFrom).To(HaveLen(4))
					Expect(containers[1].EnvFrom[0].ConfigMapRef.Name).To(Equal("not-my-config-map"))
					Expect(containers[1].EnvFrom[1].SecretRef.Name).To(Equal("not-my-secret"))
					Expect(containers[1].EnvFrom[2].ConfigMapRef.Name).To(Equal("my-config-map2-v000"))
					Expect(containers[1].EnvFrom[3].SecretRef.Name).To(Equal("not-my-secret2"))
				})
			})

			Context(".spec.template.spec.initContainers.*.envFrom.*.configMapRef.name", func() {
				BeforeEach(func() {
					resource = &unstructured.Unstructured{
						Object: map[string]interface{}{
							"kind":       "Deployment",
							"apiVersion": "apps/v1",
							"spec": map[string]interface{}{
								"template": map[string]interface{}{
									"spec": map[string]interface{}{
										"initContainers": []interface{}{
											map[string]interface{}{
												"envFrom": []interface{}{
													map[string]interface{}{
														"configMapRef": map[string]interface{}{
															"name": "not-my-config-map",
														},
													},
													map[string]interface{}{
														"secretRef": map[string]interface{}{
															"name": "not-my-secret",
														},
													},
													map[string]interface{}{
														"configMapRef": map[string]interface{}{
															"name": "my-config-map",
														},
													},
													map[string]interface{}{
														"secretRef": map[string]interface{}{
															"name": "not-my-secret2",
														},
													},
												},
											},
											map[string]interface{}{
												"envFrom": []interface{}{
													map[string]interface{}{
														"configMapRef": map[string]interface{}{
															"name": "not-my-config-map",
														},
													},
													map[string]interface{}{
														"secretRef": map[string]interface{}{
															"name": "not-my-secret",
														},
													},
													map[string]interface{}{
														"configMapRef": map[string]interface{}{
															"name": "my-config-map2",
														},
													},
													map[string]interface{}{
														"secretRef": map[string]interface{}{
															"name": "not-my-secret2",
														},
													},
												},
											},
										},
									},
								},
							},
						},
					}
				})

				It("replaces the configMap", func() {
					o := NewDeployment(resource.Object)
					containers := o.Object().Spec.Template.Spec.InitContainers
					Expect(containers).To(HaveLen(2))
					Expect(containers[0].EnvFrom).To(HaveLen(4))
					Expect(containers[0].EnvFrom[0].ConfigMapRef.Name).To(Equal("not-my-config-map"))
					Expect(containers[0].EnvFrom[1].SecretRef.Name).To(Equal("not-my-secret"))
					Expect(containers[0].EnvFrom[2].ConfigMapRef.Name).To(Equal("my-config-map-v000"))
					Expect(containers[0].EnvFrom[3].SecretRef.Name).To(Equal("not-my-secret2"))
					Expect(containers[1].EnvFrom).To(HaveLen(4))
					Expect(containers[1].EnvFrom[0].ConfigMapRef.Name).To(Equal("not-my-config-map"))
					Expect(containers[1].EnvFrom[1].SecretRef.Name).To(Equal("not-my-secret"))
					Expect(containers[1].EnvFrom[2].ConfigMapRef.Name).To(Equal("my-config-map2-v000"))
					Expect(containers[1].EnvFrom[3].SecretRef.Name).To(Equal("not-my-secret2"))
				})
			})

			Context(".spec.initContainers.*.envFrom.*.configMapRef.name", func() {
				BeforeEach(func() {
					resource = &unstructured.Unstructured{
						Object: map[string]interface{}{
							"kind":       "Pod",
							"apiVersion": "v1",
							"spec": map[string]interface{}{
								"initContainers": []interface{}{
									map[string]interface{}{
										"envFrom": []interface{}{
											map[string]interface{}{
												"configMapRef": map[string]interface{}{
													"name": "not-my-config-map",
												},
											},
											map[string]interface{}{
												"secretRef": map[string]interface{}{
													"name": "not-my-secret",
												},
											},
											map[string]interface{}{
												"configMapRef": map[string]interface{}{
													"name": "my-config-map",
												},
											},
											map[string]interface{}{
												"secretRef": map[string]interface{}{
													"name": "not-my-secret2",
												},
											},
										},
									},
									map[string]interface{}{
										"envFrom": []interface{}{
											map[string]interface{}{
												"configMapRef": map[string]interface{}{
													"name": "not-my-config-map",
												},
											},
											map[string]interface{}{
												"secretRef": map[string]interface{}{
													"name": "not-my-secret",
												},
											},
											map[string]interface{}{
												"configMapRef": map[string]interface{}{
													"name": "my-config-map2",
												},
											},
											map[string]interface{}{
												"secretRef": map[string]interface{}{
													"name": "not-my-secret2",
												},
											},
										},
									},
								},
							},
						},
					}
				})

				It("replaces the configMap", func() {
					o := NewPod(resource.Object)
					containers := o.Object().Spec.InitContainers
					Expect(containers).To(HaveLen(2))
					Expect(containers[0].EnvFrom).To(HaveLen(4))
					Expect(containers[0].EnvFrom[0].ConfigMapRef.Name).To(Equal("not-my-config-map"))
					Expect(containers[0].EnvFrom[1].SecretRef.Name).To(Equal("not-my-secret"))
					Expect(containers[0].EnvFrom[2].ConfigMapRef.Name).To(Equal("my-config-map-v000"))
					Expect(containers[0].EnvFrom[3].SecretRef.Name).To(Equal("not-my-secret2"))
					Expect(containers[1].EnvFrom).To(HaveLen(4))
					Expect(containers[1].EnvFrom[0].ConfigMapRef.Name).To(Equal("not-my-config-map"))
					Expect(containers[1].EnvFrom[1].SecretRef.Name).To(Equal("not-my-secret"))
					Expect(containers[1].EnvFrom[2].ConfigMapRef.Name).To(Equal("my-config-map2-v000"))
					Expect(containers[1].EnvFrom[3].SecretRef.Name).To(Equal("not-my-secret2"))
				})
			})
		})

		Describe("kubernetes/secret", func() {
			Context(".spec.template.spec.volumes.*.secret.secretName", func() {
				BeforeEach(func() {
					resource = &unstructured.Unstructured{
						Object: map[string]interface{}{
							"kind":       "Deployment",
							"apiVersion": "apps/v1",
							"spec": map[string]interface{}{
								"template": map[string]interface{}{
									"spec": map[string]interface{}{
										"volumes": []interface{}{
											map[string]interface{}{
												"secret": map[string]interface{}{
													"secretName": "not-my-secret",
												},
											},
											map[string]interface{}{
												"secret": map[string]interface{}{
													"secretName": "my-secret",
												},
											},
										},
									},
								},
							},
						},
					}
				})

				It("replaces the secret", func() {
					o := NewDeployment(resource.Object)
					volumes := o.Object().Spec.Template.Spec.Volumes
					Expect(volumes).To(HaveLen(2))
					Expect(volumes[0].VolumeSource.Secret.SecretName).To(Equal("not-my-secret"))
					Expect(volumes[1].VolumeSource.Secret.SecretName).To(Equal("my-secret-v000"))
				})
			})

			Context(".spec.volumes.*.secret.secretName", func() {
				BeforeEach(func() {
					resource = &unstructured.Unstructured{
						Object: map[string]interface{}{
							"kind":       "Pod",
							"apiVersion": "v1",
							"spec": map[string]interface{}{
								"volumes": []interface{}{
									map[string]interface{}{
										"secret": map[string]interface{}{
											"secretName": "not-my-secret",
										},
									},
									map[string]interface{}{
										"secret": map[string]interface{}{
											"secretName": "my-secret",
										},
									},
								},
							},
						},
					}
				})

				It("replaces the secret", func() {
					o := NewPod(resource.Object)
					volumes := o.Object().Spec.Volumes
					Expect(volumes).To(HaveLen(2))
					Expect(volumes[0].VolumeSource.Secret.SecretName).To(Equal("not-my-secret"))
					Expect(volumes[1].VolumeSource.Secret.SecretName).To(Equal("my-secret-v000"))
				})
			})

			Context(".spec.template.spec.volumes.*.projected.sources.*.secret.name", func() {
				BeforeEach(func() {
					resource = &unstructured.Unstructured{
						Object: map[string]interface{}{
							"kind":       "Deployment",
							"apiVersion": "apps/v1",
							"spec": map[string]interface{}{
								"template": map[string]interface{}{
									"spec": map[string]interface{}{
										"volumes": []interface{}{
											map[string]interface{}{
												"projected": map[string]interface{}{
													"sources": []interface{}{
														map[string]interface{}{
															"configMap": map[string]interface{}{
																"name": "not-my-config-map",
															},
														},
														map[string]interface{}{
															"secret": map[string]interface{}{
																"name": "my-secret",
															},
														},
														map[string]interface{}{
															"configMap": map[string]interface{}{
																"name": "not-my-config-map2",
															},
														},
														map[string]interface{}{
															"secret": map[string]interface{}{
																"name": "not-my-secret2",
															},
														},
													},
												},
											},
											map[string]interface{}{
												"projected": map[string]interface{}{
													"sources": []interface{}{
														map[string]interface{}{
															"configMap": map[string]interface{}{
																"name": "not-my-config-map",
															},
														},
														map[string]interface{}{
															"secret": map[string]interface{}{
																"name": "not-my-secret",
															},
														},
														map[string]interface{}{
															"configMap": map[string]interface{}{
																"name": "not-my-config-map2",
															},
														},
														map[string]interface{}{
															"secret": map[string]interface{}{
																"name": "my-secret2",
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					}
				})

				It("replaces the secret", func() {
					o := NewDeployment(resource.Object)
					volumes := o.Object().Spec.Template.Spec.Volumes
					Expect(volumes).To(HaveLen(2))
					Expect(volumes[0].VolumeSource.Projected.Sources).To(HaveLen(4))
					Expect(volumes[0].VolumeSource.Projected.Sources[0].ConfigMap.Name).To(Equal("not-my-config-map"))
					Expect(volumes[0].VolumeSource.Projected.Sources[1].Secret.Name).To(Equal("my-secret-v000"))
					Expect(volumes[0].VolumeSource.Projected.Sources[2].ConfigMap.Name).To(Equal("not-my-config-map2"))
					Expect(volumes[0].VolumeSource.Projected.Sources[3].Secret.Name).To(Equal("not-my-secret2"))
					Expect(volumes[1].VolumeSource.Projected.Sources).To(HaveLen(4))
					Expect(volumes[1].VolumeSource.Projected.Sources[0].ConfigMap.Name).To(Equal("not-my-config-map"))
					Expect(volumes[1].VolumeSource.Projected.Sources[1].Secret.Name).To(Equal("not-my-secret"))
					Expect(volumes[1].VolumeSource.Projected.Sources[2].ConfigMap.Name).To(Equal("not-my-config-map2"))
					Expect(volumes[1].VolumeSource.Projected.Sources[3].Secret.Name).To(Equal("my-secret2-v000"))
				})
			})

			Context(".spec.volumes.*.projected.sources.*.secret.name", func() {
				BeforeEach(func() {
					resource = &unstructured.Unstructured{
						Object: map[string]interface{}{
							"kind":       "Pod",
							"apiVersion": "v1",
							"spec": map[string]interface{}{
								"volumes": []interface{}{
									map[string]interface{}{
										"projected": map[string]interface{}{
											"sources": []interface{}{
												map[string]interface{}{
													"configMap": map[string]interface{}{
														"name": "not-my-config-map",
													},
												},
												map[string]interface{}{
													"secret": map[string]interface{}{
														"name": "my-secret",
													},
												},
												map[string]interface{}{
													"configMap": map[string]interface{}{
														"name": "not-my-config-map2",
													},
												},
												map[string]interface{}{
													"secret": map[string]interface{}{
														"name": "not-my-secret2",
													},
												},
											},
										},
									},
									map[string]interface{}{
										"projected": map[string]interface{}{
											"sources": []interface{}{
												map[string]interface{}{
													"configMap": map[string]interface{}{
														"name": "not-my-config-map",
													},
												},
												map[string]interface{}{
													"secret": map[string]interface{}{
														"name": "not-my-secret",
													},
												},
												map[string]interface{}{
													"configMap": map[string]interface{}{
														"name": "not-my-config-map2",
													},
												},
												map[string]interface{}{
													"secret": map[string]interface{}{
														"name": "my-secret2",
													},
												},
											},
										},
									},
								},
							},
						},
					}
				})

				It("replaces the secret", func() {
					o := NewPod(resource.Object)
					volumes := o.Object().Spec.Volumes
					Expect(volumes).To(HaveLen(2))
					Expect(volumes[0].VolumeSource.Projected.Sources).To(HaveLen(4))
					Expect(volumes[0].VolumeSource.Projected.Sources[0].ConfigMap.Name).To(Equal("not-my-config-map"))
					Expect(volumes[0].VolumeSource.Projected.Sources[1].Secret.Name).To(Equal("my-secret-v000"))
					Expect(volumes[0].VolumeSource.Projected.Sources[2].ConfigMap.Name).To(Equal("not-my-config-map2"))
					Expect(volumes[0].VolumeSource.Projected.Sources[3].Secret.Name).To(Equal("not-my-secret2"))
					Expect(volumes[1].VolumeSource.Projected.Sources).To(HaveLen(4))
					Expect(volumes[1].VolumeSource.Projected.Sources[0].ConfigMap.Name).To(Equal("not-my-config-map"))
					Expect(volumes[1].VolumeSource.Projected.Sources[1].Secret.Name).To(Equal("not-my-secret"))
					Expect(volumes[1].VolumeSource.Projected.Sources[2].ConfigMap.Name).To(Equal("not-my-config-map2"))
					Expect(volumes[1].VolumeSource.Projected.Sources[3].Secret.Name).To(Equal("my-secret2-v000"))
				})
			})

			Context(".spec.template.spec.containers.*.env.*.valueFrom.secretKeyRef.name", func() {
				BeforeEach(func() {
					resource = &unstructured.Unstructured{
						Object: map[string]interface{}{
							"kind":       "Deployment",
							"apiVersion": "apps/v1",
							"spec": map[string]interface{}{
								"template": map[string]interface{}{
									"spec": map[string]interface{}{
										"containers": []interface{}{
											map[string]interface{}{
												"env": []interface{}{
													map[string]interface{}{
														"valueFrom": map[string]interface{}{
															"configMapKeyRef": map[string]interface{}{
																"name": "not-my-config-map",
															},
														},
													},
													map[string]interface{}{
														"valueFrom": map[string]interface{}{
															"secretKeyRef": map[string]interface{}{
																"name": "my-secret",
															},
														},
													},
													map[string]interface{}{
														"valueFrom": map[string]interface{}{
															"configMapKeyRef": map[string]interface{}{
																"name": "not-my-config-map2",
															},
														},
													},
													map[string]interface{}{
														"valueFrom": map[string]interface{}{
															"secretKeyRef": map[string]interface{}{
																"name": "not-my-secret2",
															},
														},
													},
												},
											},
											map[string]interface{}{
												"env": []interface{}{
													map[string]interface{}{
														"valueFrom": map[string]interface{}{
															"configMapKeyRef": map[string]interface{}{
																"name": "not-my-config-map",
															},
														},
													},
													map[string]interface{}{
														"valueFrom": map[string]interface{}{
															"secretKeyRef": map[string]interface{}{
																"name": "not-my-secret",
															},
														},
													},
													map[string]interface{}{
														"valueFrom": map[string]interface{}{
															"configMapKeyRef": map[string]interface{}{
																"name": "not-my-config-map2",
															},
														},
													},
													map[string]interface{}{
														"valueFrom": map[string]interface{}{
															"secretKeyRef": map[string]interface{}{
																"name": "my-secret2",
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					}
				})

				It("replaces the secret", func() {
					o := NewDeployment(resource.Object)
					containers := o.Object().Spec.Template.Spec.Containers
					Expect(containers).To(HaveLen(2))
					Expect(containers[0].Env).To(HaveLen(4))
					Expect(containers[0].Env[0].ValueFrom.ConfigMapKeyRef.Name).To(Equal("not-my-config-map"))
					Expect(containers[0].Env[1].ValueFrom.SecretKeyRef.Name).To(Equal("my-secret-v000"))
					Expect(containers[0].Env[2].ValueFrom.ConfigMapKeyRef.Name).To(Equal("not-my-config-map2"))
					Expect(containers[0].Env[3].ValueFrom.SecretKeyRef.Name).To(Equal("not-my-secret2"))
					Expect(containers[1].Env).To(HaveLen(4))
					Expect(containers[1].Env[0].ValueFrom.ConfigMapKeyRef.Name).To(Equal("not-my-config-map"))
					Expect(containers[1].Env[1].ValueFrom.SecretKeyRef.Name).To(Equal("not-my-secret"))
					Expect(containers[1].Env[2].ValueFrom.ConfigMapKeyRef.Name).To(Equal("not-my-config-map2"))
					Expect(containers[1].Env[3].ValueFrom.SecretKeyRef.Name).To(Equal("my-secret2-v000"))
				})
			})

			Context(".spec.containers.*.env.*.valueFrom.secretKeyRef.name", func() {
				BeforeEach(func() {
					resource = &unstructured.Unstructured{
						Object: map[string]interface{}{
							"kind":       "Pod",
							"apiVersion": "v1",
							"spec": map[string]interface{}{
								"containers": []interface{}{
									map[string]interface{}{
										"env": []interface{}{
											map[string]interface{}{
												"valueFrom": map[string]interface{}{
													"configMapKeyRef": map[string]interface{}{
														"name": "not-my-config-map",
													},
												},
											},
											map[string]interface{}{
												"valueFrom": map[string]interface{}{
													"secretKeyRef": map[string]interface{}{
														"name": "my-secret",
													},
												},
											},
											map[string]interface{}{
												"valueFrom": map[string]interface{}{
													"configMapKeyRef": map[string]interface{}{
														"name": "not-my-config-map2",
													},
												},
											},
											map[string]interface{}{
												"valueFrom": map[string]interface{}{
													"secretKeyRef": map[string]interface{}{
														"name": "not-my-secret2",
													},
												},
											},
										},
									},
									map[string]interface{}{
										"env": []interface{}{
											map[string]interface{}{
												"valueFrom": map[string]interface{}{
													"configMapKeyRef": map[string]interface{}{
														"name": "not-my-config-map",
													},
												},
											},
											map[string]interface{}{
												"valueFrom": map[string]interface{}{
													"secretKeyRef": map[string]interface{}{
														"name": "not-my-secret",
													},
												},
											},
											map[string]interface{}{
												"valueFrom": map[string]interface{}{
													"configMapKeyRef": map[string]interface{}{
														"name": "not-my-config-map2",
													},
												},
											},
											map[string]interface{}{
												"valueFrom": map[string]interface{}{
													"secretKeyRef": map[string]interface{}{
														"name": "my-secret2",
													},
												},
											},
										},
									},
								},
							},
						},
					}
				})

				It("replaces the secret", func() {
					o := NewPod(resource.Object)
					containers := o.Object().Spec.Containers
					Expect(containers).To(HaveLen(2))
					Expect(containers[0].Env).To(HaveLen(4))
					Expect(containers[0].Env[0].ValueFrom.ConfigMapKeyRef.Name).To(Equal("not-my-config-map"))
					Expect(containers[0].Env[1].ValueFrom.SecretKeyRef.Name).To(Equal("my-secret-v000"))
					Expect(containers[0].Env[2].ValueFrom.ConfigMapKeyRef.Name).To(Equal("not-my-config-map2"))
					Expect(containers[0].Env[3].ValueFrom.SecretKeyRef.Name).To(Equal("not-my-secret2"))
					Expect(containers[1].Env).To(HaveLen(4))
					Expect(containers[1].Env[0].ValueFrom.ConfigMapKeyRef.Name).To(Equal("not-my-config-map"))
					Expect(containers[1].Env[1].ValueFrom.SecretKeyRef.Name).To(Equal("not-my-secret"))
					Expect(containers[1].Env[2].ValueFrom.ConfigMapKeyRef.Name).To(Equal("not-my-config-map2"))
					Expect(containers[1].Env[3].ValueFrom.SecretKeyRef.Name).To(Equal("my-secret2-v000"))
				})
			})

			Context(".spec.template.spec.initContainers.*.env.*.valueFrom.secretKeyRef.name", func() {
				BeforeEach(func() {
					resource = &unstructured.Unstructured{
						Object: map[string]interface{}{
							"kind":       "Deployment",
							"apiVersion": "apps/v1",
							"spec": map[string]interface{}{
								"template": map[string]interface{}{
									"spec": map[string]interface{}{
										"initContainers": []interface{}{
											map[string]interface{}{
												"env": []interface{}{
													map[string]interface{}{
														"valueFrom": map[string]interface{}{
															"configMapKeyRef": map[string]interface{}{
																"name": "not-my-config-map",
															},
														},
													},
													map[string]interface{}{
														"valueFrom": map[string]interface{}{
															"secretKeyRef": map[string]interface{}{
																"name": "my-secret",
															},
														},
													},
													map[string]interface{}{
														"valueFrom": map[string]interface{}{
															"configMapKeyRef": map[string]interface{}{
																"name": "not-my-config-map2",
															},
														},
													},
													map[string]interface{}{
														"valueFrom": map[string]interface{}{
															"secretKeyRef": map[string]interface{}{
																"name": "not-my-secret2",
															},
														},
													},
												},
											},
											map[string]interface{}{
												"env": []interface{}{
													map[string]interface{}{
														"valueFrom": map[string]interface{}{
															"configMapKeyRef": map[string]interface{}{
																"name": "not-my-config-map",
															},
														},
													},
													map[string]interface{}{
														"valueFrom": map[string]interface{}{
															"secretKeyRef": map[string]interface{}{
																"name": "not-my-secret",
															},
														},
													},
													map[string]interface{}{
														"valueFrom": map[string]interface{}{
															"configMapKeyRef": map[string]interface{}{
																"name": "not-my-config-map2",
															},
														},
													},
													map[string]interface{}{
														"valueFrom": map[string]interface{}{
															"secretKeyRef": map[string]interface{}{
																"name": "my-secret2",
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					}
				})

				It("replaces the secret", func() {
					o := NewDeployment(resource.Object)
					containers := o.Object().Spec.Template.Spec.InitContainers
					Expect(containers).To(HaveLen(2))
					Expect(containers[0].Env).To(HaveLen(4))
					Expect(containers[0].Env[0].ValueFrom.ConfigMapKeyRef.Name).To(Equal("not-my-config-map"))
					Expect(containers[0].Env[1].ValueFrom.SecretKeyRef.Name).To(Equal("my-secret-v000"))
					Expect(containers[0].Env[2].ValueFrom.ConfigMapKeyRef.Name).To(Equal("not-my-config-map2"))
					Expect(containers[0].Env[3].ValueFrom.SecretKeyRef.Name).To(Equal("not-my-secret2"))
					Expect(containers[1].Env).To(HaveLen(4))
					Expect(containers[1].Env[0].ValueFrom.ConfigMapKeyRef.Name).To(Equal("not-my-config-map"))
					Expect(containers[1].Env[1].ValueFrom.SecretKeyRef.Name).To(Equal("not-my-secret"))
					Expect(containers[1].Env[2].ValueFrom.ConfigMapKeyRef.Name).To(Equal("not-my-config-map2"))
					Expect(containers[1].Env[3].ValueFrom.SecretKeyRef.Name).To(Equal("my-secret2-v000"))
				})
			})

			Context(".spec.initContainers.*.env.*.valueFrom.secretKeyRef.name", func() {
				BeforeEach(func() {
					resource = &unstructured.Unstructured{
						Object: map[string]interface{}{
							"kind":       "Pod",
							"apiVersion": "v1",
							"spec": map[string]interface{}{
								"initContainers": []interface{}{
									map[string]interface{}{
										"env": []interface{}{
											map[string]interface{}{
												"valueFrom": map[string]interface{}{
													"configMapKeyRef": map[string]interface{}{
														"name": "not-my-config-map",
													},
												},
											},
											map[string]interface{}{
												"valueFrom": map[string]interface{}{
													"secretKeyRef": map[string]interface{}{
														"name": "my-secret",
													},
												},
											},
											map[string]interface{}{
												"valueFrom": map[string]interface{}{
													"configMapKeyRef": map[string]interface{}{
														"name": "not-my-config-map2",
													},
												},
											},
											map[string]interface{}{
												"valueFrom": map[string]interface{}{
													"secretKeyRef": map[string]interface{}{
														"name": "not-my-secret2",
													},
												},
											},
										},
									},
									map[string]interface{}{
										"env": []interface{}{
											map[string]interface{}{
												"valueFrom": map[string]interface{}{
													"configMapKeyRef": map[string]interface{}{
														"name": "not-my-config-map",
													},
												},
											},
											map[string]interface{}{
												"valueFrom": map[string]interface{}{
													"secretKeyRef": map[string]interface{}{
														"name": "not-my-secret",
													},
												},
											},
											map[string]interface{}{
												"valueFrom": map[string]interface{}{
													"configMapKeyRef": map[string]interface{}{
														"name": "not-my-config-map2",
													},
												},
											},
											map[string]interface{}{
												"valueFrom": map[string]interface{}{
													"secretKeyRef": map[string]interface{}{
														"name": "my-secret2",
													},
												},
											},
										},
									},
								},
							},
						},
					}
				})

				It("replaces the secret", func() {
					o := NewPod(resource.Object)
					containers := o.Object().Spec.InitContainers
					Expect(containers).To(HaveLen(2))
					Expect(containers[0].Env).To(HaveLen(4))
					Expect(containers[0].Env[0].ValueFrom.ConfigMapKeyRef.Name).To(Equal("not-my-config-map"))
					Expect(containers[0].Env[1].ValueFrom.SecretKeyRef.Name).To(Equal("my-secret-v000"))
					Expect(containers[0].Env[2].ValueFrom.ConfigMapKeyRef.Name).To(Equal("not-my-config-map2"))
					Expect(containers[0].Env[3].ValueFrom.SecretKeyRef.Name).To(Equal("not-my-secret2"))
					Expect(containers[1].Env).To(HaveLen(4))
					Expect(containers[1].Env[0].ValueFrom.ConfigMapKeyRef.Name).To(Equal("not-my-config-map"))
					Expect(containers[1].Env[1].ValueFrom.SecretKeyRef.Name).To(Equal("not-my-secret"))
					Expect(containers[1].Env[2].ValueFrom.ConfigMapKeyRef.Name).To(Equal("not-my-config-map2"))
					Expect(containers[1].Env[3].ValueFrom.SecretKeyRef.Name).To(Equal("my-secret2-v000"))
				})
			})

			Context(".spec.template.spec.containers.*.envFrom.*.secretRef.name", func() {
				BeforeEach(func() {
					resource = &unstructured.Unstructured{
						Object: map[string]interface{}{
							"kind":       "Deployment",
							"apiVersion": "apps/v1",
							"spec": map[string]interface{}{
								"template": map[string]interface{}{
									"spec": map[string]interface{}{
										"containers": []interface{}{
											map[string]interface{}{
												"envFrom": []interface{}{
													map[string]interface{}{
														"configMapRef": map[string]interface{}{
															"name": "not-my-config-map",
														},
													},
													map[string]interface{}{
														"secretRef": map[string]interface{}{
															"name": "my-secret",
														},
													},
													map[string]interface{}{
														"configMapRef": map[string]interface{}{
															"name": "not-my-config-map2",
														},
													},
													map[string]interface{}{
														"secretRef": map[string]interface{}{
															"name": "not-my-secret2",
														},
													},
												},
											},
											map[string]interface{}{
												"envFrom": []interface{}{
													map[string]interface{}{
														"configMapRef": map[string]interface{}{
															"name": "not-my-config-map",
														},
													},
													map[string]interface{}{
														"secretRef": map[string]interface{}{
															"name": "not-my-secret",
														},
													},
													map[string]interface{}{
														"configMapRef": map[string]interface{}{
															"name": "not-my-config-map2",
														},
													},
													map[string]interface{}{
														"secretRef": map[string]interface{}{
															"name": "my-secret2",
														},
													},
												},
											},
										},
									},
								},
							},
						},
					}
				})

				It("replaces the secret", func() {
					o := NewDeployment(resource.Object)
					containers := o.Object().Spec.Template.Spec.Containers
					Expect(containers).To(HaveLen(2))
					Expect(containers[0].EnvFrom).To(HaveLen(4))
					Expect(containers[0].EnvFrom[0].ConfigMapRef.Name).To(Equal("not-my-config-map"))
					Expect(containers[0].EnvFrom[1].SecretRef.Name).To(Equal("my-secret-v000"))
					Expect(containers[0].EnvFrom[2].ConfigMapRef.Name).To(Equal("not-my-config-map2"))
					Expect(containers[0].EnvFrom[3].SecretRef.Name).To(Equal("not-my-secret2"))
					Expect(containers[1].EnvFrom).To(HaveLen(4))
					Expect(containers[1].EnvFrom[0].ConfigMapRef.Name).To(Equal("not-my-config-map"))
					Expect(containers[1].EnvFrom[1].SecretRef.Name).To(Equal("not-my-secret"))
					Expect(containers[1].EnvFrom[2].ConfigMapRef.Name).To(Equal("not-my-config-map2"))
					Expect(containers[1].EnvFrom[3].SecretRef.Name).To(Equal("my-secret2-v000"))
				})
			})

			Context(".spec.containers.*.envFrom.*.secretRef.name", func() {
				BeforeEach(func() {
					resource = &unstructured.Unstructured{
						Object: map[string]interface{}{
							"kind":       "Pod",
							"apiVersion": "v1",
							"spec": map[string]interface{}{
								"containers": []interface{}{
									map[string]interface{}{
										"envFrom": []interface{}{
											map[string]interface{}{
												"configMapRef": map[string]interface{}{
													"name": "not-my-config-map",
												},
											},
											map[string]interface{}{
												"secretRef": map[string]interface{}{
													"name": "my-secret",
												},
											},
											map[string]interface{}{
												"configMapRef": map[string]interface{}{
													"name": "not-my-config-map2",
												},
											},
											map[string]interface{}{
												"secretRef": map[string]interface{}{
													"name": "not-my-secret2",
												},
											},
										},
									},
									map[string]interface{}{
										"envFrom": []interface{}{
											map[string]interface{}{
												"configMapRef": map[string]interface{}{
													"name": "not-my-config-map",
												},
											},
											map[string]interface{}{
												"secretRef": map[string]interface{}{
													"name": "not-my-secret",
												},
											},
											map[string]interface{}{
												"configMapRef": map[string]interface{}{
													"name": "not-my-config-map2",
												},
											},
											map[string]interface{}{
												"secretRef": map[string]interface{}{
													"name": "my-secret2",
												},
											},
										},
									},
								},
							},
						},
					}
				})

				It("replaces the secret", func() {
					o := NewPod(resource.Object)
					containers := o.Object().Spec.Containers
					Expect(containers).To(HaveLen(2))
					Expect(containers[0].EnvFrom).To(HaveLen(4))
					Expect(containers[0].EnvFrom[0].ConfigMapRef.Name).To(Equal("not-my-config-map"))
					Expect(containers[0].EnvFrom[1].SecretRef.Name).To(Equal("my-secret-v000"))
					Expect(containers[0].EnvFrom[2].ConfigMapRef.Name).To(Equal("not-my-config-map2"))
					Expect(containers[0].EnvFrom[3].SecretRef.Name).To(Equal("not-my-secret2"))
					Expect(containers[1].EnvFrom).To(HaveLen(4))
					Expect(containers[1].EnvFrom[0].ConfigMapRef.Name).To(Equal("not-my-config-map"))
					Expect(containers[1].EnvFrom[1].SecretRef.Name).To(Equal("not-my-secret"))
					Expect(containers[1].EnvFrom[2].ConfigMapRef.Name).To(Equal("not-my-config-map2"))
					Expect(containers[1].EnvFrom[3].SecretRef.Name).To(Equal("my-secret2-v000"))
				})
			})

			Context(".spec.template.spec.initContainers.*.envFrom.*.configMapRef.name", func() {
				BeforeEach(func() {
					resource = &unstructured.Unstructured{
						Object: map[string]interface{}{
							"kind":       "Deployment",
							"apiVersion": "apps/v1",
							"spec": map[string]interface{}{
								"template": map[string]interface{}{
									"spec": map[string]interface{}{
										"initContainers": []interface{}{
											map[string]interface{}{
												"envFrom": []interface{}{
													map[string]interface{}{
														"configMapRef": map[string]interface{}{
															"name": "not-my-config-map",
														},
													},
													map[string]interface{}{
														"secretRef": map[string]interface{}{
															"name": "my-secret",
														},
													},
													map[string]interface{}{
														"configMapRef": map[string]interface{}{
															"name": "not-my-config-map2",
														},
													},
													map[string]interface{}{
														"secretRef": map[string]interface{}{
															"name": "not-my-secret2",
														},
													},
												},
											},
											map[string]interface{}{
												"envFrom": []interface{}{
													map[string]interface{}{
														"configMapRef": map[string]interface{}{
															"name": "not-my-config-map",
														},
													},
													map[string]interface{}{
														"secretRef": map[string]interface{}{
															"name": "not-my-secret",
														},
													},
													map[string]interface{}{
														"configMapRef": map[string]interface{}{
															"name": "not-my-config-map2",
														},
													},
													map[string]interface{}{
														"secretRef": map[string]interface{}{
															"name": "my-secret2",
														},
													},
												},
											},
										},
									},
								},
							},
						},
					}
				})

				It("replaces the secret", func() {
					o := NewDeployment(resource.Object)
					containers := o.Object().Spec.Template.Spec.InitContainers
					Expect(containers).To(HaveLen(2))
					Expect(containers[0].EnvFrom).To(HaveLen(4))
					Expect(containers[0].EnvFrom[0].ConfigMapRef.Name).To(Equal("not-my-config-map"))
					Expect(containers[0].EnvFrom[1].SecretRef.Name).To(Equal("my-secret-v000"))
					Expect(containers[0].EnvFrom[2].ConfigMapRef.Name).To(Equal("not-my-config-map2"))
					Expect(containers[0].EnvFrom[3].SecretRef.Name).To(Equal("not-my-secret2"))
					Expect(containers[1].EnvFrom).To(HaveLen(4))
					Expect(containers[1].EnvFrom[0].ConfigMapRef.Name).To(Equal("not-my-config-map"))
					Expect(containers[1].EnvFrom[1].SecretRef.Name).To(Equal("not-my-secret"))
					Expect(containers[1].EnvFrom[2].ConfigMapRef.Name).To(Equal("not-my-config-map2"))
					Expect(containers[1].EnvFrom[3].SecretRef.Name).To(Equal("my-secret2-v000"))
				})
			})

			Context(".spec.initContainers.*.envFrom.*.secretRef.name", func() {
				BeforeEach(func() {
					resource = &unstructured.Unstructured{
						Object: map[string]interface{}{
							"kind":       "Pod",
							"apiVersion": "v1",
							"spec": map[string]interface{}{
								"initContainers": []interface{}{
									map[string]interface{}{
										"envFrom": []interface{}{
											map[string]interface{}{
												"configMapRef": map[string]interface{}{
													"name": "not-my-config-map",
												},
											},
											map[string]interface{}{
												"secretRef": map[string]interface{}{
													"name": "my-secret",
												},
											},
											map[string]interface{}{
												"configMapRef": map[string]interface{}{
													"name": "not-my-config-map2",
												},
											},
											map[string]interface{}{
												"secretRef": map[string]interface{}{
													"name": "not-my-secret2",
												},
											},
										},
									},
									map[string]interface{}{
										"envFrom": []interface{}{
											map[string]interface{}{
												"configMapRef": map[string]interface{}{
													"name": "not-my-config-map",
												},
											},
											map[string]interface{}{
												"secretRef": map[string]interface{}{
													"name": "not-my-secret",
												},
											},
											map[string]interface{}{
												"configMapRef": map[string]interface{}{
													"name": "not-my-config-map2",
												},
											},
											map[string]interface{}{
												"secretRef": map[string]interface{}{
													"name": "my-secret2",
												},
											},
										},
									},
								},
							},
						},
					}
				})

				It("replaces the secret", func() {
					o := NewPod(resource.Object)
					containers := o.Object().Spec.InitContainers
					Expect(containers).To(HaveLen(2))
					Expect(containers[0].EnvFrom).To(HaveLen(4))
					Expect(containers[0].EnvFrom[0].ConfigMapRef.Name).To(Equal("not-my-config-map"))
					Expect(containers[0].EnvFrom[1].SecretRef.Name).To(Equal("my-secret-v000"))
					Expect(containers[0].EnvFrom[2].ConfigMapRef.Name).To(Equal("not-my-config-map2"))
					Expect(containers[0].EnvFrom[3].SecretRef.Name).To(Equal("not-my-secret2"))
					Expect(containers[1].EnvFrom).To(HaveLen(4))
					Expect(containers[1].EnvFrom[0].ConfigMapRef.Name).To(Equal("not-my-config-map"))
					Expect(containers[1].EnvFrom[1].SecretRef.Name).To(Equal("not-my-secret"))
					Expect(containers[1].EnvFrom[2].ConfigMapRef.Name).To(Equal("not-my-config-map2"))
					Expect(containers[1].EnvFrom[3].SecretRef.Name).To(Equal("my-secret2-v000"))
				})
			})
		})

		Describe("kubernetes/deployment", func() {
			Context("when the manifest kind is not HorizontalPodAutoscaler", func() {
				BeforeEach(func() {
					resource = &unstructured.Unstructured{
						Object: map[string]interface{}{
							"kind":       "NotHorizontalPodAutoscaler",
							"apiVersion": "autoscaling/v2beta2",
							"spec": map[string]interface{}{
								"scaleTargetRef": map[string]interface{}{
									"kind": "deployment",
									"name": "my-deployment",
								},
							},
						},
					}
				})

				It("ignores the manifest", func() {
					o := NewHorizontalPodAutoscaler(resource.Object)
					name := o.Object().Spec.ScaleTargetRef.Name
					Expect(name).To(Equal("my-deployment"))
				})
			})

			Context("when the .spec.scaleTargetRef.kind is not a deployment", func() {
				BeforeEach(func() {
					resource = &unstructured.Unstructured{
						Object: map[string]interface{}{
							"kind":       "HorizontalPodAutoscaler",
							"apiVersion": "autoscaling/v2beta2",
							"spec": map[string]interface{}{
								"scaleTargetRef": map[string]interface{}{
									"kind": "fake",
									"name": "my-replicaset",
								},
							},
						},
					}
				})

				It("does not replace the reference", func() {
					o := NewHorizontalPodAutoscaler(resource.Object)
					name := o.Object().Spec.ScaleTargetRef.Name
					Expect(name).To(Equal("my-replicaset"))
				})
			})

			Context("when the .spec.scaleTargetRef.kind is a deployment", func() {
				BeforeEach(func() {
					resource = &unstructured.Unstructured{
						Object: map[string]interface{}{
							"kind":       "HorizontalPodAutoscaler",
							"apiVersion": "autoscaling/v2beta2",
							"spec": map[string]interface{}{
								"scaleTargetRef": map[string]interface{}{
									"kind": "Deployment",
									"name": "my-deployment",
								},
							},
						},
					}
				})

				It("does replace the reference", func() {
					o := NewHorizontalPodAutoscaler(resource.Object)
					name := o.Object().Spec.ScaleTargetRef.Name
					Expect(name).To(Equal("my-deployment-v000"))
				})
			})
		})

		Describe("kubernetes/replicaSet", func() {
			Context("when the manifest kind is not HorizontalPodAutoscaler", func() {
				BeforeEach(func() {
					resource = &unstructured.Unstructured{
						Object: map[string]interface{}{
							"kind":       "NotHorizontalPodAutoscaler",
							"apiVersion": "autoscaling/v2beta2",
							"spec": map[string]interface{}{
								"scaleTargetRef": map[string]interface{}{
									"kind": "replicaSet",
									"name": "my-replicaset",
								},
							},
						},
					}
				})

				It("ignores the manifest", func() {
					o := NewHorizontalPodAutoscaler(resource.Object)
					name := o.Object().Spec.ScaleTargetRef.Name
					Expect(name).To(Equal("my-replicaset"))
				})
			})

			Context("when the .spec.scaleTargetRef.kind is not a replicaSet", func() {
				BeforeEach(func() {
					resource = &unstructured.Unstructured{
						Object: map[string]interface{}{
							"kind":       "HorizontalPodAutoscaler",
							"apiVersion": "autoscaling/v2beta2",
							"spec": map[string]interface{}{
								"scaleTargetRef": map[string]interface{}{
									"kind": "fake",
									"name": "my-deployment",
								},
							},
						},
					}
				})

				It("does not replace the reference", func() {
					o := NewHorizontalPodAutoscaler(resource.Object)
					name := o.Object().Spec.ScaleTargetRef.Name
					Expect(name).To(Equal("my-deployment"))
				})
			})

			Context("when the .spec.scaleTargetRef.kind is a replicaSet", func() {
				BeforeEach(func() {
					resource = &unstructured.Unstructured{
						Object: map[string]interface{}{
							"kind":       "HorizontalPodAutoscaler",
							"apiVersion": "autoscaling/v2beta2",
							"spec": map[string]interface{}{
								"scaleTargetRef": map[string]interface{}{
									"kind": "ReplicaSet",
									"name": "my-replicaSet",
								},
							},
						},
					}
				})

				It("does replace the reference", func() {
					o := NewHorizontalPodAutoscaler(resource.Object)
					name := o.Object().Spec.ScaleTargetRef.Name
					Expect(name).To(Equal("my-replicaSet-v000"))
				})
			})
		})
	})

	Describe("#FindArtifacts", func() {
		var (
			resource  *unstructured.Unstructured
			artifacts []clouddriver.Artifact
		)

		BeforeEach(func() {
			resource = &unstructured.Unstructured{
				Object: map[string]interface{}{
					"kind":       "Deployment",
					"apiVersion": "apps/v1",
					"spec": map[string]interface{}{
						"template": map[string]interface{}{
							"spec": map[string]interface{}{
								"containers": []interface{}{
									map[string]interface{}{
										"name":  "test-container-name",
										"image": "gcr.io/test-project/test-container-image@some-sha",
										"env": []interface{}{
											map[string]interface{}{
												"valueFrom": map[string]interface{}{
													"configMapKeyRef": map[string]interface{}{
														"name": "my-config-map-key-ref",
													},
												},
											},
											map[string]interface{}{
												"valueFrom": map[string]interface{}{
													"secretKeyRef": map[string]interface{}{
														"name": "my-secret-key-ref",
													},
												},
											},
										},
									},
									map[string]interface{}{
										"name":  "another-test-container-name",
										"image": "gcr.io/test-project/another-test-container-image@some-sha",
										"envFrom": []interface{}{
											map[string]interface{}{
												"configMapRef": map[string]interface{}{
													"name": "my-config-map-ref",
												},
											},
											map[string]interface{}{
												"secretRef": map[string]interface{}{
													"name": "my-secret-ref",
												},
											},
										},
									},
								},
								"initContainers": []interface{}{
									map[string]interface{}{
										"name":  "test-init-container-name",
										"image": "gcr.io/test-project/test-init-container-image:some-version",
										"envFrom": []interface{}{
											map[string]interface{}{
												"configMapRef": map[string]interface{}{
													"name": "my-init-container-config-map-ref",
												},
											},
											map[string]interface{}{
												"secretRef": map[string]interface{}{
													"name": "my-init-container-secret-ref",
												},
											},
										},
									},
									map[string]interface{}{
										"name":  "another-test-init-container-name",
										"image": "gcr.io/test-project/another-test-init-container-image",
										"env": []interface{}{
											map[string]interface{}{
												"valueFrom": map[string]interface{}{
													"configMapKeyRef": map[string]interface{}{
														"name": "my-init-container-config-map-key-ref",
													},
												},
											},
											map[string]interface{}{
												"valueFrom": map[string]interface{}{
													"secretKeyRef": map[string]interface{}{
														"name": "my-init-container-secret-key-ref",
													},
												},
											},
										},
									},
								},
								"volumes": []interface{}{
									map[string]interface{}{
										"secret": map[string]interface{}{
											"secretName": "my-secret",
										},
									},
									map[string]interface{}{
										"configMap": map[string]interface{}{
											"name": "my-config-map",
										},
									},
									map[string]interface{}{
										"projected": map[string]interface{}{
											"sources": []interface{}{
												map[string]interface{}{
													"configMap": map[string]interface{}{
														"name": "my-projected-sources-config-map",
													},
												},
												map[string]interface{}{
													"secret": map[string]interface{}{
														"name": "my-projected-sources-secret",
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			}
		})

		JustBeforeEach(func() {
			artifacts = FindArtifacts(resource)
		})

		When("the manifest does not have artifacts", func() {
			BeforeEach(func() {
				resource = &unstructured.Unstructured{
					Object: map[string]interface{}{
						"kind":       "Deployment",
						"apiVersion": "apps/v1",
						"spec":       map[string]interface{}{},
					},
				}
			})

			It("returns an empty slice", func() {
				Expect(artifacts).To(HaveLen(0))
			})
		})

		When("the manifest has artifacts", func() {
			It("succeeds", func() {
				Expect(artifacts).To(HaveLen(16))
				// docker/image artifacts
				Expect(artifacts[0].Name).To(Equal("gcr.io/test-project/test-container-image"))
				Expect(artifacts[0].Reference).To(Equal("gcr.io/test-project/test-container-image@some-sha"))
				Expect(artifacts[0].Type).To(Equal(artifact.TypeDockerImage))
				Expect(artifacts[1].Name).To(Equal("gcr.io/test-project/another-test-container-image"))
				Expect(artifacts[1].Reference).To(Equal("gcr.io/test-project/another-test-container-image@some-sha"))
				Expect(artifacts[1].Type).To(Equal(artifact.TypeDockerImage))
				Expect(artifacts[2].Name).To(Equal("gcr.io/test-project/test-init-container-image"))
				Expect(artifacts[2].Reference).To(Equal("gcr.io/test-project/test-init-container-image:some-version"))
				Expect(artifacts[2].Type).To(Equal(artifact.TypeDockerImage))
				Expect(artifacts[3].Name).To(Equal("gcr.io/test-project/another-test-init-container-image"))
				Expect(artifacts[3].Reference).To(Equal("gcr.io/test-project/another-test-init-container-image"))
				Expect(artifacts[3].Type).To(Equal(artifact.TypeDockerImage))
				// kubernetes/configMap artifacts
				Expect(artifacts[4].Name).To(Equal("my-config-map"))
				Expect(artifacts[4].Reference).To(Equal("my-config-map"))
				Expect(artifacts[4].Type).To(Equal(artifact.TypeKubernetesConfigMap))
				Expect(artifacts[5].Name).To(Equal("my-projected-sources-config-map"))
				Expect(artifacts[5].Reference).To(Equal("my-projected-sources-config-map"))
				Expect(artifacts[5].Type).To(Equal(artifact.TypeKubernetesConfigMap))
				Expect(artifacts[6].Name).To(Equal("my-config-map-key-ref"))
				Expect(artifacts[6].Reference).To(Equal("my-config-map-key-ref"))
				Expect(artifacts[6].Type).To(Equal(artifact.TypeKubernetesConfigMap))
				Expect(artifacts[7].Name).To(Equal("my-init-container-config-map-key-ref"))
				Expect(artifacts[7].Reference).To(Equal("my-init-container-config-map-key-ref"))
				Expect(artifacts[7].Type).To(Equal(artifact.TypeKubernetesConfigMap))
				Expect(artifacts[8].Name).To(Equal("my-config-map-ref"))
				Expect(artifacts[8].Reference).To(Equal("my-config-map-ref"))
				Expect(artifacts[8].Type).To(Equal(artifact.TypeKubernetesConfigMap))
				Expect(artifacts[9].Name).To(Equal("my-init-container-config-map-ref"))
				Expect(artifacts[9].Reference).To(Equal("my-init-container-config-map-ref"))
				Expect(artifacts[9].Type).To(Equal(artifact.TypeKubernetesConfigMap))
				// kubernetes/secret artifacts
				Expect(artifacts[10].Name).To(Equal("my-secret"))
				Expect(artifacts[10].Reference).To(Equal("my-secret"))
				Expect(artifacts[10].Type).To(Equal(artifact.TypeKubernetesSecret))
				Expect(artifacts[11].Name).To(Equal("my-projected-sources-secret"))
				Expect(artifacts[11].Reference).To(Equal("my-projected-sources-secret"))
				Expect(artifacts[11].Type).To(Equal(artifact.TypeKubernetesSecret))
				Expect(artifacts[12].Name).To(Equal("my-secret-key-ref"))
				Expect(artifacts[12].Reference).To(Equal("my-secret-key-ref"))
				Expect(artifacts[12].Type).To(Equal(artifact.TypeKubernetesSecret))
				Expect(artifacts[13].Name).To(Equal("my-init-container-secret-key-ref"))
				Expect(artifacts[13].Reference).To(Equal("my-init-container-secret-key-ref"))
				Expect(artifacts[13].Type).To(Equal(artifact.TypeKubernetesSecret))
				Expect(artifacts[14].Name).To(Equal("my-secret-ref"))
				Expect(artifacts[14].Reference).To(Equal("my-secret-ref"))
				Expect(artifacts[14].Type).To(Equal(artifact.TypeKubernetesSecret))
				Expect(artifacts[15].Name).To(Equal("my-init-container-secret-ref"))
				Expect(artifacts[15].Reference).To(Equal("my-init-container-secret-ref"))
				Expect(artifacts[15].Type).To(Equal(artifact.TypeKubernetesSecret))
			})
		})
	})
})
