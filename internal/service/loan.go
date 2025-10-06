package service

import (
	"context"
	"errors"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/adityaokke/test-amartha/internal/entity"
	"github.com/adityaokke/test-amartha/internal/repository/db"
	"github.com/adityaokke/test-amartha/internal/repository/mail"
	"github.com/adityaokke/test-amartha/internal/repository/pdf"
	"gorm.io/gorm"
)

type LoanService interface {
	ProposeLoan(ctx context.Context, input entity.ProposeLoanInput) (result entity.Loan, err error)
	ApproveLoan(ctx context.Context, input entity.ApproveLoanInput) (result entity.Loan, err error)
	InvestLoan(ctx context.Context, input entity.InvestLoanInput) (result entity.LoanInvestment, err error)
	DisburseLoan(ctx context.Context, input entity.DisburseLoanInput) (result entity.Loan, err error)

	Loans(ctx context.Context, filter entity.LoansInput) (result []entity.Loan, err error)
	CountLoans(ctx context.Context, filter entity.LoansInput) (result int64, err error)
	Loan(ctx context.Context, filter entity.LoanInput) (result entity.Loan, err error)
	GetDraftLoanAgreementLetter(ctx context.Context, loanID int) (result string, err error)
	GetSignedLoanAgreementLetter(ctx context.Context, loanID int) (result string, err error)
}

func (s *loanService) ProposeLoan(ctx context.Context, input entity.ProposeLoanInput) (result entity.Loan, err error) {
	if input.UserID == 0 {
		err = errors.New("userId is required")
		return
	}
	if input.Amount == 0 {
		err = errors.New("amount is required")
		return
	}
	if input.Rate == 0 {
		err = errors.New("rate is required")
		return
	}
	if input.Term == 0 {
		err = errors.New("term is required")
		return
	}
	if !input.TermUnit.IsValid() {
		err = errors.New("termUnit is invalid")
		return
	}
	item := entity.Loan{
		UserID: input.UserID,
		Amount: input.Amount,
		Status: entity.LoanStatusProposed,
		Rate:   input.Rate,
		Term:   input.Term,
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
		err = errors.New("employeeId is required")
		return
	}
	if input.PhotoProofURL == "" {
		err = errors.New("photoProofUrl is required")
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

	// check if investor exists
	_, err = s.investorRepo.Investor(ctx, entity.InvestorInput{
		ID: &input.InvestorID,
	})
	if err != nil {
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

	loanInvestment, err := s.loanInvestmentRepo.LoanInvestment(ctx, entity.LoanInvestmentInput{
		LoanID:     &input.LoanID,
		InvestorID: &input.InvestorID,
	})
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return
	}
	if loanInvestment.ID != 0 {
		err = errors.New("investor has already invested in this loan")
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
		// update loan fully invested at
		fullInvestedAt := time.Now().UTC()
		loan.FullyInvestedAt = &fullInvestedAt
		err = s.loanRepo.Update(ctx, &loan)
		if err != nil {
			return
		}
		go func() {
			// generate loan agreement pdf
			var pdfRelativePath string
			pdfRelativePath, err = s.generateLoanAgreementPDF(ctx, loan)
			if err != nil {
				return
			}
			u, _ := url.Parse(os.Getenv("APP_HOST"))
			draftLoanAgreementLetterURL, err := url.JoinPath(u.String(), filepath.ToSlash(pdfRelativePath))
			if err != nil {
				return
			}
			loan.DraftLoanAgreementLetterURL = &draftLoanAgreementLetterURL
			err = s.loanRepo.Update(ctx, &loan)
			if err != nil {
				return
			}
			// send email to investors
			err = s.sendLoanAgreementEmail(ctx, loan)
			if err != nil {
				return
			}
		}()
	}
	result = item
	return
}

func (s *loanService) generateLoanAgreementPDF(ctx context.Context, loan entity.Loan) (pdfRelativePath string, err error) {
	var loanInvestments []entity.LoanInvestment
	loanInvestments, err = s.loanInvestmentRepo.LoanInvestments(ctx, entity.LoanInvestmentsInput{
		LoanID: &loan.ID,
	})
	if err != nil {
		return
	}
	investorIDs := make([]int, 0)
	investmentsMap := make(map[int]entity.LoanInvestment)
	for _, investment := range loanInvestments {
		investorIDs = append(investorIDs, investment.InvestorID)
		investmentsMap[investment.InvestorID] = investment
	}
	var investors []entity.Investor
	investors, err = s.investorRepo.Investors(ctx, entity.InvestorsInput{
		IDs: &investorIDs,
	})
	if err != nil {
		return
	}

	investorsPdf := []entity.InvestorAgreementLetterInvestor{}
	for _, investor := range investors {
		investment := investmentsMap[investor.ID]
		investorsPdf = append(investorsPdf, entity.InvestorAgreementLetterInvestor{
			Name:    investor.Email,
			Amount:  strconv.Itoa(investment.Amount),
			Percent: (float64(investment.Amount) / float64(loan.Amount)) * 100,
		})
	}
	pdfRelativePath, err = s.pdfApi.GenerateAgreementPDF(entity.InvestorAgreementLetterInput{
		AgreementNo:  strconv.Itoa(loan.ID),
		EffectiveOn:  loan.FullyInvestedAt.Format("02 Jan 2006"),
		BorrowerName: strconv.Itoa(loan.UserID),
		Amount:       strconv.Itoa(loan.Amount),
		Rate:         strconv.FormatFloat(loan.Rate, 'f', 2, 64),
		Term:         strconv.Itoa(loan.Term),
		Investors:    investorsPdf,
	})
	if err != nil {
		return
	}
	return
}

func (s *loanService) sendLoanAgreementEmail(ctx context.Context, loan entity.Loan) (err error) {
	var loanInvestments []entity.LoanInvestment
	loanInvestments, err = s.loanInvestmentRepo.LoanInvestments(ctx, entity.LoanInvestmentsInput{
		LoanID: &loan.ID,
	})
	if err != nil {
		return
	}
	investorIDs := make([]int, 0)
	investmentsMap := make(map[int]entity.LoanInvestment)
	for _, investment := range loanInvestments {
		investorIDs = append(investorIDs, investment.InvestorID)
		investmentsMap[investment.InvestorID] = investment
	}
	var investors []entity.Investor
	investors, err = s.investorRepo.Investors(ctx, entity.InvestorsInput{
		IDs: &investorIDs,
	})
	if err != nil {
		return
	}
	for _, investor := range investors {
		investment := investmentsMap[investor.ID]
		err = s.mailApi.SendInvestorAgreementMail(ctx, entity.SendInvestorAgreementMailInput{
			To:           investor.Email,
			InvestorName: investor.Email,
			InvestDate:   investment.CreatedAt.Format("02 Jan 2006"),
			Amount:       strconv.Itoa(investment.Amount),
			AgreementURL: *loan.DraftLoanAgreementLetterURL,
		})
		if err != nil {
			return
		}
	}
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

func (s *loanService) GetDraftLoanAgreementLetter(ctx context.Context, loanID int) (result string, err error) {
	loan, err := s.loanRepo.Loan(ctx, entity.LoanInput{
		ID: &loanID,
	})
	if err != nil {
		return
	}
	if loan.FullyInvestedAt == nil {
		err = errors.New("loan is not fully funded yet")
		return
	}
	if loan.DraftLoanAgreementLetterURL == nil {
		err = errors.New("loan agreement letter is not available")
		return
	}
	result = *loan.DraftLoanAgreementLetterURL
	return
}

func (s *loanService) GetSignedLoanAgreementLetter(ctx context.Context, loanID int) (result string, err error) {
	loan, err := s.loanRepo.Loan(ctx, entity.LoanInput{
		ID: &loanID,
	})
	if err != nil {
		return
	}
	if loan.Status != entity.LoanStatusDisbursed {
		err = errors.New("loan is not disbursed yet")
		return
	}
	if loan.LoanAgreementLetterURL == nil {
		err = errors.New("loan agreement letter is not available")
		return
	}
	result = *loan.LoanAgreementLetterURL
	return
}

type loanService struct {
	loanRepo           db.LoanRepository
	loanInvestmentRepo db.LoanInvestmentRepository
	investorRepo       db.InvestorRepository
	mailApi            mail.MailApi
	pdfApi             pdf.PdfApi
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

func (i InitiatorLoan) SetInvestorRepository(investorRepository db.InvestorRepository) InitiatorLoan {
	return func(s *loanService) *loanService {
		i(s).investorRepo = investorRepository
		return s
	}
}

func (i InitiatorLoan) SetMailApi(mailApi mail.MailApi) InitiatorLoan {
	return func(s *loanService) *loanService {
		i(s).mailApi = mailApi
		return s
	}
}

func (i InitiatorLoan) SetPdfApi(pdfApi pdf.PdfApi) InitiatorLoan {
	return func(s *loanService) *loanService {
		i(s).pdfApi = pdfApi
		return s
	}
}

func (i InitiatorLoan) Build() LoanService {
	return i(&loanService{})
}
