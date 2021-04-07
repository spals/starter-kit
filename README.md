[![][License img]][License]

[License]:LICENSE
[License img]:https://img.shields.io/badge/license-BSD3-blue.svg

# Go Starter Kit

Boilerplate golang code to quickly build HTTP and Grpc servers

# Introduction

The Spals Starter Kit is a collection of boilerplate server side code that allows software engineers to quickly and easily develop independently running servers in golang. Both HTTP and Grpc servers are supported (though, some Grpc features are not as mature).

Note that Starter Kit is *not* a framework or embedded server with custom features. It is literally code meant for you to copy and paste wholesale or in part into your own repository in order to get to software development as quickly as possible.

# HTTP
## Quickstart
Instructions for copying the HTTP server code:

1. Clone the Starter Kit repository : `git clone https://github.com/spals/starter-kit.git`
1. Copy the HTTP server code to your repository : `cp -R starter-kit/http /tmp/my-repo`
1. Search and replace the base package name. *NOTE* This is only tested on macOS : `LC_ALL=C find /tmp/my-repo/http -type f -exec sed -i '' -e 's/github.com\/spals\/starter-kit\/http/github.com\/my-repo\/http/g' {} +`

The code comes fully executable. Note that the instuctions below use the base package replacement from the installation instructions above. Whatever base package you choose should be used instead:

1. Navigate to the HTTP server code your repository: `cd /tmp/my-repo/http`
1. Run the code: `go run github.com/my-repo/http`

Sample output:
```
✗ go run github.com/spals/starter-kit/http          
2021/04/06 15:50:49 Initializing HTTPServer
2021/04/06 15:50:49 Parsing HTTPServerConfig
2021/04/06 15:50:49 HTTPServerConfig parsed as 
{
  "AssignRandomPort": false,
  "Port": 8080,
  "ShutdownTimeout": 1000000000,
  "LivenessConfig": {
    "MaxGoRoutines": 100
  },
  "ReadinessConfig": {}
}
2021/04/06 15:50:49 HTTPServer initialized
2021/04/06 15:50:49 Starting HTTPServer on port 8080
2021/04/06 15:50:49 HTTPServer started
```

To shutdown the server from the command line, simply ctrl-c it. Sample output:
```
2021/04/06 16:25:50 HTTPServer stopped
2021/04/06 16:25:50 Shutting down HTTPServer
2021/04/06 16:25:50 HTTPServer shutdown
```

## Configuration
The HTTP Starter Kit server is designed to be configured from the command line via environment variables and thus be easily configured within a container. Out of the box, all environment variables are prefaced with `HTTP_SERVER_`. This prefix is regsitered in the initialization code within `main.go`. The configuration variables are registered via the [go-envconfig](https://github.com/sethvargo/go-envconfig) library in the `server_config.go` file. For example, to set a custom server port number, simple set the `HTTP_SERVER_PORT` environment variable and then re-run the server.

## Tour of Code
The HTTP Starter Kit includes two basic elements: a configuration schema which allows for basic server configuration (e.g. port number) and a health check schema which allows for the registration and execution of server-side health checks (e.g. to allow integration in a container system like Kubernetes). These schemas a hand coded as Go types in `http/server/config`.

Server handlers accept HTTP requests and produce HTTP responses with JSON payloads. In the code, there are two handlers: one to retrieve current server configuration and the other to execute a server health check. Handler code is available in `http/server/handler`.

The rest of the code is pretty self-explanitory. `main.go` is the main execution entrypoint. Note that it is minimal. `server.go` contains the primary server code. Its job is to interpret configuration and run the server accordingly.

### Dependency Injection
Starter Kit relies on [Wire](https://github.com/google/wire) for dependency injection. This keeps the code organized and easily testable. Starter Kit adheres to standard Wire constructs. All schemas and handlers mentioned above use static constructors which are registered in the `wire.go` file. Any new types which require dependency injection should have a static constuctor and added to the initializer in `wire.go`.

### Adding API Endpoints
Adding an API endpoint should be straightforward. Simply create a new handler file (`http/server/handler/config.go` serves as a good example here), register the static constructor in `wire.go`, add the type as an argument in the `NewHTTPServer` constructor within the `server.go` file, and register the handler against a path and HTTP verb in that same constructor.

### Testing
The HTTP server is shipped with some tests, including an end-to-end request/response integration test (see `http/server/server_test.go`). In order to run tests:

1. Navigate to the HTTP server code your repository: `cd /tmp/my-repo/http`
1. Run the tests: `go test ./...`

# Grpc
## Quickstart
Instructions for copying the Grpc server code:

1. Clone the Starter Kit repository : `git clone https://github.com/spals/starter-kit.git`
1. Copy the Grpc server code to your repository : `cp -R starter-kit/grpc /tmp/my-repo`
1. Search and replace the base package name. *NOTE* This is only tested on macOS : `LC_ALL=C find /tmp/my-repo/grpc -type f -exec sed -i '' -e 's/github.com\/spals\/starter-kit\/grpc/github.com\/my-repo\/grpc/g' {} +`

Instructions for generating Go code from protobuf files:

1. Navigate to the Grpc server code your repository : `cd /tmp/my-repo/grpc`
1. Execute the `protoc` build command :
```
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    proto/*.proto
```

The code comes fully executable. Note that the instuctions below use the base package replacement from the installation instructions above. Whatever base package you choose should be used instead:

1. Navigate to the Grpc server code your repository: `cd /tmp/my-repo/grpc`
1. Run the code: `go run github.com/my-repo/grpc`

Sample output:
```
✗ go run github.com/spals/starter-kit/grpc          
2021/04/06 16:43:27 Initializing GrpcServer
2021/04/06 16:43:27 GrpcServer initialized
2021/04/06 16:43:27 Finding available random port
2021/04/06 16:43:27 Overwriting configured port (0) with random port (50369)
2021/04/06 16:43:27 Starting GrpcServer on port 50369
2021/04/06 16:43:27 GrpcServer started
```

To shutdown the server from the command line, simply ctrl-c it. Sample output:
```
2021/04/06 16:43:52 GrpcServer stopped
2021/04/06 16:43:52 Shutting down GrpcServer
2021/04/06 16:43:52 GrpcServer shutdown
```

## Configuration
TODO

## Tour of Code
TODO
