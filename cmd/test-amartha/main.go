package main

import (
	"fmt"
	"os"

	"github.com/adityaokke/test-amartha/internal/delivery/rest"
	"github.com/adityaokke/test-amartha/internal/repository/db/sqlite"
	"github.com/adityaokke/test-amartha/internal/repository/db/sqlite/migration"
	"github.com/adityaokke/test-amartha/internal/service"
	driver "github.com/glebarez/sqlite"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/gorm"
)

func main() {
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

	nodeEnv := os.Getenv("NODE_ENV")
	if nodeEnv == "development" {
		db = db.Debug()
	}
	loanRepo := sqlite.NewLoanRepository().
		SetDBConnection(db).
		Build()

	loanService := service.NewLoanService().
		SetRepository(loanRepo).
		Build()

	loanHandler := rest.NewLoanHandler(loanService)
	rest.Router(
		e,
		loanHandler,
	)

	host := "localhost"
	port := 3000
	e.Logger.Fatal(e.Start(fmt.Sprintf("%s:%d", host, port)))
}
