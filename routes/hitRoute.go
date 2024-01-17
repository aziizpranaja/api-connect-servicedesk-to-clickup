package routes

import (
	"api-service/controllers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func HitRoute(incomingRoutes *gin.Engine) {
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true                                             // Mengizinkan semua asal permintaan
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"} // Metode HTTP yang diizinkan
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"} // Header yang diizinkan
	incomingRoutes.Use(cors.New(config))
	user := incomingRoutes.Group("/v1/hit")
	user.POST("/", controllers.HitApi)
}
