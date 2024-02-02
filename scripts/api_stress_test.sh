#!/bin/bash

cd third_party/go-stress-testing
go run main.go -c 10 -n 1000 -u "nginx/api/v1/ad/?offset=0&limit=10&age=25&gender=M&country=US&platform=android"

# api/v1/ad?offset=0&limit=10&age=25&gender=M&country=US&platform=android