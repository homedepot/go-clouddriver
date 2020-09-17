<img src="https://github.com/billiford/go-clouddriver/blob/media/clouddriver.png" width="125" align="left">

# go-clouddriver

go-clouddriver is a rewrite of Spinnaker's [Clouddriver](https://github.com/spinnaker/clouddriver) microservice. It aims to fix severe scaling problems and operational concerns when using Clouddriver at production scale. 

It changes how clouddriver operates by providing an extended API for account onboarding (no more "dynamic accounts") and removing over-complicated strategies such as [Cache All The Stuff](https://github.com/spinnaker/clouddriver/tree/master/cats) in favor of talking directly to APIs.

Currently, go-clouddriver generates its access tokens using [arcade](https://github.com/billiford/arcade), which is meant to be used in tandem with Google's [Workload Identity](https://cloud.google.com/kubernetes-engine/docs/how-to/workload-identity) to generate your tokens in a sidecar and make them retrievable through a simple  authenticated API.

## Getting Started

### Testing

Requirements: [install Ginkgo and Gomega](https://onsi.github.io/ginkgo/)

Run from the root directory
```bash
ginkgo -r
```

### Running Locally

1) Build go-clouddriver
```bash
go build cmd/clouddriver/clouddriver.go
```

2) Run go-clouddriver
```bash
./clouddriver
```
You should see a log like `SQL config missing field - defaulting to local sqlite DB.` - this is expected when running locally. For production, you should set the env variables `DB_HOST`, `DB_NAME`, `DB_PASS`, and `DB_USER`.

3) Create your first Kubernetes provider! go-clouddriver runs on port 7002, so you'll make a POST to `localhost:7002/v1/kubernetes/providers`.
```bash
curl -XPOST localhost:7002/v1/kubernetes/providers -d '{
  "name": "test-provider",
  "host": "https://test-host",
  "caData": "test",
  "permissions": {
    "read": [
      "test-read-group"
    ],
    "write": [
      "test-write-group"
    ]
  }
}' | jq
```
And you should see the response...
```json
{
  "name": "test-provider",
  "host": "https://test-host",
  "caData": "test",
  "permissions": {
    "read": [
      "test-read-group"
    ],
    "write": [
      "test-write-group"
    ]
  }
}
```
Running the command again will return a `409 Conflict` unless you change the name of the provider.

4) List your providers by calling the `/credentials` endpoint.
```bash
curl localhost:7002/credentials | jq
```

### Verbose Request Logging

Building go-clouddriver requires a lot of reverse engineering and monitoring incoming requests.

Turn on verbose request logging by setting the environment variable `VERBOSE_REQUEST_LOGGING` to `true`. You'll now see helpful request logs.

```
REQUEST: [2020-09-17T14:26:00Z]
POST /v1/kubernetes/providers HTTP/1.1
Host: localhost:7002
Accept: */*
User-Agent: curl/7.54.0
{
  "name": "test-provider",
  "host": "https://test-host",
  "caData": "test",
  "permissions": {
    "read": [
      "test-read-group"
    ],
    "write": [
      "test-write-group"
    ]
  }
}

[GIN] 2020/09/17 - 10:24:18 | 201 |     5.19472ms |       127.0.0.1 | POST     "/v1/kubernetes/providers"
```
