package main

import (
  "log"
  "net/http"
  "github.com/gorilla/mux"
  "neovim-tips/middleware"
  "neovim-tips/handlers"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/api/login", handlers.LoginHandler).Methods("POST")
	r.HandleFunc("/api/total", handlers.TotalTipsHandler).Methods("GET")
	r.HandleFunc("/api/random", handlers.RandomTipHandler).Methods("GET")
	r.HandleFunc("/api/{id:[0-9]+}", handlers.SpecificTipHandler).Methods("GET")
	r.HandleFunc("/api/add", middleware.AuthenticateJWT(handlers.AddTipHandler)).Methods("POST")
	r.HandleFunc("/api/edit/{id:[0-9]+}", middleware.AuthenticateJWT(handlers.EditTipHandler)).Methods("PUT")
	r.HandleFunc("/api/delete/{id:[0-9]+}", middleware.AuthenticateJWT(handlers.DeleteTipHandler)).Methods("DELETE")

	log.Println("Server is running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
