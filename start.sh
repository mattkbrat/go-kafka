#!/usr/bin/env bash

set -e

docker compose up -d

# cd ./event-logger/consumer/ && air

plumber read kafka --topics events --address="localhost:9092" -f
