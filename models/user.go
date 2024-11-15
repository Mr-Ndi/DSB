package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User represents a user model for MongoDB
type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"` // MongoDB ID
	Regnumber int                `bson:"regnumber"`     // Registration number
	Username  string             `bson:"username"`      // Username
	Password  string             `bson:"password"`      // Password
}
