package main

import (
	"encoding/json"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
)

// ---- Structures ----
type RequestData struct {
	Message string `json:"message"`
}

// JSON structure for the response to the client
type ResponseData struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type User struct {
	ID    uint   `gorm:"primaryKey"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

var db *gorm.DB

// ---- Initialize Database ----
func initDatabase() {
	var err error
	dsn := "host=localhost user=postgres dbname=flash_delivery_service password=0000 sslmode=disable"
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Couldn't connect to the database:", err)
	}
	db.AutoMigrate(&User{})
	fmt.Println("Database connected and migrated successfully!")
}

// ---- Handler for POST and GET JSON requests ----
func handleJSON(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var requestData RequestData

		// Decoding JSON
		err := json.NewDecoder(r.Body).Decode(&requestData)
		if err != nil || requestData.Message == "" {
			response := ResponseData{
				Status:  "fail",
				Message: "Incorrect JSON-message",
			}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(response)
			return
		}

		// Displaying a message in the console
		fmt.Println("Message received:", requestData.Message)

		// Sending a successful response
		response := ResponseData{
			Status:  "success",
			Message: "Data received successfully",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	} else if r.Method == http.MethodGet {
		// For GET request, NOBODY to read. We can respond with a simple message.
		response := ResponseData{
			Status:  "success",
			Message: "GET request received",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// ---- Enable CORS Middleware ----
func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func renumberIDs() {
	var users []User
	// get all users
	db.Find(&users)

	// Update ID each of users
	for i, user := range users {
		user.ID = uint(i + 1) // Set new ID
		db.Save(&user)        // Save changes
	}
}

// ---- Handler to reset IDs ----
func resetIDs(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		resetSequence()
		w.Header().Set("Content-Type", "application/json")
		response := map[string]string{
			"status":  "success",
			"message": "IDs have been reset successfully.",
		}
		json.NewEncoder(w).Encode(response)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func resetSequence() {
	// Reset sequences for users
	db.Exec("SELECT setval('users_id_seq', (SELECT max(id) FROM users));")
}

// ---- CRUD Handler ----
func handleUsers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var user User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		db.Create(&user)

		//renumberIDs()
		resetSequence()

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)

	case http.MethodGet:
		var users []User
		db.Find(&users)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(users)

	case http.MethodPut: // Update User
		id := r.URL.Query().Get("id")
		var user User
		if err := db.First(&user, id).Error; err != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		var input User
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		user.Name = input.Name
		user.Email = input.Email

		//renumberIDs()
		resetSequence()

		db.Save(&user)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)

	case http.MethodDelete: // Delete User
		id := r.URL.Query().Get("id")
		if err := db.Delete(&User{}, id).Error; err != nil {
			http.Error(w, "Failed to delete user", http.StatusInternalServerError)
			return
		}
		// After remove users rename id and reset sequences
		//renumberIDs()
		resetSequence()

		// Answer about successful remove
		w.WriteHeader(http.StatusNoContent)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func main() {
	initDatabase()

	http.Handle("/api/users", enableCORS(http.HandlerFunc(handleUsers)))
	http.Handle("/api/json", enableCORS(http.HandlerFunc(handleJSON))) // JSON Processing endpoint
	http.HandleFunc("/api/users/reset", http.HandlerFunc(resetIDs))

	fmt.Println("Server is running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
