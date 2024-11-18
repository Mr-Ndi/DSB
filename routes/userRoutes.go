// routes/userRoutes.go
package routes

import (
	"dsb/controllers"

	"github.com/gorilla/mux"
)

func UserRoutes(routes *mux.Router) {
	routes.HandleFunc("/user", controllers.CreateUser).Methods("POST")
}
