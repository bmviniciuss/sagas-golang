#!/bin/bash

services=(
    "orders"
    "customers"
    "kitchen"
    "accounting"
)

for service in "${services[@]}"
do
    docker exec broker bash /bin/kafka-topics --bootstrap-server localhost:9092 --create --topic service."$service".request --partitions 1 --replication-factor 1
    docker exec broker bash /bin/kafka-topics --bootstrap-server localhost:9092 --create --topic service."$service".events --partitions 1 --replication-factor 1
done
