# ----- Build stage -----
FROM golang:alpine AS builder
RUN apk add --no-cache git

WORKDIR /go/src/http
COPY . .

RUN go get -d -v ./...
RUN go build -o /go/bin -v ./...

# ----- Final stage -----
FROM alpine:latest
RUN apk --no-cache add ca-certificates

COPY --from=builder /go/bin/http /starter-kit-http

ENTRYPOINT /starter-kit-http
ENV HTTP_SERVER_LOG_LEVEL info
ENV HTTP_SERVER_PORT 8080
EXPOSE 8080
