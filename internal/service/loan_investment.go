package service

import (
	"context"

	"github.com/adityaokke/test-amartha/internal/entity"
	"github.com/adityaokke/test-amartha/internal/repository/db"
)

type LoanInvestmentService interface {
	LoanInvestments(ctx context.Context, filter entity.LoanInvestmentsInput) (result []entity.LoanInvestment, err error)
	CountLoanInvestments(ctx context.Context, filter entity.LoanInvestmentsInput) (result int64, err error)
	LoanInvestment(ctx context.Context, filter entity.LoanInvestmentInput) (result entity.LoanInvestment, err error)
}

func (s *loanInvestmentService) LoanInvestments(ctx context.Context, filter entity.LoanInvestmentsInput) (result []entity.LoanInvestment, err error) {
	result, err = s.loanInvestmentRepo.LoanInvestments(ctx, filter)
	if err != nil {
		return
	}
	return
}

func (s *loanInvestmentService) CountLoanInvestments(ctx context.Context, filter entity.LoanInvestmentsInput) (result int64, err error) {
	result, err = s.loanInvestmentRepo.CountLoanInvestments(ctx, filter)
	if err != nil {
		return
	}
	return
}

func (s *loanInvestmentService) LoanInvestment(ctx context.Context, filter entity.LoanInvestmentInput) (result entity.LoanInvestment, err error) {
	result, err = s.loanInvestmentRepo.LoanInvestment(ctx, filter)
	if err != nil {
		return
	}
	return
}

type loanInvestmentService struct {
	loanInvestmentRepo db.LoanInvestmentRepository
}

type InitiatorLoanInvestment func(s *loanInvestmentService) *loanInvestmentService

func NewLoanInvestmentService() InitiatorLoanInvestment {
	return func(s *loanInvestmentService) *loanInvestmentService {
		return s
	}
}

func (i InitiatorLoanInvestment) SetRepository(loanInvestmentRepository db.LoanInvestmentRepository) InitiatorLoanInvestment {
	return func(s *loanInvestmentService) *loanInvestmentService {
		i(s).loanInvestmentRepo = loanInvestmentRepository
		return s
	}
}

func (i InitiatorLoanInvestment) Build() LoanInvestmentService {
	return i(&loanInvestmentService{})
}
