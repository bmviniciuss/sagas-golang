# Order Microservice

## Requirements
### Functional
- Upon receiving a async request, the order microservice should create an order with status APPROVAL_PENDING and reply asynchronously the order id.

The expected input contract is:
```json
{
  "global_id": "a6de20fa-1576-4d23-a94c-ee2464872c38",
  "event_id": "a6de20fa-1576-4d23-a94c-ee2464872c38",
  "event_type": "create_order_v1.create_order.request",
  "event_data": {
    "customer_id": "a6de20fa-1576-4d23-a94c-ee2464872c38",
    "amount": 10000,
    "currency_code": "BRL"
  },
  "saga": {
    "name": "create_order_v1",
    "reply_to": "service.order.reply",
    "step": {
      "name": "create_order",
    }
  }
}
```

The successful output contract is:
```json
{
  "global_id": "a6de20fa-1576-4d23-a94c-ee2464872c38",
  "event_id": "a6de20fa-1576-4d23-a94c-ee2464872c38",
  "event_type": "create_order_v1.create_order.success",
  "event_data": {
    "order_id": "a6de20fa-1576-4d23-a94c-ee2464872c38"
  },
  "saga": {
    "name": "create_order_v1",
    "reply_to": "service.order.reply",
    "step": {
      "name": "create_order",
    }
  }
}
```

The error output contract is:
```json
{
  "global_id": "a6de20fa-1576-4d23-a94c-ee2464872c38",
  "event_id": "a6de20fa-1576-4d23-a94c-ee2464872c38",
  "event_type": "create_order_v1.create_order.failure",
  "event_data": {},
  "saga": {
    "name": "create_order_v1",
    "reply_to": "service.order.reply",
    "step": {
      "name": "create_order",
    }
  }
}
```

### Non-functional
- The order microservice should listen for messages in the topic `service.order.request`
