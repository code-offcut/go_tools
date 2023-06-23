#!/bin/bash

RUN_NAME="go_tools"

mkdir -p output/bin

go build -gcflags "all=-N -l" -o output/bin/${RUN_NAME}