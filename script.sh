#!/bin/bash

# down docker
docker-compose -f ./platform/docker/docker-compose.yml down

# countdown 5 seconds
for i in $(seq 5 -1 1); do
    echo -n "$i..."
    sleep 1
done

# up docker
docker-compose -f ./platform/docker/docker-compose.yml up