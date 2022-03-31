package handlers

import (
	"encoding/base64"
	"encoding/json"
	"net/http"

	"com.example/database"
	"com.example/helper"
	"com.example/tokenManager"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AuthUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var params = mux.Vars(r)

	guid, err := primitive.ObjectIDFromHex(params["guid"])

	if err != nil {
		helper.GetError(err, w)
		return
	}

	user, err := database.GetUser(guid)

	if err != nil {
		helper.GetError(err, w)
		return
	}

	accessToken, err := tokenManager.GenerateAccessToken(user.GUID)

	if err != nil {
		helper.GetError(err, w)
		return
	}

	refreshToken, err := tokenManager.GenerateRefreshToken()

	if err != nil {
		helper.GetError(err, w)
		return
	}

	err = tokenManager.SaveRefreshToken(refreshToken, accessToken)
	refreshTokenBase64 := base64.StdEncoding.EncodeToString([]byte(refreshToken))

	var tokens map[string]string = make(map[string]string)
	tokens["accessToken"] = accessToken
	tokens["refreshToken"] = refreshTokenBase64

	json.NewEncoder(w).Encode(tokens)

}

func RefreshHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var tokens map[string]string = make(map[string]string)

	json.NewDecoder(r.Body).Decode(&tokens)

	accessToken, refreshToken, err := tokenManager.RefreshTokens(tokens["accessToken"], tokens["refreshToken"])

	if err != nil {
		helper.GetError(err, w)
		return
	}

	err = tokenManager.SaveRefreshToken(refreshToken, accessToken)
	if err != nil {
		helper.GetError(err, w)
		return
	}

	refreshTokenBase64 := base64.StdEncoding.EncodeToString([]byte(refreshToken))

	tokens["accessToken"] = accessToken
	tokens["refreshToken"] = refreshTokenBase64

	json.NewEncoder(w).Encode(tokens)

}
