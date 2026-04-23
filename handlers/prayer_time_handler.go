package handlers

import (
	"context"
	"net/http"
	"strconv"

	"salah-now/db"
	"salah-now/services"

	"github.com/gin-gonic/gin"
)

func GetPrayerTimes(c *gin.Context) {
	lat, _ := strconv.ParseFloat(c.Query("lat"), 64)
	lng, _ := strconv.ParseFloat(c.Query("lng"), 64)

	data, err := services.GetPrayerTimesCached(lat, lng)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch prayer times",
		})
		return
	}

	current, next := services.GetCurrentAndNext(data)

	c.JSON(http.StatusOK, gin.H{
		"prayer_times": data,
		"current":      current,
		"next":         next,
	})
}
func GetPrayerTimesAuto(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var lat, lng float64

	err := db.DB.QueryRow(context.Background(), `
		SELECT latitude, longitude
		FROM user_locations
		WHERE user_id=$1
	`, userID).Scan(&lat, &lng)

	if err != nil {
		c.JSON(400, gin.H{"error": "location not set"})
		return
	}

	data, _ := services.GetPrayerTimesCached(lat, lng)

	c.JSON(200, data)
}
