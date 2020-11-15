package postgres_repos

import (
	"context"
	"fmt"

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

// CreateMany persists new Attachments into the database
func (r AttachmentRepository) CreateMany(ctx context.Context, attachments []models.Attachment, feedbackId uint) ([]models.Attachment, error) {
	q := r.SB.Insert("attachments").Columns("name", "path", "feedback_id")

	for _, a := range attachments {
		q = q.Values(a.Name, a.Path, feedbackId)
	}
	q = q.Suffix("RETURNING id, name, path, feedback_id")

	sql, args, err := q.ToSql()
	fmt.Println(sql, args)
	if err != nil {
		return nil, err
	}

	rows, err := r.DB.Query(ctx, sql, args...)

	if err != nil {
		return nil, err
	}

	got := []models.Attachment{}
	for rows.Next() {
		var r models.Attachment
		err = rows.Scan(&r.ID, &r.Name, &r.Path, &r.FeedbackID)
		if err != nil {
			return nil, err
		}
		got = append(got, r)
	}

	return got, nil
}

func (r AttachmentRepository) GetAll(ctx context.Context) (attchs []models.Attachment, err error) {
	q := r.SB.
		Select("id", "name", "path", "feedback_id").
		From("attachments")

	sql, args, err := q.ToSql()
	fmt.Println(sql, args)
	if err != nil {
		return nil, err
	}
	rows, err := r.DB.Query(ctx, sql, args...)
	for rows.Next() {
		var a models.Attachment
		err = rows.Scan(&a.ID, &a.Name, &a.Path, &a.FeedbackID)
		if err != nil {
			return
		}
		attchs = append(attchs, a)
	}
	return attchs, err
}
