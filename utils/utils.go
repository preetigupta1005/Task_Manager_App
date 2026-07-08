package utils

import (
	"My-todo-app/models"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func ParseBody(r *http.Request, v interface{}) error {
	return json.NewDecoder(r.Body).Decode(v)
}

func EncodeJSONBody(resp http.ResponseWriter, data interface{}) error {
	return json.NewEncoder(resp).Encode(data)
}

func RespondJSON(w http.ResponseWriter, statusCode int, body interface{}) {
	w.WriteHeader(statusCode)
	if body != nil {
		if err := EncodeJSONBody(w, body); err != nil {
			log.Printf("Failed to respond JSON with error: %+v", err)
		}
	}
}

func newClientError(err error, statusCode int, messageToUser string) *models.ClientError {
	errMsg := ""
	if err != nil {
		errMsg = err.Error()
	}
	return &models.ClientError{
		MessageToUser: messageToUser,
		Err:           errMsg,
		StatusCode:    statusCode,
	}
}

func RespondError(w http.ResponseWriter, statusCode int, err error, messageToUser string) {
	log.Printf("status: %d, message: %s, err: %+v ", statusCode, messageToUser, err)
	clientError := newClientError(err, statusCode, messageToUser)
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(clientError); err != nil {
		log.Printf("status: %d, message: %s, err: %+v ", statusCode, messageToUser, err)
	}
}

func HashedPassword(password string) (string, error) {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashPassword), err
}

func CheckPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func GenerateJWT(userID string, sessionID string) (string, error) {
	claims := jwt.MapClaims{
		"userID":    userID,
		"sessionID": sessionID,
		"exp":       time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
}

func ParseJWT(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}
