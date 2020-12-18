package main

import "starter-kit/server"

func main() {
	s := server.NewHTTPServer(8080)
	s.Start()
}
