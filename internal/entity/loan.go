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
	ID     int        `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID int        `json:"userId" gorm:"index;"`
	Amount int        `json:"amount" gorm:"type:INTEGER;"`
	Status LoanStatus `json:"status" gorm:"type:VARCHAR(50);"`
	// approval info
	PhotoProofURL        *string    `json:"photoProofUrl" gorm:"type:TEXT;"`
	ApprovedByEmployeeID *int       `json:"employeeId" gorm:"index;"`
	ApprovedAt           *time.Time `json:"approvedAt" gorm:"type:DATETIME;"`
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
}

type LoanInput struct {
	ID *int
}

type WhereLoan struct {
	ID     *int
	UserID *int
}

func (w *WhereLoan) Scan(input any) {
	switch v := input.(type) {
	case LoanInput:
		w.ID = v.ID
	case LoansInput:
		w.UserID = v.UserID
	}
}

type ProposeLoanInput struct {
	UserID int
	Amount int
}

type ApproveLoanInput struct {
	ID            int
	EmployeeID    int
	PhotoProofURL string
}

type DisburseLoanInput struct {
	LoanID                         int
	LoanAgreementLetterURL         string
	DisbursedByEmployeeID          int
	AgreementCollectedByEmployeeID int
	DisbursedAt                    time.Time
}

type InvestLoanInput struct {
	ID     int
	UserID int
}
