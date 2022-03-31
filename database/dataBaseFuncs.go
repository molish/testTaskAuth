package database

import (
	"context"
	"time"

	"com.example/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var collection = ConnectDB()

func GetUser(id primitive.ObjectID) (models.User, error) {
	var user models.User
	filter := bson.M{"_id": id}
	err := collection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		return user, err
	}
	return user, nil
}

func SaveRefreshToken(refreshToken string, id primitive.ObjectID) error {
	err := saveSession(refreshToken, id, time.Now().AddDate(0, 6, 0))
	if err != nil {
		return err
	}
	return nil
}

func DeleteUserSession(id primitive.ObjectID) error {
	err := saveSession("", id, time.Now().AddDate(-10, 0, 0))
	if err != nil {
		return err
	}
	return nil
}

func saveSession(refreshToken string, id primitive.ObjectID, expiresAt time.Time) error {
	session := models.Session{
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
	}

	filter := bson.M{"_id": id}
	update := bson.D{
		{"$set", bson.D{
			{"session", session},
		}},
	}
	err := collection.FindOneAndUpdate(context.TODO(), filter, update)
	if err != nil {
		return err.Err()
	}
	return nil
}
