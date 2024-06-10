# Customers Microservice

- [Customers Microservice](#customers-microservice)
  - [Features](#features)
    - [Async Requests](#async-requests)
      - [Verify Customer](#verify-customer)
        - [Request](#request)
        - [Success Event](#success-event)
        - [Failure Event](#failure-event)


## Features
### Async Requests
The Customer Service must be able to react to async command requests from a Kafka broken in the topic: `service.customers.request`. The result of the operations produces events in the topic: `service.customers.events`

#### Verify Customer
The service must be able to handler requests with the event type of `verify_customer` with the following contract:


##### Request
```json
{
  "id": "086c2dad-22f4-4445-bfeb-c6465af4c60a",
  "type": "verify_customer",
  "date": "2024-06-10T22:56:53.888Z",
  "origin": "orchestrator",
  "correlation_id": "469cec27-106d-4767-bfbf-04c94c7f4f27",
  "data": {
    "customer_id": "018f6058-66f6-7110-82ac-8fd0034b1363"
  }
}
```

If the `customer_id` is a zero uuid the service will generate a failed validation event.

##### Success Event
```json
{
  "id": "6ffb33b9-b7d2-45da-ad6f-942967c6613c",
  "type": "customer_verified",
  "origin": "customers",
  "correlation_id": "469cec27-106d-4767-bfbf-04c94c7f4f27",
  "date": "2024-06-10T22:56:53.898Z",
  "data": {
    "customer_id": "018f6058-66f6-7110-82ac-8fd0034b1363"
  }
}
```

##### Failure Event
```json
{
  "id": "20b24b05-3e26-4692-8284-0bdb1136e3a4",
  "type": "customer_verification_failed",
  "origin": "customers",
  "correlation_id": "5c30dabc-e5d8-4e3f-a164-dc46326d6f49",
  "date": "2024-06-10T23:33:47.207Z",
  "data": {}
}
```
