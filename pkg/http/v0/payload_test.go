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

const payloadGetServerGroup = `{
  "account": "gke_github-replication-sandbox_us-central1_sandbox-us-central1-agent_smoketest-dev",
  "accountName": "gke_github-replication-sandbox_us-central1_sandbox-us-central1-agent_smoketest-dev",
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
  "createdTime": 1591791846000,
  "disabled": false,
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
      "account": "gke_github-replication-sandbox_us-central1_sandbox-us-central1-agent_smoketest-dev",
      "accountName": "gke_github-replication-sandbox_us-central1_sandbox-us-central1-agent_smoketest-dev",
      "cloudProvider": "kubernetes",
      "createdTime": 1592315226000,
      "health": [
        {
          "platform": "platform",
          "source": "Pod",
          "state": "Up",
          "type": "kubernetes/pod"
        },
        {
          "platform": "platform",
          "source": "Container frontend",
          "state": "Up",
          "type": "kuberentes/container"
        }
      ],
      "healthState": "Up",
      "humanReadableName": "pod smoketest-gke-598b749779-d8lk6",
      "key": {
        "account": "gke_github-replication-sandbox_us-central1_sandbox-us-central1-agent_smoketest-dev",
        "group": "pod",
        "kubernetesKind": "pod",
        "name": "smoketest-gke-598b749779-d8lk6",
        "namespace": "smoketest",
        "provider": "kubernetes"
      },
      "kind": "pod",
      "labels": {
        "app": "smoketest-gke",
        "app.kubernetes.io/managed-by": "spinnaker",
        "app.kubernetes.io/name": "smoketests",
        "pod-template-hash": "598b749779"
      },
      "manifest": {
        "apiVersion": "v1",
        "kind": "Pod",
        "metadata": {
          "annotations": {
            "artifact.spinnaker.io/location": "smoketest",
            "artifact.spinnaker.io/name": "smoketest-gke",
            "artifact.spinnaker.io/type": "kubernetes/deployment",
            "moniker.spinnaker.io/application": "smoketests",
            "moniker.spinnaker.io/cluster": "deployment smoketest-gke"
          },
          "creationTimestamp": "2020-06-16T13:47:06Z",
          "generateName": "smoketest-gke-598b749779-",
          "labels": {
            "app": "smoketest-gke",
            "app.kubernetes.io/managed-by": "spinnaker",
            "app.kubernetes.io/name": "smoketests",
            "pod-template-hash": "598b749779"
          },
          "name": "smoketest-gke-598b749779-d8lk6",
          "namespace": "smoketest",
          "ownerReferences": [
            {
              "apiVersion": "apps/v1",
              "blockOwnerDeletion": true,
              "controller": true,
              "kind": "ReplicaSet",
              "name": "smoketest-gke-598b749779",
              "uid": "4741c722-ab15-11ea-bb79-4201ac10010a"
            }
          ],
          "resourceVersion": "60136811",
          "selfLink": "/api/v1/namespaces/smoketest/pods/smoketest-gke-598b749779-d8lk6",
          "uid": "ddc7fbe3-afd7-11ea-bb79-4201ac10010a"
        },
        "spec": {
          "containers": [
            {
              "image": "gcr.io/github-replication-sandbox/rf:1.0.9",
              "imagePullPolicy": "IfNotPresent",
              "name": "frontend",
              "resources": {},
              "terminationMessagePath": "/dev/termination-log",
              "terminationMessagePolicy": "File",
              "volumeMounts": [
                {
                  "mountPath": "/var/run/secrets/kubernetes.io/serviceaccount",
                  "name": "default-token-q67rl",
                  "readOnly": true
                }
              ]
            }
          ],
          "dnsPolicy": "ClusterFirst",
          "enableServiceLinks": true,
          "nodeName": "gke-sandbox-us-centr-sandbox-us-centr-0ad5d721-zpdb",
          "priority": 0,
          "restartPolicy": "Always",
          "schedulerName": "default-scheduler",
          "securityContext": {},
          "serviceAccount": "default",
          "serviceAccountName": "default",
          "terminationGracePeriodSeconds": 30,
          "tolerations": [
            {
              "effect": "NoExecute",
              "key": "node.kubernetes.io/not-ready",
              "operator": "Exists",
              "tolerationSeconds": 300
            },
            {
              "effect": "NoExecute",
              "key": "node.kubernetes.io/unreachable",
              "operator": "Exists",
              "tolerationSeconds": 300
            }
          ],
          "volumes": [
            {
              "name": "default-token-q67rl",
              "secret": {
                "defaultMode": 420,
                "secretName": "default-token-q67rl"
              }
            }
          ]
        },
        "status": {
          "conditions": [
            {
              "lastTransitionTime": "2020-06-16T13:47:06Z",
              "status": "True",
              "type": "Initialized"
            },
            {
              "lastTransitionTime": "2020-06-16T13:47:08Z",
              "status": "True",
              "type": "Ready"
            },
            {
              "lastTransitionTime": "2020-06-16T13:47:08Z",
              "status": "True",
              "type": "ContainersReady"
            },
            {
              "lastTransitionTime": "2020-06-16T13:47:06Z",
              "status": "True",
              "type": "PodScheduled"
            }
          ],
          "containerStatuses": [
            {
              "containerID": "docker://f48a675fe102c49e479141dd6a0d4ba5b70b0c0378250f899970ca80153e2e76",
              "image": "gcr.io/github-replication-sandbox/rf:1.0.9",
              "imageID": "docker-pullable://gcr.io/github-replication-sandbox/rf@sha256:398db08c89829b8547956d359f27d508497e002a47fbc60095b3746373a8f4d2",
              "lastState": {},
              "name": "frontend",
              "ready": true,
              "restartCount": 0,
              "state": {
                "running": {
                  "startedAt": "2020-06-16T13:47:07Z"
                }
              }
            }
          ],
          "hostIP": "10.191.0.26",
          "phase": "Running",
          "podIP": "10.4.4.187",
          "qosClass": "BestEffort",
          "startTime": "2020-06-16T13:47:06Z"
        }
      },
      "moniker": {
        "app": "smoketests",
        "cluster": "deployment smoketest-gke"
      },
      "name": "ddc7fbe3-afd7-11ea-bb79-4201ac10010a",
      "providerType": "kubernetes",
      "region": "smoketest",
      "type": "kubernetes",
      "uid": "ddc7fbe3-afd7-11ea-bb79-4201ac10010a",
      "zone": "smoketest"
    },
    {
      "account": "gke_github-replication-sandbox_us-central1_sandbox-us-central1-agent_smoketest-dev",
      "accountName": "gke_github-replication-sandbox_us-central1_sandbox-us-central1-agent_smoketest-dev",
      "cloudProvider": "kubernetes",
      "createdTime": 1592315228000,
      "health": [
        {
          "platform": "platform",
          "source": "Pod",
          "state": "Up",
          "type": "kubernetes/pod"
        },
        {
          "platform": "platform",
          "source": "Container frontend",
          "state": "Up",
          "type": "kuberentes/container"
        }
      ],
      "healthState": "Up",
      "humanReadableName": "pod smoketest-gke-598b749779-wvrzf",
      "key": {
        "account": "gke_github-replication-sandbox_us-central1_sandbox-us-central1-agent_smoketest-dev",
        "group": "pod",
        "kubernetesKind": "pod",
        "name": "smoketest-gke-598b749779-wvrzf",
        "namespace": "smoketest",
        "provider": "kubernetes"
      },
      "kind": "pod",
      "labels": {
        "app": "smoketest-gke",
        "app.kubernetes.io/managed-by": "spinnaker",
        "app.kubernetes.io/name": "smoketests",
        "pod-template-hash": "598b749779"
      },
      "manifest": {
        "apiVersion": "v1",
        "kind": "Pod",
        "metadata": {
          "annotations": {
            "artifact.spinnaker.io/location": "smoketest",
            "artifact.spinnaker.io/name": "smoketest-gke",
            "artifact.spinnaker.io/type": "kubernetes/deployment",
            "moniker.spinnaker.io/application": "smoketests",
            "moniker.spinnaker.io/cluster": "deployment smoketest-gke"
          },
          "creationTimestamp": "2020-06-16T13:47:08Z",
          "generateName": "smoketest-gke-598b749779-",
          "labels": {
            "app": "smoketest-gke",
            "app.kubernetes.io/managed-by": "spinnaker",
            "app.kubernetes.io/name": "smoketests",
            "pod-template-hash": "598b749779"
          },
          "name": "smoketest-gke-598b749779-wvrzf",
          "namespace": "smoketest",
          "ownerReferences": [
            {
              "apiVersion": "apps/v1",
              "blockOwnerDeletion": true,
              "controller": true,
              "kind": "ReplicaSet",
              "name": "smoketest-gke-598b749779",
              "uid": "4741c722-ab15-11ea-bb79-4201ac10010a"
            }
          ],
          "resourceVersion": "60136838",
          "selfLink": "/api/v1/namespaces/smoketest/pods/smoketest-gke-598b749779-wvrzf",
          "uid": "df4f6d4e-afd7-11ea-bb79-4201ac10010a"
        },
        "spec": {
          "containers": [
            {
              "image": "gcr.io/github-replication-sandbox/rf:1.0.9",
              "imagePullPolicy": "IfNotPresent",
              "name": "frontend",
              "resources": {},
              "terminationMessagePath": "/dev/termination-log",
              "terminationMessagePolicy": "File",
              "volumeMounts": [
                {
                  "mountPath": "/var/run/secrets/kubernetes.io/serviceaccount",
                  "name": "default-token-q67rl",
                  "readOnly": true
                }
              ]
            }
          ],
          "dnsPolicy": "ClusterFirst",
          "enableServiceLinks": true,
          "nodeName": "gke-sandbox-us-centr-sandbox-us-centr-2edec7cb-4z8l",
          "priority": 0,
          "restartPolicy": "Always",
          "schedulerName": "default-scheduler",
          "securityContext": {},
          "serviceAccount": "default",
          "serviceAccountName": "default",
          "terminationGracePeriodSeconds": 30,
          "tolerations": [
            {
              "effect": "NoExecute",
              "key": "node.kubernetes.io/not-ready",
              "operator": "Exists",
              "tolerationSeconds": 300
            },
            {
              "effect": "NoExecute",
              "key": "node.kubernetes.io/unreachable",
              "operator": "Exists",
              "tolerationSeconds": 300
            }
          ],
          "volumes": [
            {
              "name": "default-token-q67rl",
              "secret": {
                "defaultMode": 420,
                "secretName": "default-token-q67rl"
              }
            }
          ]
        },
        "status": {
          "conditions": [
            {
              "lastTransitionTime": "2020-06-16T13:47:08Z",
              "status": "True",
              "type": "Initialized"
            },
            {
              "lastTransitionTime": "2020-06-16T13:47:10Z",
              "status": "True",
              "type": "Ready"
            },
            {
              "lastTransitionTime": "2020-06-16T13:47:10Z",
              "status": "True",
              "type": "ContainersReady"
            },
            {
              "lastTransitionTime": "2020-06-16T13:47:08Z",
              "status": "True",
              "type": "PodScheduled"
            }
          ],
          "containerStatuses": [
            {
              "containerID": "docker://eee7a857b1ef8593b94fe7ee686919d1bd3c750752ae50e29e2b661c051df039",
              "image": "gcr.io/github-replication-sandbox/rf:1.0.9",
              "imageID": "docker-pullable://gcr.io/github-replication-sandbox/rf@sha256:398db08c89829b8547956d359f27d508497e002a47fbc60095b3746373a8f4d2",
              "lastState": {},
              "name": "frontend",
              "ready": true,
              "restartCount": 0,
              "state": {
                "running": {
                  "startedAt": "2020-06-16T13:47:10Z"
                }
              }
            }
          ],
          "hostIP": "10.191.0.27",
          "phase": "Running",
          "podIP": "10.4.0.181",
          "qosClass": "BestEffort",
          "startTime": "2020-06-16T13:47:08Z"
        }
      },
      "moniker": {
        "app": "smoketests",
        "cluster": "deployment smoketest-gke"
      },
      "name": "df4f6d4e-afd7-11ea-bb79-4201ac10010a",
      "providerType": "kubernetes",
      "region": "smoketest",
      "type": "kubernetes",
      "uid": "df4f6d4e-afd7-11ea-bb79-4201ac10010a",
      "zone": "smoketest"
    }
  ],
  "key": {
    "account": "gke_github-replication-sandbox_us-central1_sandbox-us-central1-agent_smoketest-dev",
    "group": "replicaSet",
    "kubernetesKind": "replicaSet",
    "name": "smoketest-gke-598b749779",
    "namespace": "smoketest",
    "provider": "kubernetes"
  },
  "kind": "replicaSet",
  "labels": {
    "app": "smoketest-gke",
    "app.kubernetes.io/managed-by": "spinnaker",
    "app.kubernetes.io/name": "smoketests",
    "pod-template-hash": "598b749779"
  },
  "launchConfig": {},
  "loadBalancers": [],
  "manifest": {
    "apiVersion": "extensions/v1beta1",
    "kind": "ReplicaSet",
    "metadata": {
      "annotations": {
        "artifact.spinnaker.io/location": "smoketest",
        "artifact.spinnaker.io/name": "smoketest-gke",
        "artifact.spinnaker.io/type": "kubernetes/deployment",
        "deployment.kubernetes.io/desired-replicas": "2",
        "deployment.kubernetes.io/max-replicas": "3",
        "deployment.kubernetes.io/revision": "11",
        "deployment.kubernetes.io/revision-history": "3,5,7,9",
        "moniker.spinnaker.io/application": "smoketests",
        "moniker.spinnaker.io/cluster": "deployment smoketest-gke"
      },
      "creationTimestamp": "2020-06-10T12:24:06Z",
      "generation": 14,
      "labels": {
        "app": "smoketest-gke",
        "app.kubernetes.io/managed-by": "spinnaker",
        "app.kubernetes.io/name": "smoketests",
        "pod-template-hash": "598b749779"
      },
      "name": "smoketest-gke-598b749779",
      "namespace": "smoketest",
      "ownerReferences": [
        {
          "apiVersion": "apps/v1",
          "blockOwnerDeletion": true,
          "controller": true,
          "kind": "Deployment",
          "name": "smoketest-gke",
          "uid": "a3d93173-4e91-11ea-8db2-4201ac100109"
        }
      ],
      "resourceVersion": "60136839",
      "selfLink": "/apis/extensions/v1beta1/namespaces/smoketest/replicasets/smoketest-gke-598b749779",
      "uid": "4741c722-ab15-11ea-bb79-4201ac10010a"
    },
    "spec": {
      "replicas": 2,
      "selector": {
        "matchLabels": {
          "app": "smoketest-gke",
          "pod-template-hash": "598b749779"
        }
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
            "app.kubernetes.io/name": "smoketests",
            "pod-template-hash": "598b749779"
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
      "fullyLabeledReplicas": 2,
      "observedGeneration": 14,
      "readyReplicas": 2,
      "replicas": 2
    }
  },
  "moniker": {
    "app": "smoketests",
    "cluster": "deployment smoketest-gke",
    "sequence": 11
  },
  "name": "replicaSet smoketest-gke-598b749779",
  "providerType": "kubernetes",
  "region": "smoketest",
  "securityGroups": [],
  "serverGroupManagers": [],
  "type": "kubernetes",
  "uid": "4741c722-ab15-11ea-bb79-4201ac10010a",
  "zone": "smoketest",
  "zones": [],
  "insightActions": []
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
