# DevOps Fizzbuzz

## Pre-requisites

1. You will need the following tools installed:
   1. Runtime
      1. [Go](https://go.dev/doc/install)
   2. [Docker](https://docs.docker.com/get-docker/)
   3. [Kubectl](https://kubernetes.io/docs/tasks/tools/)
   4. [Helm](https://helm.sh/docs/intro/install/)
   5. [Kind](https://kind.sigs.k8s.io/docs/user/quick-start/#installation)

## Description

Cloud native webserver that may be plugged with Redis.

## Building

### Locally

```
make
```

Binary built as `http-server`.

### Docker

```
make build-docker
```

Can be ran via image `local/http-server`, with sample command `docker run  -p 8080:3000 --rm -it local/http-server`, service exposed on port `8080` locally.

### K8S

```
make run-k8s-kind
```

Will provision a `fizzbuzz` kind cluster cluster, install an ingress controller, a redis, and the HTTP server. By default three replicas are created with state handled via Redis.  
The service may be accessed through `localhost:9090`.

## Environment Variables

| Name                | Description                                               | Default               |
|---------------------|-----------------------------------------------------------|-----------------------|
| SERVER_PORT         | Port the web server shall be exposed with.                | 3000                  |
| HEALTHCHECK_TIMEOUT | Frequency internal healthcheck shall be ran (in seconds). | 1                     |
| HEALTHCHECK_ADDRESS | Address the healthcheck shall ping.                       | http://127.0.0.1:3000 |
| HEALTHCHECK_ENABLE  | Whether or not internal healthcheck is enabled.           | true                  |
| REDIS_ENABLE        | Whether or not Redis shall be used to handle the state.   | false                 |
| REDIS_ADDRESS       | Address of the redis server to use, if enabled.           | localhost:6379        |