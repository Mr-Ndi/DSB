package models

// User represents a user model for MongoDB
type User struct {
	Regnumber int    `bson:"regnumber"` // Registration number
	Username  string `bson:"username"`  // Username
	Password  string `bson:"password"`  // Password
}
type Authority struct {
	role     string `bson:"regnumber"` // Registration number
	Username string `bson:"username"`  // Username
	Password string `bson:"password"`  // Password
}
type LoginResponse struct {
	Regnumber int    `bson:"regnumber"`
	Username  string `bson:"username"`
	Token     string `bson:"token"`
}
