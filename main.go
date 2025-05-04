package main

import "notes-manager/server"

func init() {
	server.InitServer()
}

func main() {
	server.StartServer()
}
