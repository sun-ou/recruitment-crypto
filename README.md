# Golang code test
This program runs as C/S architecture to present wallet behaviour.

The client and server run in separate processes, and communicate via HTTP.

The reviewer should check the core logic in wallet/server_logic.go file.
I used two kinds of lock in the program:
1. Each user has a private read/write mutexes lock for account balance
2. A global exclusive lock for transfer which inside the "bank" object

I did not implement data persistence and authentication. If I have more time, I would like to save data in a json file for each 10 seconds.

# How to run

The reviewer should be able to run the tool as:

```shell
go mod tidy
go run . 
```

If you just want to run the client side.

```shell
go run . -client 
```

# Testing

```shell
go test ./...
```
