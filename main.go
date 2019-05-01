package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
)

func httphandler(w http.ResponseWriter, r *http.Request) {
	ipAddress, _, _ := net.SplitHostPort(r.RemoteAddr)
	fmt.Fprintf(w, "%s", ipAddress)
}

func main() {
	port, err := strconv.Atoi(os.Getenv("WHATISMYIP_PORT"))
	if err != nil {
		log.Fatalf("Please make sure the environment variable WHATISMYIP_PORT is defined and is a valid integer [1024-65535], error: %s", err)
	}

	listener := fmt.Sprintf(":%d", port)

	http.HandleFunc("/", httphandler)
	log.Fatal(http.ListenAndServe(listener, nil))
}
