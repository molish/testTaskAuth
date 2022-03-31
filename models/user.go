package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	GUID    primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Session Session            `json:"session,omitempty" bson:"session,omitempty"`
}
