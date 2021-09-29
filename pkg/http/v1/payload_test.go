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
							"namespace": "test-namespace",
							"permissions": {
								"read": [
									"gg_test2"
								],
								"write": [
									"gg_test2"
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
