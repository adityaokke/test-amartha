package migration

import (
	"github.com/adityaokke/test-amartha/internal/entity"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	db.AutoMigrate(&entity.Loan{}, &entity.LoanInvestment{})
}
