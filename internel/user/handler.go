package user

import (
	"errors"
	httpresponse "gotickets/internel/httpResponse"
	"gotickets/internel/user/dto"
	"net/http"

	"github.com/labstack/echo/v5"
)

type handler struct {
	service *service
}

func NewHandler(service *service) *handler {
	return &handler{
		service: service,
	}
}

func (h *handler) CreateUser(c *echo.Context) error {
	var req dto.CreateRequest // input

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.Error{
			Code:    http.StatusBadRequest,
			Details: err.Error(),
			Message: "Invalide request payload"})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.Error{
			Code:    http.StatusBadRequest,
			Details: err.Error(),
			Message: "Validation failed"})
	}

	response, err := h.service.CreateUser(req)

	if err != nil {

		if errors.Is(err, ErrorAlreadyExist) {
			return c.JSON(http.StatusConflict, httpresponse.Error{
				Code:    http.StatusConflict,
				Message: "Failed to create user",
				Details: err.Error(),
			})
		}

		return c.JSON(http.StatusInternalServerError, httpresponse.Error{
			Code:    http.StatusInternalServerError,
			Details: err.Error(),
			Message: "Failed to create user"})
	}

	return c.JSON(http.StatusCreated, response)
}

func (h *handler) LoginUser(c *echo.Context) error {
	var req dto.LoginRequest // input

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.Error{
			Code:    http.StatusBadRequest,
			Details: err.Error(),
			Message: "Invalide request payload"})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.Error{
			Code:    http.StatusBadRequest,
			Details: err.Error(),
			Message: "Validation failed"})
	}

	response, err := h.service.LoginUser(req)

	if err != nil {
		if errors.Is(err, ErrorInvalideCredentials) {
			return c.JSON(http.StatusUnauthorized, httpresponse.Error{
				Code:    http.StatusUnauthorized,
				Message: "Cannot login user",
				Details: err.Error(),
			})
		}

		return c.JSON(http.StatusInternalServerError, httpresponse.Error{
			Code:    http.StatusInternalServerError,
			Details: err.Error(),
			Message: "Failed to login user"})
	}

	return c.JSON(http.StatusOK, response)
}

func (h *handler) GetMe(c *echo.Context) error {
	userID, ok := c.Get("user_id").(uint)
	if !ok {
		return c.JSON(http.StatusUnauthorized, httpresponse.Error{
			Code:    http.StatusUnauthorized,
			Message: "Cannot get user information",
			Details: "missing user id in context",
		})
	}

	email, _ := c.Get("user_email").(string)
	name, _ := c.Get("user_name").(string)

	return c.JSON(http.StatusOK, dto.Response{
		ID:    userID,
		Name:  name,
		Email: email,
	})
}
