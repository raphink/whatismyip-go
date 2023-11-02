package main

import (
	"fmt"
	"net"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"
)

func httpHandler(w http.ResponseWriter, r *http.Request) {
	ipAddress, _, _ := net.SplitHostPort(r.RemoteAddr)
	log.WithFields(
		log.Fields{
			"ip": ipAddress,
		},
	).Info("Serving IP adress")

	allowedIP := os.Getenv("ALLOWED_IP")
	if allowedIP == "" {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(ipAddress))
		return
	}

	if ipAddress == allowedIP {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("Access granted. Your source IP (%s) matches the allowed IP.\n", ipAddress)))
		return
	}

	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte(fmt.Sprintf("Access denied. Your source IP (%s) doesn't match the allowed IP (%s)\n", ipAddress, allowedIP)))
	return
}

func main() {
	http.HandleFunc("/", httpHandler)
	log.Info("Starting server on port 8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
