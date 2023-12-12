# ETL-Sample

A application that connect to different datasource and transform to specific format and then storing to different datastore


## Pre-requisite
- Docker
- Docker compose
- Go 1.21

## How to start
`docker-compose up -d` to start postgresql db in docker  
`go mod tidy` to instart the dependencies  
`go run ./cmd/app` to start the application