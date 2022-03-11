package main

import (
	"net"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func httpHandler(w http.ResponseWriter, r *http.Request) {
	ipAddress, _, _ := net.SplitHostPort(r.RemoteAddr)
	log.WithFields(
		log.Fields{
			"ip": ipAddress,
		},
	).Info("Serving IP adress")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(ipAddress))
}

func main() {
	http.HandleFunc("/", httpHandler)
	log.Info("Starting server on port 8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
