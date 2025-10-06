package entity

type SendMailInput struct {
	To      string
	Subject string
	Body    string
}

type SendInvestorAgreementMailInput struct {
	To           string
	InvestorName string
	InvestDate   string
	Amount       string
	AgreementURL string
}
