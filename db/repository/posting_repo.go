package repository

import (
	"context"

	"crispy/domain"
)

type PostingRepository interface {
	CreatePosting(ctx context.Context, p *domain.Posting) (int64, error)
	GetPostingById(ctx context.Context, id int64) (*domain.Posting, error)
	GetPostingsByTransactionId(ctx context.Context, id int64) ([]*domain.Posting, error)
	UpdatePosting(ctx context.Context, a *domain.Posting) error
	DeletePosting(ctx context.Context, id int64) error
}
