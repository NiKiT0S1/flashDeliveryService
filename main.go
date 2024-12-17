//package main
//
//import (
//	"database/sql"
//	_ "database/sql"
//	"encoding/json"
//	"fmt"
//	_ "github.com/lib/pq"
//
//	//_ "gorm.io/gorm"
//	"log"
//	"net/http"
//)
//
//type User struct {
//	ID        uint   `gorm:"primaryKey"`
//	Name      string `gorm:"size:100"`
//	Email     string `gorm:"size:100"`
//	CreatedAt string
//	UpdatedAt string
//}
//
//type RequestData struct {
//	Message string `json:"message"`
//}
//
//type ResponseData struct {
//	Status  string `json:"status"`
//	Message string `json:"message"`
//}
//
//func handleRequest(w http.ResponseWriter, r *http.Request) {
//	if r.Method == http.MethodPost || r.Method == http.MethodGet {
//		var request RequestData
//
//		err := json.NewDecoder(r.Body).Decode(&request)
//		if err != nil || request.Message == "" {
//			response := ResponseData{
//				Status:  "error",
//				Message: "Incorrect JSON-message",
//			}
//			w.Header().Set("Content-Type", "application/json")
//			w.WriteHeader(http.StatusBadRequest)
//			json.NewEncoder(w).Encode(response)
//			return
//		}
//
//		fmt.Println("Get message: ", request.Message)
//
//		response := ResponseData{
//			Status:  "success",
//			Message: "Data received successfully",
//		}
//		w.Header().Set("Content-Type", "application/json")
//		w.WriteHeader(http.StatusOK)
//		json.NewEncoder(w).Encode(response)
//	} else {
//		http.Error(w, "Only POST and GET requests are supported", http.StatusMethodNotAllowed)
//	}
//}
//
//func main() {
//	http.HandleFunc("/api", handleRequest)
//
//	log.Println("Server started on :8080")
//	err := http.ListenAndServe(":8080", nil)
//	if err != nil {
//		log.Fatal("Server failed: ", err)
//	}
//
//	// CONNECT TO DATABASE
//	connStr := "user=postgres dbname=flash_delivery_service sslmode=disable"
//	db, err := sql.Open("postgres", connStr)
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer db.Close()
//
//	err = db.Ping()
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println("Connected to postgres")
//
//	//db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
//	//if err != nil {
//	//	log.Fatal(err)
//	//}
//	//db.AutoMigrate(&User{})
//	//
//	//newUser := User{Name: "John Doe", Email: "johndoe@gmail.com"}
//	//db.Create(&newUser)
//	//
//	//var user User
//	//db.First(&user, 1)
//	//fmt.Printf("User ID:", user)
//	//
//	//db.Model(&user).Update("Name", "John Doe")
//	//
//	//db.Delete(&user, 1)
//}

//package main
//
//import (
//	"encoding/json"
//	_ "encoding/json"
//	"fmt"
//	_ "fmt"
//	"gorm.io/driver/postgres"
//	_ "gorm.io/driver/postgres"
//	"gorm.io/gorm"
//	_ "gorm.io/gorm"
//	"log"
//	_ "log"
//	"net/http"
//	_ "net/http"
//)
//
//// ---- Structures ----
//
//// JSON structure for receiving data from the client
//type RequestData struct {
//	Message string `json:"message"`
//}
//
//// JSON structure for the response to the client
//type ResponseData struct {
//	Status  string `json:"status"`
//	Message string `json:"message"`
//}
//
//// Structure for the "users" table
//type User struct {
//	ID    uint   `gorm:"primaryKey"`
//	Name  string `json:"name"`
//	Email string `json:"email"`
//}
//
//// ---- Variables ----
//var db *gorm.DB
//
//// ---- Initializing the database ----
//func initDatabase() {
//	var err error
//	//dsn := "host=localhost user=postgres dbname=flash_delivery_service sslmode=disable"
//	dsn := "host=localhost user=postgres dbname=flash_delivery_service password=0000 sslmode=disable"
//	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
//	if err != nil {
//		log.Fatal("Couldn't connect to the database:", err)
//	}
//
//	// Automatic migration
//	db.AutoMigrate(&User{})
//	fmt.Println("Database was successfully connected and migrated!")
//}
//
//// ---- Handler for POST and GET JSON requests ----
//func handleJSON(w http.ResponseWriter, r *http.Request) {
//	if r.Method == http.MethodPost || r.Method == http.MethodGet {
//		var requestData RequestData
//
//		// Decoding JSON
//		err := json.NewDecoder(r.Body).Decode(&requestData)
//		if err != nil || requestData.Message == "" {
//			response := ResponseData{
//				Status:  "fail",
//				Message: "Incorrect JSON-message",
//			}
//			w.WriteHeader(http.StatusBadRequest)
//			json.NewEncoder(w).Encode(response)
//			return
//		}
//
//		// Displaying a message in the console
//		fmt.Println("Message received:", requestData.Message)
//
//		// Sending a successful response
//		response := ResponseData{
//			Status:  "success",
//			Message: "Data received successfully",
//		}
//		w.Header().Set("Content-Type", "application/json")
//		json.NewEncoder(w).Encode(response)
//	} else {
//		http.Error(w, "Method is unsupported", http.StatusMethodNotAllowed)
//	}
//}
//
//// Middleware для включения CORS
//func enableCORS(next http.Handler) http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		// Установка CORS-заголовков
//		w.Header().Set("Access-Control-Allow-Origin", "*") // Разрешить всем доменам
//		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
//		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
//
//		// Обработка preflight-запросов (OPTIONS)
//		if r.Method == http.MethodOptions {
//			w.WriteHeader(http.StatusOK)
//			return
//		}
//
//		// Продолжить выполнение основного обработчика
//		next.ServeHTTP(w, r)
//	})
//}
//
//// ---- CRUD operations with the users table ----
//func handleUsers(w http.ResponseWriter, r *http.Request) {
//	//w.Header().Set("Access-Control-Allow-Origin", "*")
//	//w.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE, OPTIONS")
//	//w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
//	//
//	switch r.Method {
//	//case http.MethodPost:
//	//	// Creating a new user
//	//	var user User
//	//	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
//	//		http.Error(w, "Incorrect JSON", http.StatusBadRequest)
//	//		return
//	//	}
//	//	db.Create(&user)
//	//	w.WriteHeader(http.StatusCreated)
//	//	json.NewEncoder(w).Encode(user)
//	//
//	case http.MethodGet:
//		// Getting all users
//		var users []User
//		db.Find(&users)
//		json.NewEncoder(w).Encode(users)
//
//	//case http.MethodPut:
//	//	// Updating a user by ID
//	//	var user User
//	//	id := r.URL.Query().Get("id")
//	//	if err := db.First(&user, id).Error; err != nil {
//	//		http.Error(w, "User is not found", http.StatusNotFound)
//	//		return
//	//	}
//	//	var input User
//	//	json.NewDecoder(r.Body).Decode(&input)
//	//	user.Name = input.Name // Обновляем только имя
//	//	db.Save(&user)
//	//	//json.NewDecoder(r.Body).Decode(&user)
//	//	//db.Save(&user)
//	//	json.NewEncoder(w).Encode(user)
//	//
//	//case http.MethodDelete:
//	//	// Deleting a user by ID
//	//	id := r.URL.Query().Get("id")
//	//	db.Delete(&User{}, id)
//	//	w.WriteHeader(http.StatusNoContent)
//	//
//	default:
//		http.Error(w, "Method is unsupported", http.StatusMethodNotAllowed)
//	}
//
//	if r.Method == http.MethodPost {
//		// Читаем данные из тела запроса
//		var user User
//		err := json.NewDecoder(r.Body).Decode(&user)
//		if err != nil {
//			http.Error(w, "Invalid request payload", http.StatusBadRequest)
//			return
//		}
//
//		// Логируем полученные данные (для отладки)
//		fmt.Printf("Received user: %+v\n", user)
//
//		// Возвращаем успешный ответ
//		w.Header().Set("Content-Type", "application/json")
//		w.WriteHeader(http.StatusCreated)
//		json.NewEncoder(w).Encode(map[string]string{
//			"status": "User created successfully",
//			"name":   user.Name,
//			"email":  user.Email,
//		})
//	} else {
//		// Если метод не POST, возвращаем ошибку
//		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
//	}
//}
//
//// ---- Main Function ----
//func main() {
//	// Initializing the database
//	initDatabase()
//
//	// Routes
//	http.Handle("/api/users", enableCORS(http.HandlerFunc(handleUsers)))
//	//http.HandleFunc("/api/json", handleJSON)   // JSON Processing
//	//http.HandleFunc("/api/users", handleUsers) // CRUD operations with users
//
//	// Starting the server
//	fmt.Println("The server is running on port 8080...")
//	log.Fatal(http.ListenAndServe(":8080", nil))
//}

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
		// For GET request, no body to read. We can respond with a simple message.
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
	// Получаем всех пользователей
	db.Find(&users)

	// Обновляем id каждого пользователя
	for i, user := range users {
		user.ID = uint(i + 1) // Устанавливаем новый id
		db.Save(&user)        // Сохраняем изменения
	}
}

func resetSequence() {
	// Сбросить последовательность для пользователей
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
		db.Save(&user)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)

	case http.MethodDelete: // Delete User
		id := r.URL.Query().Get("id")
		if err := db.Delete(&User{}, id).Error; err != nil {
			http.Error(w, "Failed to delete user", http.StatusInternalServerError)
			return
		}
		// После удаления пользователя перенумеруем id и сбросим последовательность
		renumberIDs()
		resetSequence()

		// Ответ об успешном удалении
		w.WriteHeader(http.StatusNoContent)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func main() {
	initDatabase()

	http.Handle("/api/users", enableCORS(http.HandlerFunc(handleUsers)))
	http.Handle("/api/json", enableCORS(http.HandlerFunc(handleJSON))) // JSON Processing endpoint

	fmt.Println("Server is running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
