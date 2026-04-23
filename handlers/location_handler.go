package handlers

import (
	"context"
	"net/http"

	"salah-now/db"

	"github.com/gin-gonic/gin"
)

func UpdateLocation(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var input struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := `
	INSERT INTO user_locations (user_id, latitude, longitude, updated_at)
	VALUES ($1, $2, $3, NOW())
	ON CONFLICT (user_id)
	DO UPDATE SET latitude=$2, longitude=$3, updated_at=NOW();
	`

	_, err := db.DB.Exec(context.Background(), query,
		userID, input.Latitude, input.Longitude,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "location updated"})
}
