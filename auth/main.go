package main

import "auth/server"

func main() {
	noteServer := server.NewServer()
	noteServer.Start(9104)
}
