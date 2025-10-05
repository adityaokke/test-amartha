package rest

import (
	echo "github.com/labstack/echo/v4"
)

func Router(
	e *echo.Echo,
	loanHandler LoanHandler,
) {
	e.POST("/loans", loanHandler.ProposeLoan)
	e.GET("/users/:id/loans", loanHandler.GetLoans)
	e.GET("/loans/:id", loanHandler.GetLoan)
}
