FROM golang:1.24.2-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

RUN apk add --no-cache protobuf protobuf-dev git build-base
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

ENV PATH="/go/bin:$PATH"

COPY . .

RUN chmod +x scripts/generate-proto.sh && ./scripts/generate-proto.sh

RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/api

FROM alpine:latest  

RUN apk --no-cache add ca-certificates postgresql-client

WORKDIR /root/

COPY --from=builder /app/main .
COPY --from=builder /go/bin/migrate /usr/local/bin/migrate
COPY ./migrations ./migrations
COPY ./seeds ./seeds
COPY scripts/docker-entrypoint.sh .

RUN chmod +x ./docker-entrypoint.sh

EXPOSE 8080
EXPOSE 3000

ENTRYPOINT ["./docker-entrypoint.sh"]