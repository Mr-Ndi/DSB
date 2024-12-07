package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Suggestion struct {
	Id        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	By        string             `bson:"by" json:"by"`
	Content   string             `bson:"suggestion" json:"content"`
	Reply     string             `bson:"reply" json:"reply"`
	Tags      []string           `bson:"tags,omitempty" json:"tags"` // Omit if nil
	Views     int                `bson:"views" json:"views"`
	Status    string             `bson:"status" json:"status"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
	Parent    string             `bson:"parent,omitempty" json:"parent"`
}

type Vote struct {
	Id           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	By           string             `bson:"by" json:"by"`
	Type         string             `bson:"type" json:"type"`
	SuggestionId string             `bson:"suggestionId" json:"suggestionId"`
}

// Response represents a response to a suggestion.
type Response struct {
	ID           string    `bson:"id" json:"id"`
	SuggestionId string    `bson:"suggestionId" json:"suggestionId"` // Reference to the suggestion
	Admin        string    `bson:"admin" json:"admin"`               // Admin responding
	Content      string    `bson:"content" json:"content"`           // Response content
	CreatedAt    time.Time `bson:"created_at" json:"created_at"`     // Timestamp of the response
}
