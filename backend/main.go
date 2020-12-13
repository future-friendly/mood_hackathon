package main

import (
	"github.com/future-friednly/mood/backend/agents"
	"github.com/future-friednly/mood/backend/auth"
	"github.com/future-friednly/mood/backend/data"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	r := mux.NewRouter()

	// auth handlers
	r.HandleFunc("/auth/signup", auth.HandleSignup).Methods("POST")
	r.HandleFunc("/auth/login", auth.HandleLogin).Methods("POST")
	r.HandleFunc("/auth/logout", auth.HandleLogout).Methods("POST")

	// agents handlers
	r.HandleFunc("/agent/get", agents.HandleGetAgents).Methods("POST")
	r.HandleFunc("/agent/add", agents.HandleCreateAgent).Methods("POST")
	r.HandleFunc("/agent/confirm", agents.HandleConfirmAgent).Methods("POST")
	r.HandleFunc("/agent/delete", agents.HandleDeleteAgent).Methods("POST")

	// data handlers
	r.HandleFunc("/data/newpage", data.HandleNewPage).Methods("POST")

	// charts
	r.HandleFunc("/chart/get", data.HandleGetChart).Methods("POST")

	// middleware
	r.Use(auth.AuthMiddleware)

	log.Println("starting backend server")
	http.ListenAndServe(":8080", r)
}