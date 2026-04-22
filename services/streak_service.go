package services

import (
	"context"
	"time"

	"salah-now/db"
)

func UpdateStreak(userID string, date time.Time) error {
	ctx := context.Background()

	var count int

	// Check if all 5 prayers completed
	err := db.DB.QueryRow(ctx, `
		SELECT COUNT(DISTINCT prayer_id)
		FROM user_prayers
		WHERE user_id=$1 AND DATE(prayed_at)=$2
	`, userID, date.Format("2006-01-02")).Scan(&count)

	if err != nil || count < 5 {
		return nil // Not a full day, no streak update
	}

	var current, longest int
	var lastDate *time.Time

	err = db.DB.QueryRow(ctx, `
		SELECT current_streak, longest_streak, last_completed_date
		FROM user_streaks WHERE user_id=$1
	`, userID).Scan(&current, &longest, &lastDate)

	// First time
	if err != nil {
		_, err = db.DB.Exec(ctx, `
			INSERT INTO user_streaks (user_id, current_streak, longest_streak, last_completed_date)
			VALUES ($1, 1, 1, $2)
		`, userID, date)

		return err
	}

	yesterday := date.AddDate(0, 0, -1)

	if lastDate != nil && lastDate.Format("2006-01-02") == yesterday.Format("2006-01-02") {
		current++
	} else if lastDate != nil && lastDate.Format("2006-01-02") == date.Format("2006-01-02") {
		return nil // already counted
	} else {
		current = 1
	}

	if current > longest {
		longest = current
	}

	_, err = db.DB.Exec(ctx, `
		UPDATE user_streaks
		SET current_streak=$1,
		    longest_streak=$2,
		    last_completed_date=$3,
		    updated_at=NOW()
		WHERE user_id=$4
	`, current, longest, date, userID)

	return err
}
