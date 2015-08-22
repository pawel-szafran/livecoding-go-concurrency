#!/usr/bin/env bash
go get golang.org/x/tools/cmd/stress
stress -p 32 go test 2>&1 | tee stress_search_tests.out
