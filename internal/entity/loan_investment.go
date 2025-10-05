package entity

type LoanInvestment struct {
	ID         int `json:"id" gorm:"primaryKey;autoIncrement"`
	LoanID     int `json:"loanId" gorm:"index;"`
	InvestorID int `json:"investorID" gorm:"index;"`
	Amount     int `json:"amount" gorm:"type:INTEGER;"`
	BaseTimeStruct
}

func (LoanInvestment) TableName() string {
	return "loan_investment"
}

type LoanInvestmentsInput struct {
	LoanID *int
}

type LoanInvestmentInput struct {
	ID *int
}

type WhereLoanInvestment struct {
	ID     *int
	LoanID *int
}

func (w *WhereLoanInvestment) Scan(input any) {
	switch v := input.(type) {
	case LoanInvestmentInput:
		w.ID = v.ID
	case LoanInvestmentsInput:
		w.LoanID = v.LoanID
	}
}

type InvestLoanInput struct {
	LoanID     int
	InvestorID int
	Amount     int
}
