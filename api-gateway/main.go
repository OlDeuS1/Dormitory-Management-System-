package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v) //nolint:errcheck
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func newProxy(rawURL string) *httputil.ReverseProxy {
	target, err := url.Parse(rawURL)
	if err != nil {
		log.Fatalf("invalid target URL %s: %v", rawURL, err)
	}
	return httputil.NewSingleHostReverseProxy(target)
}

func main() {
	studentProxy := newProxy(getEnv("STUDENT_SERVICE_URL", "http://localhost:3001"))
	roomProxy := newProxy(getEnv("ROOM_SERVICE_URL", "http://localhost:3002"))
	bookingProxy := newProxy(getEnv("BOOKING_SERVICE_URL", "http://localhost:3003"))

	mux := http.NewServeMux()

	// Health check
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok", "service": "api-gateway"})
	})

	// Proxy /students and /students/...
	mux.HandleFunc("/students", studentProxy.ServeHTTP)
	mux.HandleFunc("/students/", studentProxy.ServeHTTP)

	// Proxy /rooms and /rooms/...
	mux.HandleFunc("/rooms", roomProxy.ServeHTTP)
	mux.HandleFunc("/rooms/", roomProxy.ServeHTTP)

	// Proxy /bookings and /bookings/...
	mux.HandleFunc("/bookings", bookingProxy.ServeHTTP)
	mux.HandleFunc("/bookings/", bookingProxy.ServeHTTP)

	// 404 for everything else
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "Route not found"})
	})

	port := getEnv("PORT", "3000")
	log.Printf("API Gateway running on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
