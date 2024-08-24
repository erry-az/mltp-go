package handler

import (
	"github.com/erry-az/mltp-go/internal/server/middleware"
	"github.com/erry-az/mltp-go/internal/service"
	"golang.org/x/sync/singleflight"
)

type Handler struct {
	service *service.Service
	jwtMw   *middleware.Jwt
	sf      singleflight.Group
}

func NewHandler(service *service.Service, jwtMw *middleware.Jwt) *Handler {
	return &Handler{service: service, jwtMw: jwtMw}
}
