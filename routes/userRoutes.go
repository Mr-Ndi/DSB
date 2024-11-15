// routes/userRoutes.go
package routes

import (
	"dsb/controllers"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)

// UserRoutes registers user-related routes.
func UserRoutes(r *mux.Router, userCollection *mongo.Collection) {
	r.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		controllers.CreateUser(w, r, userCollection)
	}).Methods("POST")
}
