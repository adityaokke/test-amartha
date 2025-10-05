package main

import (
	"fmt"
	"os"

	"github.com/adityaokke/test-amartha/internal/delivery/rest"
	"github.com/adityaokke/test-amartha/internal/entity"
	"github.com/adityaokke/test-amartha/internal/repository/db/sqlite"
	"github.com/adityaokke/test-amartha/internal/repository/db/sqlite/migration"
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

	db, err := gorm.Open(driver.Open("amartha.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	migration.Migrate(db)

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

	loanService := service.NewLoanService().
		SetRepository(loanRepo).
		SetLoanInvestmentRepository(loanInvestmentRepo).
		Build()
	loanInvestmentService := service.NewLoanInvestmentService().
		SetRepository(loanInvestmentRepo).
		Build()
	fileService := service.NewFileService().
		Build()

	loanHandler := rest.NewLoanHandler(loanService)
	loanInvestmentHandler := rest.NewLoanInvestmentHandler(loanInvestmentService)
	fileHandler := rest.NewFileHandler(fileService)
	rest.Router(
		e,
		loanHandler,
		loanInvestmentHandler,
		fileHandler,
	)

	host := "localhost"
	port := 3000
	e.Logger.Fatal(e.Start(fmt.Sprintf("%s:%d", host, port)))
}
