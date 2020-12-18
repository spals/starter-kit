package main

import "crud-starter/server"

func main() {
	s := server.NewHTTPServer(8080)
	s.Start()
}
