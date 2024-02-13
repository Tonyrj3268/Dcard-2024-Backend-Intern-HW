#!/bin/bash
#請先搜尋go-stress-testing並clone
cd third_party/go-stress-testing
go run main.go -c 1000 -n 100 -u "localhost:8080/api/v1/ad?offset=0&limit=10&age=25&gender=M&country=US&platform=android"
