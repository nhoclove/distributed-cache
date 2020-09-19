# DISTRIBUTED CACHE

## Description
A simple caching system inspired by MemeCached

## Features
1. Commands: GET, SET, DEL

## How to
### Local
- Server
```shell
    go run cmd/server/main.go -p 8081
```
- Client
```shell
    go run cmd/client/main.go -servers localhost:8081
```
### Docker
- Server
```shell
    docker-compose -f docker-compose.yaml up/down
```
