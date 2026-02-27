package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"user-service/config"
	"user-service/internal/model"
	"user-service/internal/repository"
	"user-service/pkg"
)

func main() {

	_ = godotenv.Load()

	db, err := config.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(&model.User{})

	repo := &repository.UserRepository{DB: db}

	r := gin.Default()

	r.POST("/login", func(c *gin.Context) {
		var req struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}

		c.ShouldBindJSON(&req)

		user, err := repo.FindByUsername(req.Username)
		if err != nil || user.Password != req.Password {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}

		token, _ := pkg.GenerateToken(*user)

		c.JSON(http.StatusOK, gin.H{"token": token})
	})

	r.Run(":" + os.Getenv("PORT"))
}