# ----- Build stage -----
FROM golang:alpine AS builder
RUN apk add --no-cache git

WORKDIR /go/src/grpc
COPY . .

RUN go get -d -v ./...
RUN go build -o /go/bin -v ./...

# ----- Final stage -----
FROM alpine:latest
RUN apk --no-cache add ca-certificates

COPY --from=builder /go/bin/grpc /starter-kit-grpc

ENTRYPOINT /starter-kit-grpc
ENV GRPC_SERVER_PORT 8081
EXPOSE 8081
