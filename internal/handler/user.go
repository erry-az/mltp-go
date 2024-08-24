package handler

import (
	"fmt"
	"net/http"

	"github.com/erry-az/mltp-go/internal/server/middleware"

	"github.com/labstack/echo/v4"
)

type CreateUserRequest struct {
	Name     string `json:"name" validate:"required,min=2,max=20"`
	Username string `json:"username" validate:"required,username"`
}

type CreateUserResponse struct {
	Token string `json:"token"`
}

func (h *Handler) CreateUser(c echo.Context) error {
	var req CreateUserRequest

	err := c.Bind(&req)
	if err != nil {
		return err
	}

	err = c.Validate(req)
	if err != nil {
		return err
	}

	user, err := h.service.RegisterUser(c.Request().Context(), req.Username, req.Name)
	if err != nil {
		return err
	}

	token, err := h.jwtMw.NewClaims(fmt.Sprintf("%d", user.ID), user.Username, user)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, CreateUserResponse{Token: token})
}

type BalanceReadResponse struct {
	Balance int64 `json:"balance"`
}

func (h *Handler) BalanceRead(c echo.Context) error {
	claim, err := middleware.JwtParseClaims(c)
	if err != nil {
		return err
	}

	sub, err := claim.GetSubject()
	if err != nil {
		return err
	}

	balance, err := h.service.GetUser(c.Request().Context(), sub)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, BalanceReadResponse{Balance: balance.Balance})
}
