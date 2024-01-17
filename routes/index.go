package routes

import "github.com/gin-gonic/gin"

func AllRoute() {
	route := gin.New()
	route.Use(gin.Logger())

	HitRoute(route)

	route.Run()
}
