#!/bin/bash
docker run --name mongo-db -d mongo
docker run --name redis-db -d redis
docker build -t things-api .
docker run --name iothings -p 127.0.0.1:4000:4000 --link mongo-db:mongo --link redis-db:redis -d things-api
