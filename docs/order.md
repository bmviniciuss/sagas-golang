# Orders Microservice

- [Orders Microservice](#orders-microservice)
  - [Features](#features)
    - [Async Requests](#async-requests)
      - [Create Order](#create-order)
        - [Request](#request)
        - [Success Event](#success-event)
        - [Failure Event](#failure-event)
      - [Approve Order](#approve-order)
        - [Request](#request-1)
        - [Success Event](#success-event-1)
      - [Reject Order](#reject-order)
        - [Request](#request-2)
        - [Success Event](#success-event-2)
    - [API](#api)
      - [GET `v1/heath`](#get-v1heath)
        - [Response](#response)
      - [GET `v1/orders`](#get-v1orders)
        - [Response](#response-1)
      - [GET `/v1/orders/{id}`](#get-v1ordersid)
        - [Response](#response-2)


## Features
### Async Requests
The Orders Service must be able to react to async command requests from a Kafka broken in the topic: `service.orders.request`. The result of the operations produces events in the topic: `service.orders.events`

#### Create Order
The service must be able to handler requests with the event type of `create_order` with the following contract:

##### Request
```json
{
  "id": "fe97e553-29b3-416f-928b-dca857b89347",
  "type": "create_order",
  "origin": "orchestrator",
  "date": "2024-06-10T22:56:53.841Z",
  "correlation_id": "469cec27-106d-4767-bfbf-04c94c7f4f27",
  "data": {
    "amount": 140,
    "currency_code": "BRL",
    "customer_id": "018f6058-66f6-7110-82ac-8fd0034b1363",
    "items": [
      {
        "id": "018f6058-66f6-7110-82ac-8fd0034b1363",
        "quantity": 1,
        "unit_price": 140
      }
    ]
  }
}
```
Upon receiving the request, the service will create a order with the status of `APPROVAL_PENDING`.

##### Success Event
```json
{
  "id": "816f8ba4-81b1-4229-bad5-b23ec9c7f9ad",
  "type": "order_created",
  "origin": "orders",
  "correlation_id": "469cec27-106d-4767-bfbf-04c94c7f4f27",
  "date": "2024-06-10T22:56:53.873Z",
  "data": {
    "id": "469cec27-106d-4767-bfbf-04c94c7f4f27"
  }
}
```

Because this is a compensable operation, the system may produce a failure event.
##### Failure Event
```json
{
  "id": "816f8ba4-81b1-4229-bad5-b23ec9c7f9ad",
  "type": "order_creation_failed",
  "origin": "orders",
  "correlation_id": "469cec27-106d-4767-bfbf-04c94c7f4f27",
  "date": "2024-06-10T22:56:53.873Z",
  "data": {}
}
```

#### Approve Order
The service must be able to handler requests with the event type of `approve_order` with the following contract:

The service should find the order and change the status to `APPROVED`.


##### Request
```json
{
  "id": "f2f72931-b06f-4624-90e9-7b009b4bcd9f",
  "type": "approve_order",
  "date": "2024-06-10T22:56:59.696Z",
  "origin": "orchestrator",
  "correlation_id": "4c71cfd1-b2b0-4632-8075-9954d8fd1b2d",
  "data": {}
}
```

Because this is a retryable operation only e a success event is produce.
##### Success Event
```json
{
  "id": "2857a234-2d45-4c24-b96f-d1a53294d30b",
  "type": "order_approved",
  "origin": "orders",
  "correlation_id": "4c71cfd1-b2b0-4632-8075-9954d8fd1b2d",
  "date": "2024-06-10T22:56:59.706Z",
  "data": {}
}
```


#### Reject Order
The service must be able to handler requests with the event type of `reject_order` with the following contract:

The service should find the order and change the status to `REJECTED`.
##### Request
```json
{
  "id": "744648c4-b77b-4467-a227-effdb370aa0b",
  "type": "reject_order",
  "date": "2024-06-10T22:56:53.938Z",
  "origin": "orchestrator",
  "correlation_id": "469cec27-106d-4767-bfbf-04c94c7f4f27",
  "data": {}
}
```

Because this is a retryable operation only e a success event is produce.
##### Success Event
```json
{
  "id": "65f62244-c7db-405a-acfa-6fa94ea9baa7",
  "type": "order_rejected",
  "origin": "orders",
  "correlation_id": "469cec27-106d-4767-bfbf-04c94c7f4f27",
  "date": "2024-06-10T22:56:53.955Z",
  "data": {}
}
```

### API
#### GET `v1/heath`
Get a service health check
##### Response
```json
{
  "status": "ok",
  "time": "2024-06-10T23:30:11Z"
}
```
#### GET `v1/orders`
List all orders from the system
##### Response
```json
{
  "content": [
    {
      "id": "469cec27-106d-4767-bfbf-04c94c7f4f27",
      "customer_id": "018f6058-66f6-7110-82ac-8fd0034b1363",
      "amount": 140,
      "currency_code": "BRL",
      "status": "REJECTED",
      "created_at": "2024-06-10T22:56:53.859Z",
      "updated_at": "2024-06-10T22:56:53.954Z"
    },
    {
      "id": "4c71cfd1-b2b0-4632-8075-9954d8fd1b2d",
      "customer_id": "018f6058-66f6-7110-82ac-8fd0034b1363",
      "amount": 140,
      "currency_code": "BRL",
      "status": "APPROVED",
      "created_at": "2024-06-10T22:56:59.604Z",
      "updated_at": "2024-06-10T22:56:59.705Z"
    }
  ]
}
```
#### GET `/v1/orders/{id}`
Find a order by id
##### Response
```json
{
  "id": "469cec27-106d-4767-bfbf-04c94c7f4f27",
  "customer_id": "018f6058-66f6-7110-82ac-8fd0034b1363",
  "amount": 140,
  "currency_code": "BRL",
  "status": "REJECTED",
  "created_at": "2024-06-10T22:56:53.859Z",
  "updated_at": "2024-06-10T22:56:53.954Z"
}
```
