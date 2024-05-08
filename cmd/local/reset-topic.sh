#!/bin/bash

# Check if at least one argument is provided
if [ $# -eq 0 ]; then
    echo "Usage: $0 <parameter>"
    exit 1
fi

# Print the parameter
echo "Deleting topic: $1"
docker exec broker bash /bin/kafka-topics --bootstrap-server localhost:9092 --delete --topic $1

echo "Creating topic: $1"
docker exec broker bash /bin/kafka-topics --bootstrap-server localhost:9092 --create --topic $1 --partitions 1 --replication-factor 1

