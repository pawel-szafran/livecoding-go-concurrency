#!/usr/bin/env bash
go get golang.org/x/tools/cmd/stress
stress -p 32 go test
