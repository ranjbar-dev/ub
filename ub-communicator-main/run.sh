#!/usr/bin/env bash
rm -rf main
/usr/local/go/bin/go build -mod=vendor cmd/rabbit-consumer/main.go
./main