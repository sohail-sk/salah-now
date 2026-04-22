package handlers

import (
	"context"
	"net/http"

	"salah-now/db"

	"github.com/gin-gonic/gin"
)

func getWeeklyStats(userID interface{}) (map[string]interface{}, error) {
	query := `
	SELECT 
	    COUNT(*) as total,
	    COUNT(*) FILTER (WHERE mode = 'JAMAAT') as jamaat,
	    COUNT(*) FILTER (WHERE mode = 'INDIVIDUAL') as individual
	FROM user_prayers
	WHERE user_id = $1
	AND prayed_at >= NOW() - INTERVAL '7 days';
	`

	var total, jamaat, individual int

	err := db.DB.QueryRow(context.Background(), query, userID).
		Scan(&total, &jamaat, &individual)

	if err != nil {
		return nil, err
	}

	// Full days (5 prayers in a day)
	fullDaysQuery := `
	SELECT COUNT(*) FROM (
	    SELECT DATE(prayed_at)
	    FROM user_prayers
	    WHERE user_id=$1
	    AND prayed_at >= NOW() - INTERVAL '7 days'
	    GROUP BY DATE(prayed_at)
	    HAVING COUNT(DISTINCT prayer_id) = 5
	) t;
	`

	var fullDays int
	db.DB.QueryRow(context.Background(), fullDaysQuery, userID).Scan(&fullDays)

	return map[string]interface{}{
		"total_prayers": total,
		"jamaat":        jamaat,
		"individual":    individual,
		"full_days":     fullDays,
	}, nil
}

func getMonthlyStats(userID interface{}) (map[string]interface{}, error) {
	query := `
	SELECT COUNT(*) 
	FROM user_prayers
	WHERE user_id = $1
	AND date_trunc('month', prayed_at) = date_trunc('month', CURRENT_DATE);
	`

	var total int
	err := db.DB.QueryRow(context.Background(), query, userID).Scan(&total)
	if err != nil {
		return nil, err
	}

	// Total possible prayers in month till today
	totalPossibleQuery := `
	SELECT (EXTRACT(DAY FROM CURRENT_DATE)::int * 5);
	`

	var totalPossible int
	db.DB.QueryRow(context.Background(), totalPossibleQuery).Scan(&totalPossible)

	percentage := 0
	if totalPossible > 0 {
		percentage = (total * 100) / totalPossible
	}

	// Full days
	fullDaysQuery := `
	SELECT COUNT(*) FROM (
	    SELECT DATE(prayed_at)
	    FROM user_prayers
	    WHERE user_id=$1
	    AND date_trunc('month', prayed_at) = date_trunc('month', CURRENT_DATE)
	    GROUP BY DATE(prayed_at)
	    HAVING COUNT(DISTINCT prayer_id) = 5
	) t;
	`

	var fullDays int
	db.DB.QueryRow(context.Background(), fullDaysQuery, userID).Scan(&fullDays)

	return map[string]interface{}{
		"total_prayers":         total,
		"completion_percentage": percentage,
		"full_days":             fullDays,
	}, nil
}

func GetDashboard(c *gin.Context) {
	userID, _ := c.Get("user_id")

	weekly, err := getWeeklyStats(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	monthly, err := getMonthlyStats(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Streak
	var current, longest int
	db.DB.QueryRow(context.Background(), `
		SELECT current_streak, longest_streak
		FROM user_streaks
		WHERE user_id=$1
	`, userID).Scan(&current, &longest)

	c.JSON(http.StatusOK, gin.H{
		"weekly":  weekly,
		"monthly": monthly,
		"streak": gin.H{
			"current": current,
			"longest": longest,
		},
	})
}
