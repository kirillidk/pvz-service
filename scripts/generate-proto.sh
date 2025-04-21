#!/bin/sh

if ! command -v protoc > /dev/null 2>&1; then
    echo "protoc is not installed"
    exit 1
fi

mkdir -p pkg/api/proto/pvz/v1

protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    api/proto/pvz/pvz_v1/pvz.proto

echo "Proto files generated successfully!"
