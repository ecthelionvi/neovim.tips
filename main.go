package main

import (
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"neovim-tips/models"
	"net/http"
	"os"
	"strconv"
	"time"
)

var db *gorm.DB

func init() {
	err := godotenv.Load() // This will load variables from a .env file
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Replace with your database connection details
	dsn := fmt.Sprintf("host=%s user=%s dbname=%s port=%s sslmode=%s password=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
		os.Getenv("SSL_MODE"),
		os.Getenv("DB_PASS"),
	)
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto migrate the Tip model
	err = db.AutoMigrate(&model.Tip{})
	if err != nil {
		log.Fatal("Failed to migrate Tip model:", err)
	}

	// Auto migrate the User model
	err = db.AutoMigrate(&model.User{})
	if err != nil {
		log.Fatal("Failed to migrate User model:", err)
	}
}

func createSuperUser(username, password string) {
	var count int64
	db.Model(&model.User{}).Where("username = ?", username).Count(&count)
	if count == 0 {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			log.Fatal("Failed to hash password:", err)
		}
		superUser := model.User{
			Username: username,
			Password: string(hashedPassword),
			IsSuper:  true,
		}
		db.Create(&superUser)
	}
}

var jwtKey = []byte(os.Getenv("JWT_SECRET_KEY"))

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	var credentials model.User
	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var user model.User
	if result := db.Where("username = ?", credentials.Username).First(&user); result.Error != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password)); err != nil {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	expirationTime := time.Now().Add(1 * time.Hour)
	claims := &Claims{
		Username: credentials.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		http.Error(w, "Error in generating token", http.StatusInternalServerError)
		return
	}

	w.Write([]byte(tokenString))
}

func authenticateJWT(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")

		// Validate token format
		if tokenString == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		claims := &Claims{}

		// Parse the JWT string and store the result in `claims`.
		tkn, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				http.Error(w, "Invalid token signature", http.StatusUnauthorized)
				return
			}
			http.Error(w, "Invalid token", http.StatusBadRequest)
			return
		}
		if !tkn.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}

func populateTips() {
	tips := []model.Tip{}

	for _, tip := range tips {
		if !tipExists(tip.Content) {
			db.Create(&tip)
		}
	}
}

func tipExists(content string) bool {
	var count int64
	db.Model(&model.Tip{}).Where("content = ?", content).Count(&count)
	return count > 0
}

func totalTipsHandler(w http.ResponseWriter, r *http.Request) {
	var count int64
	result := db.Model(&model.Tip{}).Count(&count)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte(fmt.Sprintf("%d", count)))
}

func randomTipHandler(w http.ResponseWriter, r *http.Request) {
	var randomTip model.Tip

	// Query to select a random tip
	if result := db.Order("RANDOM()").First(&randomTip); result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte(randomTip.Content))
}

func specificTipHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var tip model.Tip
	if result := db.First(&tip, id); result.Error != nil {
		http.Error(w, "Tip not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte(tip.Content))
}

func addTipHandler(w http.ResponseWriter, r *http.Request) {
	var newTip model.Tip
	if err := json.NewDecoder(r.Body).Decode(&newTip); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Check if the tip already exists
	if tipExists(newTip.Content) {
		http.Error(w, "Tip already exists", http.StatusConflict) // Or any other appropriate status
		return
	}

	// Create new tip
	if result := db.Create(&newTip); result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(newTip)
}

func deleteTipHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// Attempt to delete the tip with the given ID
	result := db.Delete(&model.Tip{}, id)
	if result.Error != nil {
		http.Error(w, "Error deleting tip", http.StatusInternalServerError)
		return
	}

	if result.RowsAffected == 0 {
		http.Error(w, "No tip found with given ID", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Tip deleted successfully"))
}

func editTipHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var updatedContent struct {
		Content string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&updatedContent); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Check if the new content is empty
	if updatedContent.Content == "" {
		http.Error(w, "New content cannot be empty", http.StatusBadRequest)
		return
	}

	result := db.Model(&model.Tip{}).Where("id = ?", id).Update("content", updatedContent.Content)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	if result.RowsAffected == 0 {
		http.Error(w, "No tip found with given ID", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Tip updated successfully"))
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/api/login", loginHandler).Methods("POST")
	r.HandleFunc("/api/total", totalTipsHandler).Methods("GET")
	r.HandleFunc("/api/random", randomTipHandler).Methods("GET")
	r.HandleFunc("/api/{id:[0-9]+}", specificTipHandler).Methods("GET")
	r.HandleFunc("/api/add", authenticateJWT(addTipHandler)).Methods("POST")
	r.HandleFunc("/api/edit/{id:[0-9]+}", authenticateJWT(editTipHandler)).Methods("PUT")
	r.HandleFunc("/api/delete/{id:[0-9]+}", authenticateJWT(deleteTipHandler)).Methods("DELETE")

	log.Println("Server is running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
