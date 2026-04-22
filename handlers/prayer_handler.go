package handlers

import (
	"context"
	"net/http"
	"time"

	"salah-now/db"
	"salah-now/models"
	"salah-now/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func MarkPrayer(c *gin.Context) {
	var input models.UserPrayer

	userID, _ := c.Get("user_id")

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	input.ID = uuid.New()
	input.UserID, _ = uuid.Parse(userID.(string))

	if input.PrayedAt.IsZero() {
		input.PrayedAt = time.Now()
	}

	// Default mode
	if input.Mode == "" {
		input.Mode = "INDIVIDUAL"
	}

	query := `
	INSERT INTO user_prayers 
	(id, user_id, prayer_id, mosque_id, prayed_at, mode)
	VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := db.DB.Exec(context.Background(), query,
		input.ID, input.UserID, input.PrayerID,
		input.MosqueID, input.PrayedAt, input.Mode,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	services.UpdateStreak(userID.(string), input.PrayedAt)

	c.JSON(http.StatusOK, gin.H{
		"message": "Prayer marked",
	})
}

func GetTodayPrayers(c *gin.Context) {
	userID, _ := c.Get("user_id")

	query := `
	SELECT p.name, up.prayed_at,up.mode
	FROM prayers p
	LEFT JOIN user_prayers up
	ON p.id = up.prayer_id
	AND up.user_id = $1
	AND DATE(up.prayed_at) = CURRENT_DATE
	ORDER BY p.id;
	`

	rows, err := db.DB.Query(context.Background(), query, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var result []gin.H

	for rows.Next() {
		var name string
		var prayedAt *time.Time
		var mode *string
		rows.Scan(&name, &prayedAt)

		result = append(result, gin.H{
			"prayer":    name,
			"completed": prayedAt != nil,
			"time":      prayedAt,
			"mode":      mode,
		})
	}

	c.JSON(http.StatusOK, result)
}

func GetPrayerHistory(c *gin.Context) {
	userID, _ := c.Get("user_id")

	query := `
	SELECT p.name, up.prayed_at, up.mode,
	FROM user_prayers up
	JOIN prayers p ON p.id = up.prayer_id
	WHERE up.user_id = $1
	ORDER BY up.prayed_at DESC
	LIMIT 50;
	`

	rows, err := db.DB.Query(context.Background(), query, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var result []gin.H

	for rows.Next() {
		var name string
		var time time.Time
		var mode string
		rows.Scan(&name, &time, &mode)

		result = append(result, gin.H{
			"prayer": name,
			"time":   time,
			"mode":   mode,
		})
	}

	c.JSON(http.StatusOK, result)
}

func GetStreak(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var current, longest int

	err := db.DB.QueryRow(context.Background(), `
		SELECT current_streak, longest_streak
		FROM user_streaks
		WHERE user_id=$1
	`, userID).Scan(&current, &longest)

	if err != nil {
		c.JSON(200, gin.H{
			"current_streak": 0,
			"longest_streak": 0,
		})
		return
	}

	c.JSON(200, gin.H{
		"current_streak": current,
		"longest_streak": longest,
	})
}
