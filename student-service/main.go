package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
)

// Student represents a dormitory resident.
type Student struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Phone  string `json:"phone"`
	RoomID *int   `json:"roomId"`
}

var (
	mu       sync.RWMutex
	students = []Student{
		{ID: 1, Name: "Alice Smith", Email: "alice@example.com", Phone: "555-0101"},
		{ID: 2, Name: "Bob Jones", Email: "bob@example.com", Phone: "555-0102"},
	}
	nextID = 3
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
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok", "service": "student-service"})
	})

	// GET /students
	mux.HandleFunc("GET /students", func(w http.ResponseWriter, r *http.Request) {
		mu.RLock()
		defer mu.RUnlock()
		writeJSON(w, http.StatusOK, students)
	})

	// POST /students
	mux.HandleFunc("POST /students", func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Name  string `json:"name"`
			Email string `json:"email"`
			Phone string `json:"phone"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Name == "" || body.Email == "" {
			errorJSON(w, http.StatusBadRequest, "name and email are required")
			return
		}
		mu.Lock()
		s := Student{ID: nextID, Name: body.Name, Email: body.Email, Phone: body.Phone}
		nextID++
		students = append(students, s)
		mu.Unlock()
		writeJSON(w, http.StatusCreated, s)
	})

	// GET /students/{id}
	mux.HandleFunc("GET /students/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			errorJSON(w, http.StatusBadRequest, "invalid id")
			return
		}
		mu.RLock()
		defer mu.RUnlock()
		for _, s := range students {
			if s.ID == id {
				writeJSON(w, http.StatusOK, s)
				return
			}
		}
		errorJSON(w, http.StatusNotFound, "Student not found")
	})

	// PUT /students/{id}
	mux.HandleFunc("PUT /students/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			errorJSON(w, http.StatusBadRequest, "invalid id")
			return
		}
		var body struct {
			Name  string `json:"name"`
			Email string `json:"email"`
			Phone string `json:"phone"`
		}
		json.NewDecoder(r.Body).Decode(&body) //nolint:errcheck
		mu.Lock()
		defer mu.Unlock()
		for i, s := range students {
			if s.ID == id {
				if body.Name != "" {
					students[i].Name = body.Name
				}
				if body.Email != "" {
					students[i].Email = body.Email
				}
				if body.Phone != "" {
					students[i].Phone = body.Phone
				}
				writeJSON(w, http.StatusOK, students[i])
				return
			}
		}
		errorJSON(w, http.StatusNotFound, "Student not found")
	})

	// DELETE /students/{id}
	mux.HandleFunc("DELETE /students/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			errorJSON(w, http.StatusBadRequest, "invalid id")
			return
		}
		mu.Lock()
		defer mu.Unlock()
		for i, s := range students {
			if s.ID == id {
				students = append(students[:i], students[i+1:]...)
				writeJSON(w, http.StatusOK, s)
				return
			}
		}
		errorJSON(w, http.StatusNotFound, "Student not found")
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "3001"
	}
	log.Printf("Student Service running on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
