package db

import (
	"context"
	"crispy/domain"
	"database/sql"
	"fmt"
)

type SQLite_TagBuilder struct {
	SQLite_Builder
}

type SQLite_TransactionTagsBuilder struct {
	SQLite_Builder
}

func (s *SQLiteDB) TagQueryBuilder() *SQLite_TagBuilder {
	builder := &SQLite_TagBuilder{
		SQLite_Builder: *s.SQLBuilder(Table_Tag),
	}
	builder.OrIgnore = true
	return builder
}

func (s *SQLiteDB) TransactionTagsQueryBuilder() *SQLite_TransactionTagsBuilder {
	return &SQLite_TransactionTagsBuilder{
		SQLite_Builder: *s.SQLBuilder(Table_TransactionTags),
	}
}

func (s *SQLite_TagBuilder) AddColumn(c Column) *SQLite_TagBuilder {
	s.SQLite_Builder.AddColumn(c)
	return s
}

func (s *SQLite_TransactionTagsBuilder) AddColumn(c Column) *SQLite_TransactionTagsBuilder {
	s.SQLite_Builder.AddColumn(c)
	return s
}

func (s *SQLite_TagBuilder) AddColumns(columns []Column) *SQLite_TagBuilder {
	s.SQLite_Builder.AddColumns(columns)
	return s
}

func (s *SQLite_TransactionTagsBuilder) AddColumns(columns []Column) *SQLite_TransactionTagsBuilder {
	s.SQLite_Builder.AddColumns(columns)
	return s
}

func (s *SQLite_TagBuilder) SetCondition(cond Condition) *SQLite_TagBuilder {
	s.SQLite_Builder.SetCondition(cond)
	return s
}

func (s *SQLite_TransactionTagsBuilder) SetCondition(cond Condition) *SQLite_TransactionTagsBuilder {
	s.SQLite_Builder.SetCondition(cond)
	return s
}

func (s *SQLite_TagBuilder) AddJoin(tableName Table, condition string) *SQLite_TagBuilder {
	s.SQLite_Builder.AddJoin(tableName, condition)
	return s
}

func (s *SQLite_TransactionTagsBuilder) AddJoin(tableName Table, condition string) *SQLite_TransactionTagsBuilder {
	s.SQLite_Builder.AddJoin(tableName, condition)
	return s
}

func (s *SQLite_TagBuilder) AddOrderBy(c []Column, r bool) *SQLite_TagBuilder {
	s.SQLite_Builder.AddOrderBy(c, r)
	return s
}

func (s *SQLite_TransactionTagsBuilder) AddOrderBy(c []Column, r bool) *SQLite_TransactionTagsBuilder {
	s.SQLite_Builder.AddOrderBy(c, r)
	return s
}

func (s *SQLite_TagBuilder) Select(ctx context.Context, page, perPage int) ([]domain.Tag, error) {
	var id int64
	var tag string

	var tags []domain.Tag
	callback := func(nextRow *sql.Rows) error {
		err := nextRow.Scan(&id, &tag)
		if err != nil {
			return fmt.Errorf("Scan failed: %w", err)
		}
		tags = append(tags, domain.Tag{Id: id, Name: tag})
		return nil
	}
	if err := s.SQLite_Builder.Select(ctx, page, perPage, callback); err != nil {
		return nil, fmt.Errorf("SELECT in \"Tag\" failed: %w", err)
	}

	return tags, nil
}

func (s *SQLite_TransactionTagsBuilder) Select(ctx context.Context, page, perPage int) ([]int64, error) {
	var tagId int64

	var tagIdList []int64
	callback := func(nextRow *sql.Rows) error {
		err := nextRow.Scan(tagId)
		if err != nil {
			return fmt.Errorf("Scan failed: %w", err)
		}
		tagIdList = append(tagIdList, tagId)
		return nil
	}
	if err := s.SQLite_Builder.Select(ctx, page, perPage, callback); err != nil {
		return nil, fmt.Errorf("SELECT in \"Tag\" failed: %w", err)
	}

	return tagIdList, nil
}
