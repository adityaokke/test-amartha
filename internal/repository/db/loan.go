package db

import (
	"context"

	"github.com/adityaokke/test-amartha/internal/entity"
)

type LoanRepository interface {
	Create(ctx context.Context, item *entity.Loan) (err error)
	Update(ctx context.Context, item *entity.Loan) (err error)
	Delete(ctx context.Context, item *entity.Loan) (err error)

	Loans(ctx context.Context, filter entity.LoansInput) (result []entity.Loan, err error)
	CountLoans(ctx context.Context, filter entity.LoansInput) (result int64, err error)
	Loan(ctx context.Context, filter entity.LoanInput) (result entity.Loan, err error)
}
