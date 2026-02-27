package pkg

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"user-service/internal/model"
)

func GenerateToken(user model.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}