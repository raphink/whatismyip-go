package whatismyip

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"

	"cloud.google.com/go/firestore"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	log "github.com/sirupsen/logrus"
	"google.golang.org/api/iterator"
)

const (
	firestoreCollection = "whatismyip-allowed"
)

var (
	projectID string
	client    *firestore.Client
	ctx       context.Context
)

func init() {
	functions.HTTP("WhatIsMyIP", whatIsMyIP)

	if projectID = os.Getenv("GCP_PROJECT"); projectID != "" {
		ctx = context.Background()
		var err error
		client, err = firestore.NewClient(ctx, projectID)
		if err != nil {
			log.Fatalf("Failed to create Firestore client: %v", err)
		}
	}
}

func whatIsMyIP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handleGet(w, r)
	case http.MethodPost:
		handlePost(w, r)
	case http.MethodDelete:
		handleDelete(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method not allowed"))
	}
}

func getIPAddress(r *http.Request) string {
	// Try to get the IP address from the X-Forwarded-For header
	xff := r.Header.Get("X-Forwarded-For")
	if xff != "" {
		ips := strings.Split(xff, ",")
		return strings.TrimSpace(ips[0])
	}

	// Fall back to X-Real-IP header
	xri := r.Header.Get("X-Real-IP")
	if xri != "" {
		return strings.TrimSpace(xri)
	}

	// Fall back to RemoteAddr
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}

func handleGet(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/list" {
		if !checkBasicAuth(r) {
			w.WriteHeader(http.StatusUnauthorized)
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			w.Write([]byte("Unauthorized"))
			return
		}

		listAllowedIPs(w, r)
		return
	}

	ipAddress := getIPAddress(r)
	log.WithFields(
		log.Fields{
			"ip": ipAddress,
		},
	).Info("Serving IP address")

	allowedIPs, err := getAllowedIPs()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		return
	}

	if len(allowedIPs) == 0 {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(ipAddress))
		return
	}

	for _, allowedIP := range allowedIPs {
		if ipAddress == allowedIP {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(fmt.Sprintf("Access granted. Your source IP (%s) matches an allowed IP.\n", ipAddress)))
			return
		}
	}

	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte(fmt.Sprintf("Access denied. Your source IP (%s) doesn't match the allowed IPs (%s)\n", ipAddress, strings.Join(allowedIPs, ", "))))
}

func handlePost(w http.ResponseWriter, r *http.Request) {
	if !checkBasicAuth(r) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
		w.Write([]byte("Unauthorized"))
		return
	}

	ipAddress := r.FormValue("ip")
	if ipAddress == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("IP address is required"))
		return
	}

	if net.ParseIP(ipAddress) == nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid IP address format"))
		return
	}

	if projectID == "" {
		addIPToEnv(ipAddress)
	} else {
		if err := addIPToFirestore(ipAddress); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Failed to add IP address"))
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("IP address %s added successfully", ipAddress)))
}

func handleDelete(w http.ResponseWriter, r *http.Request) {
	if !checkBasicAuth(r) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
		w.Write([]byte("Unauthorized"))
		return
	}

	// Parse the form data manually
	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to parse form data"))
		return
	}

	// Read the body to get the form data
	var body []byte
	if r.Body != nil {
		body, _ = io.ReadAll(r.Body)
	}

	// Parse the body to get the IP address
	values, err := url.ParseQuery(string(body))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Failed to parse body data"))
		return
	}

	ipAddress := values.Get("ip")
	if ipAddress == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("IP address is required"))
		return
	}

	if net.ParseIP(ipAddress) == nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid IP address format"))
		return
	}

	if projectID == "" {
		removeIPFromEnv(ipAddress)
	} else {
		if err := removeIPFromFirestore(ipAddress); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Failed to remove IP address"))
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("IP address %s removed successfully", ipAddress)))
}

func getAllowedIPs() ([]string, error) {
	if projectID == "" {
		if allowedIPEnv := os.Getenv("ALLOWED_IP"); allowedIPEnv != "" {
			return strings.Split(allowedIPEnv, ","), nil
		} else {
			return []string{}, nil
		}
	}

	var allowedIPs []string
	iter := client.Collection(firestoreCollection).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		allowedIPs = append(allowedIPs, doc.Ref.ID)
	}
	return allowedIPs, nil
}

func addIPToFirestore(ip string) error {
	_, err := client.Collection(firestoreCollection).Doc(ip).Set(ctx, map[string]interface{}{})
	return err
}

func removeIPFromFirestore(ip string) error {
	_, err := client.Collection(firestoreCollection).Doc(ip).Delete(ctx)
	return err
}

func addIPToEnv(ip string) {
	allowedIPEnv := os.Getenv("ALLOWED_IP")
	allowedIPs := strings.Split(allowedIPEnv, ",")
	for _, allowedIP := range allowedIPs {
		if allowedIP == ip {
			return
		}
	}
	if allowedIPEnv == "" {
		os.Setenv("ALLOWED_IP", ip)
	} else {
		allowedIPs = append(allowedIPs, ip)
		os.Setenv("ALLOWED_IP", strings.Join(allowedIPs, ","))
	}
}

func removeIPFromEnv(ip string) {
	allowedIPs := strings.Split(os.Getenv("ALLOWED_IP"), ",")
	var updatedIPs []string
	for _, allowedIP := range allowedIPs {
		if allowedIP != ip {
			updatedIPs = append(updatedIPs, allowedIP)
		}
	}
	os.Setenv("ALLOWED_IP", strings.Join(updatedIPs, ","))
}

func listAllowedIPs(w http.ResponseWriter, r *http.Request) {
	allowedIPs, err := getAllowedIPs()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(allowedIPs); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to encode response"))
		return
	}
}

func checkBasicAuth(r *http.Request) bool {
	auth := r.Header.Get("Authorization")
	if auth == "" {
		return false
	}

	// Check if the auth method is Basic
	if !strings.HasPrefix(auth, "Basic ") {
		return false
	}

	// Decode the base64 encoded credentials
	encodedCredentials := strings.TrimPrefix(auth, "Basic ")
	decodedBytes, err := base64.StdEncoding.DecodeString(encodedCredentials)
	if err != nil {
		return false
	}
	decodedCredentials := string(decodedBytes)

	// Get the credentials from the environment variable
	expectedCredentials := os.Getenv("BASIC_AUTH")
	return decodedCredentials == expectedCredentials
}
