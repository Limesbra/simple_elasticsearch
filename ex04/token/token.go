package token

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

// структура для хранения токена
type token struct {
	T string `json:"token"`
}

// Функция для создания токена
func createToken(secret []byte) string {
	token := jwt.New(jwt.SigningMethodHS256)

	stringToken, err := token.SignedString(secret)

	if err != nil {
		log.Fatal(err)
	}

	return stringToken
}

// Функция для получения токена
func GetToken(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")

	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	var t token

	// создание токена
	t.T = createToken([]byte("SupaDupa"))

	jsonToken, err := json.MarshalIndent(t, "", "  ")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(jsonToken)
}

// Функция для валидации токена
func ValidateToken(token string) error {
	// парсинг токена
	t, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte("SupaDupa"), nil
	})

	if err != nil {
		return err
	}

	if t.Valid {
		return nil
	}

	return fmt.Errorf("unathorized")
}

// Функция для извлечения токена из запроса
func ExtractTokenFromRequest(r *http.Request) (string, error) {
	// получение значения по ключу "Authorization"
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("authorization header is missing")
	}

	// проверка на содержание "Bearer "
	parts := strings.SplitN(authHeader, " ", 2)
	if !(len(parts) == 2 && parts[0] == "Bearer") {
		return "", errors.New("invalid authorization header format")
	}

	tokenString := parts[1]
	return tokenString, nil
}
