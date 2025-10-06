package entity

import "mime/multipart"

const (
	LocalUploadPath           = "storage/uploads"
	LocalAggrementLetterPath  = "storage/agreements"
	PublicUploadPath          = "public/uploads"
	PublicAggrementLetterPath = "storage/agreements"
)

type UploadFileInput struct {
	File *multipart.FileHeader
}

type InvestorAgreementLetterInvestor struct {
	Name    string
	Amount  string
	Percent float64
}
type InvestorAgreementLetterInput struct {
	AgreementNo  string
	EffectiveOn  string
	BorrowerName string
	Amount       string
	Rate         string
	Term         string
	Investors    []InvestorAgreementLetterInvestor
}
