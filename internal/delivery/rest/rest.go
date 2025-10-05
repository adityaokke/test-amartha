package rest

import (
	echo "github.com/labstack/echo/v4"
)

func Router(
	e *echo.Echo,
	loanHandler LoanHandler,
) {
	e.POST("/loans", loanHandler.ProposeLoan)
	e.GET("loans", loanHandler.GetLoans)
	e.GET("/loans/:id", loanHandler.GetLoan)
	e.PATCH("/loans/:id", loanHandler.PatchLoan)
}
