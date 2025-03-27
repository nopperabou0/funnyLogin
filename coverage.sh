#! /usr/bin/bash

go test ./... -coverprofile="coverage.out" && go tool cover -html="coverage.out" -o="coverage.html" && xdg-open "coverage.html"