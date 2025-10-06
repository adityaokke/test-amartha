package pdf

import (
	"fmt"
	"path/filepath"

	"github.com/adityaokke/test-amartha/internal/entity"
	"github.com/google/uuid"
	"github.com/jung-kurt/gofpdf"
)

type pdfApi struct {
}

type PdfApi interface {
	GenerateAgreementPDF(input entity.InvestorAgreementLetterInput) (string, error)
}

func (r pdfApi) GenerateAgreementPDF(d entity.InvestorAgreementLetterInput) (string, error) {
	filename := uuid.New().String() + ".pdf"
	fullpath := filepath.Join(entity.LocalAggrementLetterPath, filename)

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetTitle("Loan Agreement", false)
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.CellFormat(0, 10, "LOAN AGREEMENT", "", 1, "C", false, 0, "")
	pdf.Ln(12)

	pdf.SetFont("Arial", "", 11)
	para := func(s string) { pdf.MultiCell(0, 6, s, "", "L", false); pdf.Ln(1) }

	para(fmt.Sprintf("Agreement No: %s", d.AgreementNo))
	para(fmt.Sprintf("Effective Date: %s", d.EffectiveOn))
	para(fmt.Sprintf("Borrower: %s", d.BorrowerName))
	para(fmt.Sprintf("Loan Amount (Aggregate): %s", d.Amount))
	para(fmt.Sprintf("Interest: %s   Term: %s", d.Rate, d.Term))
	pdf.Ln(2)

	pdf.SetFont("Arial", "B", 12)
	para("Schedule A - Lender List")
	pdf.SetFont("Arial", "", 11)

	// Table header
	pdf.SetFillColor(245, 246, 248)
	pdf.CellFormat(100, 8, "Lender", "1", 0, "L", true, 0, "")
	pdf.CellFormat(50, 8, "Amount", "1", 0, "R", true, 0, "")
	pdf.CellFormat(30, 8, "Share %", "1", 1, "R", true, 0, "")

	// Rows
	for _, ls := range d.Investors {
		pdf.CellFormat(100, 8, ls.Name, "1", 0, "L", false, 0, "")
		pdf.CellFormat(50, 8, ls.Amount, "1", 0, "R", false, 0, "")
		pdf.CellFormat(30, 8, fmt.Sprintf("%.2f", ls.Percent), "1", 1, "R", false, 0, "")
	}
	pdf.Ln(4)

	// Pari passu / pro-rata note
	pdf.SetFont("Arial", "", 10)
	para("All payments received from the Borrower shall be applied and distributed to the Lenders on a pari passu, pro-rata basis according to their respective shares above.")

	// Signatures
	pdf.Ln(8)
	pdf.SetFont("Arial", "", 11)
	para("Signatures:")
	pdf.Ln(8)
	pdf.CellFormat(90, 8, "Borrower: ____________________", "", 0, "L", false, 0, "")
	pdf.CellFormat(0, 8, "Agent/Lender Rep: ____________", "", 1, "L", false, 0, "")

	// Save to disk
	if err := pdf.OutputFileAndClose(fullpath); err != nil {
		return "", fmt.Errorf("write pdf: %w", err)
	}
	return fullpath, nil
}

/* -------------------------------- initiator ------------------------------- */
type initiatorPdfApi func(s *pdfApi) *pdfApi

func NewPdfApi() initiatorPdfApi {
	return func(q *pdfApi) *pdfApi {
		return q
	}
}

func (i initiatorPdfApi) Build() PdfApi {
	return i(&pdfApi{})
}
