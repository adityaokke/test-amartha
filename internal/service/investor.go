package service

import (
	"context"
	"errors"

	"github.com/adityaokke/test-amartha/internal/entity"
	"github.com/adityaokke/test-amartha/internal/repository/db"
)

type InvestorService interface {
	AddInvestor(ctx context.Context, input entity.AddInvestorInput) (result entity.Investor, err error)

	Investors(ctx context.Context, filter entity.InvestorsInput) (result []entity.Investor, err error)
	CountInvestors(ctx context.Context, filter entity.InvestorsInput) (result int64, err error)
	Investor(ctx context.Context, filter entity.InvestorInput) (result entity.Investor, err error)
}

func (s *investorService) AddInvestor(ctx context.Context, input entity.AddInvestorInput) (result entity.Investor, err error) {
	if input.Email == "" {
		err = errors.New("email is required")
		return
	}
	item := entity.Investor{
		Email: input.Email,
	}
	err = s.investorRepo.Create(ctx, &item)
	if err != nil {
		return
	}
	result = item
	return
}

func (s *investorService) Investors(ctx context.Context, filter entity.InvestorsInput) (result []entity.Investor, err error) {
	result, err = s.investorRepo.Investors(ctx, filter)
	if err != nil {
		return
	}
	return
}

func (s *investorService) CountInvestors(ctx context.Context, filter entity.InvestorsInput) (result int64, err error) {
	result, err = s.investorRepo.CountInvestors(ctx, filter)
	if err != nil {
		return
	}
	return
}

func (s *investorService) Investor(ctx context.Context, filter entity.InvestorInput) (result entity.Investor, err error) {
	result, err = s.investorRepo.Investor(ctx, filter)
	if err != nil {
		return
	}
	return
}

type investorService struct {
	investorRepo db.InvestorRepository
}

type InitiatorInvestor func(s *investorService) *investorService

func NewInvestorService() InitiatorInvestor {
	return func(s *investorService) *investorService {
		return s
	}
}

func (i InitiatorInvestor) SetRepository(investorRepository db.InvestorRepository) InitiatorInvestor {
	return func(s *investorService) *investorService {
		i(s).investorRepo = investorRepository
		return s
	}
}

func (i InitiatorInvestor) Build() InvestorService {
	return i(&investorService{})
}
