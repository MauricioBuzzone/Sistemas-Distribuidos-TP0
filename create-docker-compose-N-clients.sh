#!/bin/bash

if [ -z "$1" ]; then
    echo "it is necessary to specify the number of clients"
    exit 1
fi

if ! [[ "$1" =~ ^[0-9]+$ ]] || [ "$1" -lt 1 ]; then
    echo "the number of clients must be a positive integer."
    exit 1
fi

file_name="docker-compose-dev.yaml"

touch "$file_name"

echo "version: '3.9'
name: tp0
services:
  server:
    container_name: server
    image: server:latest
    entrypoint: python3 /main.py
    environment:
      - PYTHONUNBUFFERED=1
      - LOGGING_LEVEL=DEBUG
    volumes:
      - ./server/config.ini:/config.ini
    networks:
      - testing_net" > "$file_name"

for ((i = 1; i <= $1; i++)); do
    echo "  client$i:
    container_name: client$i
    image: client:latest
    entrypoint: /client
    environment:
      - CLI_ID=$i
      - CLI_LOG_LEVEL=DEBUG
    volumes:
      - ./client/config.yaml:/config.yaml
      - ./.data/dataset/agency-$i.csv:/agency-$i.csv
    networks:
      - testing_net
    depends_on:
      - server" >> "$file_name"
done

echo "networks:
  testing_net:
    ipam:
      driver: default
      config:
        - subnet: 172.25.125.0/24" >> "$file_name"