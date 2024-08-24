package handler

import (
	"fmt"
	"net/http"

	"github.com/erry-az/mltp-go/internal/server/middleware"
	"github.com/labstack/echo/v4"
)

type TopUpRequest struct {
	Amount int64 `json:"amount" validate:"required"`
}

func (h *Handler) TopUp(c echo.Context) error {
	claim, err := middleware.JwtParseClaims(c)
	if err != nil {
		return err
	}

	sub, err := claim.GetSubject()
	if err != nil {
		return err
	}

	var req TopUpRequest

	err = c.Bind(&req)
	if err != nil {
		return err
	}

	err = c.Validate(req)
	if err != nil {
		return err
	}

	if req.Amount <= 0 || req.Amount >= 10000000 {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid topup amount")
	}

	err = h.service.TopUp(c.Request().Context(), sub, uint64(req.Amount))
	if err != nil {
		return err
	}

	return echo.NewHTTPError(http.StatusNoContent, "topup successful")
}

type TransferRequest struct {
	ToUsername string `json:"to_username" validate:"required"`
	Amount     int64  `json:"amount" validate:"required"`
}

func (h *Handler) Transfer(c echo.Context) error {
	claim, err := middleware.JwtParseClaims(c)
	if err != nil {
		return err
	}

	sub, err := claim.GetSubject()
	if err != nil {
		return err
	}

	var req TransferRequest

	err = c.Bind(&req)
	if err != nil {
		return err
	}

	err = c.Validate(req)
	if err != nil {
		return err
	}

	if req.Amount <= 0 || req.Amount >= 10000000 {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid topup amount")
	}

	if req.ToUsername == sub {
		return echo.NewHTTPError(http.StatusBadRequest, "cannot transfer to yourself")
	}

	err = h.service.Transfer(c.Request().Context(), sub, req.ToUsername, uint64(req.Amount))
	if err != nil {
		return err
	}

	return echo.NewHTTPError(http.StatusNoContent, "transfer successful")
}

type TopUserData struct {
	Username        string `json:"username"`
	TransactedValue int64  `json:"transacted_value"`
}

func (h *Handler) TopUsers(c echo.Context) error {
	_, err := middleware.JwtParseClaims(c)
	if err != nil {
		return err
	}

	response, err, _ := h.sf.Do("top:users", func() (interface{}, error) {
		topSumTransactions, err := h.service.TopSummaryTransactions()
		if err != nil {
			return nil, err
		}

		response := make([]TopUserData, 0, len(topSumTransactions))
		for _, trx := range topSumTransactions {
			response = append(response, TopUserData{
				Username:        trx.Username,
				TransactedValue: trx.TransactedValue,
			})
		}

		return response, nil
	})
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, response)
}

type TopTransactionsData struct {
	Username string `json:"username"`
	Amount   int64  `json:"amount"`
}

func (h *Handler) TopTransactions(c echo.Context) error {
	claim, err := middleware.JwtParseClaims(c)
	if err != nil {
		return err
	}

	sub, err := claim.GetSubject()
	if err != nil {
		return err
	}

	response, err, _ := h.sf.Do(fmt.Sprintf("top:transactions:%s", sub), func() (interface{}, error) {
		topTransactions, err := h.service.TopTransactionByUsername(c.Request().Context(), sub)
		if err != nil {
			return nil, err
		}

		response := make([]TopTransactionsData, 0, len(topTransactions))
		for _, trx := range topTransactions {
			response = append(response, TopTransactionsData{
				Username: trx.Username,
				Amount:   trx.AmountValue,
			})
		}

		return response, nil
	})
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, response)
}
