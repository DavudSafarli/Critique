package repositories

import (
	"context"

	sq "github.com/Masterminds/squirrel"

	"github.com/DavudSafarli/Critique/pkg/database/postgres"
	"github.com/DavudSafarli/Critique/pkg/domain"
)

// FeedbackRepository is FeedbackRepository
type FeedbackRepository struct {
	storage *postgres.Storage
}

// GetPaginated returns records with pagination
func (r FeedbackRepository) GetPaginated(ctx context.Context, skip uint, limit uint) ([]domain.Feedback, error) {
	q := r.storage.SB.
		Select("id", "title", "body", "created_by", "extract(epoch from created_at) created_at").
		From("feedbacks").
		Offset(uint64(skip)).
		Limit(uint64(limit))

	sql, args, err := q.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.storage.DB.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}

	got := make([]domain.Feedback, 0, limit)
	for rows.Next() {
		var f domain.Feedback
		err = rows.Scan(&f.ID, &f.Title, &f.Body, &f.CreatedBy, &f.CreatedAt)
		if err != nil {
			return nil, err
		}
		got = append(got, f)
	}

	return got, nil
}

// Find finds and retrieves a single record with the given ID
func (r FeedbackRepository) Find(ctx context.Context, id uint) (f domain.Feedback, err error) {
	q := r.storage.SB.
		Select("id", "title", "body", "created_by", "extract(epoch from created_at) created_at").
		From("feedbacks").
		Where(sq.Eq{"id": id})

	sql, args, err := q.ToSql()
	if err != nil {
		return f, err
	}

	err = r.storage.DB.QueryRow(ctx, sql, args...).
		Scan(&f.ID, &f.Title, &f.Body, &f.CreatedBy, &f.CreatedAt)

	return f, err
}

// Create persists a new Feedback to the database and returns newly inserted Feedback
func (r FeedbackRepository) Create(ctx context.Context, feedback domain.Feedback) (f domain.Feedback, err error) {
	q := r.storage.SB.Insert("feedbacks").
		Columns("title", "body", "created_by", "created_at").
		SetMap(map[string]interface{}{
			"title":      feedback.Title,
			"body":       feedback.Body,
			"created_by": feedback.CreatedBy,
			"created_at": sq.Expr("TO_TIMESTAMP(?)", feedback.CreatedAt),
		})

	q = q.Suffix("RETURNING id, title, body, created_by, extract(epoch from created_at) created_at")

	sql, args, err := q.ToSql()
	if err != nil {
		return f, err
	}

	err = r.storage.DB.QueryRow(ctx, sql, args...).
		Scan(&f.ID, &f.Title, &f.Body, &f.CreatedBy, &f.CreatedAt)

	return f, err
}

// UpdateTagIDs just panics right now, but will update "tag_id"s of feedbacks from x to y.
func (r FeedbackRepository) UpdateTagIDs(ctx context.Context, tagIDFrom uint, tagIDTo uint) error {
	panic("implement me")
}
