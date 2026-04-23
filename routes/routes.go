package routes

import (
	"salah-now/handlers"
	"salah-now/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Auth
	r.POST("/signup", handlers.Signup)
	r.POST("/login", handlers.Login)

	// Protected routes
	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.POST("/mosques", handlers.CreateMosque)
		protected.GET("/mosques/nearby", handlers.GetNearbyMosques)
	}

	protected.Use(middleware.AuthMiddleware())
	{
		protected.POST("/mosques", handlers.CreateMosque)
		protected.GET("/mosques/nearby", handlers.GetNearbyMosques)

		// Prayer APIs
		protected.POST("/prayers/mark", handlers.MarkPrayer)
		protected.GET("/prayers/today", handlers.GetTodayPrayers)
		protected.GET("/prayers/history", handlers.GetPrayerHistory)
	}
	protected.GET("/streak", handlers.GetStreak)
	protected.GET("/dashboard", handlers.GetDashboard)
	protected.GET("/prayer-times", handlers.GetPrayerTimes)
	protected.POST("/user/location", handlers.UpdateLocation)
	return r
}
