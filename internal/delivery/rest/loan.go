package rest

import (
	"net/http"
	"strconv"

	"github.com/adityaokke/test-amartha/internal/entity"
	"github.com/adityaokke/test-amartha/internal/service"
	"github.com/labstack/echo/v4"
)

type LoanHandler struct {
	loanService service.LoanService
}

func NewLoanHandler(
	loanService service.LoanService,
) LoanHandler {
	return LoanHandler{
		loanService: loanService,
	}
}

func (d LoanHandler) ProposeLoan(c echo.Context) error {
	var form entity.ProposeLoanInput
	if err := c.Bind(&form); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "Invalid JSON",
		})
	}

	result, err := d.loanService.ProposeLoan(c.Request().Context(), form)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": map[string]interface{}{
			"loan": result,
		},
	})
}

func (d LoanHandler) GetLoans(c echo.Context) error {
	UserID := c.Param("id")

	var userID *int
	if UserID != "" {
		parsedUserID, err := strconv.Atoi(UserID)
		if err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{
				"error": "Invalid userId",
			})
		}
		userID = &parsedUserID
	}
	result, err := d.loanService.Loans(c.Request().Context(), entity.LoansInput{
		UserID: userID,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": map[string]interface{}{
			"loans": result,
		},
	})
}

func (d LoanHandler) GetLoan(c echo.Context) error {
	id := c.Param("id")
	parsedID, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "Invalid id",
		})
	}
	result, err := d.loanService.Loan(c.Request().Context(), entity.LoanInput{
		ID: &parsedID,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": map[string]interface{}{
			"loan": result,
		},
	})
}
