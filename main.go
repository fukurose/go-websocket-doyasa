package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
)

func handlerWebSocket(w http.ResponseWriter, r *http.Request) {
	dump, err := httputil.DumpRequest(r, true)
	if err != nil {
		http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
		return
	}
	fmt.Println(string(dump))
}

func main() {
	var httpServer http.Server
	http.Handle("/", http.FileServer(http.Dir("public")))
	http.HandleFunc("/websocket", handlerWebSocket)
	log.Println("start http listening :3000")
	httpServer.Addr = ":3000"
	log.Println(httpServer.ListenAndServe())
}
