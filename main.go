package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
)

func httpHandler(w http.ResponseWriter, r *http.Request) {
	ipAddress, _, _ := net.SplitHostPort(r.RemoteAddr)
	log.WithFields(
		log.Fields{
			"ip": ipAddress,
		},
	).Info("Serving IP adress")

	allowedIPs := os.Getenv("ALLOWED_IP")
	if allowedIPs == "" {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(ipAddress))
		return
	}

	for _, allowedIP := range strings.Split(allowedIPs, ",") {
		if ipAddress == allowedIP {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(fmt.Sprintf("Access granted. Your source IP (%s) matches an allowed IP.\n", ipAddress)))
			return
		}
	}

	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte(fmt.Sprintf("Access denied. Your source IP (%s) doesn't match the allowed IPs (%s)\n", ipAddress, allowedIPs)))
	return
}

func main() {
	http.HandleFunc("/", httpHandler)
	log.Info("Starting server on port 8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
