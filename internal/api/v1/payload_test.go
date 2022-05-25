package v1_test

const payloadBadRequest = `{
            "error": "invalid character 'd' looking for beginning of value"
          }`

const payloadRequestKubernetesProviders = `{
						"name": "test-name",
						"host": "test-host",
						"caData": "dGVzdC1jYS1kYXRhCg==",
						"permissions": {
						  "read": [
							  "gg_test"
							],
							"write": [
							  "gg_test"
							]
						}
          }`

const payloadRequestKubernetesProvidersEmptyNamespace = `{
					"name": "test-name",
					"host": "test-host",
					"caData": "dGVzdC1jYS1kYXRhCg==",
					"namespace": "  ",
					"permissions": {
						"read": [
							"gg_test"
						],
						"write": [
							"gg_test"
						]
					}
				}`

const payloadRequestKubernetesUpdateProvidersExistingNamespaces = `{
					"name": "test-name",
					"host": "test-host",
					"caData": "dGVzdC1jYS1kYXRhCg==",
					"namespace": "deprecated-field",
					"namespaces": ["n1", "n2", "n3"],
					"permissions": {
						"read": [
							"gg_test"
						],
						"write": [
							"gg_test"
						]
					}
				}`

const payloadRequestKubernetesProviderDeprecatedNamespace = `{
					"name": "test-name",
					"host": "test-host",
					"caData": "dGVzdC1jYS1kYXRhCg==",
					"namespace": "deprecated-field",
					"permissions": {
						"read": [
							"gg_test"
						],
						"write": [
							"gg_test"
						]
					}
				}`

const payloadRequestKubernetesProvidersMultipleNamespaces = `{
					"name": "test-name",
					"host": "test-host",
					"caData": "dGVzdC1jYS1kYXRhCg==",
					"namespaces": ["n1", "n2", "n3"],
					"permissions": {
						"read": [
							"gg_test"
						],
						"write": [
							"gg_test"
						]
					}
				}`

const payloadRequestKubernetesProvidersBadCAData = `{
						"name": "test-name",
						"host": "test-host",
						"caData": "test-ca-data",
						"permissions": {
						  "read": [
							  "gg_test"
							],
							"write": [
							  "gg_test"
							]
						}
          }`

const payloadRequestKubernetesProvidersMissingReadGroup = `{
					"name": "test-name",
					"host": "test-host",
					"caData": "dGVzdC1jYS1kYXRhCg==",
					"permissions": {
						"read": [
							"gg_test"
						],
						"write": [
							"gg_test2"
						]
					}
				}`

const payloadConflictRequest = `{
            "error": "provider already exists"
          }`

const payloadErrorCreatingProvider = `{
            "error": "error creating provider"
          }`

const payloadErrorCreatingReadPermission = `{
            "error": "error creating read permission"
          }`

const payloadErrorCreatingWritePermission = `{
            "error": "error creating write permission"
          }`

const payloadErrorDecodingBase64 = `{
            "error": "error decoding base64 CA data: illegal base64 data at input byte 4"
          }`

const payloadErrorGettingToken = `{
					"error": "error getting token: unsupported token provider"
				}`

const payloadErrorMissingReadGroup = `{
					"error": "error in permissions: write group 'gg_test2' must be included as a read group"
				}`

const payloadKubernetesProviderCreated = `{
            "name": "test-name",
            "host": "test-host",
            "caData": "dGVzdC1jYS1kYXRhCg==",
            "permissions": {
              "read": [
                "gg_test"
              ],
              "write": [
                "gg_test"
              ]
            }
          }`

const payloadKubernetesProviderCreatedWithDeprecatedNamespace = `{
            "name": "test-name",
            "host": "test-host",
            "caData": "dGVzdC1jYS1kYXRhCg==",
			"namespaces": ["deprecated-field"],
            "permissions": {
              "read": [
                "gg_test"
              ],
              "write": [
                "gg_test"
              ]
            }
          }`

const payloadListKubernetesProviders = `[
						{
							"name": "test-name1",
							"host": "test-host1",
							"caData": "dGVzdC1jYS1kYXRhCg==",
							"permissions": {
								"read": [
									"gg_test1"
								],
								"write": [
									"gg_test1"
								]
							}
						},
						{
							"name": "test-name2",
							"host": "test-host2",
							"caData": "dGVzdC1jYS1kYXRhCg==",
							"namespaces": ["test-namespace"],
							"permissions": {
								"read": [
									"gg_test2"
								],
								"write": [
									"gg_test2"
								]
							}
						},
						{
							"name": "test-name3",
							"host": "test-host3",
							"caData": "dGVzdC1jYS1kYXRhCg==",
							"namespaces": ["test-namespace"],
							"permissions": {
								"read": [
									"gg_test3"
								],
								"write": [
									"gg_test3"
								]
							}
						}
					]`

const payloadKubernetesProviderNotFound = `{
						"error": "provider not found"
					}`

const payloadKubernetesProviderGetGenericError = `{
            "error": "error getting provider"
          }`

const payloadKubernetesProviderDeleteGenericError = `{
            "error": "error deleting provider"
          }`

const payloadKubernetesResourcesDeleteGenericError = `{
            "error": "error deleting resources"
          }`
