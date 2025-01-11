#!/bin/bash

# down docker
docker-compose -f ./docker/docker-compose.yml down

# countdown 5 seconds
for i in $(seq 5 -1 1); do
    echo -n "$i..."
    sleep 1
done

# create network
# docker network create monitoring | true
# up docker
docker-compose -f ./docker/docker-compose.yml up


# restart docker nginx 
# docker-compose -f ./docker/docker-compose.yml restart nginx