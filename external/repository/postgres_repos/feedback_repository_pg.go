package postgres_repos

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"

	"github.com/DavudSafarli/Critique/domain/models"
)

// FeedbackRepository is FeedbackRepository
type FeedbackRepository struct {
	*Storage
}

// NewPGFeedbackRepository ..
func NewPGFeedbackRepository(connstr string) FeedbackRepository {
	storage, err := NewSingletonDbConnection(connstr)
	if err != nil {
		panic("db could not be initialized")
	}
	return FeedbackRepository{storage}
}

// GetPaginated returns records with pagination
func (r FeedbackRepository) GetPaginated(ctx context.Context, skip uint, limit uint) ([]models.Feedback, error) {
	q := r.SB.
		Select("id", "title", "body", "created_by", "extract(epoch from created_at) created_at").
		From("feedbacks").
		Offset(uint64(skip)).
		Limit(uint64(limit))

	sql, args, err := q.ToSql()
	fmt.Println(sql, args)
	if err != nil {
		return nil, err
	}

	rows, err := r.DB.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}

	got := make([]models.Feedback, 0, limit)
	for rows.Next() {
		var f models.Feedback
		err = rows.Scan(&f.ID, &f.Title, &f.Body, &f.CreatedBy, &f.CreatedAt)
		if err != nil {
			return nil, err
		}
		got = append(got, f)
	}

	return got, nil
}

// Find finds and retrieves a single record with the given ID
func (r FeedbackRepository) Find(ctx context.Context, id uint) (f models.Feedback, err error) {
	q := r.SB.
		Select("id", "title", "body", "created_by", "extract(epoch from created_at) created_at").
		From("feedbacks").
		Where(sq.Eq{"id": id})

	sql, args, err := q.ToSql()
	fmt.Println(sql, args)
	if err != nil {
		return f, err
	}

	err = r.DB.QueryRow(ctx, sql, args...).
		Scan(&f.ID, &f.Title, &f.Body, &f.CreatedBy, &f.CreatedAt)

	q2 := r.SB.
		Select("id", "name", "path", "feedback_id").
		From("attachments").
		Where(sq.Eq{"feedback_id": id})

	sql, args, err = q2.ToSql()
	fmt.Println(sql, args)
	if err != nil {
		return f, err
	}
	rows, err := r.DB.Query(ctx, sql, args...)
	for rows.Next() {
		var a models.Attachment
		err = rows.Scan(&a.ID, &a.Name, &a.Path, &a.FeedbackID)
		if err != nil {
			return
		}
		f.Attachments = append(f.Attachments, a)
	}
	return f, err
}

// Create persists a new Feedback to the database and returns newly inserted Feedback
func (r FeedbackRepository) Create(ctx context.Context, feedback models.Feedback) (f models.Feedback, err error) {
	q := r.SB.Insert("feedbacks").
		Columns("title", "body", "created_by", "created_at").
		SetMap(map[string]interface{}{
			"title":      feedback.Title,
			"body":       feedback.Body,
			"created_by": feedback.CreatedBy,
			"created_at": sq.Expr("TO_TIMESTAMP(?)", feedback.CreatedAt),
		})

	q = q.Suffix("RETURNING id, title, body, created_by, extract(epoch from created_at) created_at")

	sql, args, err := q.ToSql()
	fmt.Println(sql, args)
	if err != nil {
		return f, err
	}

	err = r.DB.QueryRow(ctx, sql, args...).
		Scan(&f.ID, &f.Title, &f.Body, &f.CreatedBy, &f.CreatedAt)

	return f, err
}

// UpdateTagIDs just panics right now, but will update "tag_id"s of feedbacks from x to y.
func (r FeedbackRepository) UpdateTagIDs(ctx context.Context, tagIDFrom uint, tagIDTo uint) error {
	panic("implement me")
}
