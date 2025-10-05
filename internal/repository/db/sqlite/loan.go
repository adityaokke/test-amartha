package sqlite

import (
	"context"

	"github.com/adityaokke/test-amartha/internal/entity"
	"github.com/adityaokke/test-amartha/internal/repository/db"
	"gorm.io/gorm"
)

type loanRepository struct {
	db *gorm.DB
}

func (r loanRepository) Create(ctx context.Context, item *entity.Loan) (err error) {
	db := r.db

	if err = db.Create(item).Error; err != nil {
		return
	}

	return
}

func (r loanRepository) Update(ctx context.Context, item *entity.Loan) (err error) {
	db := r.db

	if err = db.Save(item).Error; err != nil {
		return
	}
	return
}

func (r loanRepository) Delete(ctx context.Context, item *entity.Loan) (err error) {
	db := r.db

	if err := db.Delete(item).Error; err != nil {
		return err
	}
	return nil
}

func getWhereLoan(db *gorm.DB, filter *entity.WhereLoan) *gorm.DB {
	tableName := "loan"
	if filter.ID != nil {
		db = db.Where(tableName+".id = ?", *filter.ID)
	}
	return db
}

func (r loanRepository) Loans(ctx context.Context, filter entity.LoansInput) (result []entity.Loan, err error) {
	db := r.db

	where := entity.WhereLoan{}
	where.Scan(filter)
	db = getWhereLoan(db, &where)

	if err = db.Find(&result).Error; err != nil {
		return
	}

	return
}

func (r loanRepository) CountLoans(ctx context.Context, filter entity.LoansInput) (result int64, err error) {
	db := r.db

	where := entity.WhereLoan{}
	where.Scan(filter)
	db = getWhereLoan(db, &where)

	if err = db.Model(&entity.Loan{}).Count(&result).Error; err != nil {
		return
	}

	return
}

func (r loanRepository) Loan(ctx context.Context, filter entity.LoanInput) (result entity.Loan, err error) {
	db := r.db

	where := entity.WhereLoan{}
	where.Scan(filter)
	db = getWhereLoan(db, &where)

	if _, ok := db.Statement.Clauses["WHERE"]; !ok {
		err = gorm.ErrMissingWhereClause
		return
	}

	if err = db.First(&result).Error; err != nil {
		return
	}

	return
}

/* -------------------------------- initiator ------------------------------- */
type initiatorLoanRepository func(s *loanRepository) *loanRepository

func NewLoanRepository() initiatorLoanRepository {
	return func(q *loanRepository) *loanRepository {
		return q
	}
}

func (i initiatorLoanRepository) SetDBConnection(db *gorm.DB) initiatorLoanRepository {
	return func(s *loanRepository) *loanRepository {
		i(s).db = db
		return s
	}
}

func (i initiatorLoanRepository) Build() db.LoanRepository {
	return i(&loanRepository{})
}
