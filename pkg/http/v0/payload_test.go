package v0_test

const payloadBadRequest = `{
            "error": "invalid character 'd' looking for beginning of value"
          }`

const payloadRequestKubernetesOps = `[
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
              "accountType": "spin-cluster-account",
              "cacheThreads": 0,
              "challengeDestructiveActions": false,
              "cloudProvider": "kubernetes",
              "dockerRegistries": null,
              "enabled": false,
              "environment": "spin-cluster-account",
              "name": "spin-cluster-account",
              "namespaces": null,
              "permissions": {
                "READ": [
                  "gg_cloud_gcp_spinnaker_admins"
                ],
                "WRITE": [
                  "gg_cloud_gcp_spinnaker_admins"
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

const payloadLoadBalancers = `[
  {
    "account": "gke_github-replication-sandbox_us-central1_sandbox-us-central1-agent_smoketest-dev",
    "accountName": "gke_github-replication-sandbox_us-central1_sandbox-us-central1-agent_smoketest-dev",
    "cloudProvider": "kubernetes",
    "createdTime": 1585836001000,
    "key": {
      "account": "gke_github-replication-sandbox_us-central1_sandbox-us-central1-agent_smoketest-dev",
      "group": "service",
      "kubernetesKind": "service",
      "name": "helloworld",
      "namespace": "smoketest",
      "provider": "kubernetes"
    },
    "kind": "service",
    "labels": {
      "app.kubernetes.io/managed-by": "spinnaker",
      "app.kubernetes.io/name": "smoketests"
    },
    "manifest": {
      "apiVersion": "v1",
      "kind": "Service",
      "metadata": {
        "annotations": {
          "artifact.spinnaker.io/location": "smoketest",
          "artifact.spinnaker.io/name": "helloworld",
          "artifact.spinnaker.io/type": "kubernetes/service",
          "kubectl.kubernetes.io/last-applied-configuration": "{\"apiVersion\":\"v1\",\"kind\":\"Service\",\"metadata\":{\"annotations\":{\"artifact.spinnaker.io/location\":\"smoketest\",\"artifact.spinnaker.io/name\":\"helloworld\",\"artifact.spinnaker.io/type\":\"kubernetes/service\",\"moniker.spinnaker.io/application\":\"smoketests\",\"moniker.spinnaker.io/cluster\":\"service helloworld\"},\"labels\":{\"app.kubernetes.io/managed-by\":\"spinnaker\",\"app.kubernetes.io/name\":\"smoketests\"},\"name\":\"helloworld\",\"namespace\":\"smoketest\"},\"spec\":{\"ports\":[{\"name\":\"http\",\"port\":80,\"protocol\":\"TCP\",\"targetPort\":80}],\"selector\":{\"lb\":\"helloworld\"}}}\n",
          "moniker.spinnaker.io/application": "smoketests",
          "moniker.spinnaker.io/cluster": "service helloworld"
        },
        "creationTimestamp": "2020-04-02T14:00:01Z",
        "labels": {
          "app.kubernetes.io/managed-by": "spinnaker",
          "app.kubernetes.io/name": "smoketests"
        },
        "name": "helloworld",
        "namespace": "smoketest",
        "resourceVersion": "60136599",
        "selfLink": "/api/v1/namespaces/smoketest/services/helloworld",
        "uid": "3f10f968-74ea-11ea-b2c9-4201ac10010a"
      },
      "spec": {
        "clusterIP": "10.191.28.250",
        "ports": [
          {
            "name": "http",
            "port": 80,
            "protocol": "TCP",
            "targetPort": 80
          }
        ],
        "selector": {
          "lb": "helloworld"
        },
        "sessionAffinity": "None",
        "type": "ClusterIP"
      },
      "status": {
        "loadBalancer": {}
      }
    },
    "moniker": {
      "app": "smoketests",
      "cluster": "service helloworld"
    },
    "name": "service helloworld",
    "providerType": "kubernetes",
    "region": "smoketest",
    "serverGroups": [
      {
        "account": "gke_github-replication-sandbox_us-central1_sandbox-us-central1-agent_smoketest-dev",
        "cloudProvider": "kubernetes",
        "detachedInstances": [],
        "instances": [
          {
            "health": {
              "platform": "platform",
              "source": "Container helloworld",
              "state": "Up",
              "type": "kuberentes/container"
            },
            "id": "e89c8446-e0ca-11ea-be31-4201ac100108",
            "name": "pod helloworld-v002-2fzq2",
            "zone": "smoketest"
          }
        ],
        "isDisabled": false,
        "name": "replicaSet helloworld-v002",
        "region": "smoketest"
      },
      {
        "account": "gke_github-replication-sandbox_us-central1_sandbox-us-central1-agent_smoketest-dev",
        "cloudProvider": "kubernetes",
        "detachedInstances": [],
        "instances": [
          {
            "health": {
              "platform": "platform",
              "source": "Container helloworld",
              "state": "Up",
              "type": "kuberentes/container"
            },
            "id": "6e98d7ab-e0ca-11ea-be31-4201ac100108",
            "name": "pod helloworld-v001-89649",
            "zone": "smoketest"
          }
        ],
        "isDisabled": false,
        "name": "replicaSet helloworld-v001",
        "region": "smoketest"
      },
      {
        "account": "gke_github-replication-sandbox_us-central1_sandbox-us-central1-agent_smoketest-dev",
        "cloudProvider": "kubernetes",
        "detachedInstances": [],
        "instances": [
          {
            "health": {
              "platform": "platform",
              "source": "Container helloworld",
              "state": "Up",
              "type": "kuberentes/container"
            },
            "id": "059054f7-e0cb-11ea-be31-4201ac100108",
            "name": "pod helloworld-v003-2lchg",
            "zone": "smoketest"
          }
        ],
        "isDisabled": false,
        "name": "replicaSet helloworld-v003",
        "region": "smoketest"
      },
      {
        "account": "gke_github-replication-sandbox_us-central1_sandbox-us-central1-agent_smoketest-dev",
        "cloudProvider": "kubernetes",
        "detachedInstances": [],
        "instances": [
          {
            "health": {
              "platform": "platform",
              "source": "Container helloworld",
              "state": "Up",
              "type": "kuberentes/container"
            },
            "id": "afe8f4b9-e0b0-11ea-be31-4201ac100108",
            "name": "pod helloworld-v000-4qq7p",
            "zone": "smoketest"
          }
        ],
        "isDisabled": false,
        "name": "replicaSet helloworld-v000",
        "region": "smoketest"
      }
    ],
    "type": "kubernetes",
    "uid": "3f10f968-74ea-11ea-b2c9-4201ac10010a",
    "zone": "smoketest"
  },
  {
    "account": "gke_github-replication-sandbox_us-central1_sandbox-us-central1-agent-dev",
    "accountName": "gke_github-replication-sandbox_us-central1_sandbox-us-central1-agent-dev",
    "cloudProvider": "kubernetes",
    "createdTime": 1585836001000,
    "key": {
      "account": "gke_github-replication-sandbox_us-central1_sandbox-us-central1-agent-dev",
      "group": "service",
      "kubernetesKind": "service",
      "name": "helloworld",
      "namespace": "smoketest",
      "provider": "kubernetes"
    },
    "kind": "service",
    "labels": {
      "app.kubernetes.io/managed-by": "spinnaker",
      "app.kubernetes.io/name": "smoketests"
    },
    "manifest": {
      "apiVersion": "v1",
      "kind": "Service",
      "metadata": {
        "annotations": {
          "artifact.spinnaker.io/location": "smoketest",
          "artifact.spinnaker.io/name": "helloworld",
          "artifact.spinnaker.io/type": "kubernetes/service",
          "kubectl.kubernetes.io/last-applied-configuration": "{\"apiVersion\":\"v1\",\"kind\":\"Service\",\"metadata\":{\"annotations\":{\"artifact.spinnaker.io/location\":\"smoketest\",\"artifact.spinnaker.io/name\":\"helloworld\",\"artifact.spinnaker.io/type\":\"kubernetes/service\",\"moniker.spinnaker.io/application\":\"smoketests\",\"moniker.spinnaker.io/cluster\":\"service helloworld\"},\"labels\":{\"app.kubernetes.io/managed-by\":\"spinnaker\",\"app.kubernetes.io/name\":\"smoketests\"},\"name\":\"helloworld\",\"namespace\":\"smoketest\"},\"spec\":{\"ports\":[{\"name\":\"http\",\"port\":80,\"protocol\":\"TCP\",\"targetPort\":80}],\"selector\":{\"lb\":\"helloworld\"}}}\n",
          "moniker.spinnaker.io/application": "smoketests",
          "moniker.spinnaker.io/cluster": "service helloworld"
        },
        "creationTimestamp": "2020-04-02T14:00:01Z",
        "labels": {
          "app.kubernetes.io/managed-by": "spinnaker",
          "app.kubernetes.io/name": "smoketests"
        },
        "name": "helloworld",
        "namespace": "smoketest",
        "resourceVersion": "60136599",
        "selfLink": "/api/v1/namespaces/smoketest/services/helloworld",
        "uid": "3f10f968-74ea-11ea-b2c9-4201ac10010a"
      },
      "spec": {
        "clusterIP": "10.191.28.250",
        "ports": [
          {
            "name": "http",
            "port": 80,
            "protocol": "TCP",
            "targetPort": 80
          }
        ],
        "selector": {
          "lb": "helloworld"
        },
        "sessionAffinity": "None",
        "type": "ClusterIP"
      },
      "status": {
        "loadBalancer": {}
      }
    },
    "moniker": {
      "app": "smoketests",
      "cluster": "service helloworld"
    },
    "name": "service helloworld",
    "providerType": "kubernetes",
    "region": "smoketest",
    "serverGroups": [
      {
        "account": "gke_github-replication-sandbox_us-central1_sandbox-us-central1-agent-dev",
        "cloudProvider": "kubernetes",
        "detachedInstances": [],
        "instances": [
          {
            "health": {
              "platform": "platform",
              "source": "Container helloworld",
              "state": "Up",
              "type": "kuberentes/container"
            },
            "id": "6e98d7ab-e0ca-11ea-be31-4201ac100108",
            "name": "pod helloworld-v001-89649",
            "zone": "smoketest"
          }
        ],
        "isDisabled": false,
        "name": "replicaSet helloworld-v001",
        "region": "smoketest"
      },
      {
        "account": "gke_github-replication-sandbox_us-central1_sandbox-us-central1-agent-dev",
        "cloudProvider": "kubernetes",
        "detachedInstances": [],
        "instances": [
          {
            "health": {
              "platform": "platform",
              "source": "Container helloworld",
              "state": "Up",
              "type": "kuberentes/container"
            },
            "id": "afe8f4b9-e0b0-11ea-be31-4201ac100108",
            "name": "pod helloworld-v000-4qq7p",
            "zone": "smoketest"
          }
        ],
        "isDisabled": false,
        "name": "replicaSet helloworld-v000",
        "region": "smoketest"
      },
      {
        "account": "gke_github-replication-sandbox_us-central1_sandbox-us-central1-agent-dev",
        "cloudProvider": "kubernetes",
        "detachedInstances": [],
        "instances": [
          {
            "health": {
              "platform": "platform",
              "source": "Container helloworld",
              "state": "Up",
              "type": "kuberentes/container"
            },
            "id": "e89c8446-e0ca-11ea-be31-4201ac100108",
            "name": "pod helloworld-v002-2fzq2",
            "zone": "smoketest"
          }
        ],
        "isDisabled": false,
        "name": "replicaSet helloworld-v002",
        "region": "smoketest"
      },
      {
        "account": "gke_github-replication-sandbox_us-central1_sandbox-us-central1-agent-dev",
        "cloudProvider": "kubernetes",
        "detachedInstances": [],
        "instances": [
          {
            "health": {
              "platform": "platform",
              "source": "Container helloworld",
              "state": "Up",
              "type": "kuberentes/container"
            },
            "id": "059054f7-e0cb-11ea-be31-4201ac100108",
            "name": "pod helloworld-v003-2lchg",
            "zone": "smoketest"
          }
        ],
        "isDisabled": false,
        "name": "replicaSet helloworld-v003",
        "region": "smoketest"
      }
    ],
    "type": "kubernetes",
    "uid": "3f10f968-74ea-11ea-b2c9-4201ac10010a",
    "zone": "smoketest"
  }
]`

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
