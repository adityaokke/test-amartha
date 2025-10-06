package entity

type Investor struct {
	ID    int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Email string `json:"email" gorm:"type:VARCHAR(500);uniqueIndex;"`
	BaseTimeStruct
}

func (Investor) TableName() string {
	return "investor"
}

type InvestorsInput struct {
	IDs *[]int
}

type InvestorInput struct {
	ID    *int
	Email *string
}

type WhereInvestor struct {
	ID    *int
	Email *string
	IDs   *[]int
}

func (w *WhereInvestor) Scan(input any) {
	switch v := input.(type) {
	case InvestorInput:
		w.ID = v.ID
		w.Email = v.Email
	case InvestorsInput:
		w.IDs = v.IDs
	}
}

type AddInvestorInput struct {
	Email string
}
