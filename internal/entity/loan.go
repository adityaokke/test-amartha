package entity

import (
	"time"

	"gorm.io/gorm"
)

type LoanStatus string

const (
	LoanStatusProposed  LoanStatus = "PROPOSED"
	LoanStatusApproved  LoanStatus = "APPROVED"
	LoanStatusInvested  LoanStatus = "INVESTED"
	LoanStatusDisbursed LoanStatus = "DISBURSED"
)

func (ls LoanStatus) IsValid() bool {
	switch ls {
	case LoanStatusProposed, LoanStatusApproved, LoanStatusInvested, LoanStatusDisbursed:
		return true
	}
	return false
}

type Loan struct {
	ID             int        `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID         int        `json:"userId" gorm:"index;"`
	Amount         int        `json:"amount" gorm:"type:INTEGER;"`
	Status         LoanStatus `json:"status" gorm:"type:VARCHAR(50);"`
	InvestedAmount int        `json:"investedAmount" gorm:"type:INTEGER;default:0;"`
	Rate           float64    `json:"rate" gorm:"type:FLOAT;default:0;"`
	Term           int        `json:"term" gorm:"type:INTEGER;default:0;"`
	// approval info
	PhotoProofURL               *string    `json:"photoProofUrl" gorm:"type:TEXT;"`
	ApprovedByEmployeeID        *int       `json:"employeeId" gorm:"index;"`
	ApprovedAt                  *time.Time `json:"approvedAt" gorm:"type:DATETIME;"`
	FullyInvestedAt             *time.Time `json:"fullyInvestedAt" gorm:"type:DATETIME;"`
	DraftLoanAgreementLetterURL *string    `json:"draftLoanAgreementLetterUrl" gorm:"type:TEXT;"`
	// disbursement info
	LoanAgreementLetterURL         *string    `json:"loanAgreementLetterUrl" gorm:"type:TEXT;"`
	AgreementCollectedByEmployeeID *int       `json:"agreementCollectedByEmployeeId" gorm:"index;"`
	DisbursedByEmployeeID          *int       `json:"disbursedByEmployeeId" gorm:"index;"`
	DisbursedAt                    *time.Time `json:"disbursedAt" gorm:"type:DATETIME;"`
	BaseTimeStruct
}

func (Loan) TableName() string {
	return "loan"
}

func (l *Loan) BeforeCreate(tx *gorm.DB) (err error) {
	if !l.Status.IsValid() {
		l.Status = LoanStatusProposed
	}
	return
}

type LoansInput struct {
	UserID *int
	Status *LoanStatus
}

type LoanInput struct {
	ID *int
}

type WhereLoan struct {
	ID     *int
	UserID *int
	Status *LoanStatus
}

func (w *WhereLoan) Scan(input any) {
	switch v := input.(type) {
	case LoanInput:
		w.ID = v.ID
	case LoansInput:
		w.UserID = v.UserID
		w.Status = v.Status
	}
}

type TermUnit string

const (
	TermUnitWeek TermUnit = "WEEKLY"
)

func (e TermUnit) IsValid() bool {
	switch e {
	case TermUnitWeek:
		return true
	}
	return false
}

type ProposeLoanInput struct {
	UserID   int
	Amount   int
	Rate     float64
	Term     int
	TermUnit TermUnit
}

type PatchLoanInput struct {
	ID int
	// approval info
	EmployeeID    int
	PhotoProofURL string
	Status        LoanStatus
	// disbursement info
	LoanAgreementLetterURL         string
	DisbursedByEmployeeID          int
	AgreementCollectedByEmployeeID int
}

type ApproveLoanInput struct {
	ID            int
	EmployeeID    int
	PhotoProofURL string
}

type DisburseLoanInput struct {
	ID                             int
	DisbursedByEmployeeID          int
	LoanAgreementLetterURL         string
	AgreementCollectedByEmployeeID int
}
