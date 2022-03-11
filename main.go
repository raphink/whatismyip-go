package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
)

func httpHandler(w http.ResponseWriter, r *http.Request) {
	ipAddress, _, _ := net.SplitHostPort(r.RemoteAddr)
	fmt.Fprintf(w, "%s", ipAddress)
}

func main() {
	http.HandleFunc("/", httpHandler)
	log.Fatal(http.ListenAndServe(":8000", nil))
}
