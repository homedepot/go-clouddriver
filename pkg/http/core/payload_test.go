package core_test

const payloadBadRequest = `{
            "error": "invalid character 'd' looking for beginning of value"
          }`

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
	"reference": "%s/api/v3/repos/billiford/kubernetes-engine-samples/contents/hello-app/manifests/helloweb-deployment.yaml"
}`

const payloadRequestFetchGithubFileArtifactTestBranch = `{
  "type": "github/file",
	"reference": "%s/api/v3/repos/billiford/kubernetes-engine-samples/contents/hello-app/manifests/helloweb-deployment.yaml",
	"version": "test"
}`

const payloadRequestFetchNotImplementedArtifact = `{
  "type": "unknown/type"
}`

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
      "source": "text",
      "account": "spin-cluster-account",
      "skipExpressionEvaluation": false,
      "requiredArtifacts": []
    }
  }
]`

const payloadRequestKubernetesOpsBadManifest = `[
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
					"metadata": 1
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

const payloadServerGroupManagersResponse = `[
  {
    "account": "gke_github-replication-sandbox_us-central1_sandbox-us-central1-agent-dev",
    "accountName": "gke_github-replication-sandbox_us-central1_sandbox-us-central1-agent-dev",
    "cloudProvider": "kubernetes",
    "createdTime": 1581619801000,
    "key": {
      "account": "gke_github-replication-sandbox_us-central1_sandbox-us-central1-agent-dev",
      "group": "deployment",
      "kubernetesKind": "deployment",
      "name": "smoketest-gke",
      "namespace": "smoketest",
      "provider": "kubernetes"
    },
    "kind": "deployment",
    "labels": {
      "app": "smoketest-gke",
      "app.kubernetes.io/managed-by": "spinnaker",
      "app.kubernetes.io/name": "smoketests"
    },
    "manifest": {
      "apiVersion": "extensions/v1beta1",
      "kind": "Deployment",
      "metadata": {
        "annotations": {
          "artifact.spinnaker.io/location": "smoketest",
          "artifact.spinnaker.io/name": "smoketest-gke",
          "artifact.spinnaker.io/type": "kubernetes/deployment",
          "deployment.kubernetes.io/revision": "11",
          "kubectl.kubernetes.io/last-applied-configuration": "{\"apiVersion\":\"apps/v1\",\"kind\":\"Deployment\",\"metadata\":{\"annotations\":{\"artifact.spinnaker.io/location\":\"smoketest\",\"artifact.spinnaker.io/name\":\"smoketest-gke\",\"artifact.spinnaker.io/type\":\"kubernetes/deployment\",\"moniker.spinnaker.io/application\":\"smoketests\",\"moniker.spinnaker.io/cluster\":\"deployment smoketest-gke\"},\"labels\":{\"app\":\"smoketest-gke\",\"app.kubernetes.io/managed-by\":\"spinnaker\",\"app.kubernetes.io/name\":\"smoketests\"},\"name\":\"smoketest-gke\",\"namespace\":\"smoketest\"},\"spec\":{\"replicas\":2,\"selector\":{\"matchLabels\":{\"app\":\"smoketest-gke\"}},\"template\":{\"metadata\":{\"annotations\":{\"artifact.spinnaker.io/location\":\"smoketest\",\"artifact.spinnaker.io/name\":\"smoketest-gke\",\"artifact.spinnaker.io/type\":\"kubernetes/deployment\",\"moniker.spinnaker.io/application\":\"smoketests\",\"moniker.spinnaker.io/cluster\":\"deployment smoketest-gke\"},\"labels\":{\"app\":\"smoketest-gke\",\"app.kubernetes.io/managed-by\":\"spinnaker\",\"app.kubernetes.io/name\":\"smoketests\"}},\"spec\":{\"containers\":[{\"image\":\"gcr.io/github-replication-sandbox/rf:1.0.9\",\"name\":\"frontend\"}]}}}}\n",
          "moniker.spinnaker.io/application": "smoketests",
          "moniker.spinnaker.io/cluster": "deployment smoketest-gke"
        },
        "creationTimestamp": "2020-02-13T18:50:01Z",
        "generation": 16,
        "labels": {
          "app": "smoketest-gke",
          "app.kubernetes.io/managed-by": "spinnaker",
          "app.kubernetes.io/name": "smoketests"
        },
        "name": "smoketest-gke",
        "namespace": "smoketest",
        "resourceVersion": "60136848",
        "selfLink": "/apis/extensions/v1beta1/namespaces/smoketest/deployments/smoketest-gke",
        "uid": "a3d93173-4e91-11ea-8db2-4201ac100109"
      },
      "spec": {
        "progressDeadlineSeconds": 600,
        "replicas": 2,
        "revisionHistoryLimit": 10,
        "selector": {
          "matchLabels": {
            "app": "smoketest-gke"
          }
        },
        "strategy": {
          "rollingUpdate": {
            "maxSurge": "25%",
            "maxUnavailable": "25%"
          },
          "type": "RollingUpdate"
        },
        "template": {
          "metadata": {
            "annotations": {
              "artifact.spinnaker.io/location": "smoketest",
              "artifact.spinnaker.io/name": "smoketest-gke",
              "artifact.spinnaker.io/type": "kubernetes/deployment",
              "moniker.spinnaker.io/application": "smoketests",
              "moniker.spinnaker.io/cluster": "deployment smoketest-gke"
            },
            "labels": {
              "app": "smoketest-gke",
              "app.kubernetes.io/managed-by": "spinnaker",
              "app.kubernetes.io/name": "smoketests"
            }
          },
          "spec": {
            "containers": [
              {
                "image": "gcr.io/github-replication-sandbox/rf:1.0.9",
                "imagePullPolicy": "IfNotPresent",
                "name": "frontend",
                "resources": {},
                "terminationMessagePath": "/dev/termination-log",
                "terminationMessagePolicy": "File"
              }
            ],
            "dnsPolicy": "ClusterFirst",
            "restartPolicy": "Always",
            "schedulerName": "default-scheduler",
            "securityContext": {},
            "terminationGracePeriodSeconds": 30
          }
        }
      },
      "status": {
        "availableReplicas": 2,
        "conditions": [
          {
            "lastTransitionTime": "2020-06-02T03:09:40Z",
            "lastUpdateTime": "2020-06-02T03:09:40Z",
            "message": "Deployment has minimum availability.",
            "reason": "MinimumReplicasAvailable",
            "status": "True",
            "type": "Available"
          },
          {
            "lastTransitionTime": "2020-02-13T18:50:01Z",
            "lastUpdateTime": "2020-06-16T13:47:10Z",
            "message": "ReplicaSet \"smoketest-gke-598b749779\" has successfully progressed.",
            "reason": "NewReplicaSetAvailable",
            "status": "True",
            "type": "Progressing"
          }
        ],
        "observedGeneration": 16,
        "readyReplicas": 2,
        "replicas": 2,
        "updatedReplicas": 2
      }
    },
    "moniker": {
      "app": "smoketests",
      "cluster": "deployment smoketest-gke"
    },
    "name": "deployment smoketest-gke",
    "providerType": "kubernetes",
    "region": "smoketest",
    "serverGroups": [
      {
        "account": "gke_github-replication-sandbox_us-central1_sandbox-us-central1-agent-dev",
        "moniker": {
          "app": "smoketests",
          "cluster": "deployment smoketest-gke",
          "sequence": 11
        },
        "name": "replicaSet smoketest-gke-598b749779",
        "namespace": "smoketest",
        "region": "smoketest"
      }
    ],
    "type": "kubernetes",
    "uid": "a3d93173-4e91-11ea-8db2-4201ac100109",
    "zone": "smoketest"
  },
  {
    "account": "gke_github-replication-sandbox_us-central1_sandbox-us-central1-agent_smoketest-dev",
    "accountName": "gke_github-replication-sandbox_us-central1_sandbox-us-central1-agent_smoketest-dev",
    "cloudProvider": "kubernetes",
    "createdTime": 1581619801000,
    "key": {
      "account": "gke_github-replication-sandbox_us-central1_sandbox-us-central1-agent_smoketest-dev",
      "group": "deployment",
      "kubernetesKind": "deployment",
      "name": "smoketest-gke",
      "namespace": "smoketest",
      "provider": "kubernetes"
    },
    "kind": "deployment",
    "labels": {
      "app": "smoketest-gke",
      "app.kubernetes.io/managed-by": "spinnaker",
      "app.kubernetes.io/name": "smoketests"
    },
    "manifest": {
      "apiVersion": "extensions/v1beta1",
      "kind": "Deployment",
      "metadata": {
        "annotations": {
          "artifact.spinnaker.io/location": "smoketest",
          "artifact.spinnaker.io/name": "smoketest-gke",
          "artifact.spinnaker.io/type": "kubernetes/deployment",
          "deployment.kubernetes.io/revision": "11",
          "kubectl.kubernetes.io/last-applied-configuration": "{\"apiVersion\":\"apps/v1\",\"kind\":\"Deployment\",\"metadata\":{\"annotations\":{\"artifact.spinnaker.io/location\":\"smoketest\",\"artifact.spinnaker.io/name\":\"smoketest-gke\",\"artifact.spinnaker.io/type\":\"kubernetes/deployment\",\"moniker.spinnaker.io/application\":\"smoketests\",\"moniker.spinnaker.io/cluster\":\"deployment smoketest-gke\"},\"labels\":{\"app\":\"smoketest-gke\",\"app.kubernetes.io/managed-by\":\"spinnaker\",\"app.kubernetes.io/name\":\"smoketests\"},\"name\":\"smoketest-gke\",\"namespace\":\"smoketest\"},\"spec\":{\"replicas\":2,\"selector\":{\"matchLabels\":{\"app\":\"smoketest-gke\"}},\"template\":{\"metadata\":{\"annotations\":{\"artifact.spinnaker.io/location\":\"smoketest\",\"artifact.spinnaker.io/name\":\"smoketest-gke\",\"artifact.spinnaker.io/type\":\"kubernetes/deployment\",\"moniker.spinnaker.io/application\":\"smoketests\",\"moniker.spinnaker.io/cluster\":\"deployment smoketest-gke\"},\"labels\":{\"app\":\"smoketest-gke\",\"app.kubernetes.io/managed-by\":\"spinnaker\",\"app.kubernetes.io/name\":\"smoketests\"}},\"spec\":{\"containers\":[{\"image\":\"gcr.io/github-replication-sandbox/rf:1.0.9\",\"name\":\"frontend\"}]}}}}\n",
          "moniker.spinnaker.io/application": "smoketests",
          "moniker.spinnaker.io/cluster": "deployment smoketest-gke"
        },
        "creationTimestamp": "2020-02-13T18:50:01Z",
        "generation": 16,
        "labels": {
          "app": "smoketest-gke",
          "app.kubernetes.io/managed-by": "spinnaker",
          "app.kubernetes.io/name": "smoketests"
        },
        "name": "smoketest-gke",
        "namespace": "smoketest",
        "resourceVersion": "60136848",
        "selfLink": "/apis/extensions/v1beta1/namespaces/smoketest/deployments/smoketest-gke",
        "uid": "a3d93173-4e91-11ea-8db2-4201ac100109"
      },
      "spec": {
        "progressDeadlineSeconds": 600,
        "replicas": 2,
        "revisionHistoryLimit": 10,
        "selector": {
          "matchLabels": {
            "app": "smoketest-gke"
          }
        },
        "strategy": {
          "rollingUpdate": {
            "maxSurge": "25%",
            "maxUnavailable": "25%"
          },
          "type": "RollingUpdate"
        },
        "template": {
          "metadata": {
            "annotations": {
              "artifact.spinnaker.io/location": "smoketest",
              "artifact.spinnaker.io/name": "smoketest-gke",
              "artifact.spinnaker.io/type": "kubernetes/deployment",
              "moniker.spinnaker.io/application": "smoketests",
              "moniker.spinnaker.io/cluster": "deployment smoketest-gke"
            },
            "labels": {
              "app": "smoketest-gke",
              "app.kubernetes.io/managed-by": "spinnaker",
              "app.kubernetes.io/name": "smoketests"
            }
          },
          "spec": {
            "containers": [
              {
                "image": "gcr.io/github-replication-sandbox/rf:1.0.9",
                "imagePullPolicy": "IfNotPresent",
                "name": "frontend",
                "resources": {},
                "terminationMessagePath": "/dev/termination-log",
                "terminationMessagePolicy": "File"
              }
            ],
            "dnsPolicy": "ClusterFirst",
            "restartPolicy": "Always",
            "schedulerName": "default-scheduler",
            "securityContext": {},
            "terminationGracePeriodSeconds": 30
          }
        }
      },
      "status": {
        "availableReplicas": 2,
        "conditions": [
          {
            "lastTransitionTime": "2020-06-02T03:09:40Z",
            "lastUpdateTime": "2020-06-02T03:09:40Z",
            "message": "Deployment has minimum availability.",
            "reason": "MinimumReplicasAvailable",
            "status": "True",
            "type": "Available"
          },
          {
            "lastTransitionTime": "2020-02-13T18:50:01Z",
            "lastUpdateTime": "2020-06-16T13:47:10Z",
            "message": "ReplicaSet \"smoketest-gke-598b749779\" has successfully progressed.",
            "reason": "NewReplicaSetAvailable",
            "status": "True",
            "type": "Progressing"
          }
        ],
        "observedGeneration": 16,
        "readyReplicas": 2,
        "replicas": 2,
        "updatedReplicas": 2
      }
    },
    "moniker": {
      "app": "smoketests",
      "cluster": "deployment smoketest-gke"
    },
    "name": "deployment smoketest-gke",
    "providerType": "kubernetes",
    "region": "smoketest",
    "serverGroups": [
      {
        "account": "gke_github-replication-sandbox_us-central1_sandbox-us-central1-agent_smoketest-dev",
        "moniker": {
          "app": "smoketests",
          "cluster": "deployment smoketest-gke",
          "sequence": 11
        },
        "name": "replicaSet smoketest-gke-598b749779",
        "namespace": "smoketest",
        "region": "smoketest"
      }
    ],
    "type": "kubernetes",
    "uid": "a3d93173-4e91-11ea-8db2-4201ac100109",
    "zone": "smoketest"
  }
]`

const payloadServerGroupManagers = `[
            {
              "account": "account1",
              "accountName": "account1",
              "cloudProvider": "kubernetes",
              "createdTime": 1581603123000,
              "key": {
                "account": "account1",
                "group": "deployment",
                "kubernetesKind": "deployment",
                "name": "test-deployment1",
                "namespace": "test-namespace1",
                "provider": "kubernetes"
              },
              "kind": "deployment",
              "labels": {
                "label1": "test-label1"
              },
              "manifest": {
                "apiVersion": "apps/v1",
                "kind": "Deployment",
                "metadata": {
                  "creationTimestamp": "2020-02-13T14:12:03Z",
                  "labels": {
                    "label1": "test-label1"
                  },
                  "name": "test-deployment1",
                  "namespace": "test-namespace1",
                  "uid": "cec15437-4e6a-11ea-9788-4201ac100006"
                }
              },
              "moniker": {
                "app": "test-application",
                "cluster": "deployment test-deployment1"
              },
              "name": "deployment test-deployment1",
              "namespace": "test-namespace1",
              "providerType": "kubernetes",
              "region": "test-namespace1",
              "serverGroups": [
                {
                  "account": "account1",
                  "moniker": {
                    "app": "test-application",
                    "cluster": "deployment test-deployment1",
                    "sequence": 236
                  },
                  "name": "replicaSet test-rs1",
                  "namespace": "test-namespace1",
                  "region": "test-namespace1"
                }
              ],
              "type": "kubernetes",
              "uid": "cec15437-4e6a-11ea-9788-4201ac100006",
              "zone": "test-application"
            },
            {
              "account": "account1",
              "accountName": "account1",
              "cloudProvider": "kubernetes",
              "createdTime": 1581603123000,
              "key": {
                "account": "account1",
                "group": "deployment",
                "kubernetesKind": "deployment",
                "name": "test-deployment2",
                "namespace": "test-namespace2",
                "provider": "kubernetes"
              },
              "kind": "deployment",
              "labels": {
                "label1": "test-label2"
              },
              "manifest": {
                "apiVersion": "apps/v1",
                "kind": "Deployment",
                "metadata": {
                  "creationTimestamp": "2020-02-13T14:12:03Z",
                  "labels": {
                    "label1": "test-label2"
                  },
                  "name": "test-deployment2",
                  "namespace": "test-namespace2",
                  "uid": "cec15437-4e6a-11ea-9788-4201ac100006"
                }
              },
              "moniker": {
                "app": "test-application",
                "cluster": "deployment test-deployment2"
              },
              "name": "deployment test-deployment2",
              "namespace": "test-namespace2",
              "providerType": "kubernetes",
              "region": "test-namespace2",
              "serverGroups": [],
              "type": "kubernetes",
              "uid": "cec15437-4e6a-11ea-9788-4201ac100006",
              "zone": "test-application"
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
              "kind": "",
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
              "kind": "",
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
              "kind": "",
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

const serverGroupsResponse = `[
  {
    "account": "gke_github-replication-sandbox_us-central1_sandbox-us-central1-agent_smoketest-dev",
    "buildInfo": {
      "images": [
        "gcr.io/github-replication-sandbox/rf:1.0.9"
      ]
    },
    "capacity": {
      "desired": 2,
      "pinned": false
    },
    "cloudProvider": "kubernetes",
    "cluster": "deployment smoketest-gke",
    "createdTime": 1591791846000,
    "instanceCounts": {
      "down": 0,
      "outOfService": 0,
      "starting": 0,
      "total": 2,
      "unknown": 0,
      "up": 2
    },
    "instances": [
      {
        "availabilityZone": "smoketest",
        "health": [
          {
            "state": "Up",
            "type": "kubernetes/pod"
          },
          {
            "state": "Up",
            "type": "kuberentes/container"
          }
        ],
        "healthState": "Up",
        "id": "ddc7fbe3-afd7-11ea-bb79-4201ac10010a",
        "name": "pod smoketest-gke-598b749779-d8lk6"
      },
      {
        "availabilityZone": "smoketest",
        "health": [
          {
            "state": "Up",
            "type": "kubernetes/pod"
          },
          {
            "state": "Up",
            "type": "kuberentes/container"
          }
        ],
        "healthState": "Up",
        "id": "df4f6d4e-afd7-11ea-bb79-4201ac10010a",
        "name": "pod smoketest-gke-598b749779-wvrzf"
      }
    ],
    "isDisabled": false,
    "labels": {
      "app": "smoketest-gke",
      "app.kubernetes.io/managed-by": "spinnaker",
      "app.kubernetes.io/name": "smoketests",
      "pod-template-hash": "598b749779"
    },
    "loadBalancers": [],
    "moniker": {
      "app": "smoketests",
      "cluster": "deployment smoketest-gke",
      "sequence": 11
    },
    "name": "replicaSet smoketest-gke-598b749779",
    "region": "smoketest",
    "securityGroups": [],
    "serverGroupManagers": [
      {
        "account": "gke_github-replication-sandbox_us-central1_sandbox-us-central1-agent_smoketest-dev",
        "location": "smoketest",
        "name": "smoketest-gke"
      }
    ],
    "type": "kubernetes"
  },
  {
    "account": "gke_github-replication-sandbox_us-central1_sandbox-us-central1-agent_smoketest-dev",
    "buildInfo": {
      "images": [
        "gcr.io/google-samples/hello-app:1.0"
      ]
    },
    "capacity": {
      "desired": 1,
      "pinned": false
    },
    "cloudProvider": "kubernetes",
    "cluster": "replicaSet helloworld",
    "createdTime": 1597697267000,
    "instanceCounts": {
      "down": 0,
      "outOfService": 0,
      "starting": 0,
      "total": 1,
      "unknown": 0,
      "up": 1
    },
    "instances": [
      {
        "availabilityZone": "smoketest",
        "health": [
          {
            "state": "Up",
            "type": "kubernetes/pod"
          },
          {
            "state": "Up",
            "type": "kuberentes/container"
          }
        ],
        "healthState": "Up",
        "id": "e89c8446-e0ca-11ea-be31-4201ac100108",
        "name": "pod helloworld-v002-2fzq2"
      }
    ],
    "isDisabled": false,
    "labels": {
      "app.kubernetes.io/managed-by": "spinnaker",
      "app.kubernetes.io/name": "smoketests",
      "moniker.spinnaker.io/sequence": "2",
      "tier": "helloworld"
    },
    "loadBalancers": [
      "service helloworld"
    ],
    "moniker": {
      "app": "smoketests",
      "cluster": "replicaSet helloworld",
      "sequence": 2
    },
    "name": "replicaSet helloworld-v002",
    "region": "smoketest",
    "securityGroups": [],
    "serverGroupManagers": [],
    "type": "kubernetes"
  },
  {
    "account": "gke_github-replication-sandbox_us-central1_sandbox-us-central1-agent_smoketest-dev",
    "buildInfo": {
      "images": [
        "gcr.io/google-samples/hello-app:1.0"
      ]
    },
    "capacity": {
      "desired": 1,
      "pinned": false
    },
    "cloudProvider": "kubernetes",
    "cluster": "replicaSet helloworld",
    "createdTime": 1597697316000,
    "instanceCounts": {
      "down": 0,
      "outOfService": 0,
      "starting": 0,
      "total": 1,
      "unknown": 0,
      "up": 1
    },
    "instances": [
      {
        "availabilityZone": "smoketest",
        "health": [
          {
            "state": "Up",
            "type": "kubernetes/pod"
          },
          {
            "state": "Up",
            "type": "kuberentes/container"
          }
        ],
        "healthState": "Up",
        "id": "059054f7-e0cb-11ea-be31-4201ac100108",
        "name": "pod helloworld-v003-2lchg"
      }
    ],
    "isDisabled": false,
    "labels": {
      "app.kubernetes.io/managed-by": "spinnaker",
      "app.kubernetes.io/name": "smoketests",
      "moniker.spinnaker.io/sequence": "3",
      "tier": "helloworld"
    },
    "loadBalancers": [
      "service helloworld"
    ],
    "moniker": {
      "app": "smoketests",
      "cluster": "replicaSet helloworld",
      "sequence": 3
    },
    "name": "replicaSet helloworld-v003",
    "region": "smoketest",
    "securityGroups": [],
    "serverGroupManagers": [],
    "type": "kubernetes"
  },
  {
    "account": "gke_github-replication-sandbox_us-central1_sandbox-us-central1-agent_smoketest-dev",
    "buildInfo": {
      "images": [
        "gcr.io/google-samples/hello-app:1.0"
      ]
    },
    "capacity": {
      "desired": 1,
      "pinned": false
    },
    "cloudProvider": "kubernetes",
    "cluster": "replicaSet helloworld",
    "createdTime": 1597697063000,
    "instanceCounts": {
      "down": 0,
      "outOfService": 0,
      "starting": 0,
      "total": 1,
      "unknown": 0,
      "up": 1
    },
    "instances": [
      {
        "availabilityZone": "smoketest",
        "health": [
          {
            "state": "Up",
            "type": "kubernetes/pod"
          },
          {
            "state": "Up",
            "type": "kuberentes/container"
          }
        ],
        "healthState": "Up",
        "id": "6e98d7ab-e0ca-11ea-be31-4201ac100108",
        "name": "pod helloworld-v001-89649"
      }
    ],
    "isDisabled": false,
    "labels": {
      "app.kubernetes.io/managed-by": "spinnaker",
      "app.kubernetes.io/name": "smoketests",
      "moniker.spinnaker.io/sequence": "1",
      "tier": "helloworld"
    },
    "loadBalancers": [
      "service helloworld"
    ],
    "moniker": {
      "app": "smoketests",
      "cluster": "replicaSet helloworld",
      "sequence": 1
    },
    "name": "replicaSet helloworld-v001",
    "region": "smoketest",
    "securityGroups": [],
    "serverGroupManagers": [],
    "type": "kubernetes"
  },
  {
    "account": "gke_github-replication-sandbox_us-central1_sandbox-us-central1-agent_smoketest-dev",
    "buildInfo": {
      "images": [
        "gcr.io/google-samples/hello-app:1.0"
      ]
    },
    "capacity": {
      "desired": 1,
      "pinned": false
    },
    "cloudProvider": "kubernetes",
    "cluster": "replicaSet helloworld",
    "createdTime": 1597686005000,
    "instanceCounts": {
      "down": 0,
      "outOfService": 0,
      "starting": 0,
      "total": 1,
      "unknown": 0,
      "up": 1
    },
    "instances": [
      {
        "availabilityZone": "smoketest",
        "health": [
          {
            "state": "Up",
            "type": "kubernetes/pod"
          },
          {
            "state": "Up",
            "type": "kuberentes/container"
          }
        ],
        "healthState": "Up",
        "id": "afe8f4b9-e0b0-11ea-be31-4201ac100108",
        "name": "pod helloworld-v000-4qq7p"
      }
    ],
    "isDisabled": false,
    "labels": {
      "app.kubernetes.io/managed-by": "spinnaker",
      "app.kubernetes.io/name": "smoketests",
      "moniker.spinnaker.io/sequence": "0",
      "tier": "helloworld"
    },
    "loadBalancers": [
      "service helloworld"
    ],
    "moniker": {
      "app": "smoketests",
      "cluster": "replicaSet helloworld",
      "sequence": 0
    },
    "name": "replicaSet helloworld-v000",
    "region": "smoketest",
    "securityGroups": [],
    "serverGroupManagers": [],
    "type": "kubernetes"
  },
  {
    "account": "gke_github-replication-sandbox_us-central1_sandbox-us-central1-agent-dev",
    "buildInfo": {
      "images": [
        "gcr.io/google-samples/hello-app:1.0"
      ]
    },
    "capacity": {
      "desired": 1,
      "pinned": false
    },
    "cloudProvider": "kubernetes",
    "cluster": "replicaSet helloworld",
    "createdTime": 1597697063000,
    "instanceCounts": {
      "down": 0,
      "outOfService": 0,
      "starting": 0,
      "total": 1,
      "unknown": 0,
      "up": 1
    },
    "instances": [
      {
        "availabilityZone": "smoketest",
        "health": [
          {
            "state": "Up",
            "type": "kubernetes/pod"
          },
          {
            "state": "Up",
            "type": "kuberentes/container"
          }
        ],
        "healthState": "Up",
        "id": "6e98d7ab-e0ca-11ea-be31-4201ac100108",
        "name": "pod helloworld-v001-89649"
      }
    ],
    "isDisabled": false,
    "labels": {
      "app.kubernetes.io/managed-by": "spinnaker",
      "app.kubernetes.io/name": "smoketests",
      "moniker.spinnaker.io/sequence": "1",
      "tier": "helloworld"
    },
    "loadBalancers": [
      "service helloworld"
    ],
    "moniker": {
      "app": "smoketests",
      "cluster": "replicaSet helloworld",
      "sequence": 1
    },
    "name": "replicaSet helloworld-v001",
    "region": "smoketest",
    "securityGroups": [],
    "serverGroupManagers": [],
    "type": "kubernetes"
  },
  {
    "account": "gke_github-replication-sandbox_us-central1_sandbox-us-central1-agent-dev",
    "buildInfo": {
      "images": [
        "gcr.io/google-samples/hello-app:1.0"
      ]
    },
    "capacity": {
      "desired": 1,
      "pinned": false
    },
    "cloudProvider": "kubernetes",
    "cluster": "replicaSet helloworld",
    "createdTime": 1597686005000,
    "instanceCounts": {
      "down": 0,
      "outOfService": 0,
      "starting": 0,
      "total": 1,
      "unknown": 0,
      "up": 1
    },
    "instances": [
      {
        "availabilityZone": "smoketest",
        "health": [
          {
            "state": "Up",
            "type": "kubernetes/pod"
          },
          {
            "state": "Up",
            "type": "kuberentes/container"
          }
        ],
        "healthState": "Up",
        "id": "afe8f4b9-e0b0-11ea-be31-4201ac100108",
        "name": "pod helloworld-v000-4qq7p"
      }
    ],
    "isDisabled": false,
    "labels": {
      "app.kubernetes.io/managed-by": "spinnaker",
      "app.kubernetes.io/name": "smoketests",
      "moniker.spinnaker.io/sequence": "0",
      "tier": "helloworld"
    },
    "loadBalancers": [
      "service helloworld"
    ],
    "moniker": {
      "app": "smoketests",
      "cluster": "replicaSet helloworld",
      "sequence": 0
    },
    "name": "replicaSet helloworld-v000",
    "region": "smoketest",
    "securityGroups": [],
    "serverGroupManagers": [],
    "type": "kubernetes"
  },
  {
    "account": "gke_github-replication-sandbox_us-central1_sandbox-us-central1-agent-dev",
    "buildInfo": {
      "images": [
        "gcr.io/google-samples/hello-app:1.0"
      ]
    },
    "capacity": {
      "desired": 1,
      "pinned": false
    },
    "cloudProvider": "kubernetes",
    "cluster": "replicaSet helloworld",
    "createdTime": 1597697316000,
    "instanceCounts": {
      "down": 0,
      "outOfService": 0,
      "starting": 0,
      "total": 1,
      "unknown": 0,
      "up": 1
    },
    "instances": [
      {
        "availabilityZone": "smoketest",
        "health": [
          {
            "state": "Up",
            "type": "kubernetes/pod"
          },
          {
            "state": "Up",
            "type": "kuberentes/container"
          }
        ],
        "healthState": "Up",
        "id": "059054f7-e0cb-11ea-be31-4201ac100108",
        "name": "pod helloworld-v003-2lchg"
      }
    ],
    "isDisabled": false,
    "labels": {
      "app.kubernetes.io/managed-by": "spinnaker",
      "app.kubernetes.io/name": "smoketests",
      "moniker.spinnaker.io/sequence": "3",
      "tier": "helloworld"
    },
    "loadBalancers": [
      "service helloworld"
    ],
    "moniker": {
      "app": "smoketests",
      "cluster": "replicaSet helloworld",
      "sequence": 3
    },
    "name": "replicaSet helloworld-v003",
    "region": "smoketest",
    "securityGroups": [],
    "serverGroupManagers": [],
    "type": "kubernetes"
  },
  {
    "account": "gke_github-replication-sandbox_us-central1_sandbox-us-central1-agent-dev",
    "buildInfo": {
      "images": [
        "gcr.io/google-samples/hello-app:1.0"
      ]
    },
    "capacity": {
      "desired": 1,
      "pinned": false
    },
    "cloudProvider": "kubernetes",
    "cluster": "replicaSet helloworld",
    "createdTime": 1597697267000,
    "instanceCounts": {
      "down": 0,
      "outOfService": 0,
      "starting": 0,
      "total": 1,
      "unknown": 0,
      "up": 1
    },
    "instances": [
      {
        "availabilityZone": "smoketest",
        "health": [
          {
            "state": "Up",
            "type": "kubernetes/pod"
          },
          {
            "state": "Up",
            "type": "kuberentes/container"
          }
        ],
        "healthState": "Up",
        "id": "e89c8446-e0ca-11ea-be31-4201ac100108",
        "name": "pod helloworld-v002-2fzq2"
      }
    ],
    "isDisabled": false,
    "labels": {
      "app.kubernetes.io/managed-by": "spinnaker",
      "app.kubernetes.io/name": "smoketests",
      "moniker.spinnaker.io/sequence": "2",
      "tier": "helloworld"
    },
    "loadBalancers": [
      "service helloworld"
    ],
    "moniker": {
      "app": "smoketests",
      "cluster": "replicaSet helloworld",
      "sequence": 2
    },
    "name": "replicaSet helloworld-v002",
    "region": "smoketest",
    "securityGroups": [],
    "serverGroupManagers": [],
    "type": "kubernetes"
  },
  {
    "account": "gke_github-replication-sandbox_us-central1_sandbox-us-central1-agent-dev",
    "buildInfo": {
      "images": [
        "gcr.io/github-replication-sandbox/rf:1.0.9"
      ]
    },
    "capacity": {
      "desired": 2,
      "pinned": false
    },
    "cloudProvider": "kubernetes",
    "cluster": "deployment smoketest-gke",
    "createdTime": 1591791846000,
    "instanceCounts": {
      "down": 0,
      "outOfService": 0,
      "starting": 0,
      "total": 2,
      "unknown": 0,
      "up": 2
    },
    "instances": [
      {
        "availabilityZone": "smoketest",
        "health": [
          {
            "state": "Up",
            "type": "kubernetes/pod"
          },
          {
            "state": "Up",
            "type": "kuberentes/container"
          }
        ],
        "healthState": "Up",
        "id": "df4f6d4e-afd7-11ea-bb79-4201ac10010a",
        "name": "pod smoketest-gke-598b749779-wvrzf"
      },
      {
        "availabilityZone": "smoketest",
        "health": [
          {
            "state": "Up",
            "type": "kubernetes/pod"
          },
          {
            "state": "Up",
            "type": "kuberentes/container"
          }
        ],
        "healthState": "Up",
        "id": "ddc7fbe3-afd7-11ea-bb79-4201ac10010a",
        "name": "pod smoketest-gke-598b749779-d8lk6"
      }
    ],
    "isDisabled": false,
    "labels": {
      "app": "smoketest-gke",
      "app.kubernetes.io/managed-by": "spinnaker",
      "app.kubernetes.io/name": "smoketests",
      "pod-template-hash": "598b749779"
    },
    "loadBalancers": [],
    "moniker": {
      "app": "smoketests",
      "cluster": "deployment smoketest-gke",
      "sequence": 11
    },
    "name": "replicaSet smoketest-gke-598b749779",
    "region": "smoketest",
    "securityGroups": [],
    "serverGroupManagers": [
      {
        "account": "gke_github-replication-sandbox_us-central1_sandbox-us-central1-agent-dev",
        "location": "smoketest",
        "name": "smoketest-gke"
      }
    ],
    "type": "kubernetes"
  }
]`

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

const payloadListApplications = `[
  {
    "attributes": {},
    "clusterNames": {
      "gae_github-replication-sandbox-dev": [
        "20190717t160707"
      ],
      "gae_pr-developer-tools-pr": [
        "20190717t160707"
      ]
    },
    "name": "20190717t160707"
  },
  {
    "attributes": {
      "name": "application"
    },
    "clusterNames": {
      "gke_github-replication-sandbox-dev": [
        "Alertmanager.monitoring.coreos.com spin-prometheus-operator-alertmanager",
        "deployment spin-kube-state-metrics",
        "deployment spin-prometheus-operator-operator"
      ],
      "spin-cluster-account": [
        "Alertmanager.monitoring.coreos.com spin-prometheus-operator-alertmanager",
        "deployment spin-kube-state-metrics",
        "deployment spin-prometheus-operator-operator"
      ]
    },
    "name": "application"
  },
  {
    "attributes": {},
    "clusterNames": {
      "gae_github-replication-sandbox-dev": [
        "aps-dev"
      ],
      "gae_pr-developer-tools-pr": [
        "aps-dev"
      ]
    },
    "name": "aps"
  },
  {
    "attributes": {
      "name": "armory-kubesvc"
    },
    "clusterNames": {
      "gke_github-replication-sandbox-dev": [
        "deployment kubesvc-buddy",
        "service kubesvc-metrics",
        "deployment spin-kubesvc",
        "service buddy"
      ],
      "spin-cluster-account": [
        "deployment kubesvc-buddy",
        "service kubesvc-metrics",
        "deployment spin-kubesvc",
        "service buddy"
      ]
    },
    "name": "armory-kubesvc"
  },
  {
    "attributes": {
      "name": "billy"
    },
    "clusterNames": {
      "gke_github-replication-sandbox-dev": [
        "service locksmith",
        "deployment thdca-qp-digital-p2p-qp-digital-p2p",
        "deployment nginx-test",
        "deployment nginx-deployment"
      ],
      "gke_github-replication-sandbox_us-central1-c_test-1-dev": [
        "deployment hello",
        "replicaSet frontend-rs",
        "deployment nginx-test",
        "deployment nginx-deployment"
      ],
      "gke_github-replication-sandbox_us-central1_sandbox-us-central1-agent-dev": [
        "service qp-digital-p2p",
        "deployment smoketest-gke"
      ],
      "gke_github-replication-sandbox_us-central1_sandbox-us-central1-agent_default-dev": [
        "service qp-digital-p2p"
      ],
      "spin-cluster-account": [
        "service locksmith",
        "deployment thdca-qp-digital-p2p-qp-digital-p2p",
        "deployment nginx-test",
        "deployment nginx-deployment"
      ]
    },
    "name": "billy"
  },
  {
    "attributes": {
      "name": "captain-hook"
    },
    "clusterNames": {
      "gke_github-replication-sandbox-dev": [
        "deployment captain-hook",
        "service captain-hook"
      ],
      "spin-cluster-account": [
        "deployment captain-hook",
        "service captain-hook"
      ]
    },
    "name": "captain-hook"
  },
  {
    "attributes": {
      "name": "cd-nginx-ingress-controller"
    },
    "clusterNames": {
      "gke_github-replication-sandbox-dev": [
        "deployment nginx-ingress-controller",
        "service nginx-ingress-default-backend",
        "deployment nginx-ingress-default-backend",
        "service nginx-ingress-controller"
      ],
      "spin-cluster-account": [
        "deployment nginx-ingress-controller",
        "service nginx-ingress-default-backend",
        "deployment nginx-ingress-default-backend",
        "service nginx-ingress-controller"
      ]
    },
    "name": "cd-nginx-ingress-controller"
  },
  {
    "attributes": {
      "name": "cd-skipper"
    },
    "clusterNames": {
      "gke_github-replication-sandbox-dev": [
        "service skipper",
        "deployment skipper"
      ],
      "spin-cluster-account": [
        "service skipper",
        "deployment skipper"
      ]
    },
    "name": "cd-skipper"
  },
  {
    "attributes": {
      "name": "cd-vault"
    },
    "clusterNames": {
      "gke_github-replication-sandbox-dev": [
        "service vault",
        "ingress vault",
        "service vault-ui",
        "service vault-internal",
        "deployment vault-agent-injector",
        "service vault-standby",
        "service vault-active",
        "service vault-agent-injector-svc",
        "statefulSet vault"
      ],
      "spin-cluster-account": [
        "service vault-ui",
        "ingress vault",
        "service vault",
        "deployment vault-agent-injector",
        "service vault-internal",
        "service vault-standby",
        "service vault-active",
        "service vault-agent-injector-svc",
        "statefulSet vault"
      ]
    },
    "name": "cd-vault"
  },
  {
    "attributes": {
      "name": "dxr05"
    },
    "clusterNames": {
      "gke_github-replication-sandbox_us-central1-c_test-1-dev": [
        "service hello-app"
      ],
      "gke_github-replication-sandbox_us-central1_sandbox-us-central1-agent-dev": [
        "deployment dxr05-nginx",
        "service hello-app",
        "deployment hello-pt173287190-1"
      ],
      "gke_github-replication-sandbox_us-central1_sandbox-us-central1-agent_default-dev": [
        "deployment dxr05-nginx",
        "deployment hello-pt173287190-1",
        "service hello-app"
      ]
    },
    "name": "dxr05"
  },
  {
    "attributes": {},
    "clusterNames": {
      "gae_github-replication-sandbox-dev": [
        "dxr05-default",
        "dxr05"
      ],
      "gae_np-platforms-cd-thd-dev": [
        "dxr05-default",
        "dxr05"
      ],
      "gae_pr-developer-tools-pr": [
        "dxr05-default",
        "dxr05"
      ]
    },
    "name": "dxr05"
  },
  {
    "attributes": {},
    "clusterNames": {
      "gae_github-replication-sandbox-dev": [
        "gaeexamplepipelines-default"
      ],
      "gae_pr-developer-tools-pr": [
        "gaeexamplepipelines-default"
      ]
    },
    "name": "gaeexamplepipelines"
  },
  {
    "attributes": {},
    "clusterNames": {
      "gae_github-replication-sandbox-dev": [
        "gaesandboxtesting-default"
      ],
      "gae_pr-developer-tools-pr": [
        "gaesandboxtesting-default"
      ]
    },
    "name": "gaesandboxtesting"
  },
  {
    "attributes": {},
    "clusterNames": {
      "gae_np-platforms-cd-thd-dev": [
        "gaetestingpython-default"
      ]
    },
    "name": "gaetestingpython"
  },
  {
    "attributes": {
      "name": "helm-test"
    },
    "clusterNames": {
      "gke_github-replication-sandbox-dev": [
        "service cert-manager-webhook",
        "service cert-manager"
      ],
      "spin-cluster-account": [
        "service cert-manager-webhook",
        "service cert-manager"
      ]
    },
    "name": "helm-test"
  },
  {
    "attributes": {},
    "clusterNames": {
      "gae_github-replication-sandbox-dev": [
        "myapp-billy"
      ],
      "gae_pr-developer-tools-pr": [
        "myapp-billy"
      ]
    },
    "name": "myapp"
  },
  {
    "attributes": {
      "name": "nate"
    },
    "clusterNames": {
      "gke_github-replication-sandbox-dev": [
        "deployment nginx-deployment"
      ],
      "spin-cluster-account": [
        "deployment nginx-deployment"
      ]
    },
    "name": "nate"
  },
  {
    "attributes": {
      "name": "poc-konduktor"
    },
    "clusterNames": {
      "gke_github-replication-sandbox_us-central1_sandbox-us-central1-agent-dev": [
        "deployment konduktor"
      ],
      "gke_github-replication-sandbox_us-central1_sandbox-us-central1-agent_konduktor-dev": [
        "deployment konduktor"
      ]
    },
    "name": "poc-konduktor"
  },
  {
    "attributes": {},
    "clusterNames": {
      "gae_github-replication-sandbox-dev": [
        "pt174045607"
      ],
      "gae_pr-developer-tools-pr": [
        "pt174045607"
      ]
    },
    "name": "pt174045607"
  },
  {
    "attributes": {
      "name": "scsre"
    },
    "clusterNames": {
      "gke_github-replication-sandbox_us-central1_sandbox-us-central1-agent-dev": [
        "deployment hello-app",
        "service hello-app"
      ],
      "gke_github-replication-sandbox_us-central1_sandbox-us-central1-agent_smoketest-dev": [
        "deployment hello-app",
        "service hello-app"
      ]
    },
    "name": "scsre"
  },
  {
    "attributes": {
      "name": "smoketests"
    },
    "clusterNames": {
      "gke_github-replication-sandbox_us-central1_sandbox-us-central1-agent-dev": [
        "deployment smoketest-gke",
        "service helloworld"
      ],
      "gke_github-replication-sandbox_us-central1_sandbox-us-central1-agent_smoketest-dev": [
        "deployment smoketest-gke",
        "service helloworld"
      ]
    },
    "name": "smoketests"
  },
  {
    "attributes": {},
    "clusterNames": {
      "gae_github-replication-sandbox-dev": [
        "smoketests-default",
        "smoketests-go-helloworld"
      ],
      "gae_pr-developer-tools-pr": [
        "smoketests-default",
        "smoketests-go-helloworld"
      ]
    },
    "name": "smoketests"
  },
  {
    "attributes": {
      "name": "spin"
    },
    "clusterNames": {
      "gke_github-replication-sandbox-dev": [
        "igor",
        "dinghy",
        "clouddriver",
        "orca",
        "front50",
        "deck",
        "rosco",
        "fiat",
        "echo",
        "gate"
      ],
      "spin-cluster-account": [
        "dinghy",
        "igor",
        "clouddriver",
        "orca",
        "deck",
        "front50",
        "rosco",
        "fiat",
        "echo",
        "gate"
      ]
    },
    "name": "spin"
  },
  {
    "attributes": {
      "name": "spinnaker-agent-operations"
    },
    "clusterNames": {
      "gke_github-replication-sandbox-dev": [
        "deployment cleanup-operator"
      ],
      "gke_github-replication-sandbox_us-central1_sandbox-us-central1-agent-dev": [
        "deployment cleanup-operator"
      ],
      "gke_github-replication-sandbox_us-central1_sandbox-us-central1-agent_default-dev": [
        "deployment cleanup-operator"
      ],
      "spin-cluster-account": [
        "deployment cleanup-operator"
      ]
    },
    "name": "spinnaker-agent-operations"
  },
  {
    "attributes": {
      "name": "spinnaker-blackbeard"
    },
    "clusterNames": {
      "gke_github-replication-sandbox-dev": [
        "service blackbeard-service",
        "deployment blackbeard"
      ],
      "spin-cluster-account": [
        "service blackbeard-service",
        "deployment blackbeard"
      ]
    },
    "name": "spinnaker-blackbeard"
  },
  {
    "attributes": {
      "name": "spinnaker-examples"
    },
    "clusterNames": {
      "gke_github-replication-sandbox-dev": [
        "service hello-app-red-black",
        "ingress hello-app-red-black"
      ],
      "gke_github-replication-sandbox_us-central1_sandbox-us-central1-agent-dev": [
        "deployment hello-app",
        "service hello-app"
      ],
      "spin-cluster-account": [
        "service hello-app-red-black",
        "ingress hello-app-red-black"
      ]
    },
    "name": "spinnaker-examples"
  },
  {
    "attributes": {
      "name": "spinnaker-monitoring"
    },
    "clusterNames": {
      "gke_github-replication-sandbox-dev": [
        "Prometheus.monitoring.coreos.com spin-prometheus-operator-prometheus",
        "service spin-kube-state-metrics",
        "service spin-grafana",
        "service spin-prometheus-operator-operator",
        "service spin-prometheus-operator-kube-controller-manager",
        "service spin-prometheus-operator-alertmanager",
        "daemonSet spin-prometheus-node-exporter",
        "deployment spin-grafana",
        "ingress spin-grafana",
        "deployment spin-kube-state-metrics",
        "service grafana-image-renderer",
        "deployment grafana-image-renderer",
        "deployment spin-prometheus-operator-operator",
        "service spin-prometheus-operator-kube-etcd",
        "service spin-prometheus-operator-coredns",
        "service spin-prometheus-node-exporter",
        "service spin-prometheus-operator-kube-proxy",
        "service spin-prometheus-operator-kube-scheduler",
        "service spin-prometheus-operator-prometheus"
      ],
      "spin-cluster-account": [
        "Prometheus.monitoring.coreos.com spin-prometheus-operator-prometheus",
        "service spin-kube-state-metrics",
        "service spin-grafana",
        "service spin-prometheus-operator-operator",
        "service spin-prometheus-operator-kube-controller-manager",
        "service spin-prometheus-operator-alertmanager",
        "daemonSet spin-prometheus-node-exporter",
        "deployment spin-grafana",
        "ingress spin-grafana",
        "deployment spin-kube-state-metrics",
        "service grafana-image-renderer",
        "deployment grafana-image-renderer",
        "deployment spin-prometheus-operator-operator",
        "service spin-prometheus-operator-kube-etcd",
        "service spin-prometheus-operator-coredns",
        "service spin-prometheus-node-exporter",
        "service spin-prometheus-operator-kube-proxy",
        "service spin-prometheus-operator-kube-scheduler",
        "service spin-prometheus-operator-prometheus"
      ]
    },
    "name": "spinnaker-monitoring"
  },
  {
    "attributes": {
      "name": "spinnaker-opa"
    },
    "clusterNames": {
      "gke_github-replication-sandbox-dev": [
        "service opa",
        "deployment opa"
      ],
      "spin-cluster-account": [
        "service opa",
        "deployment opa"
      ]
    },
    "name": "spinnaker-opa"
  },
  {
    "attributes": {
      "name": "test"
    },
    "clusterNames": {
      "gke_github-replication-sandbox-dev": [
        "service hello-app-red-black",
        "service hello-app-red-black43"
      ],
      "spin-cluster-account": [
        "service hello-app-red-black",
        "service hello-app-red-black43"
      ]
    },
    "name": "test"
  },
  {
    "attributes": {
      "name": "thd-public"
    },
    "clusterNames": {
      "gke_github-replication-sandbox-dev": [
        "deployment thd-public-hello-app"
      ],
      "spin-cluster-account": [
        "deployment thd-public-hello-app"
      ]
    },
    "name": "thd-public"
  },
  {
    "attributes": {},
    "clusterNames": {
      "gae_github-replication-sandbox-dev": [
        "thdcacomllc-billy"
      ],
      "gae_pr-developer-tools-pr": [
        "thdcacomllc-billy"
      ]
    },
    "name": "thdcacomllc"
  },
  {
    "attributes": {
      "name": "vishnu"
    },
    "clusterNames": {
      "gke_github-replication-sandbox_us-central1_sandbox-us-central1-agent-dev": [
        "deployment smoketest-gke"
      ],
      "gke_github-replication-sandbox_us-central1_sandbox-us-central1-agent_default-dev": [
        "deployment smoketest-gke"
      ]
    },
    "name": "vishnu"
  }
]`

const payloadDeleteManifest = `[
   {
      "deleteManifest":{
         "manifestName":"ReplicaSet frontend234-65cbfd6b8f",
         "cloudProvider":"kubernetes",
         "options":{
            "orphanDependants":false,
            "gracePeriodSeconds":15
         },
         "location":"default",
         "user":"me@me.com",
         "account":"spin-cluster-account"
      }
   }
]`

const payloadUndoRolloutManifest = `[
		{
			 "undoRolloutManifest":{
					"manifestName":"deployment spin-clouddriver",
					"cloudProvider":"kubernetes",
					"location":"spinnaker",
					"user":"me@me.com",
					"account":"spin-cluster-account",
					"revision":"196"
			 }
		}
]`

const payloadRollingRestartManfiest = `[
   {
      "rollingRestartManifest":{
         "cloudProvider":"kubernetes",
         "manifestName":"deployment spin-clouddriver",
         "location":"spinnaker",
         "user":"me@demo.me.com",
         "account":"spin-cluster-account"
      }
   }
]`

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
