# Order Microservice

## Requirements
### Functional
- Upon receiving a async request, the customer microservice should verify if the customer can make an order and return a response.

The expected input contract is:
```json
{
  "global_id": "5049b639-bd70-4eed-921a-e09422e38eb7",
  "event_id": "1ffaa156-ae0c-4c48-bc46-cb25fc59c506",
  "event_type": "create_order_v1.verify_customer.request",
  "saga": {
    "reply_channel": "saga.create_order_v1.response"
  },
  "event_data": {
    "customer_id": "135f7d8d-11da-4441-9860-b290ed0a8252"
  }
}
```

The expected successful response is:
```json
{
  "global_id": "5049b639-bd70-4eed-921a-e09422e38eb7",
  "event_id": "1ffaa156-ae0c-4c48-bc46-cb25fc59c506",
  "event_type": "create_order_v1.verify_customer.success",
  "saga": {
    "reply_channel": "saga.create_order_v1.response"
  },
  "event_data": {
    "customer_id": "135f7d8d-11da-4441-9860-b290ed0a8252"
  }
}
```

The expected failure response is:
```json
{
  "global_id": "5049b639-bd70-4eed-921a-e09422e38eb7",
  "event_id": "1ffaa156-ae0c-4c48-bc46-cb25fc59c506",
  "event_type": "create_order_v1.verify_customer.failure",
  "saga": {
    "reply_channel": "saga.create_order_v1.response"
  },
  "event_data": {
    "customer_id": "135f7d8d-11da-4441-9860-b290ed0a8252"
  }
}
```
