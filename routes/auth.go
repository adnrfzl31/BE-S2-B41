package routes

import (
	"BE-S2-B41/handlers"
	"BE-S2-B41/pkg/middleware"
	"BE-S2-B41/pkg/mysql"
	"BE-S2-B41/repositories"

	"github.com/gorilla/mux"
)

func AuthRoutes(r *mux.Router) {
	userRepository := repositories.RepositoryUser(mysql.DB)
	h := handlers.HandlerAuth(userRepository)

	r.HandleFunc("/register", middleware.UploadFile(h.Register)).Methods("POST")
	r.HandleFunc("/login", h.Login).Methods("POST")
	r.HandleFunc("/check-auth", middleware.Auth(h.CheckAuth)).Methods("GET")
}
