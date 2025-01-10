#!/bin/bash

# Import variables and functions
source "$(dirname "$0")/00_var&func.sh"


# /opt/dackbok
#   ├── docker-compose.yml -> frontend, backend, traefik, portainer
#   └── config
#       └── .env
#       └── firebase-cedentials.json
#   └── data
#   └── logs


# backedn -> otel sdk -> otel-collector -> uptrace