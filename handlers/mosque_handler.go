package handlers

import (
	"context"
	"net/http"
	"strconv"

	"salah-now/db"
	"salah-now/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CreateMosque(c *gin.Context) {
	var m models.Mosque

	if err := c.ShouldBindJSON(&m); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	m.ID = uuid.New()

	query := `
	INSERT INTO mosques (id, name, address, latitude, longitude, location)
	VALUES ($1, $2, $3, $4, $5, ST_SetSRID(ST_MakePoint($5, $4), 4326))
	`

	_, err := db.DB.Exec(context.Background(), query,
		m.ID, m.Name, m.Address, m.Latitude, m.Longitude,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, m)
}

func GetNearbyMosques(c *gin.Context) {
	lat, _ := strconv.ParseFloat(c.Query("lat"), 64)
	lng, _ := strconv.ParseFloat(c.Query("lng"), 64)

	query := `
	SELECT id, name, address, latitude, longitude,
	       ST_DistanceSphere(location, ST_MakePoint($1, $2)) AS distance
	FROM mosques
	ORDER BY distance
	LIMIT 10;
	`

	rows, err := db.DB.Query(context.Background(), query, lng, lat)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var result []gin.H

	for rows.Next() {
		var m models.Mosque
		var distance float64

		rows.Scan(&m.ID, &m.Name, &m.Address, &m.Latitude, &m.Longitude, &distance)

		result = append(result, gin.H{
			"id":       m.ID,
			"name":     m.Name,
			"distance": distance,
		})
	}

	c.JSON(http.StatusOK, result)
}