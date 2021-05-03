# ----- Build stage -----
FROM golang:alpine AS builder
RUN apk add --no-cache git

WORKDIR /go/src/http
COPY . .

RUN go get -d -v ./...
RUN go build -o /go/bin -v ./...

# ----- Final stage -----
FROM alpine:latest
ARG GIT_HEAD=unknown
RUN apk --no-cache add ca-certificates

COPY --from=builder /go/bin/http /starterkit-http

ENTRYPOINT /starterkit-http
ENV HTTP_SERVER_PORT 8080
EXPOSE 8080

# See https://github.com/opencontainers/image-spec/blob/master/annotations.md
LABEL authors=Spals
LABEL licenses=BSD-3-Clause
LABEL name=starterkit-http
LABEL revision=${GIT_HEAD}
LABEL version=0.0.1
