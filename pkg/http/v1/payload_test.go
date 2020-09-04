package v1_test

const payloadBadRequest = `{
            "error": "invalid character 'd' looking for beginning of value"
          }`

const payloadRequestKubernetesProviders = `{
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

const payloadKubernetesProviderCreated = `{
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
