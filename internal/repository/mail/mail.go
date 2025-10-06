package mail

import (
	"bytes"
	"context"
	"embed"
	"errors"
	"strconv"
	"text/template"

	"github.com/adityaokke/test-amartha/internal/entity"
	pkgMail "github.com/adityaokke/test-amartha/internal/pkg/mail"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

//go:embed template/*.html
var FS embed.FS

type mailApi struct {
	mailer pkgMail.Mailer
	tt     *template.Template
}

type MailApi interface {
	SendInvestorAgreementMail(ctx context.Context, input entity.SendInvestorAgreementMailInput) (err error)
}

func (r mailApi) SendInvestorAgreementMail(ctx context.Context, input entity.SendInvestorAgreementMailInput) (err error) {
	amount, err := strconv.Atoi(input.Amount)
	if err != nil {
		return
	}
	input.Amount = message.NewPrinter(language.Indonesian).Sprint(amount)

	bodyT := r.tt.Lookup("investment-agreement.html")
	if bodyT == nil {
		err = errors.New("template not found")
		return
	}
	var bodyBuf bytes.Buffer
	err = bodyT.Execute(&bodyBuf, input)
	if err != nil {
		return
	}
	err = r.mailer.SendMail(input.To, "Your investment agreement (PDF link inside)", bodyBuf.String())
	if err != nil {
		return
	}
	return
}

/* -------------------------------- initiator ------------------------------- */
type initiatorMailApi func(s *mailApi) *mailApi

func NewMailApi() initiatorMailApi {
	return func(q *mailApi) *mailApi {
		return q
	}
}

func (i initiatorMailApi) SetMailer(mailer *pkgMail.Mailer) initiatorMailApi {
	return func(s *mailApi) *mailApi {
		s.mailer = *mailer
		return s
	}
}

func (i initiatorMailApi) Build() MailApi {
	return i(&mailApi{
		tt: template.Must(template.ParseFS(FS, "template/*.html")),
	})
}
