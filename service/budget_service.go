package service

import (
	"context"
	"crispy/db"
	"crispy/domain"
	"fmt"
	"strconv"
)

func (s *Service) CreateBudget(ctx context.Context, newBudget *domain.Budget) error {
	const errorMsg = "CreateBudget failed: %w"
	if err := newBudget.Validate(); err != nil {
		return fmt.Errorf(errorMsg, err)
	}

	if err := s.repo.BeginTx(ctx); err != nil {
		return fmt.Errorf(errorMsg, err)
	}

	// Budget
	bqb = s.repo.BudgetQueryBuilder().
		AddColumns([]db.Column{
			db.Column_Name,
			db.Column_Description,
			db.Column_Target,
			db.Column_TagId,
			db.Column_Frequency,
			db.Column_Interval,
		})
	var lastId int64
	switch newBudget.Frequency() {
	case domain.FreqWeekly:
		lastId, err = bqb.
			AddColumn(db.Column_DayOfWeek).
			Insert(ctx, newBudget.Name(), newBudget.Description(), newBudget.Target(), newBudget.TagId(), newBudget.Frequency(), newBudget.Interval(), newBudget.DayOfWeek())
	case domain.FreqMonthly:
		lastId, err = bqb.
			AddColumn(db.Column_DayOfMonth).
			Insert(ctx, newBudget.Name(), newBudget.Description(), newBudget.Target(), newBudget.TagId(), newBudget.Frequency(), newBudget.Interval(), newBudget.DayOfMonth())
	}
	if err != nil {
		s.repo.Rollback()
		return fmt.Errorf(errorMsg, err)
	}

	if err := s.repo.Commit(); err != nil {
		s.repo.Rollback()
		return fmt.Errorf(errorMsg, err)
	}

	return nil
}

func (s *Service) GetAllBudgets(ctx context.Context, page, perPage int) ([]*domain.Budget, error) {
	const errorMsg = "GetAllBudgets failed: %w"

	b_arr, err := s.repo.BudgetQueryBuilder().
		Select(ctx, page, perPage)
	if err != nil {
		return nil, fmt.Errorf(errorMsg, err)
	}

	for _, b := range b_arr {
		r, err := s.repo.RecurrenceQueryBuilder().
			SetCondition(db.NewWhere(db.Column_Id, db.Equal, strconv.FormatInt(b.)))
			Select(ctx, 0, 1)
	}
	return b_arr, nil
}

func (s *Service) UpdateBudget(ctx context.Context, newBudget *domain.Budget) error {
	const errorMsg = "UpdateBudget failed: %w"
	if err := newBudget.Validate(); err != nil {
		return fmt.Errorf(errorMsg, err)
	}
	if err := s.repo.BeginTx(ctx); err != nil {
		return fmt.Errorf(errorMsg, err)
	}

	err := s.repo.BudgetQueryBuilder().
		AddColumns([]db.Column{
			db.Column_Name,
			db.Column_Description,
			db.Column_Target,
			db.Column_TagId,
		}).
		SetCondition(
			db.NewWhere(db.Column_Id, db.Equal, strconv.FormatInt(newBudget.ID(), 10)),
		).
		Update(ctx, newBudget.Name(), newBudget.Description(), newBudget.Target(), newBudget.TagId())
	if err != nil {
		s.repo.Rollback()
		return fmt.Errorf(errorMsg, err)
	}

	if err := s.repo.Commit(); err != nil {
		s.repo.Rollback()
		return fmt.Errorf(errorMsg, err)
	}
	return nil
}
