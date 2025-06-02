package utils

import (
	"encoding/json"
	"fmt"
	"github.com/form3tech-oss/jwt-go"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
	"time"
)

func EncodeJSONBody(res http.ResponseWriter, body interface{}) error {
	err := json.NewEncoder(res).Encode(body)
	if err != nil {
		return err
	}
	return nil
}

func DecodeJSONBody(req *http.Request, out interface{}) error {
	err := json.NewDecoder(req.Body).Decode(&out)

	if err != nil {
		return err
	}
	return nil
}
func RespondJSON(w http.ResponseWriter, statusCode int, body interface{}) {
	w.WriteHeader(statusCode)
	if body != nil {
		if err := EncodeJSONBody(w, body); err != nil {
			logrus.Errorf("Failed to respond JSON with error: %+v", err)
		}
	}
}
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w /n", err)
	}
	return string(hashedPassword), nil
}

func CheckPassword(password string, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))

}
func CreateToken(userId string) (string, error) {
	godotenv.Load()
	var secretKey = []byte(os.Getenv("SECRET_KEY"))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"userId": userId,
			"exp":    time.Now().Add(time.Hour * 24).Unix(),
		})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func VerifyToken(tokenString string) error {
	godotenv.Load()
	var secretKey = []byte(os.Getenv("SECRET_KEY"))
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return err
	}

	if !token.Valid {
		return fmt.Errorf("invalid token")
	}

	return nil
}
