#!/bin/sh

set -e

go test -cover ./cmd/integration
go test -cover ./internal/handler
go test -cover ./internal/middleware
go test -cover ./internal/repository
go test -cover ./internal/service/auth
go test -cover ./internal/service/grpc
go test -cover ./internal/service/product
go test -cover ./internal/service/pvz
go test -cover ./internal/service/reception