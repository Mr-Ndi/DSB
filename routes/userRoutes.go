// routes/userRoutes.go
package routes

import (
	"dsb/controllers"

	"github.com/gorilla/mux"
)

// UserRoutes registers user-related routes.
// func UserRoutes(r *mux.Router, userCollection *mongo.Collection) {
// 	r.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
// 		controllers.CreateUser(w, r, userCollection)
// 	}).Methods("POST")
// }

func UserRoutes(routes *mux.Router) {

	routes.HandleFunc("/user", controllers.CreateUser).Methods("POST")
	routes.HandleFunc("/login", controllers.HandleLogin).Methods("POST")
}
