#!/bin/bash

PROJECT_COMPOSE_FILE=$(dirname "$0")/../docker-compose.yaml

docker compose -f $PROJECT_COMPOSE_FILE run --rm certbot
docker compose -f $PROJECT_COMPOSE_FILE exec nginx nginx -s reload