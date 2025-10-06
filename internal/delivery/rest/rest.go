package rest

import (
	"fmt"

	"github.com/adityaokke/test-amartha/internal/entity"
	echo "github.com/labstack/echo/v4"
)

func Router(
	e *echo.Echo,
	loanHandler LoanHandler,
	loanInvestmentHandler LoanInvestmentHandler,
	fileHandler FileHandler,
	InvestorHandler InvestorHandler,
) {
	e.POST("/files", fileHandler.Upload)
	e.POST("/loans", loanHandler.ProposeLoan)
	e.GET("loans", loanHandler.GetLoans)
	e.GET("/loans/:id", loanHandler.GetLoan)
	e.PATCH("/loans/:id", loanHandler.PatchLoan)
	e.POST("/loans/:id/investments", loanHandler.InvestLoan)
	e.GET("/loans/:id/investments", loanInvestmentHandler.GetLoanInvestments)
	e.Static(fmt.Sprintf("/%s", entity.PublicUploadPath), entity.LocalUploadPath)
	e.Static(fmt.Sprintf("/%s", entity.PublicAggrementLetterPath), entity.LocalAggrementLetterPath)
	e.POST("/investors", InvestorHandler.AddInvestor)
	e.GET("/investors", InvestorHandler.GetInvestors)
	e.GET("/loans/:id/agreement/contents", loanHandler.GetAgreementLetter)
}
