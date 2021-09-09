# Webhook

> A golang webhook server comply with at least once deliver.

## 🔨 Installation

```bash
go get github.com/si3nloong/webhook
```

## ⚙️ Configuration

```yaml
# RESTful server
enabled: true
port: "3000"

# gRPC server
grpc:
  enabled: true
  port: "9000"

message_queue:
  engine: ""
```

## ✨ Features

- Support [YAML](https://yaml.org/) and [env](https://en.wikipedia.org/wiki/Env) configuration
- Support retry send webhook if the response is fail.
- [RESTful](https://en.wikipedia.org/wiki/Representational_state_transfer) API ready
- Support [gRPC](https://grpc.io/) protocol
- Support tracing, [Jaeger](https://github.com/jaegertracing/jaeger), [OpenCensus](https://opencensus.io/)
- Allow to send a webhook using [cURL](https://curl.se/) command
- Support Redis, NATS, NSQ as [message queue](https://en.wikipedia.org/wiki/Message_queue) engine
- CLI ready
- Dockerize
- Configurable
- Kubernetes ready

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

Copyright 2019 SianLoong

Licensed under the MIT License.
