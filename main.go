package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
)

func handler(w http.ResponseWriter, r *http.Request) {
	ipAddress, _, _ := net.SplitHostPort(r.RemoteAddr)
	fmt.Fprintf(w, "%s", ipAddress)
}

func main() {
	listener := fmt.Sprintf(":%d", os.Setenv("WHATISMYIP_PORT", "8000"))

	http.HandleFunc("/", handler)
	http.ListenAndServe(listener, nil)
}
