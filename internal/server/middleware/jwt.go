package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/erry-az/mltp-go/internal/config"
	"github.com/golang-jwt/jwt/v5"
	echoJwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

var ErrSigningMethodNotFound = fmt.Errorf("signing method not found")

const JwtCtxKey = "jwt"

type Jwt struct {
	signMethod jwt.SigningMethod
	cfg        config.JwtConfig
}

type CustomClaims[T any] struct {
	Data T `json:"data,omitempty"`
	jwt.RegisteredClaims
}

func NewJWT(config config.JwtConfig) (*Jwt, error) {
	signMethod := jwt.GetSigningMethod(config.SigningMethod)
	if signMethod == nil {
		return nil, ErrSigningMethodNotFound
	}

	return &Jwt{
		signMethod: signMethod,
		cfg:        config,
	}, nil
}

func (j *Jwt) NewClaims(id string, username string, custom any) (string, error) {
	claims := CustomClaims[any]{
		Data: custom,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    j.cfg.Issuer,
			Subject:   username,
			ID:        id,
		},
	}

	return jwt.NewWithClaims(j.signMethod, claims).SignedString([]byte(j.cfg.SigningKey))
}

func (j *Jwt) Middleware() echo.MiddlewareFunc {
	return echoJwt.WithConfig(echoJwt.Config{
		ContextKey:    JwtCtxKey,
		SigningKey:    []byte(j.cfg.SigningKey),
		SigningMethod: j.cfg.SigningMethod,
	})
}

func JwtParseClaims(c echo.Context) (jwt.MapClaims, error) {
	token := c.Get(JwtCtxKey).(*jwt.Token)
	if token == nil {
		return nil, echo.NewHTTPError(http.StatusUnauthorized, "unauthorized user, token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || claims == nil {
		return nil, echo.NewHTTPError(http.StatusUnauthorized, "unauthorized user, claims")
	}

	return claims, nil
}
