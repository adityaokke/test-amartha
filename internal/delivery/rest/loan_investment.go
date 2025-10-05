package rest

import (
	"net/http"
	"strconv"

	"github.com/adityaokke/test-amartha/internal/entity"
	"github.com/adityaokke/test-amartha/internal/service"
	"github.com/labstack/echo/v4"
)

type LoanInvestmentHandler struct {
	loanInvestmentService service.LoanInvestmentService
}

func NewLoanInvestmentHandler(
	loanInvestmentService service.LoanInvestmentService,
) LoanInvestmentHandler {
	return LoanInvestmentHandler{
		loanInvestmentService: loanInvestmentService,
	}
}

func (d LoanInvestmentHandler) GetLoanInvestments(c echo.Context) error {
	id := c.Param("id")
	parsedID, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "Invalid id",
		})
	}

	result, err := d.loanInvestmentService.LoanInvestments(c.Request().Context(), entity.LoanInvestmentsInput{
		LoanID: &parsedID,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": map[string]interface{}{
			"loan_investments": result,
		},
	})
}
