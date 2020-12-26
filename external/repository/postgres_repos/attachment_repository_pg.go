package postgres_repos

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"

	"github.com/DavudSafarli/Critique/domain/models"
)

// AttachmentRepository is AttachmentRepository
type AttachmentRepository struct {
	database database
}

// NewPGAttachmentRepository ..
func NewPGAttachmentRepository(database database) AttachmentRepository {
	return AttachmentRepository{database}
}

func (r AttachmentRepository) GetAttachmentsByFeedbackID(ctx context.Context, feedbackID uint) ([]models.Attachment, error) {
	db := getDB(ctx, r.database.DB)
	q := r.database.SB.
		Select("id", "name", "path", "feedback_id").
		From("attachments").
		Where(sq.Eq{"feedback_id": feedbackID})
	sql, args, err := q.ToSql()
	rows, err := db.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	got, err := r.scan(rows)
	if err != nil {
		return nil, err
	}
	//r.close(db)
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
	if err := rows.Err(); err != nil {
		return got, err
	}
	return
}

// CreateMany persists new Attachments into the database
func (r AttachmentRepository) CreateManyAttachments(ctx context.Context, attachments []models.Attachment, feedbackID uint) error {
	db := getDB(ctx, r.database.DB)
	q := r.database.SB.Insert("attachments").Columns("name", "path", "feedback_id")

	for _, a := range attachments {
		q = q.Values(a.Name, a.Path, feedbackID)
	}
	q = q.Suffix("RETURNING id, name, path, feedback_id")

	sql, args, err := q.ToSql()
	if err != nil {
		return err
	}

	rows, err := db.Query(ctx, sql, args...)

	if err != nil {
		return err
	}

	got, err := r.scan(rows)
	if err != nil {
		return err
	}
	copy(attachments, got)
	//r.close(db)
	return nil

}

func (r AttachmentRepository) GetAllAttachments(ctx context.Context) ([]models.Attachment, error) {
	db := getDB(ctx, r.database.DB)
	q := r.database.SB.
		Select("id", "name", "path", "feedback_id").
		From("attachments")

	sql, args, err := q.ToSql()
	if err != nil {
		return nil, err
	}
	rows, err := db.Query(ctx, sql, args...)

	got, err := r.scan(rows)
	if err != nil {
		return nil, err
	}
	//r.close(db)
	return got, err
}
