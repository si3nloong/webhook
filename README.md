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
- Allow to send a webhook using [cURL](https://curl.se/) command
- Support Redis, NATS, NSQ as [message queue](https://en.wikipedia.org/wiki/Message_queue) engine
- CLI ready
- Dockerize
- Configurable
- Kubernetes ready

## ⚡️ RESTful API

- **POST** `/v1/webhook/send`

| Name    | Data Type           | Description        | Required |
| ------- | ------------------- | ------------------ | :------: |
| url     | `string`            | URI                |    ✅    |
| headers | `map[string]string` | HTTP headers       |    ❌    |
| body    | `string`            | HTTP body          |    ❌    |
| retry   | `uint`              | Maximum of retries |    ❌    |

## 💡 gRPC API

Please refer to [proto](/grpc/api) files.

## ⚠️ Disclaimer

This project still under development, don't use this in production!

## License

Copyright 2019 SianLoong

Licensed under the MIT License.
