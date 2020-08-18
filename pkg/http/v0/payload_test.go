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
