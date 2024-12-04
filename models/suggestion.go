package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Suggestion struct {
	Id         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	By         int                `bson:"by" json:"by"`
	Suggestion string             `bson:"suggestion" json:"suggestion"`
	Reply      string             `bson:"reply" json:"reply"`
	Votes      int                `bson:"votes" json:"votes"`
	Tags       []string           `bson:"tags,omitempty" json:"tags,omitempty"` // Omit if nil
	Views      int                `bson:"views" json:"views"`
	Status     string             `bson:"status" json:"status"`
	CreatedAt  time.Time          `bson:"createdAt" json:"createdAt"`
}
