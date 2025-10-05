package rest

import (
	echo "github.com/labstack/echo/v4"
)

func Router(
	e *echo.Echo,
	loanHandler LoanHandler,
	loanInvestmentHandler LoanInvestmentHandler,
) {
	e.POST("/loans", loanHandler.ProposeLoan)
	e.GET("loans", loanHandler.GetLoans)
	e.GET("/loans/:id", loanHandler.GetLoan)
	e.PATCH("/loans/:id", loanHandler.PatchLoan)
	e.POST("/loans/:id/investments", loanHandler.InvestLoan)
	e.GET("/loans/:id/investments", loanInvestmentHandler.GetLoanInvestments)
}
