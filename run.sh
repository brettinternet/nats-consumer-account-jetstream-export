#!/usr/bin/env bash

go run ./setup.go

docker-compose up -d --force-recreate

docker-compose logs nats

sleep 3

go run ./main.go
