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

- Support [RESTful](https://en.wikipedia.org/wiki/Representational_state_transfer) API
- Support [gRPC](https://grpc.io/) protocol
- Allow to send a webhook using [cURL](https://curl.se/) command
- Support Redis, NATS, NSQ as [message queue](https://en.wikipedia.org/wiki/Message_queue) engine
- Dockerize
- Configurable
- Kubernetes ready

## ⚠️ Disclaimer

This project still under development, don't use this in production!

## License

[GPL 3.0](https://www.gnu.org/licenses/gpl-3.0.html)
