<p align="center">
  <img src="https://user-images.githubusercontent.com/28108597/136213335-8eb3bff5-cda2-4758-a723-2fce660892af.png" width="140px">
</p>

# Webhook Server

> A golang webhook server comply with at least once deliver.

[![Build](https://github.com/si3nloong/webhook/workflows/Testing/badge.svg?branch=master)](https://github.com/si3nloong/webhook/actions/workflows/test.yml)
[![Release](https://img.shields.io/github/v/tag/si3nloong/webhook)](https://github.com/si3nloong/webhook/releases)
[![Go Report](https://goreportcard.com/badge/github.com/si3nloong/webhook)](https://goreportcard.com/report/github.com/si3nloong/webhook)
[![Go Coverage](https://codecov.io/gh/si3nloong/webhook/branch/master/graph/badge.svg)](https://codecov.io/gh/si3nloong/webhook)
[![LICENSE](https://img.shields.io/github/license/si3nloong/webhook)](https://github.com/si3nloong/webhook/blob/master/LICENSE)

## 🔨 Installation

```bash
go get github.com/si3nloong/webhook
```

## ⚙️ Configuration

```yaml
# HTTP server
enabled: true
port: 3000
no_of_worker: 2
max_pending_webhook: 10000

# gRPC server
grpc:
  enabled: true # enable gRPC server
  api_key: "abcd"
  port: 9000

message_queue:
  engine: "redis" # possible value is redis, nats, nsq
  topic: "webhook"
  queue_group: "webhook"
  redis:
    cluster: false
    addr: "127.0.0.1:6379"
    password: ""
    db: 1
  nats:
    js: true # indicate whether use jetstream or not
```

## ✨ Features

- Support [YAML](https://yaml.org/) and [env](https://en.wikipedia.org/wiki/Env) configuration.
- Automatically re-send webhook if the response is fail.
- [RESTful](https://en.wikipedia.org/wiki/Representational_state_transfer) API ready.
- [GraphQL](https://graphql.org/) API ready.
- Support [gRPC](https://grpc.io/) protocol.
- Allow to send a webhook using [cURL](https://curl.se/) command
- Support [Redis](https://redis.io/), [NATS](https://nats.io/), NSQ as [Message Queue](https://en.wikipedia.org/wiki/Message_queue) engine.
- Support [Elasticsearch](https://www.elastic.co/) as Persistent Volume engine.
- Dockerize.
- Configurable.
- Kubernetes ready.
  <!-- - CLI ready. -->
  <!-- - Support tracing, [Jaeger](https://github.com/jaegertracing/jaeger), [OpenCensus](https://opencensus.io/) -->

## ⚡️ RESTful APIs

Please refer to [here](/http/README.md).

## 💡 gRPC API

Please refer to [proto](/grpc/api) files.

## ⚠️ Disclaimer

This project still under development, don't use this in production!

## 🎉 Big Thanks To

Thanks to these awesome companies for their support of Open Source developers ❤

[![GitHub](https://jstools.dev/img/badges/github.svg)](https://github.com/open-source)
[![NPM](https://jstools.dev/img/badges/npm.svg)](https://www.npmjs.com/)

## License

Copyright 2021 SianLoong

Licensed under the [MIT License](/LICENSE).
