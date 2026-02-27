package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
)

// Room represents a dormitory room.
type Room struct {
	ID        int    `json:"id"`
	Number    string `json:"number"`
	Type      string `json:"type"`
	Capacity  int    `json:"capacity"`
	Available bool   `json:"available"`
}

var (
	mu    sync.RWMutex
	rooms = []Room{
		{ID: 1, Number: "101", Type: "single", Capacity: 1, Available: true},
		{ID: 2, Number: "102", Type: "double", Capacity: 2, Available: true},
		{ID: 3, Number: "201", Type: "double", Capacity: 2, Available: true},
	}
	nextID = 4
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
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok", "service": "room-service"})
	})

	// GET /rooms  (supports ?available=true|false)
	mux.HandleFunc("GET /rooms", func(w http.ResponseWriter, r *http.Request) {
		mu.RLock()
		defer mu.RUnlock()
		availParam := r.URL.Query().Get("available")
		if availParam == "" {
			writeJSON(w, http.StatusOK, rooms)
			return
		}
		want := availParam == "true"
		filtered := make([]Room, 0)
		for _, rm := range rooms {
			if rm.Available == want {
				filtered = append(filtered, rm)
			}
		}
		writeJSON(w, http.StatusOK, filtered)
	})

	// POST /rooms
	mux.HandleFunc("POST /rooms", func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Number   string `json:"number"`
			Type     string `json:"type"`
			Capacity int    `json:"capacity"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Number == "" || body.Type == "" || body.Capacity == 0 {
			errorJSON(w, http.StatusBadRequest, "number, type and capacity are required")
			return
		}
		mu.Lock()
		rm := Room{ID: nextID, Number: body.Number, Type: body.Type, Capacity: body.Capacity, Available: true}
		nextID++
		rooms = append(rooms, rm)
		mu.Unlock()
		writeJSON(w, http.StatusCreated, rm)
	})

	// GET /rooms/{id}
	mux.HandleFunc("GET /rooms/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			errorJSON(w, http.StatusBadRequest, "invalid id")
			return
		}
		mu.RLock()
		defer mu.RUnlock()
		for _, rm := range rooms {
			if rm.ID == id {
				writeJSON(w, http.StatusOK, rm)
				return
			}
		}
		errorJSON(w, http.StatusNotFound, "Room not found")
	})

	// PUT /rooms/{id}
	mux.HandleFunc("PUT /rooms/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			errorJSON(w, http.StatusBadRequest, "invalid id")
			return
		}
		var body struct {
			Number    string `json:"number"`
			Type      string `json:"type"`
			Capacity  int    `json:"capacity"`
			Available *bool  `json:"available"`
		}
		json.NewDecoder(r.Body).Decode(&body) //nolint:errcheck
		mu.Lock()
		defer mu.Unlock()
		for i, rm := range rooms {
			if rm.ID == id {
				if body.Number != "" {
					rooms[i].Number = body.Number
				}
				if body.Type != "" {
					rooms[i].Type = body.Type
				}
				if body.Capacity != 0 {
					rooms[i].Capacity = body.Capacity
				}
				if body.Available != nil {
					rooms[i].Available = *body.Available
				}
				writeJSON(w, http.StatusOK, rooms[i])
				return
			}
		}
		errorJSON(w, http.StatusNotFound, "Room not found")
	})

	// DELETE /rooms/{id}
	mux.HandleFunc("DELETE /rooms/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			errorJSON(w, http.StatusBadRequest, "invalid id")
			return
		}
		mu.Lock()
		defer mu.Unlock()
		for i, rm := range rooms {
			if rm.ID == id {
				rooms = append(rooms[:i], rooms[i+1:]...)
				writeJSON(w, http.StatusOK, rm)
				return
			}
		}
		errorJSON(w, http.StatusNotFound, "Room not found")
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "3002"
	}
	log.Printf("Room Service running on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
