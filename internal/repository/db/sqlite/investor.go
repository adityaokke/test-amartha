package sqlite

import (
	"context"

	"github.com/adityaokke/test-amartha/internal/entity"
	"github.com/adityaokke/test-amartha/internal/repository/db"
	"gorm.io/gorm"
)

type investorRepository struct {
	db *gorm.DB
}

func (r investorRepository) Create(ctx context.Context, item *entity.Investor) (err error) {
	db := r.db

	if err = db.Create(item).Error; err != nil {
		return
	}

	return
}

func (r investorRepository) Update(ctx context.Context, item *entity.Investor) (err error) {
	db := r.db

	if err = db.Save(item).Error; err != nil {
		return
	}
	return
}

func (r investorRepository) Delete(ctx context.Context, item *entity.Investor) (err error) {
	db := r.db

	if err := db.Delete(item).Error; err != nil {
		return err
	}
	return nil
}

func getWhereInvestor(db *gorm.DB, filter *entity.WhereInvestor) *gorm.DB {
	tableName := entity.Investor{}.TableName()
	if filter.ID != nil {
		db = db.Where(tableName+".id = ?", *filter.ID)
	}
	if filter.IDs != nil {
		if len(*filter.IDs) > 0 {
			db = db.Where(tableName+".id IN (?)", *filter.IDs)
		} else {
			db = db.Where("1 = 0")
		}
	}
	if filter.Email != nil {
		db = db.Where(tableName+".email = ?", *filter.Email)
	}
	return db
}

func (r investorRepository) Investors(ctx context.Context, filter entity.InvestorsInput) (result []entity.Investor, err error) {
	db := r.db

	where := entity.WhereInvestor{}
	where.Scan(filter)
	db = getWhereInvestor(db, &where)

	if err = db.Find(&result).Error; err != nil {
		return
	}

	return
}

func (r investorRepository) CountInvestors(ctx context.Context, filter entity.InvestorsInput) (result int64, err error) {
	db := r.db

	where := entity.WhereInvestor{}
	where.Scan(filter)
	db = getWhereInvestor(db, &where)

	if err = db.Model(&entity.Investor{}).Count(&result).Error; err != nil {
		return
	}

	return
}

func (r investorRepository) Investor(ctx context.Context, filter entity.InvestorInput) (result entity.Investor, err error) {
	db := r.db

	where := entity.WhereInvestor{}
	where.Scan(filter)
	db = getWhereInvestor(db, &where)

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
type initiatorInvestorRepository func(s *investorRepository) *investorRepository

func NewInvestorRepository() initiatorInvestorRepository {
	return func(q *investorRepository) *investorRepository {
		return q
	}
}

func (i initiatorInvestorRepository) SetDBConnection(db *gorm.DB) initiatorInvestorRepository {
	return func(s *investorRepository) *investorRepository {
		i(s).db = db
		return s
	}
}

func (i initiatorInvestorRepository) Build() db.InvestorRepository {
	return i(&investorRepository{})
}
