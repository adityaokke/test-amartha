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
	input := entity.LoansInput{}

	userID := c.QueryParam("userId")
	if userID != "" {
		userIDParsed, err := strconv.Atoi(userID)
		if err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{
				"error": "Invalid userId",
			})
		}
		input.UserID = &userIDParsed
	}

	status := entity.LoanStatus(c.QueryParam("status"))
	if status != "" {
		if !status.IsValid() {
			return c.JSON(http.StatusBadRequest, echo.Map{
				"error": "Invalid status",
			})
		}
		input.Status = &status
	}

	result, err := d.loanService.Loans(c.Request().Context(), input)
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

func (d LoanHandler) PatchLoan(c echo.Context) error {
	id := c.Param("id")
	parsedID, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "Invalid id",
		})
	}

	var form entity.PatchLoanInput
	if err := c.Bind(&form); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "Invalid JSON",
		})
	}
	form.ID = parsedID
	var result entity.Loan
	switch form.Status {
	case entity.LoanStatusApproved:
		result, err = d.loanService.ApproveLoan(c.Request().Context(), entity.ApproveLoanInput{
			ID:            form.ID,
			EmployeeID:    form.EmployeeID,
			PhotoProofURL: form.PhotoProofURL,
		})
	case entity.LoanStatusDisbursed:
		result, err = d.loanService.DisburseLoan(c.Request().Context(), entity.DisburseLoanInput{
			ID:                             form.ID,
			LoanAgreementLetterURL:         form.LoanAgreementLetterURL,
			DisbursedByEmployeeID:          form.DisbursedByEmployeeID,
			AgreementCollectedByEmployeeID: form.AgreementCollectedByEmployeeID,
		})
	default:
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "Invalid status",
		})
	}
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": map[string]interface{}{
			"loan": result,
		},
	})
}

func (d LoanHandler) InvestLoan(c echo.Context) error {
	id := c.Param("id")
	parsedID, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "Invalid id",
		})
	}
	var form entity.InvestLoanInput
	if err := c.Bind(&form); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "Invalid JSON",
		})
	}
	form.LoanID = parsedID
	result, err := d.loanService.InvestLoan(c.Request().Context(), form)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": map[string]interface{}{
			"loan_investment": result,
		},
	})
}

func (d LoanHandler) GetAgreementLetter(c echo.Context) error {
	id := c.Param("id")
	parsedID, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "Invalid id",
		})
	}

	variant := c.QueryParam("variant")
	result := ""
	switch variant {
	case entity.AggrementLetterVariantDraft:
		result, err = d.loanService.GetDraftLoanAgreementLetter(c.Request().Context(), parsedID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
		}
	case entity.AggrementLetterVariantSign:
		result, err = d.loanService.GetSignedLoanAgreementLetter(c.Request().Context(), parsedID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
		}
	default:
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "Invalid variant",
		})
	}

	return c.Redirect(http.StatusFound, result)
}

func (d LoanHandler) GetLoanQuotes(c echo.Context) error {
	id := c.Param("id")
	parsedID, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "Invalid id",
		})
	}

	result, err := d.loanService.GetLoanQuote(c.Request().Context(), parsedID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": map[string]interface{}{
			"loan_quotes": result,
		},
	})
}
