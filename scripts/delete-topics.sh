#!/bin/bash

services=(
    "orders"
    "customers"
    "accounting"
)

for service in "${services[@]}"
do
    docker exec broker bash /bin/kafka-topics --bootstrap-server localhost:9092 --delete --topic service."$service".request
    docker exec broker bash /bin/kafka-topics --bootstrap-server localhost:9092 --delete --topic service."$service".events
done
