package postgres_repos

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"

	"github.com/DavudSafarli/Critique/domain/models"
)

// AttachmentRepository is AttachmentRepository
type AttachmentRepository struct {
	*Storage
}

// NewPGAttachmentRepository ..
func NewPGAttachmentRepository(connstr string) AttachmentRepository {
	storage, err := NewSingletonDbConnection(connstr)
	if err != nil {
		panic("db could not be initialized")
	}
	return AttachmentRepository{storage}
}

func (r AttachmentRepository) GetByFeedbackID(ctx context.Context, feedbackID uint) ([]models.Attachment, error) {
	db := r.getDB(ctx)
	q := r.SB.
		Select("id", "name", "path", "feedback_id").
		From("attachments").
		Where(sq.Eq{"feedback_id": feedbackID})
	sql, args, _ := q.ToSql()
	rows, err := db.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	got, err := r.scan(rows)
	if err != nil {
		return nil, err
	}
	return got, err
}
func (r AttachmentRepository) scan(rows pgx.Rows) (got []models.Attachment, err error) {
	got = []models.Attachment{}
	for rows.Next() {
		var r models.Attachment
		err = rows.Scan(&r.ID, &r.Name, &r.Path, &r.FeedbackID)
		if err != nil {
			return nil, err
		}
		got = append(got, r)
	}
	return
}

// CreateMany persists new Attachments into the database
func (r AttachmentRepository) CreateMany(ctx context.Context, attachments []models.Attachment, feedbackID uint) ([]models.Attachment, error) {
	db := r.getDB(ctx)
	q := r.SB.Insert("attachments").Columns("name", "path", "feedback_id")

	for _, a := range attachments {
		q = q.Values(a.Name, a.Path, feedbackID)
	}
	q = q.Suffix("RETURNING id, name, path, feedback_id")

	sql, args, err := q.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := db.Query(ctx, sql, args...)

	if err != nil {
		return nil, err
	}

	got, err := r.scan(rows)
	if err != nil {
		return nil, err
	}
	return got, err

	return got, nil
}

func (r AttachmentRepository) GetAll(ctx context.Context) ([]models.Attachment, error) {
	q := r.SB.
		Select("id", "name", "path", "feedback_id").
		From("attachments")

	sql, args, err := q.ToSql()
	if err != nil {
		return nil, err
	}
	rows, err := r.DB.Query(ctx, sql, args...)

	got, err := r.scan(rows)
	if err != nil {
		return nil, err
	}
	return got, err
}
