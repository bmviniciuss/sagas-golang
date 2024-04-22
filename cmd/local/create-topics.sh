#/bin/bash

docker exec broker bash /bin/kafka-topics --bootstrap-server localhost:9092 --list
docker exec broker bash /bin/kafka-topics --bootstrap-server localhost:9092 --create --topic saga.create_order_v1.response --partitions 1 --replication-factor 1
docker exec broker bash /bin/kafka-topics --bootstrap-server localhost:9092 --create --topic service.order.request --partitions 1 --replication-factor 1
docker exec broker bash /bin/kafka-topics --bootstrap-server localhost:9092 --create --topic service.customer.request --partitions 1 --replication-factor 1
docker exec broker bash /bin/kafka-topics --bootstrap-server localhost:9092 --list
