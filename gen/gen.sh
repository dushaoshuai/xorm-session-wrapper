#!/usr/bin/env bash

go run ./...
gofmt -w -r 'interface{} -> any' ../session.go