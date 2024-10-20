# Golang code test
This program implemented RESTful API to present wallet behaviour.

The reviewer should check the core logic in wallet/server_logic.go file.

The ./coverage file is a coverage report. I used unit test for most of files, except main.go and wallet/server_logic.go.
I'd rather use integration test for these files.

I did not implement authentication. If I have more time, I would like to introduce jwt for authentication.

# How to run

The reviewer should be able to run the application as:

```shell
docker-compose up
or
docker-compose -f [path_of_this_project]/docker-compose.yml up
```

# Example for the API

See postman_collection.json

# Testing
Please set environment variable 'pq_host' which point to your postgres host, after you launch the application with docker-compose.

```shell
export pq_host=192.168.56.127:5432
go test ./... -race -cover
```
