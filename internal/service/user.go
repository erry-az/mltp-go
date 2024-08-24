package service

import (
	"context"
	"database/sql"
	"errors"
	"net/http"

	"github.com/erry-az/mltp-go/db/query"
	"github.com/labstack/echo/v4"
)

func (s *Service) RegisterUser(ctx context.Context, username, fullname string) (query.User, error) {
	user, err := s.queries.GetUserByUsername(ctx, username)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return query.User{}, err
	}

	if user.ID != 0 {
		return query.User{}, echo.NewHTTPError(http.StatusConflict, "username already exists")
	}

	return s.queries.CreateUser(ctx, query.CreateUserParams{
		Username: username,
		Fullname: fullname,
	})
}

func (s *Service) GetUser(ctx context.Context, username string) (query.User, error) {
	return s.queries.GetUserByUsername(ctx, username)
}
