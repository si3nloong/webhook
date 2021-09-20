# RESTful API

- [/health]()
- [/v1/webhook/send](/http/README.md#)

## Send Webhook

URL: /v1/webhook/send

Method: **POST**

| Name    | Data Type           | Description        | Required |
| ------- | ------------------- | ------------------ | :------: |
| url     | `string`            | URI                |    ✅    |
| headers | `map[string]string` | HTTP headers       |    ❌    |
| body    | `string`            | HTTP body          |    ❌    |
| retry   | `uint`              | Maximum of retries |    ❌    |
