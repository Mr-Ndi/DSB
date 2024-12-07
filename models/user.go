package models

type User struct {
	Regnumber       int    `bson:"regnumber"`
	Username        string `bson:"username"`
	Password        string `bson:"password"`
	SuggestionCount int    `bson:"suggestion_count"`
}

type LoginResponse struct {
	Regnumber int    `bson:"regnumber"`
	Username  string `bson:"username"`
	Token     string `bson:"token"`
}
type Admin struct {
	Role     string `bson:"role"`
	Password string `bson:"password"`
}
