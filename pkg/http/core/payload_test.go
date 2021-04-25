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
    "deleteManifest": {}
  }
]`

const payloadRequestKubernetesOpsCleanupArtifacts = `[
  {
    "cleanupArtifacts": {
      "manifests": [
        {}
      ]
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
    "patchManifest": {}
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
              "kind": "DaemonSet",
              "labels": null,
              "loadBalancers": null,
              "manifest": null,
              "moniker": {
                "app": "test-deployment1",
                "cluster": "deployment test-deployment1",
                "sequence": 19
              },
              "name": "DaemonSet test-ds1",
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
              "kind": "ReplicaSet",
              "labels": null,
              "loadBalancers": null,
              "manifest": null,
              "moniker": {
                "app": "test-deployment1",
                "cluster": "deployment test-deployment1",
                "sequence": 19
              },
              "name": "ReplicaSet test-rs1",
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
              "kind": "StatefulSet",
              "labels": null,
              "loadBalancers": null,
              "manifest": null,
              "moniker": {
                "app": "test-deployment1",
                "cluster": "deployment test-deployment1",
                "sequence": 19
              },
              "name": "StatefulSet test-rs1",
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
              "kind": "DaemonSet",
              "labels": null,
              "loadBalancers": null,
              "manifest": null,
              "moniker": {
                "app": "test-deployment1",
                "cluster": "deployment test-deployment1",
                "sequence": 19
              },
              "name": "DaemonSet test-ds1",
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
              "kind": "ReplicaSet",
              "labels": null,
              "loadBalancers": null,
              "manifest": null,
              "moniker": {
                "app": "test-deployment1",
                "cluster": "deployment test-deployment1",
                "sequence": 19
              },
              "name": "ReplicaSet test-rs1",
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
              "kind": "StatefulSet",
              "labels": null,
              "loadBalancers": null,
              "manifest": null,
              "moniker": {
                "app": "test-deployment1",
                "cluster": "deployment test-deployment1",
                "sequence": 19
              },
              "name": "StatefulSet test-rs1",
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
              "kind": "DaemonSet",
              "labels": null,
              "loadBalancers": null,
              "manifest": null,
              "moniker": {
                "app": "test-deployment1",
                "cluster": "deployment test-deployment1",
                "sequence": 19
              },
              "name": "DaemonSet test-ds1",
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
              "kind": "StatefulSet",
              "labels": null,
              "loadBalancers": null,
              "manifest": null,
              "moniker": {
                "app": "test-deployment1",
                "cluster": "deployment test-deployment1",
                "sequence": 19
              },
              "name": "StatefulSet test-rs1",
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
              "cloudProvider": "kubernetes",
              "labels": {
                "label1": "test-label1"
              },
              "moniker": {
                "app": "test-application",
                "cluster": "ingress test-ingress1"
              },
              "name": "ingress test-ingress1",
              "displayName": "test-ingress1",
              "region": "test-namespace1",
              "serverGroups": null,
              "type": "kubernetes",
              "accountName": "account1",
              "createdTime": 1581603123000,
              "key": {
                "account": "account1",
                "group": "networking.k8s.io",
                "kubernetesKind": "ingress",
                "name": "ingress test-ingress1",
                "namespace": "test-namespace1",
                "provider": "kubernetes"
              },
              "kind": "ingress",
              "manifest": {
                "apiVersion": "networking.k8s.io/v1beta1",
                "kind": "Ingress",
                "metadata": {
                  "creationTimestamp": "2020-02-13T14:12:03Z",
                  "labels": {
                    "label1": "test-label1"
                  },
                  "name": "test-ingress1",
                  "namespace": "test-namespace1",
                  "uid": "cec15437-4e6a-11ea-9788-4201ac100006"
                }
              },
              "providerType": "kubernetes",
              "uid": "cec15437-4e6a-11ea-9788-4201ac100006",
              "zone": "test-application"
            },
            {
              "account": "account1",
              "cloudProvider": "kubernetes",
              "moniker": {
                "app": "test-application",
                "cluster": "service test-service1"
              },
              "name": "service test-service1",
              "displayName": "test-service1",
              "region": "test-namespace1",
              "serverGroups": null,
              "type": "kubernetes",
              "accountName": "account1",
              "createdTime": -62135596800000,
              "key": {
                "account": "account1",
                "group": "",
                "kubernetesKind": "service",
                "name": "service test-service1",
                "namespace": "test-namespace1",
                "provider": "kubernetes"
              },
              "kind": "service",
              "manifest": {
                "apiVersion": "v1",
                "kind": "Service",
                "metadata": {
                  "name": "test-service1",
                  "namespace": "test-namespace1"
                }
              },
              "providerType": "kubernetes",
              "zone": "test-application"
            }
          ]`

const payloadListLoadBalancersSorted = `[
            {
              "account": "account1",
              "cloudProvider": "kubernetes",
              "labels": {
                "label1": "test-label1"
              },
              "moniker": {
                "app": "test-application",
                "cluster": "ingress test-ingress1"
              },
              "name": "ingress test-ingress1",
              "displayName": "test-ingress1",
              "region": "test-namespace1",
              "serverGroups": null,
              "type": "kubernetes",
              "accountName": "account1",
              "createdTime": 1581603123000,
              "key": {
                "account": "account1",
                "group": "networking.k8s.io",
                "kubernetesKind": "ingress",
                "name": "ingress test-ingress1",
                "namespace": "test-namespace1",
                "provider": "kubernetes"
              },
              "kind": "ingress",
              "manifest": {
                "apiVersion": "networking.k8s.io/v1beta1",
                "kind": "Ingress",
                "metadata": {
                  "creationTimestamp": "2020-02-13T14:12:03Z",
                  "labels": {
                    "label1": "test-label1"
                  },
                  "name": "test-ingress1",
                  "namespace": "test-namespace1",
                  "uid": "cec15437-4e6a-11ea-9788-4201ac100006"
                }
              },
              "providerType": "kubernetes",
              "uid": "cec15437-4e6a-11ea-9788-4201ac100006",
              "zone": "test-application"
            },
            {
              "account": "account1",
              "cloudProvider": "kubernetes",
              "moniker": {
                "app": "test-application",
                "cluster": "service test-service1"
              },
              "name": "service test-service1",
              "displayName": "test-service1",
              "region": "test-namespace1",
              "serverGroups": null,
              "type": "kubernetes",
              "accountName": "account1",
              "createdTime": -62135596800000,
              "key": {
                "account": "account1",
                "group": "",
                "kubernetesKind": "service",
                "name": "service test-service1",
                "namespace": "test-namespace1",
                "provider": "kubernetes"
              },
              "kind": "service",
              "manifest": {
                "apiVersion": "v1",
                "kind": "Service",
                "metadata": {
                  "name": "test-service1",
                  "namespace": "test-namespace1"
                }
              },
              "providerType": "kubernetes",
              "zone": "test-application"
            },
            {
              "account": "account1",
              "cloudProvider": "kubernetes",
              "labels": {
                "label1": "test-label1"
              },
              "moniker": {
                "app": "test-application",
                "cluster": "ingress test-ingress2"
              },
              "name": "ingress test-ingress2",
              "displayName": "test-ingress2",
              "region": "test-namespace2",
              "serverGroups": null,
              "type": "kubernetes",
              "accountName": "account1",
              "createdTime": 1581603123000,
              "key": {
                "account": "account1",
                "group": "networking.k8s.io",
                "kubernetesKind": "ingress",
                "name": "ingress test-ingress2",
                "namespace": "test-namespace2",
                "provider": "kubernetes"
              },
              "kind": "ingress",
              "manifest": {
                "apiVersion": "networking.k8s.io/v1beta1",
                "kind": "Ingress",
                "metadata": {
                  "creationTimestamp": "2020-02-13T14:12:03Z",
                  "labels": {
                    "label1": "test-label1"
                  },
                  "name": "test-ingress2",
                  "namespace": "test-namespace2",
                  "uid": "cec15437-4e6a-11ea-9788-4201ac100006"
                }
              },
              "providerType": "kubernetes",
              "uid": "cec15437-4e6a-11ea-9788-4201ac100006",
              "zone": "test-application"
            },
            {
              "account": "account1",
              "cloudProvider": "kubernetes",
              "labels": {
                "label1": "test-label1"
              },
              "moniker": {
                "app": "test-application",
                "cluster": "ingress test-ingress3"
              },
              "name": "ingress test-ingress3",
              "displayName": "test-ingress3",
              "region": "test-namespace2",
              "serverGroups": null,
              "type": "kubernetes",
              "accountName": "account1",
              "createdTime": 1581603123000,
              "key": {
                "account": "account1",
                "group": "networking.k8s.io",
                "kubernetesKind": "ingress",
                "name": "ingress test-ingress3",
                "namespace": "test-namespace2",
                "provider": "kubernetes"
              },
              "kind": "ingress",
              "manifest": {
                "apiVersion": "networking.k8s.io/v1beta1",
                "kind": "Ingress",
                "metadata": {
                  "creationTimestamp": "2020-02-13T14:12:03Z",
                  "labels": {
                    "label1": "test-label1"
                  },
                  "name": "test-ingress3",
                  "namespace": "test-namespace2",
                  "uid": "cec15437-4e6a-11ea-9788-4201ac100006"
                }
              },
              "providerType": "kubernetes",
              "uid": "cec15437-4e6a-11ea-9788-4201ac100006",
              "zone": "test-application"
            },
            {
              "account": "account1",
              "cloudProvider": "kubernetes",
              "moniker": {
                "app": "test-application",
                "cluster": "service test-service2"
              },
              "name": "service test-service2",
              "displayName": "test-service2",
              "region": "test-namespace2",
              "serverGroups": null,
              "type": "kubernetes",
              "accountName": "account1",
              "createdTime": -62135596800000,
              "key": {
                "account": "account1",
                "group": "",
                "kubernetesKind": "service",
                "name": "service test-service2",
                "namespace": "test-namespace2",
                "provider": "kubernetes"
              },
              "kind": "service",
              "manifest": {
                "apiVersion": "v1",
                "kind": "Service",
                "metadata": {
                  "name": "test-service2",
                  "namespace": "test-namespace2"
                }
              },
              "providerType": "kubernetes",
              "zone": "test-application"
            },
            {
              "account": "account1",
              "cloudProvider": "kubernetes",
              "moniker": {
                "app": "test-application",
                "cluster": "service test-service3"
              },
              "name": "service test-service3",
              "displayName": "test-service3",
              "region": "test-namespace2",
              "serverGroups": null,
              "type": "kubernetes",
              "accountName": "account1",
              "createdTime": -62135596800000,
              "key": {
                "account": "account1",
                "group": "",
                "kubernetesKind": "service",
                "name": "service test-service3",
                "namespace": "test-namespace2",
                "provider": "kubernetes"
              },
              "kind": "service",
              "manifest": {
                "apiVersion": "v1",
                "kind": "Service",
                "metadata": {
                  "name": "test-service3",
                  "namespace": "test-namespace2"
                }
              },
              "providerType": "kubernetes",
              "zone": "test-application"
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
              "group": "ReplicaSet",
              "kubernetesKind": "ReplicaSet",
              "name": "test-rs1",
              "namespace": "test-namespace1",
              "provider": "kubernetes"
            },
            "kind": "ReplicaSet",
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
                  "moniker.spinnaker.io/application": "test-deployment2",
                  "moniker.spinnaker.io/cluster": "deployment test-deployment1"
                },
                "creationTimestamp": "2020-02-13T14:12:03Z",
                "name": "test-rs1",
                "namespace": "test-namespace1"
              },
              "spec": {
                "replicas": 1,
                "template": {
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
              "app": "test-deployment2",
              "cluster": "deployment test-deployment1",
              "sequence": 19
            },
            "name": "ReplicaSet test-rs1",
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

const payloadSearchEmptyResponseWithPageSizeZero = `[
            {
              "pageNumber": 1,
              "pageSize": 0,
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
                  "name": "pod test-name1",
                  "namespace": "default",
                  "provider": "kubernetes",
                  "region": "default",
                  "type": "instances"
                },
                {
                  "account": "account1",
                  "group": "pod",
                  "kubernetesKind": "pod",
                  "name": "pod test-name2",
                  "namespace": "default",
                  "provider": "kubernetes",
                  "region": "default",
                  "type": "instances"
                },
                {
                  "account": "account2",
                  "group": "pod",
                  "kubernetesKind": "pod",
                  "name": "pod test-name1",
                  "namespace": "default",
                  "provider": "kubernetes",
                  "region": "default",
                  "type": "instances"
                },
                {
                  "account": "account2",
                  "group": "pod",
                  "kubernetesKind": "pod",
                  "name": "pod test-name2",
                  "namespace": "default",
                  "provider": "kubernetes",
                  "region": "default",
                  "type": "instances"
                }
              ],
              "totalMatches": 4
            }
          ]`

const payloadSearchWithPageSizeThree = `[
            {
              "pageNumber": 1,
              "pageSize": 3,
              "query": "default",
              "results": [
                {
                  "account": "account1",
                  "group": "pod",
                  "kubernetesKind": "pod",
                  "name": "pod test-name1",
                  "namespace": "default",
                  "provider": "kubernetes",
                  "region": "default",
                  "type": "instances"
                },
                {
                  "account": "account1",
                  "group": "pod",
                  "kubernetesKind": "pod",
                  "name": "pod test-name2",
                  "namespace": "default",
                  "provider": "kubernetes",
                  "region": "default",
                  "type": "instances"
                },
                {
                  "account": "account2",
                  "group": "pod",
                  "kubernetesKind": "pod",
                  "name": "pod test-name1",
                  "namespace": "default",
                  "provider": "kubernetes",
                  "region": "default",
                  "type": "instances"
                }
              ],
              "totalMatches": 3
            }
          ]`

const payloadManifestCoordinates = `{
            "kind": "test-kind",
            "name": "test-name",
            "namespace": "test-namespace"
          }`
