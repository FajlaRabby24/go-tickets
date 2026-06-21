package main

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string `json:"name" validate:"required,min=3" gorm:"type:varchar(100);not null"`
	Email    string `json:"email" validate:"required,email" gorm:"type:varchar(100);uniqueIndex;not null"`
	Password string `json:"password" validate:"required,min=6" gorm:"type:varchar(100);not null"`
}

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i any) error {
	if err := cv.validator.Struct(i); err != nil {
		return echo.ErrBadRequest.Wrap(err)
	}

	return nil
}

func main() {
	dsn := "host=localhost user=postgres password=12345 dbname=gotickets port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		TranslateError: true,
	})

	if err != nil {
		panic("failed to connect database")
	} else {
		println("connected to database")
	}

	db.AutoMigrate(&User{})

	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}

	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())

	e.GET("/", func(c *echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"message": "Hello, World!"})
	})

	e.POST("/users", func(c *echo.Context) error {
		newUser := new(User)

		// binding the user data
		if err := c.Bind(newUser); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]any{"error": err.Error()})
		}

		// validating the user data
		if err := c.Validate(newUser); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]any{"error": err.Error()})
		}

		// save to db
		result := db.Create(newUser)
		if result.Error != nil {
			return c.JSON(http.StatusInternalServerError, map[string]any{"error": result.Error.Error()})
		}
		return c.JSON(http.StatusCreated, newUser)
	})

	if err := e.Start(":5000"); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}
