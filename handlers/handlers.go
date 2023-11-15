package handlers

import (
  "net/http"
  "github.com/gorilla/mux"
  "neovim-tips/database"
  "neovim-tips/models"
  "encoding/json"
  "fmt"
  "strconv"
  "neovim-tips/utils"
  "time"
  "github.com/golang-jwt/jwt"
  "golang.org/x/crypto/bcrypt"
  "neovim-tips/middleware"
  )


func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var credentials model.User
	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var user model.User
	if result := database.DB.Where("username = ?", credentials.Username).First(&user); result.Error != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password)); err != nil {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	expirationTime := time.Now().Add(1 * time.Hour)
	claims := &middleware.Claims{
		Username: credentials.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(middleware.JWTKey)
	if err != nil {
		http.Error(w, "Error in generating token", http.StatusInternalServerError)
		return
	}

	w.Write([]byte(tokenString))
}

func TotalTipsHandler(w http.ResponseWriter, r *http.Request) {
	var count int64
	result := database.DB.Model(&model.Tip{}).Count(&count)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte(fmt.Sprintf("%d", count)))
}

func RandomTipHandler(w http.ResponseWriter, r *http.Request) {
	var randomTip model.Tip

	if result := database.DB.Order("RANDOM()").First(&randomTip); result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte(randomTip.Content))
}

func SpecificTipHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var tip model.Tip
	if result := database.DB.First(&tip, id); result.Error != nil {
		http.Error(w, "Tip not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte(tip.Content))
}

func AddTipHandler(w http.ResponseWriter, r *http.Request) {
	var newTip model.Tip
	if err := json.NewDecoder(r.Body).Decode(&newTip); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if utils.TipExists(newTip.Content) {
		http.Error(w, "Tip already exists", http.StatusConflict)
		return
	}

	if result := database.DB.Create(&newTip); result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(newTip)
}

func DeleteTipHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	result := database.DB.Delete(&model.Tip{}, id)
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

func EditTipHandler(w http.ResponseWriter, r *http.Request) {
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

	if updatedContent.Content == "" {
		http.Error(w, "New content cannot be empty", http.StatusBadRequest)
		return
	}

	result := database.DB.Model(&model.Tip{}).Where("id = ?", id).Update("content", updatedContent.Content)
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
