package v1_test

const payloadBadRequest = `{
            "error": "invalid character 'd' looking for beginning of value"
          }`

const payloadRequestKubernetesProviders = `{
						"name": "test-name",
						"host": "test-host",
						"caData": "test-ca-data"
          }`

const payloadConflictRequest = `{
            "error": "provider already exists"
          }`

const payloadErrorCreatingProvider = `{
            "error": "error creating provider"
          }`

const payloadKubernetesProviderCreated = `{
            "name": "test-name",
            "host": "test-host",
            "caData": "test-ca-data"
          }`
