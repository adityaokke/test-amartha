package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/adityaokke/test-amartha/internal/entity"
	"github.com/adityaokke/test-amartha/internal/repository/db"
)

type LoanService interface {
	ProposeLoan(ctx context.Context, input entity.ProposeLoanInput) (result entity.Loan, err error)
	ApproveLoan(ctx context.Context, input entity.ApproveLoanInput) (result entity.Loan, err error)
	InvestLoan(ctx context.Context, input entity.InvestLoanInput) (result entity.LoanInvestment, err error)
	DisburseLoan(ctx context.Context, input entity.DisburseLoanInput) (result entity.Loan, err error)

	Loans(ctx context.Context, filter entity.LoansInput) (result []entity.Loan, err error)
	CountLoans(ctx context.Context, filter entity.LoansInput) (result int64, err error)
	Loan(ctx context.Context, filter entity.LoanInput) (result entity.Loan, err error)
}

func (s *loanService) ProposeLoan(ctx context.Context, input entity.ProposeLoanInput) (result entity.Loan, err error) {
	if input.UserID == 0 {
		err = errors.New("user_id is required")
		return
	}
	item := entity.Loan{
		UserID: input.UserID,
		Amount: input.Amount,
		Status: entity.LoanStatusProposed,
	}
	err = s.loanRepo.Create(ctx, &item)
	if err != nil {
		return
	}
	result = item
	return
}

func (s *loanService) ApproveLoan(ctx context.Context, input entity.ApproveLoanInput) (result entity.Loan, err error) {
	if input.ID == 0 {
		err = errors.New("id is required")
		return
	}
	if input.EmployeeID == 0 {
		err = errors.New("employee_id is required")
		return
	}
	if input.PhotoProofURL == "" {
		err = errors.New("photo_proof_url is required")
		return
	}
	currentItem, err := s.loanRepo.Loan(ctx, entity.LoanInput{
		ID: &input.ID,
	})
	if err != nil {
		return
	}

	if currentItem.Status != entity.LoanStatusProposed {
		err = errors.New("only proposed loan can be approved")
		return
	}

	currentItem.ApprovedByEmployeeID = &input.EmployeeID
	currentItem.Status = entity.LoanStatusApproved
	trimmedURL := strings.TrimSpace(input.PhotoProofURL)
	currentItem.PhotoProofURL = &trimmedURL
	approvedAt := time.Now().UTC()
	currentItem.ApprovedAt = &approvedAt
	err = s.loanRepo.Update(ctx, &currentItem)
	if err != nil {
		return
	}
	result = currentItem
	return
}

func (s *loanService) InvestLoan(ctx context.Context, input entity.InvestLoanInput) (result entity.LoanInvestment, err error) {
	if input.LoanID == 0 {
		err = errors.New("loanId is required")
		return
	}
	if input.InvestorID == 0 {
		err = errors.New("investorId is required")
		return
	}
	if input.Amount == 0 {
		err = errors.New("amount is required")
		return
	}

	loan, err := s.loanRepo.Loan(ctx, entity.LoanInput{
		ID: &input.LoanID,
	})
	if err != nil {
		return
	}
	if loan.Status != entity.LoanStatusApproved {
		err = errors.New("only approved loan can be invested")
		return
	}
	if loan.InvestedAmount >= loan.Amount {
		err = errors.New("loan is already fully funded")
		return
	}
	if loan.InvestedAmount+input.Amount > loan.Amount {
		err = errors.New("investment would exceed loan amount")
		return
	}

	item := entity.LoanInvestment{
		LoanID:     input.LoanID,
		InvestorID: input.InvestorID,
		Amount:     input.Amount,
	}
	err = s.loanInvestmentRepo.InvestLoan(ctx, &item)
	if err != nil {
		return
	}

	loan, err = s.loanRepo.Loan(ctx, entity.LoanInput{
		ID: &input.LoanID,
	})
	if err != nil {
		return
	}
	if loan.InvestedAmount == loan.Amount {
		// TODO: send email to investor
	}
	result = item
	return
}

func (s *loanService) DisburseLoan(ctx context.Context, input entity.DisburseLoanInput) (result entity.Loan, err error) {
	if input.ID == 0 {
		err = errors.New("id is required")
		return
	}
	if input.DisbursedByEmployeeID == 0 {
		err = errors.New("disbursedByEmployeeId is required")
		return
	}
	if input.LoanAgreementLetterURL == "" {
		err = errors.New("loanAgreementLetterUrl is required")
		return
	}
	if input.AgreementCollectedByEmployeeID == 0 {
		err = errors.New("agreementCollectedByEmployeeId is required")
		return
	}
	currentItem, err := s.loanRepo.Loan(ctx, entity.LoanInput{
		ID: &input.ID,
	})
	if err != nil {
		return
	}

	if currentItem.Status != entity.LoanStatusApproved {
		err = errors.New("only approved loan can be disbursed")
		return
	}
	if currentItem.InvestedAmount < currentItem.Amount {
		err = errors.New("only fully invested loan can be disbursed")
		return
	}
	loanInvestments, err := s.loanInvestmentRepo.LoanInvestments(ctx, entity.LoanInvestmentsInput{
		LoanID: &input.ID,
	})
	if err != nil {
		return
	}
	totalInvestment := 0
	for _, investment := range loanInvestments {
		totalInvestment += investment.Amount
	}
	if totalInvestment != currentItem.Amount {
		err = errors.New("total investment does not match loan amount, please contact admin/cs")
		return
	}

	currentItem.DisbursedByEmployeeID = &input.DisbursedByEmployeeID
	trimmedURL := strings.TrimSpace(input.LoanAgreementLetterURL)
	currentItem.LoanAgreementLetterURL = &trimmedURL
	currentItem.AgreementCollectedByEmployeeID = &input.AgreementCollectedByEmployeeID
	disbursedAt := time.Now().UTC()
	currentItem.DisbursedAt = &disbursedAt
	currentItem.Status = entity.LoanStatusDisbursed
	err = s.loanRepo.Update(ctx, &currentItem)
	if err != nil {
		return
	}
	result = currentItem
	return
}

func (s *loanService) Loans(ctx context.Context, filter entity.LoansInput) (result []entity.Loan, err error) {
	result, err = s.loanRepo.Loans(ctx, filter)
	if err != nil {
		return
	}
	return
}

func (s *loanService) CountLoans(ctx context.Context, filter entity.LoansInput) (result int64, err error) {
	result, err = s.loanRepo.CountLoans(ctx, filter)
	if err != nil {
		return
	}
	return
}

func (s *loanService) Loan(ctx context.Context, filter entity.LoanInput) (result entity.Loan, err error) {
	result, err = s.loanRepo.Loan(ctx, filter)
	if err != nil {
		return
	}
	return
}

type loanService struct {
	loanRepo           db.LoanRepository
	loanInvestmentRepo db.LoanInvestmentRepository
}

type InitiatorLoan func(s *loanService) *loanService

func NewLoanService() InitiatorLoan {
	return func(s *loanService) *loanService {
		return s
	}
}

func (i InitiatorLoan) SetRepository(loanRepository db.LoanRepository) InitiatorLoan {
	return func(s *loanService) *loanService {
		i(s).loanRepo = loanRepository
		return s
	}
}

func (i InitiatorLoan) SetLoanInvestmentRepository(loanInvestmentRepository db.LoanInvestmentRepository) InitiatorLoan {
	return func(s *loanService) *loanService {
		i(s).loanInvestmentRepo = loanInvestmentRepository
		return s
	}
}

func (i InitiatorLoan) Build() LoanService {
	return i(&loanService{})
}
