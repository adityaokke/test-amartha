package db

import (
	"context"

	"github.com/adityaokke/test-amartha/internal/entity"
)

type LoanInvestmentRepository interface {
	InvestLoan(ctx context.Context, item *entity.LoanInvestment) (err error)

	LoanInvestments(ctx context.Context, filter entity.LoanInvestmentsInput) (result []entity.LoanInvestment, err error)
	CountLoanInvestments(ctx context.Context, filter entity.LoanInvestmentsInput) (result int64, err error)
	LoanInvestment(ctx context.Context, filter entity.LoanInvestmentInput) (result entity.LoanInvestment, err error)
}
