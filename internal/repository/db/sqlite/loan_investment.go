package sqlite

import (
	"context"
	"errors"

	"github.com/adityaokke/test-amartha/internal/entity"
	"github.com/adityaokke/test-amartha/internal/repository/db"
	"gorm.io/gorm"
)

type loanInvestmentRepository struct {
	db *gorm.DB
}

func (r loanInvestmentRepository) InvestLoan(ctx context.Context, item *entity.LoanInvestment) (err error) {
	err = r.db.Transaction(func(tx *gorm.DB) (errTx error) {
		if errTx = tx.Create(item).Error; errTx != nil {
			return
		}

		res := tx.Model(&entity.Loan{}).Where("id = ? AND invested_amount + ? <= amount", item.LoanID, item.Amount).UpdateColumn("invested_amount", gorm.Expr("invested_amount + ?", item.Amount))
		errTx = res.Error
		if errTx != nil {
			return
		}
		if res.RowsAffected == 0 {
			errTx = errors.New("failed to update loan invested amount, possibly exceeding loan amount")
			return
		}
		return
	})
	return
}

func getWhereLoanInvestment(db *gorm.DB, filter *entity.WhereLoanInvestment) *gorm.DB {
	tableName := entity.LoanInvestment{}.TableName()
	if filter.ID != nil {
		db = db.Where(tableName+".id = ?", *filter.ID)
	}
	if filter.LoanID != nil {
		db = db.Where(tableName+".loan_id = ?", *filter.LoanID)
	}
	return db
}

func (r loanInvestmentRepository) LoanInvestments(ctx context.Context, filter entity.LoanInvestmentsInput) (result []entity.LoanInvestment, err error) {
	db := r.db

	where := entity.WhereLoanInvestment{}
	where.Scan(filter)
	db = getWhereLoanInvestment(db, &where)

	if err = db.Find(&result).Error; err != nil {
		return
	}

	return
}

func (r loanInvestmentRepository) CountLoanInvestments(ctx context.Context, filter entity.LoanInvestmentsInput) (result int64, err error) {
	db := r.db

	where := entity.WhereLoanInvestment{}
	where.Scan(filter)
	db = getWhereLoanInvestment(db, &where)

	if err = db.Model(&entity.LoanInvestment{}).Count(&result).Error; err != nil {
		return
	}

	return
}

func (r loanInvestmentRepository) LoanInvestment(ctx context.Context, filter entity.LoanInvestmentInput) (result entity.LoanInvestment, err error) {
	db := r.db

	where := entity.WhereLoanInvestment{}
	where.Scan(filter)
	db = getWhereLoanInvestment(db, &where)

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
type initiatorLoanInvestmentRepository func(s *loanInvestmentRepository) *loanInvestmentRepository

func NewLoanInvestmentRepository() initiatorLoanInvestmentRepository {
	return func(q *loanInvestmentRepository) *loanInvestmentRepository {
		return q
	}
}

func (i initiatorLoanInvestmentRepository) SetDBConnection(db *gorm.DB) initiatorLoanInvestmentRepository {
	return func(s *loanInvestmentRepository) *loanInvestmentRepository {
		i(s).db = db
		return s
	}
}

func (i initiatorLoanInvestmentRepository) Build() db.LoanInvestmentRepository {
	return i(&loanInvestmentRepository{})
}
