package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

// Booking represents a room reservation.
type Booking struct {
	ID        int    `json:"id"`
	StudentID int    `json:"studentId"`
	RoomID    int    `json:"roomId"`
	CheckIn   string `json:"checkIn"`
	CheckOut  string `json:"checkOut"`
	Status    string `json:"status"`
	CreatedAt string `json:"createdAt"`
}

var (
	mu       sync.RWMutex
	bookings = []Booking{}
	nextID   = 1
)

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v) //nolint:errcheck
}

func errorJSON(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func main() {
	mux := http.NewServeMux()

	// Health check
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok", "service": "booking-service"})
	})

	// GET /bookings
	mux.HandleFunc("GET /bookings", func(w http.ResponseWriter, r *http.Request) {
		mu.RLock()
		defer mu.RUnlock()
		writeJSON(w, http.StatusOK, bookings)
	})

	// POST /bookings
	mux.HandleFunc("POST /bookings", func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			StudentID int    `json:"studentId"`
			RoomID    int    `json:"roomId"`
			CheckIn   string `json:"checkIn"`
			CheckOut  string `json:"checkOut"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil ||
			body.StudentID == 0 || body.RoomID == 0 || body.CheckIn == "" || body.CheckOut == "" {
			errorJSON(w, http.StatusBadRequest, "studentId, roomId, checkIn and checkOut are required")
			return
		}
		mu.Lock()
		b := Booking{
			ID:        nextID,
			StudentID: body.StudentID,
			RoomID:    body.RoomID,
			CheckIn:   body.CheckIn,
			CheckOut:  body.CheckOut,
			Status:    "active",
			CreatedAt: time.Now().UTC().Format(time.RFC3339),
		}
		nextID++
		bookings = append(bookings, b)
		mu.Unlock()
		writeJSON(w, http.StatusCreated, b)
	})

	// GET /bookings/{id}
	mux.HandleFunc("GET /bookings/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			errorJSON(w, http.StatusBadRequest, "invalid id")
			return
		}
		mu.RLock()
		defer mu.RUnlock()
		for _, b := range bookings {
			if b.ID == id {
				writeJSON(w, http.StatusOK, b)
				return
			}
		}
		errorJSON(w, http.StatusNotFound, "Booking not found")
	})

	// PUT /bookings/{id}
	mux.HandleFunc("PUT /bookings/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			errorJSON(w, http.StatusBadRequest, "invalid id")
			return
		}
		var body struct {
			CheckIn  string `json:"checkIn"`
			CheckOut string `json:"checkOut"`
			Status   string `json:"status"`
		}
		json.NewDecoder(r.Body).Decode(&body) //nolint:errcheck
		mu.Lock()
		defer mu.Unlock()
		for i, b := range bookings {
			if b.ID == id {
				if body.CheckIn != "" {
					bookings[i].CheckIn = body.CheckIn
				}
				if body.CheckOut != "" {
					bookings[i].CheckOut = body.CheckOut
				}
				if body.Status != "" {
					bookings[i].Status = body.Status
				}
				writeJSON(w, http.StatusOK, bookings[i])
				return
			}
		}
		errorJSON(w, http.StatusNotFound, "Booking not found")
	})

	// DELETE /bookings/{id}  (soft cancel)
	mux.HandleFunc("DELETE /bookings/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			errorJSON(w, http.StatusBadRequest, "invalid id")
			return
		}
		mu.Lock()
		defer mu.Unlock()
		for i, b := range bookings {
			if b.ID == id {
				bookings[i].Status = "cancelled"
				writeJSON(w, http.StatusOK, bookings[i])
				return
			}
		}
		errorJSON(w, http.StatusNotFound, "Booking not found")
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "3003"
	}
	log.Printf("Booking Service running on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
