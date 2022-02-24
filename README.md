<img src="https://github.com/homedepot/go-clouddriver/blob/media/clouddriver.png" width="125" align="left">

# Go Clouddriver

Go Clouddriver is a rewrite of Spinnaker's [Clouddriver](https://github.com/spinnaker/clouddriver) microservice. It has an observed 95%+ decrease in CPU and memory load for Kubernetes operations.

Go Clouddriver brings many features to the table which allow it to perform better than Clouddriver OSS for production loads:
- it does not rely on `kubectl` and instead interfaces directly with the Kubernetes API for all operations
- it utilizes an in-memory cache store for Kubernetes API discovery
- it stores Kubernetes providers in a database, fronted by a simple CRUD API
- it removes over-complicated strategies such as [Cache All The Stuff](https://github.com/spinnaker/clouddriver/tree/master/cats), instead making live calls for all operations

Go Clouddriver is *not* an API complete implementation of Clouddriver OSS and only handles Kubernetes providers and operations. It is meant to be run in tandem with Clouddriver OSS. Visit [the wiki](https://github.com/homedepot/go-clouddriver/wiki) for feature support and intallation instructions.

## Getting Started

### Testing

Run from the root directory
```bash
make tools test
```

### Building

Run from the root directory
```bash
make build
```

### Running Locally

1) Go Clouddriver generates its access tokens using [Arcade](https://github.com/billiford/arcade) as a sidecar, so a working instance of Arcade will need to be running locally in order for Go Clouddriver to talk to Kubernetes clusters.

2) Export the Arcade API key (the same one you set up in step 1).
```bash
export ARCADE_API_KEY=test
```

3) Run Go Clouddriver.
```bash
make run
```

4) Create your first Kubernetes provider! Go Clouddriver runs on port 7002, so you'll make a POST to `localhost:7002/v1/kubernetes/providers`.
```bash
curl -XPOST localhost:7002/v1/kubernetes/providers -d '{
  "name": "test-provider",
  "host": "https://test-host",
  "caData": "test",
  "permissions": {
    "read": [
      "test-group"
    ],
    "write": [
      "test-group"
    ]
  }
}' | jq
```
And you should see the response
```json
{
  "name": "test-provider",
  "host": "https://test-host",
  "caData": "test",
  "permissions": {
    "read": [
      "test-group"
    ],
    "write": [
      "test-group"
    ]
  }
}
```
Running the command again will return a `409 Conflict` unless you change the name of the provider.

5) List your providers by calling the `/credentials` endpoint.
```bash
curl localhost:7002/credentials | jq
```

### Configuration

| Environment Variable | Description | Notes | Default Value |
|----------|:-------------:|-----------:|------------:|
| `ARCADE_API_KEY` | Needed to talk to [Arcade](https://github.com/billiford/arcade). | Required for most operations. ||
| `ARTIFACTS_CREDENTIALS_CONFIG_DIR` | Sets the directory for artifacts configuration. | Optional. Leave unset to use OSS Clouddriver's Artifacts API. ||
| `KUBERNETES_USE_DISK_CACHE` | Stores Kubernetes API discovery on disk instead of in-memory. || `false` |
| `DB_HOST` | Used to connect to MySQL database. | If not set will default to local SQLite database. ||
| `DB_NAME` | Used to connect to MySQL database. | If not set will default to local SQLite database. ||
| `DB_PASS` | Used to connect to MySQL database. | If not set will default to local SQLite database. ||
| `DB_USER` | Used to connect to MySQL database. | If not set will default to local SQLite database. ||
| `VERBOSE_REQUEST_LOGGING` | Logs all incoming request information. | Should only be used in non-production for testing. | `false` |
