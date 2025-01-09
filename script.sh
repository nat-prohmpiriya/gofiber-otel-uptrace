#!/bin/bash

docker compose down 

# # # countdown 5 seconds
for i in $(seq 5 -1 1); do
    echo -n "$i..."
    sleep 1
done

docker compose up

# docker-compose logs uptrace
# docker-compose logs -f  otelc

# docker-compose exec otelc sh
# telnet uptrace 14317

# docker-compose logs uptrace | grep 14317

# ติดตั้ง grpcurl
# brew install grpcurl

# ทดสอบ connection
# grpcurl -plaintext localhost:14317 list