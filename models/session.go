package models

import (
	"time"
)

type Session struct {
	RefreshToken string    `json:"refreshToken,omitempty" bson:"refreshToken,omitempty"`
	ExpiresAt    time.Time `json:"expiresAt,omitempty" bson:"expiresAT,omitempty"`
}
