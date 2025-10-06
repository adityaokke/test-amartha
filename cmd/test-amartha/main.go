package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/adityaokke/test-amartha/internal/delivery/rest"
	"github.com/adityaokke/test-amartha/internal/entity"
	pkgMail "github.com/adityaokke/test-amartha/internal/pkg/mail"
	"github.com/adityaokke/test-amartha/internal/repository/db/sqlite"
	"github.com/adityaokke/test-amartha/internal/repository/db/sqlite/migration"
	"github.com/adityaokke/test-amartha/internal/repository/mail"
	"github.com/adityaokke/test-amartha/internal/repository/pdf"
	"github.com/adityaokke/test-amartha/internal/service"
	driver "github.com/glebarez/sqlite"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/gorm"
)

func main() {
	godotenv.Load(".env")

	// ensure dir
	os.MkdirAll(entity.LocalUploadPath, 0o755)
	os.MkdirAll(entity.LocalAggrementLetterPath, 0o755)

	db, err := gorm.Open(driver.Open("amartha.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	migration.Migrate(db)

	mailFrom := os.Getenv("MAIL_FROM")
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPortEnv := os.Getenv("SMTP_PORT")
	smtpPort := 0
	if smtpPortEnv != "" {
		smtpPort, err = strconv.Atoi(smtpPortEnv)
		if err != nil {
			panic("invalid SMTP_PORT")
		}
	}
	smtpUser := os.Getenv("SMTP_USER")
	smtpPass := os.Getenv("SMTP_PASS")
	mailer := pkgMail.NewMailer(mailFrom, smtpHost, smtpPort, smtpUser, smtpPass)

	// initialize echo
	e := echo.New()
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${remote_ip} ${time_rfc3339_nano} \"${method} ${path}\" ${status} ${bytes_out} \"${referer}\" \"${user_agent}\"\n",
	}))
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		// AllowOrigins:     []string{"*"},
		AllowOrigins:     []string{"https://*", "http://*"},
		AllowMethods:     []string{echo.GET, echo.POST, echo.OPTIONS},
		AllowCredentials: true,
	}))

	goEnv := os.Getenv("GO_ENV")
	if goEnv == "development" {
		db = db.Debug()
	}
	loanRepo := sqlite.NewLoanRepository().
		SetDBConnection(db).
		Build()
	loanInvestmentRepo := sqlite.NewLoanInvestmentRepository().
		SetDBConnection(db).
		Build()
	investorRepo := sqlite.NewInvestorRepository().
		SetDBConnection(db).
		Build()
	mailApi := mail.NewMailApi().
		SetMailer(&mailer).
		Build()
	pdfApi := pdf.NewPdfApi().
		Build()

	loanService := service.NewLoanService().
		SetRepository(loanRepo).
		SetLoanInvestmentRepository(loanInvestmentRepo).
		SetInvestorRepository(investorRepo).
		SetMailApi(mailApi).
		SetPdfApi(pdfApi).
		Build()
	loanInvestmentService := service.NewLoanInvestmentService().
		SetRepository(loanInvestmentRepo).
		Build()
	fileService := service.NewFileService().
		Build()
	investorService := service.NewInvestorService().
		SetRepository(investorRepo).
		Build()

	loanHandler := rest.NewLoanHandler(loanService)
	loanInvestmentHandler := rest.NewLoanInvestmentHandler(loanInvestmentService)
	fileHandler := rest.NewFileHandler(fileService)
	InvestorHandler := rest.NewInvestorHandler(investorService)
	rest.Router(
		e,
		loanHandler,
		loanInvestmentHandler,
		fileHandler,
		InvestorHandler,
	)

	host := "localhost"
	port := 3000
	e.Logger.Fatal(e.Start(fmt.Sprintf("%s:%d", host, port)))
}
