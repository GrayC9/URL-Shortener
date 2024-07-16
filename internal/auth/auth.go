package auth

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"
	"url_shortener/internal/handlers"
	"url_shortener/internal/storage"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

var JWTSecretKey []byte

type Claims struct {
	UserID int `json:"user_id"`
	jwt.StandardClaims
}

func SignUp(db storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		username := r.FormValue("username")
		password := r.FormValue("hash_password")

		id, err := db.CreateUser(username, Hashing(password))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		handlers.WriteJSON(w, r, http.StatusOK, map[string]interface{}{
			"id": id,
		})
	}
}

func Login(db storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		username := r.FormValue("username")
		password := r.FormValue("hash_password")

		id, err := db.EnterUser(username, password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		token, err := GenerateToken(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		handlers.WriteJSON(w, r, http.StatusOK, map[string]interface{}{
			"token": token,
		})
	}
}

func GenerateToken(userID int) (string, error) {
	claims := &Claims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(JWTSecretKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return signedToken, nil
}

func ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return JWTSecretKey, nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		claims, err := ParseToken(tokenString)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		r = r.WithContext(context.WithValue(r.Context(), "userID", claims.UserID))

		next.ServeHTTP(w, r)
	})
}

func GetContextUserID(r *http.Request) (int, error) {
	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		return 0, errors.New("user ID not found in context")
	}
	return userID, nil
}

func Hashing(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return ""
	}
	return string(hash)
}
