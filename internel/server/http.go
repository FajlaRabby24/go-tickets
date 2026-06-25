package server

import (
	"fmt"
	"gotickets/internel/booking"
	"gotickets/internel/config"
	"gotickets/internel/event"
	"gotickets/internel/user"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"gorm.io/gorm"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i any) error {
	if err := cv.validator.Struct(i); err != nil {
		return echo.ErrBadRequest.Wrap(err)
	}

	return nil
}

func Start(db *gorm.DB, cfg *config.Config) {
	db.AutoMigrate(&user.User{}, &event.Event{}, &booking.Booking{})

	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	e.Use(middleware.RequestLogger())

	e.GET("/", func(c *echo.Context) error {
		return c.String(http.StatusOK, "Gotickets server is running")
	})

	// user route register
	user.RegisterRoutes(e, db, cfg)
	event.RegisterRoutes(e, db)
	booking.RegisterRoutes(e, db, cfg)

	port := fmt.Sprintf(":%s", cfg.Port)
	if err := e.Start(port); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}
