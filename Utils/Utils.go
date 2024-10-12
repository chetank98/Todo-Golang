package Utils

import (
	"encoding/json"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
	"time"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPassword(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

}

func Health(w http.ResponseWriter, r *http.Request) {
	RespondJSON(w, http.StatusOK, struct {
		Status string `json:"status"`
	}{Status: "server is running"})
}

func RespondError(w http.ResponseWriter, statusCode int, err error, message string) {
	logrus.Errorf("status: %d message: %s err: %+v ", statusCode, err, message)
	w.WriteHeader(statusCode)
}

func EncodeJSONBody(resp http.ResponseWriter, data interface{}) error {
	return json.NewEncoder(resp).Encode(data)
}

func RespondJSON(w http.ResponseWriter, statusCode int, body interface{}) {
	w.WriteHeader(statusCode)
	if body != nil {
		if err := EncodeJSONBody(w, body); err != nil {
			logrus.Errorf("Failed to respond with JSON %+v", err)
		}
	}
}

func GenerateJWT(userID, email, name, sessionID string) (string, error) {
	claims := jwt.MapClaims{
		"userID":    userID,
		"email":     email,
		"name":      name,
		"sessionID": sessionID,
		"iat":       time.Now().Unix(),
		"exp":       time.Now().Add(time.Hour * 2).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}
