package core_test

const payloadRequestKubernetesOpsDeployManifest = `[
  {
    "deployManifest": {
      "enableTraffic": true,
      "namespaceOverride": "default",
      "optionalArtifacts": [],
      "cloudProvider": "kubernetes",
      "manifests": [
        {
          "metadata": {
            "name": "rss-site",
            "labels": {
              "app": "web"
            }
          },
          "apiVersion": "v1",
          "kind": "Pod",
          "spec": {
            "containers": [
              {
                "image": "nginx",
                "name": "front-end",
                "ports": [
                  {
                    "containerPort": 80
                  }
                ]
              },
              {
                "image": "nickchase/rss-php-nginx:vasdf1",
                "name": "rss-reader",
                "ports": [
                  {
                    "containerPort": 88
                  }
                ]
              }
            ]
          }
        }
      ],
      "trafficManagement": {
        "options": {
          "enableTraffic": true,
          "namespace": "default",
          "services": [
            "service hello-app-red-black"
          ],
          "strategy": "redblack"
        },
        "enabled": false
      },
      "moniker": {
        "app": "test"
      },
      "source": "text",
      "account": "spin-cluster-account",
      "skipExpressionEvaluation": false,
      "requiredArtifacts": []
    }
  }
]`

const payloadRequestFetchHelmArtifact = `{
  "name": "test-chart-name",
  "type": "helm/chart",
	"version": "1.0.0"
}`

const payloadRequestFetchBase64ArtifactBadReference = `{
  "type": "embedded/base64",
	"reference": "not-base-64!"
}`

const payloadRequestFetchBase64Artifact = `{
  "type": "embedded/base64",
	"reference": "aGVsbG93b3JsZAo="
}`

const payloadRequestFetchGCSObjetArtifact = `{
  "type": "gcs/object",
	"reference": "gs://fake-bucket/fake-path/fake-file.txt"
}`

const payloadRequestFetchGCSObjetArtifactBadReference = `{
  "type": "gcs/object",
	"reference": "not-gcs-format"
}`

const payloadRequestFetchGCSObjetArtifactBadFileGeneration = `{
  "type": "gcs/object",
	"reference": "gs://fake-bucket/fake-path/fake-file.txt#v1"
}`

const payloadRequestFetchGCSObjetArtifactNotFound = `{
  "type": "gcs/object",
	"reference": "gs://fake-bucket/fake-path/not-found.txt"
}`

const payloadRequestFetchGithubFileArtifact = `{
  "type": "github/file",
	"reference": "%s/api/v3/repos/homedepot/kubernetes-engine-samples/contents/hello-app/manifests/helloweb-deployment.yaml"
}`

const payloadRequestFetchGitRepoArtifact = `{
  "type": "git/repo",
	"reference": "%s/git-repo"
}`

const payloadRequestFetchGitRepoArtifactBranch = `{
  "type": "git/repo",
  "reference": "%s/git-repo",
  "version": "test"
}`

const payloadRequestFetchGitRepoArtifactSubPath = `{
  "type": "git/repo",
  "reference": "%s/git-repo",
  "location": "kustomize"
}`

const payloadRequestFetchHTTPFileArtifact = `{
  "type": "http/file",
	"reference": "%s/hello"
}`

const payloadRequestFetchGithubFileArtifactTestBranch = `{
  "type": "github/file",
	"reference": "%s/api/v3/repos/homedepot/kubernetes-engine-samples/contents/hello-app/manifests/helloweb-deployment.yaml",
	"version": "test"
}`

const payloadRequestFetchNotImplementedArtifact = `{
  "type": "unknown/type"
}`

const payloadRequestKubernetesOpsDeleteManifest = `[
  {
    "deleteManifest": {
      "account": "spin-cluster-account"
		}
  }
]`

const payloadRequestKubernetesOpsDisableManifest = `[
  {
    "disableManifest": {
      "account": "spin-cluster-account"
		}
  }
]`

const payloadRequestKubernetesOpsEnableManifest = `[
  {
    "enableManifest": {
      "account": "spin-cluster-account"
		}
  }
]`

const payloadRequestKubernetesOpsCleanupArtifacts = `[
  {
    "cleanupArtifacts": {
      "account": "spin-cluster-account",
      "manifests": "asdf"
		}
  }
]`

const payloadRequestKubernetesOpsScaleManifest = `[
  {
    "scaleManifest": {
      "enableTraffic": true,
      "namespaceOverride": "default",
      "optionalArtifacts": [],
      "cloudProvider": "kubernetes",
      "manifests": [
        {
          "metadata": {
            "name": "rss-site",
            "labels": {
              "app": "web"
            }
          },
          "apiVersion": "v1",
          "kind": "Pod",
          "spec": {
            "containers": [
              {
                "image": "nginx",
                "name": "front-end",
                "ports": [
                  {
                    "containerPort": 80
                  }
                ]
              },
              {
                "image": "nickchase/rss-php-nginx:vasdf1",
                "name": "rss-reader",
                "ports": [
                  {
                    "containerPort": 88
                  }
                ]
              }
            ]
          }
        }
      ],
      "trafficManagement": {
        "options": {
          "enableTraffic": true,
          "namespace": "default",
          "services": [
            "service hello-app-red-black"
          ],
          "strategy": "redblack"
        },
        "enabled": false
      },
      "moniker": {
        "app": "test"
      },
      "source": "text",
      "account": "spin-cluster-account",
      "skipExpressionEvaluation": false,
      "requiredArtifacts": []
    }
  }
]`

const payloadRequestKubernetesOpsRollingRestartManifest = `[
  {
    "rollingRestartManifest": {
      "enableTraffic": true,
      "namespaceOverride": "default",
      "optionalArtifacts": [],
      "cloudProvider": "kubernetes",
      "manifests": [
        {
          "metadata": {
            "name": "rss-site",
            "labels": {
              "app": "web"
            }
          },
          "apiVersion": "v1",
          "kind": "Pod",
          "spec": {
            "containers": [
              {
                "image": "nginx",
                "name": "front-end",
                "ports": [
                  {
                    "containerPort": 80
                  }
                ]
              },
              {
                "image": "nickchase/rss-php-nginx:vasdf1",
                "name": "rss-reader",
                "ports": [
                  {
                    "containerPort": 88
                  }
                ]
              }
            ]
          }
        }
      ],
      "trafficManagement": {
        "options": {
          "enableTraffic": true,
          "namespace": "default",
          "services": [
            "service hello-app-red-black"
          ],
          "strategy": "redblack"
        },
        "enabled": false
      },
      "moniker": {
        "app": "test"
      },
			"manifestName": "deployment test-deployment",
      "source": "text",
      "account": "spin-cluster-account",
      "skipExpressionEvaluation": false,
      "requiredArtifacts": []
    }
  }
]`

const payloadRequestKubernetesOpsRunJob = `[
  {
    "runJob": {
      "enableTraffic": true,
      "namespaceOverride": "default",
      "optionalArtifacts": [],
      "cloudProvider": "kubernetes",
      "manifests": [
        {
          "metadata": {
            "name": "rss-site",
            "labels": {
              "app": "web"
            }
          },
          "apiVersion": "v1",
          "kind": "Pod",
          "spec": {
            "containers": [
              {
                "image": "nginx",
                "name": "front-end",
                "ports": [
                  {
                    "containerPort": 80
                  }
                ]
              },
              {
                "image": "nickchase/rss-php-nginx:vasdf1",
                "name": "rss-reader",
                "ports": [
                  {
                    "containerPort": 88
                  }
                ]
              }
            ]
          }
        }
      ],
      "trafficManagement": {
        "options": {
          "enableTraffic": true,
          "namespace": "default",
          "services": [
            "service hello-app-red-black"
          ],
          "strategy": "redblack"
        },
        "enabled": false
      },
      "moniker": {
        "app": "test"
      },
      "source": "text",
      "account": "spin-cluster-account",
      "skipExpressionEvaluation": false,
      "requiredArtifacts": []
    }
  }
]`

const payloadRequestKubernetesOpsUndoRolloutManifest = `[
  {
    "undoRolloutManifest": {
      "enableTraffic": true,
      "namespaceOverride": "default",
      "optionalArtifacts": [],
      "cloudProvider": "kubernetes",
      "manifests": [
        {
          "metadata": {
            "name": "rss-site",
            "labels": {
              "app": "web"
            }
          },
          "apiVersion": "v1",
          "kind": "Pod",
          "spec": {
            "containers": [
              {
                "image": "nginx",
                "name": "front-end",
                "ports": [
                  {
                    "containerPort": 80
                  }
                ]
              },
              {
                "image": "nickchase/rss-php-nginx:vasdf1",
                "name": "rss-reader",
                "ports": [
                  {
                    "containerPort": 88
                  }
                ]
              }
            ]
          }
        }
      ],
      "trafficManagement": {
        "options": {
          "enableTraffic": true,
          "namespace": "default",
          "services": [
            "service hello-app-red-black"
          ],
          "strategy": "redblack"
        },
        "enabled": false
      },
      "moniker": {
        "app": "test"
      },
			"manifestName": "deployment test-deployment",
      "source": "text",
      "account": "spin-cluster-account",
      "skipExpressionEvaluation": false,
      "requiredArtifacts": []
    }
  }
]`

const payloadRequestKubernetesOpsPatchManifest = `[
  {
    "patchManifest": {
		  "account": "spin-cluster-account"
		}
  }
]`

const payloadCredentials = `[
              {
                "accountType": "provider1",
                "cacheThreads": 0,
                "challengeDestructiveActions": false,
                "cloudProvider": "kubernetes",
                "dockerRegistries": null,
                "enabled": false,
                "environment": "provider1",
                "name": "provider1",
                "namespaces": null,
                "permissions": {
                  "READ": [
                    "gg_test"
                  ],
                  "WRITE": [
                    "gg_test"
                  ]
                },
                "primaryAccount": false,
                "providerVersion": "v2",
                "requiredGroupMembership": [],
                "skin": "v2",
                "spinnakerKindMap": null,
                "type": "kubernetes"
              },
              {
                "accountType": "provider2",
                "cacheThreads": 0,
                "challengeDestructiveActions": false,
                "cloudProvider": "kubernetes",
                "dockerRegistries": null,
                "enabled": false,
                "environment": "provider2",
                "name": "provider2",
                "namespaces": null,
                "permissions": {
                  "READ": [
                    "gg_test2"
                  ],
                  "WRITE": [
                    "gg_test2"
                  ]
                },
                "primaryAccount": false,
                "providerVersion": "v2",
                "requiredGroupMembership": [],
                "skin": "v2",
                "spinnakerKindMap": null,
                "type": "kubernetes"
              }
            ]`

const payloadArtifactCredentials = `[
            {
              "name": "helm-stable",
              "types": [
                "helm/chart"
              ]
            },
            {
              "name": "embedded-artifact",
              "types": [
                "embedded/base64"
              ]
            }
          ]`

const payloadListHelmArtifactAccountNames = `[
            "minecraft",
            "prometheus-operator"
          ]`

const payloadListHelmArtifactAccountVersions = `[
            "1.0.0",
            "1.1.0"
          ]`

const payloadApplications = `[
            {
              "attributes": {
                "name": "test-spinnaker-app1"
              },
              "clusterNames": {
                "test-account1": [
                  "test-kind1 test-name1"
                ]
              },
              "name": "test-spinnaker-app1"
            },
            {
              "attributes": {
                "name": "test-spinnaker-app2"
              },
              "clusterNames": {
                "test-account2": [
                  "test-kind2 test-name2"
                ],
                "test-account3": [
                  "test-kind3 test-name3"
                ]
              },
              "name": "test-spinnaker-app2"
            }
          ]`

const payloadApplicationsSorted = `[
            {
              "attributes": {
                "name": "a"
              },
              "clusterNames": {
                "test-account3": [
                  "test-kind3 test-name3"
                ]
              },
              "name": "a"
            },
            {
              "attributes": {
                "name": "b"
              },
              "clusterNames": {
                "test-account2": [
                  "test-kind2 test-name2"
                ]
              },
              "name": "b"
            },
            {
              "attributes": {
                "name": "c"
              },
              "clusterNames": {
                "test-account1": [
                  "test-kind1 test-name1"
                ]
              },
              "name": "c"
            }
          ]`

const payloadGetAccountCredentials = `{
            "accountType": "test-account",
            "cacheThreads": 0,
            "challengeDestructiveActions": false,
            "cloudProvider": "kubernetes",
            "dockerRegistries": null,
            "enabled": false,
            "environment": "test-account",
            "name": "test-account",
            "namespaces": null,
            "permissions": {
              "READ": null,
              "WRITE": null
            },
            "primaryAccount": false,
            "providerVersion": "v2",
            "requiredGroupMembership": [],
            "skin": "v2",
            "spinnakerKindMap": {
              "apiService": "unclassified",
              "clusterRole": "unclassified",
              "clusterRoleBinding": "unclassified",
              "configMap": "configs",
              "controllerRevision": "unclassified",
              "cronJob": "serverGroups",
              "customResourceDefinition": "unclassified",
              "daemonSet": "serverGroups",
              "deployment": "serverGroupManagers",
              "event": "unclassified",
              "horizontalpodautoscaler": "unclassified",
              "ingress": "loadBalancers",
              "job": "serverGroups",
              "limitRange": "unclassified",
              "mutatingWebhookConfiguration": "unclassified",
              "namespace": "unclassified",
              "networkPolicy": "securityGroups",
              "persistentVolume": "configs",
              "persistentVolumeClaim": "configs",
              "pod": "instances",
              "podDisruptionBudget": "unclassified",
              "podPreset": "unclassified",
              "podSecurityPolicy": "unclassified",
              "replicaSet": "serverGroups",
              "role": "unclassified",
              "roleBinding": "unclassified",
              "secret": "configs",
              "service": "loadBalancers",
              "serviceAccount": "unclassified",
              "statefulSet": "serverGroups",
              "storageClass": "unclassified",
              "validatingWebhookConfiguration": "unclassified"
            },
            "type": "kubernetes"
          }`

const payloadCredentialsExpandTrue = `[
              {
                "accountType": "provider1",
                "cacheThreads": 0,
                "challengeDestructiveActions": false,
                "cloudProvider": "kubernetes",
                "dockerRegistries": null,
                "enabled": false,
                "environment": "provider1",
                "name": "provider1",
                "namespaces": [
                  "namespace1",
                  "namespace2"
                ],
                "permissions": {
                  "READ": [
                    "gg_test"
                  ],
                  "WRITE": [
                    "gg_test"
                  ]
                },
                "primaryAccount": false,
                "providerVersion": "v2",
                "requiredGroupMembership": [],
                "skin": "v2",
                "spinnakerKindMap": {
                  "apiService": "unclassified",
                  "clusterRole": "unclassified",
                  "clusterRoleBinding": "unclassified",
                  "configMap": "configs",
                  "controllerRevision": "unclassified",
                  "cronJob": "serverGroups",
                  "customResourceDefinition": "unclassified",
                  "daemonSet": "serverGroups",
                  "deployment": "serverGroupManagers",
                  "event": "unclassified",
                  "horizontalpodautoscaler": "unclassified",
                  "ingress": "loadBalancers",
                  "job": "serverGroups",
                  "limitRange": "unclassified",
                  "mutatingWebhookConfiguration": "unclassified",
                  "namespace": "unclassified",
                  "networkPolicy": "securityGroups",
                  "persistentVolume": "configs",
                  "persistentVolumeClaim": "configs",
                  "pod": "instances",
                  "podDisruptionBudget": "unclassified",
                  "podPreset": "unclassified",
                  "podSecurityPolicy": "unclassified",
                  "replicaSet": "serverGroups",
                  "role": "unclassified",
                  "roleBinding": "unclassified",
                  "secret": "configs",
                  "service": "loadBalancers",
                  "serviceAccount": "unclassified",
                  "statefulSet": "serverGroups",
                  "storageClass": "unclassified",
                  "validatingWebhookConfiguration": "unclassified"
                },
                "type": "kubernetes"
              },
              {
                "accountType": "provider2",
                "cacheThreads": 0,
                "challengeDestructiveActions": false,
                "cloudProvider": "kubernetes",
                "dockerRegistries": null,
                "enabled": false,
                "environment": "provider2",
                "name": "provider2",
                "namespaces": [
                  "namespace1",
                  "namespace2"
                ],
                "permissions": {
                  "READ": [
                    "gg_test2"
                  ],
                  "WRITE": [
                    "gg_test2"
                  ]
                },
                "primaryAccount": false,
                "providerVersion": "v2",
                "requiredGroupMembership": [],
                "skin": "v2",
                "spinnakerKindMap": {
                  "apiService": "unclassified",
                  "clusterRole": "unclassified",
                  "clusterRoleBinding": "unclassified",
                  "configMap": "configs",
                  "controllerRevision": "unclassified",
                  "cronJob": "serverGroups",
                  "customResourceDefinition": "unclassified",
                  "daemonSet": "serverGroups",
                  "deployment": "serverGroupManagers",
                  "event": "unclassified",
                  "horizontalpodautoscaler": "unclassified",
                  "ingress": "loadBalancers",
                  "job": "serverGroups",
                  "limitRange": "unclassified",
                  "mutatingWebhookConfiguration": "unclassified",
                  "namespace": "unclassified",
                  "networkPolicy": "securityGroups",
                  "persistentVolume": "configs",
                  "persistentVolumeClaim": "configs",
                  "pod": "instances",
                  "podDisruptionBudget": "unclassified",
                  "podPreset": "unclassified",
                  "podSecurityPolicy": "unclassified",
                  "replicaSet": "serverGroups",
                  "role": "unclassified",
                  "roleBinding": "unclassified",
                  "secret": "configs",
                  "service": "loadBalancers",
                  "serviceAccount": "unclassified",
                  "statefulSet": "serverGroups",
                  "storageClass": "unclassified",
                  "validatingWebhookConfiguration": "unclassified"
                },
                "type": "kubernetes"
              }
            ]`

const payloadCredentialsExpandTrueNoNamespaces = `[
              {
                "accountType": "provider1",
                "cacheThreads": 0,
                "challengeDestructiveActions": false,
                "cloudProvider": "kubernetes",
                "dockerRegistries": null,
                "enabled": false,
                "environment": "provider1",
                "name": "provider1",
                "namespaces": [],
                "permissions": {
                  "READ": [
                    "gg_test"
                  ],
                  "WRITE": [
                    "gg_test"
                  ]
                },
                "primaryAccount": false,
                "providerVersion": "v2",
                "requiredGroupMembership": [],
                "skin": "v2",
                "spinnakerKindMap": {
                  "apiService": "unclassified",
                  "clusterRole": "unclassified",
                  "clusterRoleBinding": "unclassified",
                  "configMap": "configs",
                  "controllerRevision": "unclassified",
                  "cronJob": "serverGroups",
                  "customResourceDefinition": "unclassified",
                  "daemonSet": "serverGroups",
                  "deployment": "serverGroupManagers",
                  "event": "unclassified",
                  "horizontalpodautoscaler": "unclassified",
                  "ingress": "loadBalancers",
                  "job": "serverGroups",
                  "limitRange": "unclassified",
                  "mutatingWebhookConfiguration": "unclassified",
                  "namespace": "unclassified",
                  "networkPolicy": "securityGroups",
                  "persistentVolume": "configs",
                  "persistentVolumeClaim": "configs",
                  "pod": "instances",
                  "podDisruptionBudget": "unclassified",
                  "podPreset": "unclassified",
                  "podSecurityPolicy": "unclassified",
                  "replicaSet": "serverGroups",
                  "role": "unclassified",
                  "roleBinding": "unclassified",
                  "secret": "configs",
                  "service": "loadBalancers",
                  "serviceAccount": "unclassified",
                  "statefulSet": "serverGroups",
                  "storageClass": "unclassified",
                  "validatingWebhookConfiguration": "unclassified"
                },
                "type": "kubernetes"
              },
              {
                "accountType": "provider2",
                "cacheThreads": 0,
                "challengeDestructiveActions": false,
                "cloudProvider": "kubernetes",
                "dockerRegistries": null,
                "enabled": false,
                "environment": "provider2",
                "name": "provider2",
                "namespaces": [],
                "permissions": {
                  "READ": [
                    "gg_test2"
                  ],
                  "WRITE": [
                    "gg_test2"
                  ]
                },
                "primaryAccount": false,
                "providerVersion": "v2",
                "requiredGroupMembership": [],
                "skin": "v2",
                "spinnakerKindMap": {
                  "apiService": "unclassified",
                  "clusterRole": "unclassified",
                  "clusterRoleBinding": "unclassified",
                  "configMap": "configs",
                  "controllerRevision": "unclassified",
                  "cronJob": "serverGroups",
                  "customResourceDefinition": "unclassified",
                  "daemonSet": "serverGroups",
                  "deployment": "serverGroupManagers",
                  "event": "unclassified",
                  "horizontalpodautoscaler": "unclassified",
                  "ingress": "loadBalancers",
                  "job": "serverGroups",
                  "limitRange": "unclassified",
                  "mutatingWebhookConfiguration": "unclassified",
                  "namespace": "unclassified",
                  "networkPolicy": "securityGroups",
                  "persistentVolume": "configs",
                  "persistentVolumeClaim": "configs",
                  "pod": "instances",
                  "podDisruptionBudget": "unclassified",
                  "podPreset": "unclassified",
                  "podSecurityPolicy": "unclassified",
                  "replicaSet": "serverGroups",
                  "role": "unclassified",
                  "roleBinding": "unclassified",
                  "secret": "configs",
                  "service": "loadBalancers",
                  "serviceAccount": "unclassified",
                  "statefulSet": "serverGroups",
                  "storageClass": "unclassified",
                  "validatingWebhookConfiguration": "unclassified"
                },
                "type": "kubernetes"
              }
            ]`

const payloadCredentialsExpandTrueNamespaceScopedProvider = `[
						{
							"accountType": "provider1",
							"cacheThreads": 0,
							"challengeDestructiveActions": false,
							"cloudProvider": "kubernetes",
							"dockerRegistries": null,
							"enabled": false,
							"environment": "provider1",
							"name": "provider1",
							"namespaces": [
								"namespace1"
							],
							"permissions": {
								"READ": [
									"gg_test"
								],
								"WRITE": [
									"gg_test"
								]
							},
							"primaryAccount": false,
							"providerVersion": "v2",
							"requiredGroupMembership": [],
							"skin": "v2",
							"spinnakerKindMap": {
								"apiService": "unclassified",
								"clusterRole": "unclassified",
								"clusterRoleBinding": "unclassified",
								"configMap": "configs",
								"controllerRevision": "unclassified",
								"cronJob": "serverGroups",
								"customResourceDefinition": "unclassified",
								"daemonSet": "serverGroups",
								"deployment": "serverGroupManagers",
								"event": "unclassified",
								"horizontalpodautoscaler": "unclassified",
								"ingress": "loadBalancers",
								"job": "serverGroups",
								"limitRange": "unclassified",
								"mutatingWebhookConfiguration": "unclassified",
								"namespace": "unclassified",
								"networkPolicy": "securityGroups",
								"persistentVolume": "configs",
								"persistentVolumeClaim": "configs",
								"pod": "instances",
								"podDisruptionBudget": "unclassified",
								"podPreset": "unclassified",
								"podSecurityPolicy": "unclassified",
								"replicaSet": "serverGroups",
								"role": "unclassified",
								"roleBinding": "unclassified",
								"secret": "configs",
								"service": "loadBalancers",
								"serviceAccount": "unclassified",
								"statefulSet": "serverGroups",
								"storageClass": "unclassified",
								"validatingWebhookConfiguration": "unclassified"
							},
							"type": "kubernetes"
						}
					]`

const payloadServerGroupManagers = `[
            {
              "account": "account1",
              "apiVersion": "apps/v1",
              "cloudProvider": "kubernetes",
              "createdTime": 1581603123000,
              "kind": "deployment",
              "labels": {
                "label1": "test-label1"
              },
              "moniker": {
                "app": "test-application",
                "cluster": "deployment test-deployment1"
              },
              "name": "deployment test-deployment1",
              "displayName": "test-deployment1",
              "namespace": "test-namespace1",
              "region": "test-namespace1",
              "serverGroups": [
                {
                  "account": "account1",
                  "moniker": {
                    "app": "test-application",
                    "cluster": "deployment test-rs1",
                    "sequence": 236
                  },
                  "name": "replicaSet test-rs1",
                  "namespace": "test-namespace1",
                  "region": "test-namespace1"
                }
              ]
            },
            {
              "account": "account1",
              "apiVersion": "apps/v1",
              "cloudProvider": "kubernetes",
              "createdTime": 1581603123000,
              "kind": "deployment",
              "labels": {
                "label1": "test-label2"
              },
              "moniker": {
                "app": "test-application",
                "cluster": "deployment test-deployment2"
              },
              "name": "deployment test-deployment2",
              "displayName": "test-deployment2",
              "namespace": "test-namespace2",
              "region": "test-namespace2",
              "serverGroups": []
            }
          ]`

const payloadServerGroupManagersSorted = `[
            {
              "account": "account1",
              "apiVersion": "apps/v1",
              "cloudProvider": "kubernetes",
              "createdTime": 1581603123000,
              "kind": "deployment",
              "labels": {
                "label1": "test-label1"
              },
              "moniker": {
                "app": "test-application",
                "cluster": "deployment test-deployment1"
              },
              "name": "deployment test-deployment1",
              "displayName": "test-deployment1",
              "namespace": "test-namespace1",
              "region": "test-namespace1",
              "serverGroups": [
                {
                  "account": "account1",
                  "moniker": {
                    "app": "test-application",
                    "cluster": "deployment test-rs1",
                    "sequence": 236
                  },
                  "name": "replicaSet test-rs1",
                  "namespace": "test-namespace1",
                  "region": "test-namespace1"
                }
              ]
            },
            {
              "account": "account1",
              "apiVersion": "apps/v1",
              "cloudProvider": "kubernetes",
              "createdTime": 1581603123000,
              "kind": "deployment",
              "labels": {
                "label1": "test-label2"
              },
              "moniker": {
                "app": "test-application",
                "cluster": "deployment test-deployment2"
              },
              "name": "deployment test-deployment2",
              "displayName": "test-deployment2",
              "namespace": "test-namespace2",
              "region": "test-namespace2",
              "serverGroups": [
                {
                  "account": "account1",
                  "moniker": {
                    "app": "test-application",
                    "cluster": "deployment test-rs2",
                    "sequence": 236
                  },
                  "name": "replicaSet test-rs4",
                  "namespace": "test-namespace2",
                  "region": "test-namespace2"
                },
                {
                  "account": "account1",
                  "moniker": {
                    "app": "test-application",
                    "cluster": "deployment test-rs2",
                    "sequence": 236
                  },
                  "name": "replicaSet test-rs2",
                  "namespace": "test-namespace2",
                  "region": "test-namespace2"
                }
              ]
            },
            {
              "account": "account1",
              "apiVersion": "apps/v1",
              "cloudProvider": "kubernetes",
              "createdTime": 1581603123000,
              "kind": "deployment",
              "labels": {
                "label1": "test-label2"
              },
              "moniker": {
                "app": "test-application",
                "cluster": "deployment test-deployment4"
              },
              "name": "deployment test-deployment4",
              "displayName": "test-deployment4",
              "namespace": "test-namespace2",
              "region": "test-namespace2",
              "serverGroups": [
                {
                  "account": "account1",
                  "moniker": {
                    "app": "test-application",
                    "cluster": "deployment test-rs2",
                    "sequence": 236
                  },
                  "name": "replicaSet test-rs4",
                  "namespace": "test-namespace2",
                  "region": "test-namespace2"
                },
                {
                  "account": "account1",
                  "moniker": {
                    "app": "test-application",
                    "cluster": "deployment test-rs2",
                    "sequence": 236
                  },
                  "name": "replicaSet test-rs2",
                  "namespace": "test-namespace2",
                  "region": "test-namespace2"
                }
              ]
            },
            {
              "account": "account1",
              "apiVersion": "apps/v1",
              "cloudProvider": "kubernetes",
              "createdTime": 1581603123000,
              "kind": "deployment",
              "labels": {
                "label1": "test-label3"
              },
              "moniker": {
                "app": "test-application",
                "cluster": "deployment test-deployment3"
              },
              "name": "deployment test-deployment3",
              "displayName": "test-deployment3",
              "namespace": "test-namespace3",
              "region": "test-namespace3",
              "serverGroups": [
                {
                  "account": "account1",
                  "moniker": {
                    "app": "test-application",
                    "cluster": "deployment test-rs3",
                    "sequence": 236
                  },
                  "name": "replicaSet test-rs3",
                  "namespace": "test-namespace3",
                  "region": "test-namespace3"
                }
              ]
            }
          ]`

const payloadListServerGroupsDisabled = `[
            {
              "account": "account1",
              "accountName": "",
              "buildInfo": {
                "images": [
                  "test-image1",
                  "test-image2"
                ]
              },
              "capacity": {
                "desired": 2,
                "pinned": false
              },
              "cloudProvider": "kubernetes",
              "cluster": "deployment test-deployment1",
              "createdTime": 1581603123000,
              "disabled": false,
              "displayName": "test-ds1",
              "instanceCounts": {
                "down": 0,
                "outOfService": 0,
                "starting": 0,
                "total": 2,
                "unknown": 0,
                "up": 1
              },
              "instances": [
                {
                  "availabilityZone": "test-namespace1",
                  "health": [
                    {
                      "state": "Down",
                      "type": "kubernetes/pod"
                    },
                    {
                      "state": "Down",
                      "type": "kubernetes/container"
                    }
                  ],
                  "healthState": "Down",
                  "id": "cec15437-4e6a-11ea-9788-4201ac100006",
                  "key": {
                    "account": "",
                    "group": "",
                    "kubernetesKind": "",
                    "name": "",
                    "namespace": "",
                    "provider": ""
                  },
                  "moniker": {
                    "app": "",
                    "cluster": ""
                  },
                  "name": "pod test-pod1"
                }
              ],
              "isDisabled": false,
              "key": {
                "account": "",
                "group": "",
                "kubernetesKind": "",
                "name": "",
                "namespace": "",
                "provider": ""
              },
              "kind": "daemonSet",
              "labels": null,
              "loadBalancers": null,
              "manifest": null,
              "moniker": {
                "app": "test-application",
                "cluster": "deployment test-deployment1",
                "sequence": 19
              },
              "name": "daemonSet test-ds1",
              "namespace": "test-namespace1",
              "providerType": "",
              "region": "test-namespace1",
              "securityGroups": null,
              "serverGroupManagers": [],
              "type": "kubernetes",
              "uid": "",
              "zone": "",
              "zones": null,
              "insightActions": null
            },
            {
              "account": "account1",
              "accountName": "",
              "buildInfo": {
                "images": [
                  "test-image1",
                  "test-image2"
                ]
              },
              "capacity": {
                "desired": 1,
                "pinned": false
              },
              "cloudProvider": "kubernetes",
              "cluster": "deployment test-deployment1",
              "createdTime": 1581603123000,
              "disabled": false,
              "displayName": "test-rs1",
              "instanceCounts": {
                "down": 0,
                "outOfService": 0,
                "starting": 0,
                "total": 1,
                "unknown": 0,
                "up": 0
              },
              "instances": [
                {
                  "availabilityZone": "test-namespace1",
                  "health": [
                    {
                      "state": "Down",
                      "type": "kubernetes/pod"
                    },
                    {
                      "state": "Down",
                      "type": "kubernetes/container"
                    }
                  ],
                  "healthState": "Down",
                  "id": "cec15437-4e6a-11ea-9788-4201ac100006",
                  "key": {
                    "account": "",
                    "group": "",
                    "kubernetesKind": "",
                    "name": "",
                    "namespace": "",
                    "provider": ""
                  },
                  "moniker": {
                    "app": "",
                    "cluster": ""
                  },
                  "name": "pod test-pod1"
                }
              ],
              "isDisabled": true,
              "key": {
                "account": "",
                "group": "",
                "kubernetesKind": "",
                "name": "",
                "namespace": "",
                "provider": ""
              },
              "kind": "replicaSet",
              "labels": null,
              "loadBalancers": [
							  "service test-svc1",
							  "service my-managed-service1",
							  "service my-managed-service2"
							],
              "manifest": null,
              "moniker": {
                "app": "test-application",
                "cluster": "deployment test-deployment1",
                "sequence": 19
              },
              "name": "replicaSet test-rs1",
              "namespace": "test-namespace1",
              "providerType": "",
              "region": "test-namespace1",
              "securityGroups": null,
              "serverGroupManagers": [
                {
                  "account": "account1",
                  "location": "test-namespace1",
                  "name": "test-deployment1"
                }
              ],
              "type": "kubernetes",
              "uid": "",
              "zone": "",
              "zones": null,
              "insightActions": null
            },
            {
              "account": "account1",
              "accountName": "",
              "buildInfo": {
                "images": [
                  "test-image1",
                  "test-image2"
                ]
              },
              "capacity": {
                "desired": 1,
                "pinned": false
              },
              "cloudProvider": "kubernetes",
              "cluster": "deployment test-deployment1",
              "createdTime": 1581603123000,
              "disabled": false,
              "displayName": "test-sts1",
              "instanceCounts": {
                "down": 0,
                "outOfService": 0,
                "starting": 0,
                "total": 1,
                "unknown": 0,
                "up": 0
              },
              "instances": [
                {
                  "availabilityZone": "test-namespace1",
                  "health": [
                    {
                      "state": "Down",
                      "type": "kubernetes/pod"
                    },
                    {
                      "state": "Down",
                      "type": "kubernetes/container"
                    }
                  ],
                  "healthState": "Down",
                  "id": "cec15437-4e6a-11ea-9788-4201ac100006",
                  "key": {
                    "account": "",
                    "group": "",
                    "kubernetesKind": "",
                    "name": "",
                    "namespace": "",
                    "provider": ""
                  },
                  "moniker": {
                    "app": "",
                    "cluster": ""
                  },
                  "name": "pod test-pod1"
                }
              ],
              "isDisabled": false,
              "key": {
                "account": "",
                "group": "",
                "kubernetesKind": "",
                "name": "",
                "namespace": "",
                "provider": ""
              },
              "kind": "statefulSet",
              "labels": null,
              "loadBalancers": [
							  "service test-svc2"
							],
              "manifest": null,
              "moniker": {
                "app": "test-application",
                "cluster": "deployment test-deployment1",
                "sequence": 19
              },
              "name": "statefulSet test-sts1",
              "namespace": "test-namespace1",
              "providerType": "",
              "region": "test-namespace1",
              "securityGroups": null,
              "serverGroupManagers": [],
              "type": "kubernetes",
              "uid": "",
              "zone": "",
              "zones": null,
              "insightActions": null
            }
          ]`

const payloadListServerGroups = `[
            {
              "account": "account1",
              "accountName": "",
              "buildInfo": {
                "images": [
                  "test-image1",
                  "test-image2"
                ]
              },
              "capacity": {
                "desired": 2,
                "pinned": false
              },
              "cloudProvider": "kubernetes",
              "cluster": "deployment test-deployment1",
              "createdTime": 1581603123000,
              "disabled": false,
              "displayName": "test-ds1",
              "instanceCounts": {
                "down": 0,
                "outOfService": 0,
                "starting": 0,
                "total": 2,
                "unknown": 0,
                "up": 1
              },
              "instances": [
                {
                  "availabilityZone": "test-namespace1",
                  "health": [
                    {
                      "state": "Down",
                      "type": "kubernetes/pod"
                    },
                    {
                      "state": "Down",
                      "type": "kubernetes/container"
                    }
                  ],
                  "healthState": "Down",
                  "id": "cec15437-4e6a-11ea-9788-4201ac100006",
                  "key": {
                    "account": "",
                    "group": "",
                    "kubernetesKind": "",
                    "name": "",
                    "namespace": "",
                    "provider": ""
                  },
                  "moniker": {
                    "app": "",
                    "cluster": ""
                  },
                  "name": "pod test-pod1"
                }
              ],
              "isDisabled": false,
              "key": {
                "account": "",
                "group": "",
                "kubernetesKind": "",
                "name": "",
                "namespace": "",
                "provider": ""
              },
              "kind": "daemonSet",
              "labels": null,
              "loadBalancers": null,
              "manifest": null,
              "moniker": {
                "app": "test-application",
                "cluster": "deployment test-deployment1",
                "sequence": 19
              },
              "name": "daemonSet test-ds1",
              "namespace": "test-namespace1",
              "providerType": "",
              "region": "test-namespace1",
              "securityGroups": null,
              "serverGroupManagers": [],
              "type": "kubernetes",
              "uid": "",
              "zone": "",
              "zones": null,
              "insightActions": null
            },
            {
              "account": "account1",
              "accountName": "",
              "buildInfo": {
                "images": [
                  "test-image1",
                  "test-image2"
                ]
              },
              "capacity": {
                "desired": 1,
                "pinned": false
              },
              "cloudProvider": "kubernetes",
              "cluster": "deployment test-deployment1",
              "createdTime": 1581603123000,
              "disabled": false,
              "displayName": "test-rs1",
              "instanceCounts": {
                "down": 0,
                "outOfService": 0,
                "starting": 0,
                "total": 1,
                "unknown": 0,
                "up": 0
              },
              "instances": [
                {
                  "availabilityZone": "test-namespace1",
                  "health": [
                    {
                      "state": "Down",
                      "type": "kubernetes/pod"
                    },
                    {
                      "state": "Down",
                      "type": "kubernetes/container"
                    }
                  ],
                  "healthState": "Down",
                  "id": "cec15437-4e6a-11ea-9788-4201ac100006",
                  "key": {
                    "account": "",
                    "group": "",
                    "kubernetesKind": "",
                    "name": "",
                    "namespace": "",
                    "provider": ""
                  },
                  "moniker": {
                    "app": "",
                    "cluster": ""
                  },
                  "name": "pod test-pod1"
                }
              ],
              "isDisabled": false,
              "key": {
                "account": "",
                "group": "",
                "kubernetesKind": "",
                "name": "",
                "namespace": "",
                "provider": ""
              },
              "kind": "replicaSet",
              "labels": null,
              "loadBalancers": [
							  "service test-svc1"
							],
              "manifest": null,
              "moniker": {
                "app": "test-application",
                "cluster": "deployment test-deployment1",
                "sequence": 19
              },
              "name": "replicaSet test-rs1",
              "namespace": "test-namespace1",
              "providerType": "",
              "region": "test-namespace1",
              "securityGroups": null,
              "serverGroupManagers": [
                {
                  "account": "account1",
                  "location": "test-namespace1",
                  "name": "test-deployment1"
                }
              ],
              "type": "kubernetes",
              "uid": "",
              "zone": "",
              "zones": null,
              "insightActions": null
            },
            {
              "account": "account1",
              "accountName": "",
              "buildInfo": {
                "images": [
                  "test-image1",
                  "test-image2"
                ]
              },
              "capacity": {
                "desired": 1,
                "pinned": false
              },
              "cloudProvider": "kubernetes",
              "cluster": "deployment test-deployment1",
              "createdTime": 1581603123000,
              "disabled": false,
              "displayName": "test-sts1",
              "instanceCounts": {
                "down": 0,
                "outOfService": 0,
                "starting": 0,
                "total": 1,
                "unknown": 0,
                "up": 0
              },
              "instances": [
                {
                  "availabilityZone": "test-namespace1",
                  "health": [
                    {
                      "state": "Down",
                      "type": "kubernetes/pod"
                    },
                    {
                      "state": "Down",
                      "type": "kubernetes/container"
                    }
                  ],
                  "healthState": "Down",
                  "id": "cec15437-4e6a-11ea-9788-4201ac100006",
                  "key": {
                    "account": "",
                    "group": "",
                    "kubernetesKind": "",
                    "name": "",
                    "namespace": "",
                    "provider": ""
                  },
                  "moniker": {
                    "app": "",
                    "cluster": ""
                  },
                  "name": "pod test-pod1"
                }
              ],
              "isDisabled": false,
              "key": {
                "account": "",
                "group": "",
                "kubernetesKind": "",
                "name": "",
                "namespace": "",
                "provider": ""
              },
              "kind": "statefulSet",
              "labels": null,
              "loadBalancers": [
							  "service test-svc2"
							],
              "manifest": null,
              "moniker": {
                "app": "test-application",
                "cluster": "deployment test-deployment1",
                "sequence": 19
              },
              "name": "statefulSet test-sts1",
              "namespace": "test-namespace1",
              "providerType": "",
              "region": "test-namespace1",
              "securityGroups": null,
              "serverGroupManagers": [],
              "type": "kubernetes",
              "uid": "",
              "zone": "",
              "zones": null,
              "insightActions": null
            }
          ]`

const payloadListServerGroupsUniquePodInstances = `[
            {
              "account": "account1",
              "accountName": "",
              "buildInfo": {
                "images": [
                  "test-image1",
                  "test-image2"
                ]
              },
              "capacity": {
                "desired": 2,
                "pinned": false
              },
              "cloudProvider": "kubernetes",
              "cluster": "deployment test-deployment1",
              "createdTime": 1581603123000,
              "disabled": false,
              "displayName": "test-ds1",
              "instanceCounts": {
                "down": 0,
                "outOfService": 0,
                "starting": 0,
                "total": 2,
                "unknown": 0,
                "up": 1
              },
              "instances": [],
              "isDisabled": false,
              "key": {
                "account": "",
                "group": "",
                "kubernetesKind": "",
                "name": "",
                "namespace": "",
                "provider": ""
              },
              "kind": "daemonSet",
              "labels": null,
              "loadBalancers": null,
              "manifest": null,
              "moniker": {
                "app": "test-application",
                "cluster": "deployment test-deployment1",
                "sequence": 19
              },
              "name": "daemonSet test-ds1",
              "namespace": "test-namespace1",
              "providerType": "",
              "region": "test-namespace1",
              "securityGroups": null,
              "serverGroupManagers": [],
              "type": "kubernetes",
              "uid": "",
              "zone": "",
              "zones": null,
              "insightActions": null
            },
            {
              "account": "account1",
              "accountName": "",
              "buildInfo": {
                "images": [
                  "test-image1",
                  "test-image2"
                ]
              },
              "capacity": {
                "desired": 1,
                "pinned": false
              },
              "cloudProvider": "kubernetes",
              "cluster": "deployment test-deployment1",
              "createdTime": 1581603123000,
              "disabled": false,
              "displayName": "test-rs1",
              "instanceCounts": {
                "down": 0,
                "outOfService": 0,
                "starting": 0,
                "total": 1,
                "unknown": 0,
                "up": 0
              },
              "instances": [
                {
                  "availabilityZone": "test-namespace1",
                  "health": [
                    {
                      "state": "Down",
                      "type": "kubernetes/pod"
                    },
                    {
                      "state": "Down",
                      "type": "kubernetes/container"
                    }
                  ],
                  "healthState": "Down",
                  "id": "cec15437-4e6a-11ea-9788-4201ac100000",
                  "key": {
                    "account": "",
                    "group": "",
                    "kubernetesKind": "",
                    "name": "",
                    "namespace": "",
                    "provider": ""
                  },
                  "moniker": {
                    "app": "",
                    "cluster": ""
                  },
                  "name": "pod test-pod3"
                }
              ],
              "isDisabled": false,
              "key": {
                "account": "",
                "group": "",
                "kubernetesKind": "",
                "name": "",
                "namespace": "",
                "provider": ""
              },
              "kind": "replicaSet",
              "labels": null,
              "loadBalancers": null,
              "manifest": null,
              "moniker": {
                "app": "test-application",
                "cluster": "deployment test-deployment1",
                "sequence": 19
              },
              "name": "replicaSet test-rs1",
              "namespace": "test-namespace1",
              "providerType": "",
              "region": "test-namespace1",
              "securityGroups": null,
              "serverGroupManagers": [
                {
                  "account": "account1",
                  "location": "test-namespace1",
                  "name": "test-deployment1"
                }
              ],
              "type": "kubernetes",
              "uid": "",
              "zone": "",
              "zones": null,
              "insightActions": null
            },
            {
              "account": "account1",
              "accountName": "",
              "buildInfo": {
                "images": [
                  "test-image1",
                  "test-image2"
                ]
              },
              "capacity": {
                "desired": 1,
                "pinned": false
              },
              "cloudProvider": "kubernetes",
              "cluster": "deployment test-deployment1",
              "createdTime": 1581603123000,
              "disabled": false,
              "displayName": "test-sts1",
              "instanceCounts": {
                "down": 0,
                "outOfService": 0,
                "starting": 0,
                "total": 1,
                "unknown": 0,
                "up": 0
              },
              "instances": [],
              "isDisabled": false,
              "key": {
                "account": "",
                "group": "",
                "kubernetesKind": "",
                "name": "",
                "namespace": "",
                "provider": ""
              },
              "kind": "statefulSet",
              "labels": null,
              "loadBalancers": [
                "service test-svc2"
              ],
              "manifest": null,
              "moniker": {
                "app": "test-application",
                "cluster": "deployment test-deployment1",
                "sequence": 19
              },
              "name": "statefulSet test-sts1",
              "namespace": "test-namespace1",
              "providerType": "",
              "region": "test-namespace1",
              "securityGroups": null,
              "serverGroupManagers": [],
              "type": "kubernetes",
              "uid": "",
              "zone": "",
              "zones": null,
              "insightActions": null
            }
          ]`

const payloadListServerGroupsSorted = `[
            {
              "account": "account1",
              "accountName": "",
              "buildInfo": {
                "images": [
                  "test-image1",
                  "test-image2"
                ]
              },
              "capacity": {
                "desired": 2,
                "pinned": false
              },
              "cloudProvider": "kubernetes",
              "cluster": "deployment test-deployment1",
              "createdTime": 1581603123000,
              "disabled": false,
              "displayName": "test-ds1",
              "instanceCounts": {
                "down": 0,
                "outOfService": 0,
                "starting": 0,
                "total": 2,
                "unknown": 0,
                "up": 1
              },
              "instances": [
                {
                  "availabilityZone": "test-namespace1",
                  "health": [
                    {
                      "state": "Down",
                      "type": "kubernetes/pod"
                    },
                    {
                      "state": "Down",
                      "type": "kubernetes/container"
                    }
                  ],
                  "healthState": "Down",
                  "id": "cec15437-4e6a-11ea-9788-4201ac100004",
                  "key": {
                    "account": "",
                    "group": "",
                    "kubernetesKind": "",
                    "name": "",
                    "namespace": "",
                    "provider": ""
                  },
                  "moniker": {
                    "app": "",
                    "cluster": ""
                  },
                  "name": "pod test-pod2"
                },
                {
                  "availabilityZone": "test-namespace1",
                  "health": [
                    {
                      "state": "Down",
                      "type": "kubernetes/pod"
                    },
                    {
                      "state": "Down",
                      "type": "kubernetes/container"
                    }
                  ],
                  "healthState": "Down",
                  "id": "cec15437-4e6a-11ea-9788-4201ac100003",
                  "key": {
                    "account": "",
                    "group": "",
                    "kubernetesKind": "",
                    "name": "",
                    "namespace": "",
                    "provider": ""
                  },
                  "moniker": {
                    "app": "",
                    "cluster": ""
                  },
                  "name": "pod test-pod3"
                },
                {
                  "availabilityZone": "test-namespace1",
                  "health": [
                    {
                      "state": "Down",
                      "type": "kubernetes/pod"
                    },
                    {
                      "state": "Down",
                      "type": "kubernetes/container"
                    }
                  ],
                  "healthState": "Down",
                  "id": "cec15437-4e6a-11ea-9788-4201ac100002",
                  "key": {
                    "account": "",
                    "group": "",
                    "kubernetesKind": "",
                    "name": "",
                    "namespace": "",
                    "provider": ""
                  },
                  "moniker": {
                    "app": "",
                    "cluster": ""
                  },
                  "name": "pod test-pod4"
                }
              ],
              "isDisabled": false,
              "key": {
                "account": "",
                "group": "",
                "kubernetesKind": "",
                "name": "",
                "namespace": "",
                "provider": ""
              },
              "kind": "daemonSet",
              "labels": null,
              "loadBalancers": null,
              "manifest": null,
              "moniker": {
                "app": "test-application",
                "cluster": "deployment test-deployment1",
                "sequence": 19
              },
              "name": "daemonSet test-ds1",
              "namespace": "test-namespace1",
              "providerType": "",
              "region": "test-namespace1",
              "securityGroups": null,
              "serverGroupManagers": [],
              "type": "kubernetes",
              "uid": "",
              "zone": "",
              "zones": null,
              "insightActions": null
            },
            {
              "account": "account1",
              "accountName": "",
              "buildInfo": {
                "images": [
                  "test-image1",
                  "test-image2"
                ]
              },
              "capacity": {
                "desired": 1,
                "pinned": false
              },
              "cloudProvider": "kubernetes",
              "cluster": "deployment test-deployment1",
              "createdTime": 1581603123000,
              "disabled": false,
              "displayName": "test-rs1",
              "instanceCounts": {
                "down": 0,
                "outOfService": 0,
                "starting": 0,
                "total": 1,
                "unknown": 0,
                "up": 0
              },
              "instances": [
                {
                  "availabilityZone": "test-namespace1",
                  "health": [
                    {
                      "state": "Down",
                      "type": "kubernetes/pod"
                    },
                    {
                      "state": "Down",
                      "type": "kubernetes/container"
                    }
                  ],
                  "healthState": "Down",
                  "id": "cec15437-4e6a-11ea-9788-4201ac100001",
                  "key": {
                    "account": "",
                    "group": "",
                    "kubernetesKind": "",
                    "name": "",
                    "namespace": "",
                    "provider": ""
                  },
                  "moniker": {
                    "app": "",
                    "cluster": ""
                  },
                  "name": "pod test-pod3"
                }
              ],
              "isDisabled": false,
              "key": {
                "account": "",
                "group": "",
                "kubernetesKind": "",
                "name": "",
                "namespace": "",
                "provider": ""
              },
              "kind": "replicaSet",
              "labels": null,
              "loadBalancers": null,
              "manifest": null,
              "moniker": {
                "app": "test-application",
                "cluster": "deployment test-deployment1",
                "sequence": 19
              },
              "name": "replicaSet test-rs1",
              "namespace": "test-namespace1",
              "providerType": "",
              "region": "test-namespace1",
              "securityGroups": null,
              "serverGroupManagers": [
                {
                  "account": "account1",
                  "location": "test-namespace1",
                  "name": "test-deployment1"
                }
              ],
              "type": "kubernetes",
              "uid": "",
              "zone": "",
              "zones": null,
              "insightActions": null
            },
            {
              "account": "account1",
              "accountName": "",
              "buildInfo": {
                "images": [
                  "test-image1",
                  "test-image2"
                ]
              },
              "capacity": {
                "desired": 1,
                "pinned": false
              },
              "cloudProvider": "kubernetes",
              "cluster": "deployment test-deployment1",
              "createdTime": 1581603123000,
              "disabled": false,
              "displayName": "test-rs1",
              "instanceCounts": {
                "down": 0,
                "outOfService": 0,
                "starting": 0,
                "total": 1,
                "unknown": 0,
                "up": 0
              },
              "instances": [
                {
                  "availabilityZone": "test-namespace1",
                  "health": [
                    {
                      "state": "Down",
                      "type": "kubernetes/pod"
                    },
                    {
                      "state": "Down",
                      "type": "kubernetes/container"
                    }
                  ],
                  "healthState": "Down",
                  "id": "cec15437-4e6a-11ea-9788-4201ac100005",
                  "key": {
                    "account": "",
                    "group": "",
                    "kubernetesKind": "",
                    "name": "",
                    "namespace": "",
                    "provider": ""
                  },
                  "moniker": {
                    "app": "",
                    "cluster": ""
                  },
                  "name": "pod test-pod1"
                }
              ],
              "isDisabled": false,
              "key": {
                "account": "",
                "group": "",
                "kubernetesKind": "",
                "name": "",
                "namespace": "",
                "provider": ""
              },
              "kind": "statefulSet",
              "labels": null,
              "loadBalancers": null,
              "manifest": null,
              "moniker": {
                "app": "test-application",
                "cluster": "deployment test-deployment1",
                "sequence": 19
              },
              "name": "statefulSet test-rs1",
              "namespace": "test-namespace1",
              "providerType": "",
              "region": "test-namespace1",
              "securityGroups": null,
              "serverGroupManagers": [],
              "type": "kubernetes",
              "uid": "",
              "zone": "",
              "zones": null,
              "insightActions": null
            },
            {
              "account": "account1",
              "accountName": "",
              "buildInfo": {
                "images": [
                  "test-image1",
                  "test-image2"
                ]
              },
              "capacity": {
                "desired": 2,
                "pinned": false
              },
              "cloudProvider": "kubernetes",
              "cluster": "deployment test-deployment1",
              "createdTime": 1581603123000,
              "disabled": false,
              "displayName": "test-ds1",
              "instanceCounts": {
                "down": 0,
                "outOfService": 0,
                "starting": 0,
                "total": 2,
                "unknown": 0,
                "up": 1
              },
              "instances": [
                {
                  "availabilityZone": "test-namespace1",
                  "health": [
                    {
                      "state": "Down",
                      "type": "kubernetes/pod"
                    },
                    {
                      "state": "Down",
                      "type": "kubernetes/container"
                    }
                  ],
                  "healthState": "Down",
                  "id": "cec15437-4e6a-11ea-9788-4201ac100004",
                  "key": {
                    "account": "",
                    "group": "",
                    "kubernetesKind": "",
                    "name": "",
                    "namespace": "",
                    "provider": ""
                  },
                  "moniker": {
                    "app": "",
                    "cluster": ""
                  },
                  "name": "pod test-pod2"
                },
                {
                  "availabilityZone": "test-namespace1",
                  "health": [
                    {
                      "state": "Down",
                      "type": "kubernetes/pod"
                    },
                    {
                      "state": "Down",
                      "type": "kubernetes/container"
                    }
                  ],
                  "healthState": "Down",
                  "id": "cec15437-4e6a-11ea-9788-4201ac100003",
                  "key": {
                    "account": "",
                    "group": "",
                    "kubernetesKind": "",
                    "name": "",
                    "namespace": "",
                    "provider": ""
                  },
                  "moniker": {
                    "app": "",
                    "cluster": ""
                  },
                  "name": "pod test-pod3"
                },
                {
                  "availabilityZone": "test-namespace1",
                  "health": [
                    {
                      "state": "Down",
                      "type": "kubernetes/pod"
                    },
                    {
                      "state": "Down",
                      "type": "kubernetes/container"
                    }
                  ],
                  "healthState": "Down",
                  "id": "cec15437-4e6a-11ea-9788-4201ac100002",
                  "key": {
                    "account": "",
                    "group": "",
                    "kubernetesKind": "",
                    "name": "",
                    "namespace": "",
                    "provider": ""
                  },
                  "moniker": {
                    "app": "",
                    "cluster": ""
                  },
                  "name": "pod test-pod4"
                }
              ],
              "isDisabled": false,
              "key": {
                "account": "",
                "group": "",
                "kubernetesKind": "",
                "name": "",
                "namespace": "",
                "provider": ""
              },
              "kind": "daemonSet",
              "labels": null,
              "loadBalancers": null,
              "manifest": null,
              "moniker": {
                "app": "test-application",
                "cluster": "deployment test-deployment1",
                "sequence": 19
              },
              "name": "daemonSet test-ds1",
              "namespace": "test-namespace2",
              "providerType": "",
              "region": "test-namespace2",
              "securityGroups": null,
              "serverGroupManagers": [],
              "type": "kubernetes",
              "uid": "",
              "zone": "",
              "zones": null,
              "insightActions": null
            },
            {
              "account": "account1",
              "accountName": "",
              "buildInfo": {
                "images": [
                  "test-image1",
                  "test-image2"
                ]
              },
              "capacity": {
                "desired": 1,
                "pinned": false
              },
              "cloudProvider": "kubernetes",
              "cluster": "deployment test-deployment1",
              "createdTime": 1581603123000,
              "disabled": false,
              "displayName": "test-rs1",
              "instanceCounts": {
                "down": 0,
                "outOfService": 0,
                "starting": 0,
                "total": 1,
                "unknown": 0,
                "up": 0
              },
              "instances": [
                {
                  "availabilityZone": "test-namespace1",
                  "health": [
                    {
                      "state": "Down",
                      "type": "kubernetes/pod"
                    },
                    {
                      "state": "Down",
                      "type": "kubernetes/container"
                    }
                  ],
                  "healthState": "Down",
                  "id": "cec15437-4e6a-11ea-9788-4201ac100005",
                  "key": {
                    "account": "",
                    "group": "",
                    "kubernetesKind": "",
                    "name": "",
                    "namespace": "",
                    "provider": ""
                  },
                  "moniker": {
                    "app": "",
                    "cluster": ""
                  },
                  "name": "pod test-pod1"
                }
              ],
              "isDisabled": false,
              "key": {
                "account": "",
                "group": "",
                "kubernetesKind": "",
                "name": "",
                "namespace": "",
                "provider": ""
              },
              "kind": "statefulSet",
              "labels": null,
              "loadBalancers": null,
              "manifest": null,
              "moniker": {
                "app": "test-application",
                "cluster": "deployment test-deployment1",
                "sequence": 19
              },
              "name": "statefulSet test-rs1",
              "namespace": "test-namespace2",
              "providerType": "",
              "region": "test-namespace2",
              "securityGroups": null,
              "serverGroupManagers": [],
              "type": "kubernetes",
              "uid": "",
              "zone": "",
              "zones": null,
              "insightActions": null
            }
          ]`

const payloadListLoadBalancers = `[
            {
              "account": "account1",
              "apiVersion": "networking.k8s.io/v1beta1",
              "cloudProvider": "kubernetes",
              "createdTime": 1581603123000,
              "displayName": "test-ingress1",
              "kind": "ingress",
              "labels": {
                "label1": "test-label1"
              },
              "moniker": {
                "app": "test-application",
                "cluster": "ingress test-ingress1"
              },
              "name": "ingress test-ingress1",
              "namespace": "test-namespace1",
              "region": "test-namespace1",
              "serverGroups": [],
              "type": "kubernetes"
            },
            {
              "account": "account1",
              "apiVersion": "v1",
              "cloudProvider": "kubernetes",
              "createdTime": -62135596800000,
              "displayName": "test-service1",
              "kind": "service",
              "moniker": {
                "app": "test-application",
                "cluster": "service test-service1"
              },
              "name": "service test-service1",
              "namespace": "test-namespace1",
              "region": "test-namespace1",
              "serverGroups": [
                {
                  "account": "account1",
                  "cloudProvider": "kubernetes",
                  "detachedInstances": [],
                  "instances": [
                    {
                      "health": {
                        "platform": "platform",
                        "source": "Container TODO",
                        "state": "Down",
                        "type": "kubernetes/container"
                      },
                      "id": "cec15437-4e6a-11ea-9788-4201ac100006",
                      "name": "pod test-pod1",
                      "zone": "test-namespace1"
                    }
                  ],
                  "isDisabled": false,
                  "name": "replicaSet test-rs1",
                  "region": "test-namespace1"
                }
              ],
              "type": "kubernetes"
            }
          ]`

const payloadListLoadBalancersSorted = `[
						{
							"account": "account1",
							"apiVersion": "networking.k8s.io/v1beta1",
							"cloudProvider": "kubernetes",
							"createdTime": 1581603123000,
							"displayName": "test-ingress1",
							"kind": "ingress",
							"labels": {
								"label1": "test-label1"
							},
							"moniker": {
								"app": "test-application",
								"cluster": "ingress test-ingress1"
							},
							"name": "ingress test-ingress1",
							"namespace": "test-namespace1",
							"region": "test-namespace1",
							"serverGroups": [],
							"type": "kubernetes"
						},
						{
							"account": "account1",
							"apiVersion": "v1",
							"cloudProvider": "kubernetes",
							"createdTime": -62135596800000,
							"displayName": "test-service1",
							"kind": "service",
							"moniker": {
								"app": "test-application",
								"cluster": "service test-service1"
							},
							"name": "service test-service1",
							"namespace": "test-namespace1",
							"region": "test-namespace1",
							"serverGroups": [
								{
									"account": "account1",
									"cloudProvider": "kubernetes",
									"detachedInstances": [],
									"instances": [
										{
											"health": {
												"platform": "platform",
												"source": "Container TODO",
												"state": "Down",
												"type": "kubernetes/container"
											},
											"id": "cec15437-4e6a-11ea-9788-4201ac100006",
											"name": "pod test-pod1",
											"zone": "test-namespace1"
										}
									],
									"isDisabled": false,
									"name": "statefulSet test-sts1",
									"region": "test-namespace1"
								}
							],
							"type": "kubernetes"
						},
						{
							"account": "account1",
							"apiVersion": "networking.k8s.io/v1beta1",
							"cloudProvider": "kubernetes",
							"createdTime": 1581603123000,
							"displayName": "test-ingress2",
							"kind": "ingress",
							"labels": {
								"label1": "test-label1"
							},
							"moniker": {
								"app": "test-application",
								"cluster": "ingress test-ingress2"
							},
							"name": "ingress test-ingress2",
							"namespace": "test-namespace2",
							"region": "test-namespace2",
							"serverGroups": [],
							"type": "kubernetes"
						},
						{
							"account": "account1",
							"apiVersion": "networking.k8s.io/v1beta1",
							"cloudProvider": "kubernetes",
							"createdTime": 1581603123000,
							"displayName": "test-ingress3",
							"kind": "ingress",
							"labels": {
								"label1": "test-label1"
							},
							"moniker": {
								"app": "test-application",
								"cluster": "ingress test-ingress3"
							},
							"name": "ingress test-ingress3",
							"namespace": "test-namespace2",
							"region": "test-namespace2",
							"serverGroups": [],
							"type": "kubernetes"
						},
						{
							"account": "account1",
							"apiVersion": "v1",
							"cloudProvider": "kubernetes",
							"createdTime": -62135596800000,
							"displayName": "test-service2",
							"kind": "service",
							"moniker": {
								"app": "test-application",
								"cluster": "service test-service2"
							},
							"name": "service test-service2",
							"namespace": "test-namespace2",
							"region": "test-namespace2",
							"serverGroups": [
								{
									"account": "account1",
									"cloudProvider": "kubernetes",
									"detachedInstances": [],
									"instances": [
										{
											"health": {
												"platform": "platform",
												"source": "Container TODO",
												"state": "Down",
												"type": "kubernetes/container"
											},
											"id": "cec15437-4e6a-11ea-9788-4201ac100006",
											"name": "pod test-pod1",
											"zone": "test-namespace1"
										}
									],
									"isDisabled": false,
									"name": "replicaSet test-rs1",
									"region": "test-namespace2"
								}
							],
							"type": "kubernetes"
						},
						{
							"account": "account1",
							"apiVersion": "v1",
							"cloudProvider": "kubernetes",
							"createdTime": -62135596800000,
							"displayName": "test-service3",
							"kind": "service",
							"moniker": {
								"app": "test-application",
								"cluster": "service test-service3"
							},
							"name": "service test-service3",
							"namespace": "test-namespace2",
							"region": "test-namespace2",
							"serverGroups": [],
							"type": "kubernetes"
						}
					]`

const payloadListClusters = `{
            "test-account1": [
              "test-kind1 test-name1"
            ],
            "test-account2": [
              "test-kind2 test-name2",
              "test-kind3 test-name3"
            ]
          }`

const payloadListClusters2 = `{
            "test-account1": [
              "test-kind1 test-name1"
            ],
            "test-account2": [
              "test-kind3 test-name3"
            ]
          }`

const payloadListClustersByName = `{
          "accountName": "test-account",
          "application": "test-application",
          "loadBalancers": [],
          "moniker": {
            "app": "test-application",
            "cluster": "deployment test-deployment"
          },
          "name": "deployment test-deployment",
          "serverGroups": [
            {
              "account": "test-account",
              "apiVersion": "apps/v1",
              "cloudProvider": "kubernetes",
              "displayName": "test-rs1",
              "kind": "ReplicaSet",
              "labels": {
                "labelKey1": "labelValue1",
                "labelKey2": "labelValue2"
              },
              "moniker": {
                "app": "test-application",
                "cluster": "deployment test-deployment",
                "sequence": 19
              },
              "name": "replicaSet test-rs1",
              "namespace": "test-namespace1",
              "region": "test-namespace1",
              "serverGroupManagers": [
                {
                  "account": "test-account",
                  "location": "test-namespace1",
                  "name": "test-deployment"
                }
              ],
              "type": "kubernetes",
              "zones": []
            },
            {
              "account": "test-account",
              "apiVersion": "apps/v1",
              "cloudProvider": "kubernetes",
              "displayName": "test-rs4",
              "kind": "ReplicaSet",
              "labels": {
                "labelKey1": "labelValue1",
                "labelKey2": "labelValue2"
              },
              "moniker": {
                "app": "test-application",
                "cluster": "deployment test-deployment",
                "sequence": 19
              },
              "name": "replicaSet test-rs4",
              "namespace": "test-namespace1",
              "region": "test-namespace1",
              "serverGroupManagers": [
                {
                  "account": "test-account",
                  "location": "test-namespace1",
                  "name": "test-deployment"
                }
              ],
              "type": "kubernetes",
              "zones": []
            }
          ],
          "type": "kubernetes"
        }`

const payloadGetServerGroupManagedLoadBalancersMalformed = `{
              "account": "test-account",
              "accountName": "test-account",
              "buildInfo": {
                "images": [
                  "test-image3",
                  "test-image4"
                ]
              },
              "capacity": {
                "desired": 1,
                "pinned": false
              },
              "cloudProvider": "kubernetes",
              "createdTime": 1581603123000,
              "disabled": false,
              "displayName": "test-rs1",
              "instanceCounts": {
                "down": 0,
                "outOfService": 0,
                "starting": 0,
                "total": 1,
                "unknown": 0,
                "up": 0
              },
              "instances": [],
              "isDisabled": false,
              "key": {
                "account": "test-account",
                "group": "replicaSet",
                "kubernetesKind": "replicaSet",
                "name": "test-rs1",
                "namespace": "test-namespace1",
                "provider": "kubernetes"
              },
              "kind": "replicaSet",
              "labels": null,
              "loadBalancers": [],
              "manifest": {
                "apiVersion": "apps/v1",
                "kind": "ReplicaSet",
                "metadata": {
                  "annotations": {
                    "artifact.spinnaker.io/location": "test-namespace1",
                    "artifact.spinnaker.io/name": "test-deployment1",
                    "artifact.spinnaker.io/type": "kubernetes/deployment",
                    "moniker.spinnaker.io/application": "test-application",
                    "moniker.spinnaker.io/cluster": "deployment test-deployment1",
                    "moniker.spinnaker.io/sequence": "19",
                    "traffic.spinnaker.io/load-balancers": "[malformed]"
                  },
                  "creationTimestamp": "2020-02-13T14:12:03Z",
                  "name": "test-rs1",
                  "namespace": "test-namespace1",
                  "ownerReferences": [
                    {
                      "kind": "Deployment",
                      "name": "test-deployment1",
                      "uid": "test-uid3"
                    }
                  ],
                  "uid": "test-uid1"
                },
                "spec": {
                  "replicas": 1,
                  "template": {
                    "metadata": {
                      "labels": {
                        "test": "label"
                      }
                    },
                    "spec": {
                      "containers": [
                        {
                          "image": "test-image3"
                        },
                        {
                          "image": "test-image4"
                        }
                      ]
                    }
                  }
                },
                "status": {
                  "readyReplicas": 0,
                  "replicas": 1
                }
              },
              "moniker": {
                "app": "test-application",
                "cluster": "deployment test-deployment1",
                "sequence": 19
              },
              "name": "replicaSet test-rs1",
              "namespace": "test-namespace1",
              "providerType": "kubernetes",
              "region": "test-namespace1",
              "securityGroups": [],
              "serverGroupManagers": [],
              "type": "kubernetes",
              "uid": "test-uid1",
              "zone": "test-namespace1",
              "zones": [],
              "insightActions": []
            }`

const payloadGetServerGroupManagedLoadBalancersNoKind = `{
              "account": "test-account",
              "accountName": "test-account",
              "buildInfo": {
                "images": [
                  "test-image3",
                  "test-image4"
                ]
              },
              "capacity": {
                "desired": 1,
                "pinned": false
              },
              "cloudProvider": "kubernetes",
              "createdTime": 1581603123000,
              "disabled": false,
              "displayName": "test-rs1",
              "instanceCounts": {
                "down": 0,
                "outOfService": 0,
                "starting": 0,
                "total": 1,
                "unknown": 0,
                "up": 0
              },
              "instances": [],
              "isDisabled": false,
              "key": {
                "account": "test-account",
                "group": "replicaSet",
                "kubernetesKind": "replicaSet",
                "name": "test-rs1",
                "namespace": "test-namespace1",
                "provider": "kubernetes"
              },
              "kind": "replicaSet",
              "labels": null,
              "loadBalancers": [],
              "manifest": {
                "apiVersion": "apps/v1",
                "kind": "ReplicaSet",
                "metadata": {
                  "annotations": {
                    "artifact.spinnaker.io/location": "test-namespace1",
                    "artifact.spinnaker.io/name": "test-deployment1",
                    "artifact.spinnaker.io/type": "kubernetes/deployment",
                    "moniker.spinnaker.io/application": "test-application",
                    "moniker.spinnaker.io/cluster": "deployment test-deployment1",
                    "moniker.spinnaker.io/sequence": "19",
                    "traffic.spinnaker.io/load-balancers": "[\"my-managed-service-no-kind\"]"
                  },
                  "creationTimestamp": "2020-02-13T14:12:03Z",
                  "name": "test-rs1",
                  "namespace": "test-namespace1",
                  "ownerReferences": [
                    {
                      "kind": "Deployment",
                      "name": "test-deployment1",
                      "uid": "test-uid3"
                    }
                  ],
                  "uid": "test-uid1"
                },
                "spec": {
                  "replicas": 1,
                  "template": {
                    "metadata": {
                      "labels": {
                        "test": "label"
                      }
                    },
                    "spec": {
                      "containers": [
                        {
                          "image": "test-image3"
                        },
                        {
                          "image": "test-image4"
                        }
                      ]
                    }
                  }
                },
                "status": {
                  "readyReplicas": 0,
                  "replicas": 1
                }
              },
              "moniker": {
                "app": "test-application",
                "cluster": "deployment test-deployment1",
                "sequence": 19
              },
              "name": "replicaSet test-rs1",
              "namespace": "test-namespace1",
              "providerType": "kubernetes",
              "region": "test-namespace1",
              "securityGroups": [],
              "serverGroupManagers": [],
              "type": "kubernetes",
              "uid": "test-uid1",
              "zone": "test-namespace1",
              "zones": [],
              "insightActions": []
            }`

const payloadGetServerGroupManagedLoadBalancersWrongKind = `{
              "account": "test-account",
              "accountName": "test-account",
              "buildInfo": {
                "images": [
                  "test-image3",
                  "test-image4"
                ]
              },
              "capacity": {
                "desired": 1,
                "pinned": false
              },
              "cloudProvider": "kubernetes",
              "createdTime": 1581603123000,
              "disabled": false,
              "displayName": "test-rs1",
              "instanceCounts": {
                "down": 0,
                "outOfService": 0,
                "starting": 0,
                "total": 1,
                "unknown": 0,
                "up": 0
              },
              "instances": [],
              "isDisabled": false,
              "key": {
                "account": "test-account",
                "group": "replicaSet",
                "kubernetesKind": "replicaSet",
                "name": "test-rs1",
                "namespace": "test-namespace1",
                "provider": "kubernetes"
              },
              "kind": "replicaSet",
              "labels": null,
              "loadBalancers": [],
              "manifest": {
                "apiVersion": "apps/v1",
                "kind": "ReplicaSet",
                "metadata": {
                  "annotations": {
                    "artifact.spinnaker.io/location": "test-namespace1",
                    "artifact.spinnaker.io/name": "test-deployment1",
                    "artifact.spinnaker.io/type": "kubernetes/deployment",
                    "moniker.spinnaker.io/application": "test-application",
                    "moniker.spinnaker.io/cluster": "deployment test-deployment1",
                    "moniker.spinnaker.io/sequence": "19",
                    "traffic.spinnaker.io/load-balancers": "[\"ingress my-managed-ingress\"]"
                  },
                  "creationTimestamp": "2020-02-13T14:12:03Z",
                  "name": "test-rs1",
                  "namespace": "test-namespace1",
                  "ownerReferences": [
                    {
                      "kind": "Deployment",
                      "name": "test-deployment1",
                      "uid": "test-uid3"
                    }
                  ],
                  "uid": "test-uid1"
                },
                "spec": {
                  "replicas": 1,
                  "template": {
                    "metadata": {
                      "labels": {
                        "test": "label"
                      }
                    },
                    "spec": {
                      "containers": [
                        {
                          "image": "test-image3"
                        },
                        {
                          "image": "test-image4"
                        }
                      ]
                    }
                  }
                },
                "status": {
                  "readyReplicas": 0,
                  "replicas": 1
                }
              },
              "moniker": {
                "app": "test-application",
                "cluster": "deployment test-deployment1",
                "sequence": 19
              },
              "name": "replicaSet test-rs1",
              "namespace": "test-namespace1",
              "providerType": "kubernetes",
              "region": "test-namespace1",
              "securityGroups": [],
              "serverGroupManagers": [],
              "type": "kubernetes",
              "uid": "test-uid1",
              "zone": "test-namespace1",
              "zones": [],
              "insightActions": []
            }`

const payloadGetServerGroupManagedLoadBalancersErrorGettingService = `{
              "account": "test-account",
              "accountName": "test-account",
              "buildInfo": {
                "images": [
                  "test-image3",
                  "test-image4"
                ]
              },
              "capacity": {
                "desired": 1,
                "pinned": false
              },
              "cloudProvider": "kubernetes",
              "createdTime": 1581603123000,
              "disabled": false,
              "displayName": "test-rs1",
              "instanceCounts": {
                "down": 0,
                "outOfService": 0,
                "starting": 0,
                "total": 1,
                "unknown": 0,
                "up": 0
              },
              "instances": [],
              "isDisabled": false,
              "key": {
                "account": "test-account",
                "group": "replicaSet",
                "kubernetesKind": "replicaSet",
                "name": "test-rs1",
                "namespace": "test-namespace1",
                "provider": "kubernetes"
              },
              "kind": "replicaSet",
              "labels": null,
              "loadBalancers": [],
              "manifest": {
                "apiVersion": "apps/v1",
                "kind": "ReplicaSet",
                "metadata": {
                  "annotations": {
                    "artifact.spinnaker.io/location": "test-namespace1",
                    "artifact.spinnaker.io/name": "test-deployment1",
                    "artifact.spinnaker.io/type": "kubernetes/deployment",
                    "moniker.spinnaker.io/application": "test-application",
                    "moniker.spinnaker.io/cluster": "deployment test-deployment1",
                    "moniker.spinnaker.io/sequence": "19",
                    "traffic.spinnaker.io/load-balancers": "[\"service my-managed-service\"]"
                  },
                  "creationTimestamp": "2020-02-13T14:12:03Z",
                  "name": "test-rs1",
                  "namespace": "test-namespace1",
                  "ownerReferences": [
                    {
                      "kind": "Deployment",
                      "name": "test-deployment1",
                      "uid": "test-uid3"
                    }
                  ],
                  "uid": "test-uid1"
                },
                "spec": {
                  "replicas": 1,
                  "template": {
                    "metadata": {
                      "labels": {
                        "test": "label"
                      }
                    },
                    "spec": {
                      "containers": [
                        {
                          "image": "test-image3"
                        },
                        {
                          "image": "test-image4"
                        }
                      ]
                    }
                  }
                },
                "status": {
                  "readyReplicas": 0,
                  "replicas": 1
                }
              },
              "moniker": {
                "app": "test-application",
                "cluster": "deployment test-deployment1",
                "sequence": 19
              },
              "name": "replicaSet test-rs1",
              "namespace": "test-namespace1",
              "providerType": "kubernetes",
              "region": "test-namespace1",
              "securityGroups": [],
              "serverGroupManagers": [],
              "type": "kubernetes",
              "uid": "test-uid1",
              "zone": "test-namespace1",
              "zones": [],
              "insightActions": []
            }`

const payloadGetServerGroupManagedLoadBalancersDisabled = `{
              "account": "test-account",
              "accountName": "test-account",
              "buildInfo": {
                "images": [
                  "test-image3",
                  "test-image4"
                ]
              },
              "capacity": {
                "desired": 1,
                "pinned": false
              },
              "cloudProvider": "kubernetes",
              "createdTime": 1581603123000,
              "disabled": true,
              "displayName": "test-rs1",
              "instanceCounts": {
                "down": 0,
                "outOfService": 0,
                "starting": 0,
                "total": 1,
                "unknown": 0,
                "up": 0
              },
              "instances": [],
              "isDisabled": false,
              "key": {
                "account": "test-account",
                "group": "replicaSet",
                "kubernetesKind": "replicaSet",
                "name": "test-rs1",
                "namespace": "test-namespace1",
                "provider": "kubernetes"
              },
              "kind": "replicaSet",
              "labels": null,
              "loadBalancers": [
                "service my-managed-service"
							],
              "manifest": {
                "apiVersion": "apps/v1",
                "kind": "ReplicaSet",
                "metadata": {
                  "annotations": {
                    "artifact.spinnaker.io/location": "test-namespace1",
                    "artifact.spinnaker.io/name": "test-deployment1",
                    "artifact.spinnaker.io/type": "kubernetes/deployment",
                    "moniker.spinnaker.io/application": "test-application",
                    "moniker.spinnaker.io/cluster": "deployment test-deployment1",
                    "moniker.spinnaker.io/sequence": "19",
                    "traffic.spinnaker.io/load-balancers": "[\"service my-managed-service\"]"
                  },
                  "creationTimestamp": "2020-02-13T14:12:03Z",
                  "name": "test-rs1",
                  "namespace": "test-namespace1",
                  "ownerReferences": [
                    {
                      "kind": "Deployment",
                      "name": "test-deployment1",
                      "uid": "test-uid3"
                    }
                  ],
                  "uid": "test-uid1"
                },
                "spec": {
                  "replicas": 1,
                  "template": {
                    "metadata": {
                      "labels": {
                        "test": "label"
                      }
                    },
                    "spec": {
                      "containers": [
                        {
                          "image": "test-image3"
                        },
                        {
                          "image": "test-image4"
                        }
                      ]
                    }
                  }
                },
                "status": {
                  "readyReplicas": 0,
                  "replicas": 1
                }
              },
              "moniker": {
                "app": "test-application",
                "cluster": "deployment test-deployment1",
                "sequence": 19
              },
              "name": "replicaSet test-rs1",
              "namespace": "test-namespace1",
              "providerType": "kubernetes",
              "region": "test-namespace1",
              "securityGroups": [],
              "serverGroupManagers": [],
              "type": "kubernetes",
              "uid": "test-uid1",
              "zone": "test-namespace1",
              "zones": [],
              "insightActions": []
            }`

const payloadGetServerGroupManagedLoadBalancers = `{
              "account": "test-account",
              "accountName": "test-account",
              "buildInfo": {
                "images": [
                  "test-image3",
                  "test-image4"
                ]
              },
              "capacity": {
                "desired": 1,
                "pinned": false
              },
              "cloudProvider": "kubernetes",
              "createdTime": 1581603123000,
              "disabled": false,
              "displayName": "test-rs1",
              "instanceCounts": {
                "down": 0,
                "outOfService": 0,
                "starting": 0,
                "total": 1,
                "unknown": 0,
                "up": 0
              },
              "instances": [],
              "isDisabled": false,
              "key": {
                "account": "test-account",
                "group": "replicaSet",
                "kubernetesKind": "replicaSet",
                "name": "test-rs1",
                "namespace": "test-namespace1",
                "provider": "kubernetes"
              },
              "kind": "replicaSet",
              "labels": null,
              "loadBalancers": [
                "service my-managed-service"
							],
              "manifest": {
                "apiVersion": "apps/v1",
                "kind": "ReplicaSet",
                "metadata": {
                  "annotations": {
                    "artifact.spinnaker.io/location": "test-namespace1",
                    "artifact.spinnaker.io/name": "test-deployment1",
                    "artifact.spinnaker.io/type": "kubernetes/deployment",
                    "moniker.spinnaker.io/application": "test-application",
                    "moniker.spinnaker.io/cluster": "deployment test-deployment1",
                    "moniker.spinnaker.io/sequence": "19",
                    "traffic.spinnaker.io/load-balancers": "[\"service my-managed-service\"]"
                  },
                  "creationTimestamp": "2020-02-13T14:12:03Z",
                  "name": "test-rs1",
                  "namespace": "test-namespace1",
                  "ownerReferences": [
                    {
                      "kind": "Deployment",
                      "name": "test-deployment1",
                      "uid": "test-uid3"
                    }
                  ],
                  "uid": "test-uid1"
                },
                "spec": {
                  "replicas": 1,
                  "template": {
                    "metadata": {
                      "labels": {
                        "test": "label"
                      }
                    },
                    "spec": {
                      "containers": [
                        {
                          "image": "test-image3"
                        },
                        {
                          "image": "test-image4"
                        }
                      ]
                    }
                  }
                },
                "status": {
                  "readyReplicas": 0,
                  "replicas": 1
                }
              },
              "moniker": {
                "app": "test-application",
                "cluster": "deployment test-deployment1",
                "sequence": 19
              },
              "name": "replicaSet test-rs1",
              "namespace": "test-namespace1",
              "providerType": "kubernetes",
              "region": "test-namespace1",
              "securityGroups": [],
              "serverGroupManagers": [],
              "type": "kubernetes",
              "uid": "test-uid1",
              "zone": "test-namespace1",
              "zones": [],
              "insightActions": []
            }`

const payloadGetServerGroup = `{
            "account": "test-account",
            "accountName": "test-account",
            "buildInfo": {
              "images": [
                "test-image3",
                "test-image4"
              ]
            },
            "capacity": {
              "desired": 1,
              "pinned": false
            },
            "cloudProvider": "kubernetes",
            "createdTime": 1581603123000,
            "disabled": false,
            "displayName": "test-rs1",
            "instanceCounts": {
              "down": 0,
              "outOfService": 0,
              "starting": 0,
              "total": 1,
              "unknown": 0,
              "up": 0
            },
            "instances": [
              {
                "account": "test-account",
                "accountName": "test-account",
                "availabilityZone": "test-namespace1",
                "cloudProvider": "kubernetes",
                "createdTime": 1581603123000,
                "health": [
                  {
                    "state": "Down",
                    "type": "kubernetes/pod"
                  },
                  {
                    "state": "Down",
                    "type": "kubernetes/container"
                  }
                ],
                "healthState": "Down",
                "humanReadableName": "pod test-pod1",
                "id": "cec15437-4e6a-11ea-9788-4201ac100006",
                "key": {
                  "account": "test-account",
                  "group": "pod",
                  "kubernetesKind": "pod",
                  "name": "test-pod1",
                  "namespace": "test-namespace1",
                  "provider": "kubernetes"
                },
                "kind": "pod",
                "labels": {
                  "label1": "test-label1"
                },
                "manifest": {
                  "apiVersion": "v1",
                  "kind": "Pod",
                  "metadata": {
										"annotations": {
											"moniker.spinnaker.io/application": "test-application"
										},
                    "creationTimestamp": "2020-02-13T14:12:03Z",
                    "labels": {
                      "label1": "test-label1"
                    },
                    "name": "test-pod1",
                    "namespace": "test-namespace1",
                    "ownerReferences": [
                      {
                        "name": "test-rs1"
                      }
                    ],
                    "uid": "cec15437-4e6a-11ea-9788-4201ac100006"
                  }
                },
                "moniker": {
                  "app": "test-application",
                  "cluster": ""
                },
                "name": "pod test-pod1",
                "providerType": "kubernetes",
                "region": "test-namespace1",
                "type": "kubernetes",
                "uid": "cec15437-4e6a-11ea-9788-4201ac100006",
                "zone": "test-namespace1"
              },
              {
                "account": "test-account",
                "accountName": "test-account",
                "availabilityZone": "test-namespace2",
                "cloudProvider": "kubernetes",
                "createdTime": 1581603123000,
                "health": [
                  {
                    "state": "Down",
                    "type": "kubernetes/pod"
                  },
                  {
                    "state": "Down",
                    "type": "kubernetes/container"
                  }
                ],
                "healthState": "Down",
                "humanReadableName": "pod test-pod2",
                "id": "cec15437-4e6a-11ea-9788-4201ac100006",
                "key": {
                  "account": "test-account",
                  "group": "pod",
                  "kubernetesKind": "pod",
                  "name": "test-pod2",
                  "namespace": "test-namespace2",
                  "provider": "kubernetes"
                },
                "kind": "pod",
                "labels": {
                  "label1": "test-label1"
                },
                "manifest": {
                  "apiVersion": "v1",
                  "kind": "Pod",
                  "metadata": {
										"annotations": {
											"moniker.spinnaker.io/application": "test-application"
										},
                    "creationTimestamp": "2020-02-13T14:12:03Z",
                    "labels": {
                      "label1": "test-label1"
                    },
                    "name": "test-pod2",
                    "namespace": "test-namespace2",
                    "ownerReferences": [
                      {
                        "name": "test-rs1"
                      }
                    ],
                    "uid": "cec15437-4e6a-11ea-9788-4201ac100006"
                  }
                },
                "moniker": {
                  "app": "test-application",
                  "cluster": ""
                },
                "name": "pod test-pod2",
                "providerType": "kubernetes",
                "region": "test-namespace2",
                "type": "kubernetes",
                "uid": "cec15437-4e6a-11ea-9788-4201ac100006",
                "zone": "test-namespace2"
              }
            ],
						"isDisabled": false,
            "key": {
              "account": "test-account",
              "group": "replicaSet",
              "kubernetesKind": "replicaSet",
              "name": "test-rs1",
              "namespace": "test-namespace1",
              "provider": "kubernetes"
            },
            "kind": "replicaSet",
            "labels": null,
            "loadBalancers": [],
            "manifest": {
              "apiVersion": "apps/v1",
              "kind": "ReplicaSet",
              "metadata": {
                "annotations": {
                  "artifact.spinnaker.io/location": "test-namespace2",
                  "artifact.spinnaker.io/name": "test-deployment2",
                  "artifact.spinnaker.io/type": "kubernetes/deployment",
                  "deployment.kubernetes.io/revision": "19",
                  "moniker.spinnaker.io/application": "test-application",
                  "moniker.spinnaker.io/cluster": "deployment test-deployment1"
                },
                "creationTimestamp": "2020-02-13T14:12:03Z",
                "name": "test-rs1",
                "namespace": "test-namespace1"
              },
							"spec": {
                "replicas": 1,
                "template": {
                  "metadata": {
                    "labels": {
                      "selectorKey1": "selectorValue1"
                    }
                  },
                  "spec": {
                    "containers": [
                      {
                        "image": "test-image3"
                      },
                      {
                        "image": "test-image4"
                      }
                    ]
                  }
                }
              },
              "status": {
                "readyReplicas": 0,
                "replicas": 1
              }
            },
            "moniker": {
              "app": "test-application",
              "cluster": "deployment test-deployment1",
              "sequence": 19
            },
            "name": "replicaSet test-rs1",
            "namespace": "test-namespace1",
            "providerType": "kubernetes",
            "region": "test-namespace1",
            "securityGroups": [],
            "serverGroupManagers": [],
            "type": "kubernetes",
            "uid": "",
            "zone": "test-namespace1",
            "zones": [],
            "insightActions": []
          }`

const payloadGetJob = `{
            "account": "test-account",
            "completionDetails": {
              "exitCode": "",
              "message": "",
              "reason": "",
              "signal": ""
            },
            "createdTime": 1581603123000,
            "jobState": "Running",
            "location": "test-namespace",
            "name": "test-job1",
            "pods": [],
            "provider": "kubernetes"
          }`

const payloadSearchDefault = `[
							 {
								 "pageNumber": 1,
								 "pageSize": 500,
								 "platform": "aws",
								 "query": "default",
								 "results": [],
								 "totalMatches": 0
							 }
						 ]`

const payloadSearchEmptyResponse = `[
							 {
								 "pageNumber": 1,
								 "pageSize": 500,
								 "query": "default",
								 "results": [],
								 "totalMatches": 0
							 }
						 ]`

const payloadSearch = `[
          {
            "pageNumber": 1,
            "pageSize": 500,
            "query": "default",
            "results": [
              {
                "account": "account1",
                "group": "pod",
                "kubernetesKind": "pod",
                "name": "pod test-name",
                "namespace": "default",
                "provider": "kubernetes",
                "region": "default",
                "type": "instances"
              }
            ],
            "totalMatches": 1
          }
        ]`

const payloadManifestCoordinates = `{
            "kind": "test-kind",
            "name": "test-name",
            "namespace": "test-namespace"
          }`

const payloadManifestCoordinatesList = `[
            {
              "kind": "replicaSet",
              "name": "rs2-v000",
              "namespace": "test-namespace"
            },
            {
              "kind": "replicaSet",
              "name": "rs2-v001",
              "namespace": "test-namespace"
            }
          ]`

const payloadManifestClusterRoleNoRules = `{
            "account": "test-account",
						"artifacts": [],
            "events": [],
            "location": "test-namespace",
            "manifest": {
              "apiVersion": "rbac.authorization.k8s.io/v1",
              "kind": "ClusterRole",
              "metadata": {
                "annotations": {
                  "artifact.spinnaker.io/location": "",
                  "artifact.spinnaker.io/name": "test-cluster-role",
                  "artifact.spinnaker.io/type": "kubernetes/clusterRole",
                  "artifact.spinnaker.io/version": "",
                  "moniker.spinnaker.io/application": "test-application",
                  "moniker.spinnaker.io/cluster": "clusterRole test-cluster-role"
                },
                "creationTimestamp": "2021-10-20T15:29:26Z",
                "labels": {
                  "app.kubernetes.io/managed-by": "spinnaker",
                  "app.kubernetes.io/name": "test-application"
                },
                "name": "test-cluster-role",
                "resourceVersion": "53990465",
                "uid": "d1f1ab80-1320-4e2d-8d12-893c326af416"
              }
            },
            "metrics": [],
            "moniker": {
              "app": "test-application",
              "cluster": "clusterRole test-cluster-role"
            },
            "name": "clusterRole test-cluster-role",
            "status": {
              "available": {
                "state": true,
                "message": ""
              },
              "failed": {
                "state": false,
                "message": ""
              },
              "paused": {
                "state": false,
                "message": ""
              },
              "stable": {
                "state": true,
                "message": ""
              }
            },
            "warnings": []
          }`

const payloadManifestWithArtifacts = `{
        "account": "test-account",
        "artifacts": [
          {
            "customKind": false,
            "metadata": {},
            "name": "gcr.io/test-project/test-container-image",
            "reference": "gcr.io/test-project/test-container-image",
            "type": "docker/image"
          },
          {
            "customKind": false,
            "metadata": {},
            "name": "gcr.io/test-project/another-test-container-image",
            "reference": "gcr.io/test-project/another-test-container-image",
            "type": "docker/image"
          }
        ],
        "events": [],
        "location": "test-namespace",
        "manifest": {
          "apiVersion": "apps/v1",
          "kind": "Deployment",
          "spec": {
            "template": {
              "spec": {
                "containers": [
                  {
                    "image": "gcr.io/test-project/test-container-image",
                    "name": "test-container-name"
                  },
                  {
                    "image": "gcr.io/test-project/another-test-container-image",
                    "name": "another-test-container-name"
                  }
                ]
              }
            }
          }
        },
        "metrics": [],
        "moniker": {
          "app": "unknown",
          "cluster": "deployment test-deployment"
        },
        "name": "deployment test-deployment",
        "status": {
          "available": {
            "state": true,
            "message": ""
          },
          "failed": {
            "state": false,
            "message": ""
          },
          "paused": {
            "state": false,
            "message": ""
          },
          "stable": {
            "state": true,
            "message": ""
          }
        },
        "warnings": []
      }`

const payloadManifestNoEvents = `{
            "account": "test-account",
            "artifacts": [],
            "events": [],
            "location": "test-namespace",
            "manifest": {},
            "metrics": [],
            "moniker": {
              "app": "unknown",
              "cluster": "pod test-pod"
            },
            "name": "pod test-pod",
            "status": {
              "available": {
                "state": true,
                "message": ""
              },
              "failed": {
                "state": false,
                "message": ""
              },
              "paused": {
                "state": false,
                "message": ""
              },
              "stable": {
                "state": true,
                "message": ""
              }
            },
            "warnings": []
          }`

const payloadManifestIncludeEvents = `{
              "account": "test-account",
              "artifacts": [],
              "events": [
                {
                  "kind": "test-kind",
                  "apiVersion": "test-api-version",
                  "metadata": {
                    "name": "test-event-name",
                    "generateName": "test-event-generate-name",
                    "namespace": "test-event-namespace",
                    "creationTimestamp": null
                  },
                  "involvedObject": {
                    "kind": "test-kind",
                    "namespace": "test-namespace",
                    "name": "test-pod"
                  },
                  "reason": "test reason",
                  "message": "test message",
                  "source": {},
                  "firstTimestamp": null,
                  "lastTimestamp": null,
                  "count": 1,
                  "eventTime": null,
                  "reportingComponent": "",
                  "reportingInstance": ""
                },
                {
                  "kind": "test-kind",
                  "apiVersion": "test-api-version",
                  "metadata": {
                    "name": "test-event-name2",
                    "generateName": "test-event-generate-name",
                    "namespace": "test-event-namespace",
                    "creationTimestamp": null
                  },
                  "involvedObject": {
                    "kind": "test-kind",
                    "namespace": "test-namespace",
                    "name": "test-pod"
                  },
                  "reason": "test reason",
                  "message": "test message",
                  "source": {},
                  "firstTimestamp": null,
                  "lastTimestamp": null,
                  "count": 2,
                  "eventTime": null,
                  "reportingComponent": "",
                  "reportingInstance": ""
                }
              ],
              "location": "test-namespace",
              "manifest": {},
              "metrics": [],
              "moniker": {
                "app": "unknown",
                "cluster": "pod test-pod"
              },
              "name": "pod test-pod",
              "status": {
                "available": {
                  "state": true,
                  "message": ""
                },
                "failed": {
                  "state": false,
                  "message": ""
                },
                "paused": {
                  "state": false,
                  "message": ""
                },
                "stable": {
                  "state": true,
                  "message": ""
                }
              },
              "warnings": []
            }`

const payloadGetInstance = `{
            "account": "test-account",
            "apiVersion": "v1",
            "cloudProvider": "kubernetes",
            "createdTime": 1581603123000,
            "displayName": "test-pod1",
            "health": [
              {
                "platform": "platform",
                "source": "Pod",
                "state": "Up",
                "type": "kubernetes/pod"
              },
              {
                "platform": "platform",
                "source": "Container test-container-name",
                "state": "Up",
                "type": "kubernetes/container"
              }
            ],
            "healthState": "Up",
            "humanReadableName": "pod test-pod1",
            "kind": "pod",
            "labels": {
              "label1": "test-label1"
            },
            "moniker": {
              "app": "test-application",
              "cluster": "test cluster"
            },
            "name": "cec15437-4e6a-11ea-9788-4201ac100006",
            "namespace": "test-namespace1",
            "providerType": "kubernetes",
            "zone": "test-namespace1"
          }`

const payloadGetDownInstance = `{
            "account": "test-account",
            "apiVersion": "v1",
            "cloudProvider": "kubernetes",
            "createdTime": 1581603123000,
            "displayName": "test-pod1",
            "health": [
              {
                "platform": "platform",
                "source": "Pod",
                "state": "Down",
                "type": "kubernetes/pod"
              },
              {
                "platform": "platform",
                "source": "Container test-container-name",
                "state": "Up",
                "type": "kubernetes/container"
              }
            ],
            "healthState": "Down",
            "humanReadableName": "pod test-pod1",
            "kind": "pod",
            "labels": {
              "label1": "test-label1"
            },
            "moniker": {
              "app": "test-application",
              "cluster": "test cluster"
            },
            "name": "cec15437-4e6a-11ea-9788-4201ac100006",
            "namespace": "test-namespace1",
            "providerType": "kubernetes",
            "zone": "test-namespace1"
          }`

const payloadGetInstanceConsole = `{
            "output": [
              {
                "name": "test-container-name",
                "output": "log output"
              },
              {
                "name": "test-init-container-name",
                "output": "log output"
              }
            ]
          }`

const payloadListProjectClustersNoMatches = `[
	{
		"account": "test-account-1",
		"applications": [
			{
				"application": "fake-application",
				"clusters": [],
				"lastPush": 0
			}
		],
		"detail": "*",
		"instanceCounts": {
			"down": 0,
			"outOfService": 0,
			"starting": 0,
			"total": 0,
			"unknown": 0,
			"up": 0
		},
		"stack": "*"
	}
]`

const payloadListProjectClusters = `[
	{
		"account": "test-account-1",
		"applications": [
			{
				"application": "test-application-1",
				"clusters": [
					{
						"builds": [
							{
								"buildNumber": "0",
								"deployed": 1581603123000,
								"images": [
									"test-image-1"
								]
							}
						],
						"instanceCounts": {
							"down": 0,
							"outOfService": 0,
							"starting": 0,
							"total": 4,
							"unknown": 0,
							"up": 2
						},
						"lastPush": 1581603123000,
						"region": "test-namespace-1"
					}
				],
				"lastPush": 1581603123000
			}
		],
		"detail": "test-detail",
		"instanceCounts": {
			"down": 0,
			"outOfService": 0,
			"starting": 0,
			"total": 4,
			"unknown": 0,
			"up": 2
		},
		"stack": "test-stack"
	},
	{
		"account": "test-account-1",
		"applications": [
			{
				"application": "test-application-1",
				"clusters": [
					{
						"builds": [
							{
								"buildNumber": "0",
								"deployed": 1581603123000,
								"images": [
									"test-image-1",
									"test-image-2"
								]
							}
						],
						"instanceCounts": {
							"down": 0,
							"outOfService": 0,
							"starting": 0,
							"total": 8,
							"unknown": 0,
							"up": 4
						},
						"lastPush": 1581603123000,
						"region": "test-namespace-1"
					}
				],
				"lastPush": 1581603123000
			}
		],
		"detail": "",
		"instanceCounts": {
			"down": 0,
			"outOfService": 0,
			"starting": 0,
			"total": 8,
			"unknown": 0,
			"up": 4
		},
		"stack": ""
	},
	{
		"account": "test-account-1",
		"applications": [
			{
				"application": "test-application-1",
				"clusters": [
					{
						"builds": [
							{
								"buildNumber": "0",
								"deployed": 1581603123000,
								"images": [
									"test-image-1",
									"test-image-2"
								]
							}
						],
						"instanceCounts": {
							"down": 0,
							"outOfService": 0,
							"starting": 0,
							"total": 12,
							"unknown": 0,
							"up": 6
						},
						"lastPush": 1581603123000,
						"region": "test-namespace-1"
					},
					{
						"builds": [
							{
								"buildNumber": "0",
								"deployed": 1581603123000,
								"images": [
									"test-image-3"
								]
							}
						],
						"instanceCounts": {
							"down": 0,
							"outOfService": 0,
							"starting": 0,
							"total": 16,
							"unknown": 0,
							"up": 8
						},
						"lastPush": 1581603123000,
						"region": "test-namespace-2"
					}
				],
				"lastPush": 1581603123000
			},
			{
				"application": "test-application-2",
				"clusters": [],
				"lastPush": 0
			}
		],
		"detail": "*",
		"instanceCounts": {
			"down": 0,
			"outOfService": 0,
			"starting": 0,
			"total": 28,
			"unknown": 0,
			"up": 14
		},
		"stack": "*"
	}
]`

const payloadTaskIncomplete = `{
              "id": "task-id",
              "resultObjects": [
                {
                  "boundArtifacts": [
                    {
                      "customKind": false,
                      "location": "test-namespace",
                      "metadata": {
                        "account": "test-account-name"
                      },
                      "name": "",
                      "reference": "test-deployment",
                      "type": "kubernetes/deployment"
                    }
                  ],
                  "createdArtifacts": [
                    {
                      "customKind": false,
                      "location": "test-namespace",
                      "metadata": {
                        "account": "test-account-name"
                      },
                      "name": "",
                      "reference": "test-deployment",
                      "type": "kubernetes/deployment"
                    }
                  ],
                  "deployedNamesByLocation": {
                    "test-namespace": [
                      "Deployment test-deployment"
                    ]
                  },
                  "manifestNamesByNamespace": {
                    "test-namespace": [
                      "Deployment test-deployment"
                    ]
                  },
                  "manifestNamesByNamespaceToRefresh": {
                    "test-namespace": [
                      "Deployment test-deployment"
                    ]
                  },
                  "manifests": [
                    {
                      "apiVersion": "apps/v1",
                      "kind": "Deployment",
                      "metadata": {
                        "name": "test-deployment",
                        "namespace": "test-namespace"
                      }
                    }
                  ]
                }
              ],
              "status": {
                "complete": false,
                "completed": false,
                "failed": false,
                "phase": "ORCHESTRATION",
                "retryable": false,
                "status": "Orchestration in progress."
              }
            }`

const payloadTaskComplete = `{
              "id": "task-id",
              "resultObjects": [
                {
                  "boundArtifacts": [
                    {
                      "customKind": false,
                      "location": "test-namespace",
                      "metadata": {
                        "account": "test-account-name"
                      },
                      "name": "",
                      "reference": "test-deployment",
                      "type": "kubernetes/deployment"
                    }
                  ],
                  "createdArtifacts": [
                    {
                      "customKind": false,
                      "location": "test-namespace",
                      "metadata": {
                        "account": "test-account-name"
                      },
                      "name": "",
                      "reference": "test-deployment",
                      "type": "kubernetes/deployment"
                    }
                  ],
                  "deployedNamesByLocation": {
                    "test-namespace": [
                      "Deployment test-deployment"
                    ]
                  },
                  "manifestNamesByNamespace": {
                    "test-namespace": [
                      "Deployment test-deployment"
                    ]
                  },
                  "manifestNamesByNamespaceToRefresh": {
                    "test-namespace": [
                      "Deployment test-deployment"
                    ]
                  },
                  "manifests": [
                    {}
                  ]
                }
              ],
              "status": {
                "complete": true,
                "completed": true,
                "failed": false,
                "phase": "ORCHESTRATION",
                "retryable": false,
                "status": "Orchestration completed."
              }
            }`
