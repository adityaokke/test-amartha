package db

import (
	"context"

	"github.com/adityaokke/test-amartha/internal/entity"
)

type InvestorRepository interface {
	Create(ctx context.Context, item *entity.Investor) (err error)
	Update(ctx context.Context, item *entity.Investor) (err error)
	Delete(ctx context.Context, item *entity.Investor) (err error)

	Investors(ctx context.Context, filter entity.InvestorsInput) (result []entity.Investor, err error)
	CountInvestors(ctx context.Context, filter entity.InvestorsInput) (result int64, err error)
	Investor(ctx context.Context, filter entity.InvestorInput) (result entity.Investor, err error)
}
