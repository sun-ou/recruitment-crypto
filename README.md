# Golang code test
This program implemented RESTful API to present wallet behaviour.

The reviewer should check the core logic in wallet/server_logic.go file.
I used two kinds of lock in the program:
1. Each user has a private read/write mutexes lock for account balance
2. A global exclusive lock for transfer which inside the "bank" object

I did not implement data persistence and authentication. If I have more time, I would like to save data to a database.

# How to run

The reviewer should be able to run the application as:

```shell
docker-compose up
```

# Example for the API

See postman_collection.json

# Testing

```shell
go test ./... -race -cover
```
