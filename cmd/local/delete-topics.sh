#/bin/bash

docker exec broker bash /bin/kafka-topics --bootstrap-server localhost:9092 --delete --topic saga.create_order_v1.response
docker exec broker bash /bin/kafka-topics --bootstrap-server localhost:9092 --delete --topic service.order.request
docker exec broker bash /bin/kafka-topics --bootstrap-server localhost:9092 --delete --topic service.customer.request
docker exec broker bash /bin/kafka-topics --bootstrap-server localhost:9092 --delete --topic service.kitchen.request

docker exec broker bash /bin/kafka-topics --bootstrap-server localhost:9092 --list
