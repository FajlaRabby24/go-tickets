package user

import (
	"gotickets/internel/auth"
	"gotickets/internel/config"
	middlewares "gotickets/internel/middleware"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

func RegisterRoutes(e *echo.Echo, db *gorm.DB, cfg *config.Config) {
	userRepository := NewRepository(db)
	jwtService := auth.NewJwtService(cfg.JwtSecret) // you can pass a custom secret key here
	userService := NewService(userRepository, jwtService)
	userHandler := NewHandler(userService)

	api := e.Group("/api/v1/auth")

	api.POST("/register", userHandler.CreateUser)
	api.POST("/login", userHandler.LoginUser)
	api.GET("/me", userHandler.GetMe, middlewares.AuthMiddleware(jwtService))
}
