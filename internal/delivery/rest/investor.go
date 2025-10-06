package rest

import (
	"net/http"

	"github.com/adityaokke/test-amartha/internal/entity"
	"github.com/adityaokke/test-amartha/internal/service"
	"github.com/labstack/echo/v4"
)

type InvestorHandler struct {
	investorService service.InvestorService
}

func NewInvestorHandler(
	investorService service.InvestorService,
) InvestorHandler {
	return InvestorHandler{
		investorService: investorService,
	}
}

func (d InvestorHandler) AddInvestor(c echo.Context) error {
	var form entity.AddInvestorInput
	if err := c.Bind(&form); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "Invalid JSON",
		})
	}

	result, err := d.investorService.AddInvestor(c.Request().Context(), form)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": map[string]interface{}{
			"investor": result,
		},
	})
}

func (d InvestorHandler) GetInvestors(c echo.Context) error {
	var filter entity.InvestorsInput
	result, err := d.investorService.Investors(c.Request().Context(), filter)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": map[string]interface{}{
			"investors": result,
		},
	})
}
