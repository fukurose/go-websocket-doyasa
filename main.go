package main

import (
	"log"
	"net/http"
)

func main() {
	var httpServer http.Server
	http.Handle("/", http.FileServer(http.Dir("public")))
	log.Println("start http listening :3000")
	httpServer.Addr = ":3000"
	log.Println(httpServer.ListenAndServe())
}
