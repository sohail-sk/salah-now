package handlers

import (
	"context"
	"net/http"

	"salah-now/db"
	"salah-now/models"
	"salah-now/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func Signup(c *gin.Context) {
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	user.ID = uuid.New()

	query := `INSERT INTO users (id, name, email, password)
	          VALUES ($1, $2, $3, $4)`

	_, err := db.DB.Exec(context.Background(), query,
		user.ID, user.Name, user.Email, string(hashedPassword),
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User already exists"})
		return
	}

	token, _ := utils.GenerateToken(user.ID.String())

	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}

func Login(c *gin.Context) {
	var input models.User
	var dbUser models.User

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := `SELECT id, password FROM users WHERE email=$1`

	err := db.DB.QueryRow(context.Background(), query, input.Email).
		Scan(&dbUser.ID, &dbUser.Password)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(input.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	token, _ := utils.GenerateToken(dbUser.ID.String())

	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}
