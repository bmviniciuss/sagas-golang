# Accounting Microservice

- [Accounting Microservice](#accounting-microservice)
  - [Features](#features)
    - [Async Requests](#async-requests)
      - [Authorize Card Purchase](#authorize-card-purchase)
        - [Request](#request)
        - [Success Event](#success-event)
        - [Failure Event](#failure-event)


## Features
### Async Requests
The Customer Service must be able to react to async command requests from a Kafka broken in the topic: `service.accounting.request`. The result of the operations produces events in the topic: `service.accounting.events`

#### Authorize Card Purchase
The service must be able to handler requests with the event type of `authorize_card` to authorize a card payment with the following contract:


##### Request
```json
{
  "id": "9cc50b18-8ebe-45fb-8e61-7e35148d246f",
  "type": "authorize_card",
  "origin": "orchestrator",
  "date": "2024-06-10T22:56:53.916Z",
  "correlation_id": "469cec27-106d-4767-bfbf-04c94c7f4f27",
  "data": {
    "amount": 140,
    "card": "0000000000000000"
  }
}
```

If the `data.card` is a zero value the service will generate a failed authorization event.

##### Success Event
```json
{
  "id": "7863df51-800b-41c8-8511-28df5491a091",
  "type": "card_authorized",
  "origin": "accounting",
  "correlation_id": "4c71cfd1-b2b0-4632-8075-9954d8fd1b2d",
  "date": "2024-06-10T22:56:59.658Z",
  "data": {}
}
```

##### Failure Event
```json
{
  "id": "531db911-8463-46cd-bc4f-f33fdcb6bf57",
  "type": "card_authorization_failed",
  "origin": "accounting",
  "correlation_id": "469cec27-106d-4767-bfbf-04c94c7f4f27",
  "date": "2024-06-10T22:56:53.924Z",
  "data": {}
}
```
