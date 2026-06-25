package booking

import (
	"gotickets/internel/auth"
	"gotickets/internel/config"
	"gotickets/internel/event"
	middlewares "gotickets/internel/middleware"

	"github.com/labstack/echo/v5"
	"gorm.io/gorm"
)

func RegisterRoutes(e *echo.Echo, db *gorm.DB, cfg *config.Config) {
	bookingRepo := NewRepository(db)
	eventRepo := event.NewRepository(db)

	svc := NewService(bookingRepo, eventRepo)
	handler := NewHandler(svc)

	jwtService := auth.NewJwtService(cfg.JwtSecret)

	api := e.Group("/api/v1/bookings", middlewares.AuthMiddleware(jwtService))

	api.POST("", handler.CreateBooking)
	api.GET("/me", handler.GetMyBookings)
}
