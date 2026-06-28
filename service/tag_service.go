package service

import (
	"context"
	"crispy/db"
	"crispy/domain"
	"fmt"
	"strconv"
)

func (s *Service) UpsertTags(ctx context.Context, tags []string) error {
	const errorMsg = "UpsertTags failed: %w"

	if err := s.repo.BeginTx(ctx); err != nil {
		return fmt.Errorf(errorMsg, err)
	}

	for _, tag := range tags {
		if err := s.insertTag(ctx, tag); err != nil {
			s.repo.Rollback()
			return fmt.Errorf(errorMsg, err)
		}
	}
	if err := s.repo.Commit(); err != nil {
		s.repo.Rollback()
		return fmt.Errorf(errorMsg, err)
	}
	return nil
}

func (s *Service) insertTag(ctx context.Context, tag string) error {
	_, err := s.repo.TagQueryBuilder().
		AddColumn(db.Column_Name).
		Insert(ctx, tag)
	return err
}

func (s *Service) GetAllTags(ctx context.Context, page, perPage int) ([]domain.Tag, error) {
	const errorMsg = "GetAllTags failed: %w"

	tagList, err := s.repo.TagQueryBuilder().
		Select(ctx, page, perPage)
	if err != nil {
		return nil, fmt.Errorf(errorMsg, err)
	}

	return tagList, nil
}

func (s *Service) GetTagById(ctx context.Context, id int64) (string, error) {
	const errorMsg = "GetTagById failed: %w"

	tag, err := s.repo.TagQueryBuilder().
		SetCondition(db.NewWhere(db.Column_Id, db.Equal, strconv.FormatInt(id, 10))).
		Select(ctx, 0, 1)
	if err != nil {
		return "", fmt.Errorf(errorMsg, err)
	}

	return tag[0].Name, nil
}

func (s *Service) GetTagIDByName(ctx context.Context, name string) (int64, error) {
	const errorMsg = "GetTagIDByName for name %s failed: %w"

	tag, err := s.repo.TagQueryBuilder().
		SetCondition(db.NewWhere(db.Column_Name, db.Equal, name)).
		Select(ctx, 0, 1)
	if err != nil {
		return -1, fmt.Errorf(errorMsg, name, err)
	}
	return tag[0].Id, nil
}

func (s *Service) UpdateTag(ctx context.Context, id int64, newName string) error {
	const errorMsg = "UpdateTag for id %d failed: %w"

	if err := s.repo.BeginTx(ctx); err != nil {
		return fmt.Errorf(errorMsg, id, err)
	}
	err := s.repo.TagQueryBuilder().
		AddColumn(db.Column_Name).
		SetCondition(db.NewWhere(db.Column_Id, db.Equal, strconv.FormatInt(id, 10))).
		Update(ctx, newName)
	if err != nil {
		s.repo.Rollback()
		return fmt.Errorf(errorMsg, id, err)
	}
	if err := s.repo.Commit(); err != nil {
		s.repo.Rollback()
		return fmt.Errorf(errorMsg, id, err)
	}
	return nil
}

func (s *Service) DeleteTag(ctx context.Context, id int64) error {
	const errorMsg = "DeleteTag on ID %d failed: %w"
	if err := s.repo.BeginTx(ctx); err != nil {
		return fmt.Errorf(errorMsg, id, err)
	}
	err := s.repo.TransactionQueryBuilder().
		SetCondition(
			db.NewWhere(db.Column_Id, db.Equal, strconv.FormatInt(id, 10)),
		).
		Delete(ctx)
	if err != nil {
		s.repo.Rollback()
		return fmt.Errorf(errorMsg, id, err)
	}
	if err := s.repo.Commit(); err != nil {
		s.repo.Rollback()
		return fmt.Errorf(errorMsg, id, err)
	}
	return nil
}

func (s *Service) updateTransactionTags(ctx context.Context, transactionId int64, tagList []string) error {
	for _, tag := range tagList {
		if err := s.insertTag(ctx, tag); err != nil {
			return err
		}
	}

	existing, err := s.getTags(ctx, transactionId)
	if err != nil {
		return err
	}

	toAdd, toRemove := diffIDs(existing, tagList)
	for _, tag := range toRemove {
		tagId, err := s.GetTagIDByName(ctx, tag)
		if err != nil {
			return err
		}
		err = s.repo.TransactionTagsQueryBuilder().
			SetCondition(db.AndCondition{Conditions: []db.Condition{
				db.NewWhere(db.Column_TransactionId, db.Equal, strconv.FormatInt(transactionId, 10)),
				db.NewWhere(db.Column_TagId, db.Equal, strconv.FormatInt(tagId, 10)),
			}}).
			Delete(ctx)
		if err != nil {
			return err
		}
	}

	for _, tag := range toAdd {
		fmt.Println("adding: %s", tag)
		id, err := s.GetTagIDByName(ctx, tag)
		if err != nil {
			return err
		}
		_, err = s.repo.TransactionTagsQueryBuilder().
			AddColumn(db.Column_TransactionId).
			AddColumn(db.Column_TagId).
			Insert(ctx, transactionId, id)
		if err != nil {
			return err
		}
	}
	fmt.Println("Done")
	return nil
}

func diffIDs(existing, desired []string) (toAdd, toRemove []string) {
	existSet := make(map[string]struct{}, len(existing))
	for _, id := range existing {
		existSet[id] = struct{}{}
	}
	desiredSet := make(map[string]struct{}, len(desired))
	for _, id := range desired {
		desiredSet[id] = struct{}{}
	}

	for id := range desiredSet {
		if _, ok := existSet[id]; !ok {
			toAdd = append(toAdd, id)
		}
	}
	for id := range existSet {
		if _, ok := desiredSet[id]; !ok {
			toRemove = append(toRemove, id)
		}
	}
	return
}

func (s *Service) RemoveTagFromTransaction(ctx context.Context, transactionId int64, tagId int64) error {
	const errorMsg = "RemoveTagFromTransaction on transactionID %d and tagID %d failed: %w"
	if err := s.repo.BeginTx(ctx); err != nil {
		return fmt.Errorf(errorMsg, transactionId, tagId, err)
	}

	err := s.repo.TransactionTagsQueryBuilder().
		SetCondition(db.AndCondition{
			Conditions: []db.Condition{
				db.NewWhere(db.Column_TransactionId, db.Equal, strconv.FormatInt(transactionId, 10)),
				db.NewWhere(db.Column_TagId, db.Equal, strconv.FormatInt(tagId, 10)),
			},
		}).
		Delete(ctx)
	if err != nil {
		s.repo.Rollback()
		return fmt.Errorf(errorMsg, transactionId, tagId, err)
	}

	if err := s.repo.Commit(); err != nil {
		s.repo.Rollback()
		return fmt.Errorf(errorMsg, transactionId, tagId, err)
	}

	return nil
}

func (s *Service) getTags(ctx context.Context, transactionId int64) ([]string, error) {
	tag_arr, err := s.repo.TagQueryBuilder().
		SetCondition(
			db.NewWhere(db.Column_TransactionId, db.Equal, strconv.FormatInt(transactionId, 10)),
		).
		AddJoin(db.Table_TransactionTags, "tag_id = id").
		Select(ctx, -1, -1)
	if err != nil {
		return nil, fmt.Errorf("Failed to get tags for transaction %d: %w", transactionId, err)
	}
	names := make([]string, len(tag_arr))
	for i, t := range tag_arr {
		names[i] = t.Name
	}
	return names, nil
}
