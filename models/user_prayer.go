package models

import (
	"time"

	"github.com/google/uuid"
)

type UserPrayer struct {
	ID       uuid.UUID `json:"id"`
	UserID   uuid.UUID `json:"user_id"`
	PrayerID int       `json:"prayer_id"`
	MosqueID uuid.UUID `json:"mosque_id"`
	PrayedAt time.Time `json:"prayed_at"`
	Mode     string    `json:"mode"` // JAMAAT / INDIVIDUAL
}
