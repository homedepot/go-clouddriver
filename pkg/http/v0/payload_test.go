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
