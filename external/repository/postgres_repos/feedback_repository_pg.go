package postgres_repos

import (
	"context"

	"github.com/jackc/pgx/v4"

	sq "github.com/Masterminds/squirrel"

	"github.com/DavudSafarli/Critique/domain/models"
)

// FeedbackRepository is FeedbackRepository
type FeedbackRepository struct {
	database database
}

// NewPGFeedbackRepository ..
func NewPGFeedbackRepository(database database) FeedbackRepository {
	return FeedbackRepository{database: database}
}

// GetFeedbacksPaginated returns records with pagination
func (r FeedbackRepository) GetFeedbacksPaginated(ctx context.Context, skip uint, limit uint) ([]models.Feedback, error) {
	db := getDB(ctx, r.database.DB)
	q := r.database.SB.
		Select("id", "title", "body", "created_by", "extract(epoch from created_at) created_at").
		From("feedbacks").
		Offset(uint64(skip)).
		Limit(uint64(limit))

	sql, args, err := q.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := db.Query(ctx, sql, args...)
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
func (r FeedbackRepository) FindFeedback(ctx context.Context, id uint) (f models.Feedback, err error) {
	db := getDB(ctx, r.database.DB)
	q := r.database.SB.
		Select("id", "title", "body", "created_by", "extract(epoch from created_at) created_at").
		From("feedbacks").
		Where(sq.Eq{"id": id})

	sql, args, err := q.ToSql()
	if err != nil {
		return f, err
	}

	err = db.QueryRow(ctx, sql, args...).
		Scan(&f.ID, &f.Title, &f.Body, &f.CreatedBy, &f.CreatedAt)

	if err != nil {
		return f, convertPgErrorToDomainError(err)
	}

	q2 := r.database.SB.
		Select("id", "name", "path", "feedback_id").
		From("attachments").
		Where(sq.Eq{"feedback_id": id})

	sql, args, err = q2.ToSql()
	if err != nil {
		return f, err
	}
	rows, err := db.Query(ctx, sql, args...)
	for rows.Next() {
		var a models.Attachment
		err = rows.Scan(&a.ID, &a.Name, &a.Path, &a.FeedbackID)
		if err != nil {
			return
		}
		f.Attachments = append(f.Attachments, a)
	}
	if err == pgx.ErrNoRows {
		err = nil
	}

	return f, err
}

// Create persists a new Feedback to the database and returns newly inserted Feedback
func (r FeedbackRepository) CreateFeedback(ctx context.Context, feedback *models.Feedback) (err error) {
	db := getDB(ctx, r.database.DB)
	q := r.database.SB.Insert("feedbacks").
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
		return err
	}

	err = db.QueryRow(ctx, sql, args...).
		Scan(&feedback.ID, &feedback.Title, &feedback.Body, &feedback.CreatedBy, &feedback.CreatedAt)

	return err
}

// UpdateTagIDs just panics right now, but will update "tag_id"s of feedbacks from x to y.
func (r FeedbackRepository) UpdateTagIDs(ctx context.Context, tagIDFrom uint, tagIDTo uint) error {
	panic("implement me")
}

func (r FeedbackRepository) GetAll(ctx context.Context) ([]models.Feedback, error) {
	db := getDB(ctx, r.database.DB)
	q := r.database.SB.Select("id", "title", "body", "created_by", "extract(epoch from created_at) created_at").From("feedbacks")
	sql, _, err := q.ToSql()
	if err != nil {
		return nil, err
	}
	rows, err := db.Query(ctx, sql)
	if err != nil {
		return nil, err
	}
	got, err := r.scan(rows)
	if err != nil {
		return nil, err
	}
	return got, err
}

func (r FeedbackRepository) scan(rows pgx.Rows) (got []models.Feedback, err error) {
	got = []models.Feedback{}
	for rows.Next() {
		var f models.Feedback
		err = rows.Scan(&f.ID, &f.Title, &f.Body, &f.CreatedBy, &f.CreatedAt)
		if err != nil {
			return nil, err
		}
		got = append(got, f)
	}
	if err := rows.Err(); err != nil {
		return got, err
	}
	return
}
