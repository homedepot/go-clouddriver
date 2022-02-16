<img src="https://github.com/homedepot/go-clouddriver/blob/media/clouddriver.png" width="125" align="left">

# Go Clouddriver

Go Clouddriver is a rewrite of Spinnaker's [Clouddriver](https://github.com/spinnaker/clouddriver) microservice. It aims to fix severe scaling problems and operational concerns when using Clouddriver at production scale.

It changes how clouddriver operates by providing an extended API for account onboarding (no more "dynamic accounts") and removing over-complicated strategies such as [Cache All The Stuff](https://github.com/spinnaker/clouddriver/tree/master/cats) in favor of talking directly to APIs.

Go Clouddriver generates its access tokens using [Arcade](https://github.com/billiford/arcade), which is meant to be used in tandem with Google's [Workload Identity](https://cloud.google.com/kubernetes-engine/docs/how-to/workload-identity) to generate your tokens in a sidecar and make them retrievable through a simple authenticated API.

Visit [the wiki](https://github.com/homedepot/go-clouddriver/wiki) for feature support.

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

1) Clone and run [arcade](https://github.com/billiford/arcade).

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

| Environment Variable | Description | Notes |
|----------|:-------------:|-----------:|
| `ARCADE_API_KEY` | Needed to talk to [Arcade](https://github.com/billiford/arcade). | Required for most operations. |
| `ARTIFACTS_CREDENTIALS_CONFIG_DIR` | Sets the directory for artifacts configuration. | Not recommended - use OSS Clouddriver's artifacts API instead. |
| `KUBERNETES_USE_DISK_CACHE` | Stores Kubernetes API discovery on disk instead of in-memory. | |
| `DB_HOST` | Used to connect to MySQL datastore. | If not set will default to local SQLite DB. |
| `DB_NAME` | Used to connect to MySQL datastore. | If not set will default to local SQLite DB. |
| `DB_PASS` | Used to connect to MySQL datastore. | If not set will default to local SQLite DB. |
| `DB_USER` | Used to connect to MySQL datastore. | If not set will default to local SQLite DB. |
| `VERBOSE_REQUEST_LOGGING` | Logs all incoming request information. | Should only be used in non-production for testing. |
