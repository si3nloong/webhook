# RESTful API

- Health check - [/health]()
- Get webhook list - [/v1/webhooks](/app/http/restful/README.md#get-webhooks)
- Find webhook by id - [/v1/webhook/:id](/app/http/restful/README.md#find-webhook)
- Send a webhook - [/v1/webhook/send](/app/http/restful/README.md#send-webhook)

## Get Webhooks

URL: /v1/webhooks

Method: **GET**

## Find Webhook by ID

URL: /v1/webhook/:id

Method: **GET**

## Send Webhook

URL: /v1/webhook/send

Method: **POST**

| Name    | Data Type           | Description        | Required |
| ------- | ------------------- | ------------------ | :------: |
| url     | `string`            | URI                |    ✅    |
| headers | `map[string]string` | HTTP headers       |    ❌    |
| body    | `string`            | HTTP body          |    ❌    |
| retry   | `uint`              | Maximum of retries |    ❌    |
