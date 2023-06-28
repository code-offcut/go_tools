#!/bin/bash

RUN_NAME="go_tools"

mkdir -p output/bin

go build -gcflags "all=-N -l" -o output/bin/${RUN_NAME}

# Windows
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o output/bin/${RUN_NAME}_win.exe
# Linux
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o output/bin/${RUN_NAME}_linux
# MacOS
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o output/bin/${RUN_NAME}_mac