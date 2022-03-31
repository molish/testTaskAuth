package tokenManager

import (
	"encoding/base64"
	"errors"
	"math/rand"
	"strings"
	"time"

	"com.example/database"
	"com.example/models"
	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type userClaim struct {
	GUID primitive.ObjectID `json:"guid"`
	jwt.StandardClaims
}

var mySigningKey = []byte("298f8294948x29989#8M(&#&(*#@H<)@)#<_F*#_FH#<B#*#&*#*#&$9834hf,8#*3f8034f7309#H3,H0")
var salt = []byte("8072h578h428f527g24h08j57f268j2d982-4f702h8457f0422hf6807240j745hf245278fh24j289f02478f248")

func GenerateAccessToken(guid primitive.ObjectID) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, userClaim{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Duration(15) * time.Minute).Unix(),
		},
		GUID: guid,
	})

	return token.SignedString(mySigningKey)
}

func GenerateRefreshToken() (string, error) {
	token := make([]byte, 64)
	rndSource := rand.NewSource(time.Now().Unix())
	rndRes := rand.New(rndSource)

	_, err := rndRes.Read(token)
	if err != nil {
		return string(token), err
	}
	return string(token), nil
}

func SaveRefreshToken(refreshToken string, accessToken string) error {

	guid, err := getGuidFromAccessToken(accessToken)
	if err != nil {
		return err
	}

	refreshTokenForHAshing := refreshToken + accessToken + string(salt)

	refreshTokenHash, err := bcrypt.GenerateFromPassword([]byte(refreshTokenForHAshing), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	err = database.SaveRefreshToken(string(refreshTokenHash), guid)
	if err != nil {
		return err
	}

	return nil
}

func RefreshTokens(accessToken string, refreshToken string) (string, string, error) {

	guid, err := getGuidFromAccessToken(accessToken)
	if err != nil {
		return "", "", err
	}

	refreshTokenBytes, err := base64.StdEncoding.DecodeString(refreshToken)

	var user models.User
	user, err = database.GetUser(guid)

	if time.Now().After(user.Session.ExpiresAt) {
		return "", "", errors.New("refresh token Expired")
	}

	refreshTokenForCompare := string(refreshTokenBytes) + accessToken + string(salt)
	err = bcrypt.CompareHashAndPassword([]byte(user.Session.RefreshToken), []byte(refreshTokenForCompare))
	if err != nil {
		database.DeleteUserSession(guid)
		return "", "", errors.New("refresh token is invalid")
	}

	newAccesToken, err := GenerateAccessToken(guid)
	if err != nil {
		return "", "", err
	}
	newRefreshToken, err := GenerateRefreshToken()
	if err != nil {
		return "", "", err
	}

	return newAccesToken, newRefreshToken, nil
}

func getGuidFromAccessToken(accessToken string) (primitive.ObjectID, error) {
	segmentWithGUID := strings.Split(accessToken, ".")[1]
	decodedSegmentWithGUID, err := jwt.DecodeSegment(segmentWithGUID)
	if err != nil {
		return primitive.NilObjectID, err
	}
	GUID := strings.Split(string(decodedSegmentWithGUID), "\"")[3]
	guid, err := primitive.ObjectIDFromHex(GUID)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return guid, nil
}
