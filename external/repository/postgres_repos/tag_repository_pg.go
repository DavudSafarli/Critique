package postgres_repos

import (
	"context"

	"github.com/DavudSafarli/Critique/domain/models"

	sq "github.com/Masterminds/squirrel"
)

// TagRepository is TagRepository
type TagRepository struct {
	*Storage
}

// NewPGTagRepository ..
func NewPGTagRepository(storage *Storage) TagRepository {
	return TagRepository{storage}
}

// CreateMany persists new Tags into the database
func (r TagRepository) CreateMany(ctx context.Context, tags []models.Tag) error {
	db := r.getDB(ctx)
	q := r.SB.Insert("tags").Columns("name")

	for _, tag := range tags {
		q = q.Values(tag.Name)
	}
	q = q.Suffix("RETURNING id, name")

	sql, args, err := q.ToSql()
	if err != nil {
		return err
	}

	rows, err := db.Query(ctx, sql, args...)

	got := []models.Tag{}
	for rows.Next() {
		var r models.Tag
		err = rows.Scan(&r.ID, &r.Name)
		if err != nil {
			return err
		}
		got = append(got, r)
	}
	if err = rows.Err(); err != nil {
		return err
	}
	copy(tags, got)
	return nil
}

// Get returns all Tags
func (r TagRepository) Get(ctx context.Context) ([]models.Tag, error) {
	db := r.getDB(ctx)
	q := r.SB.Select("*").From("tags")

	sql, args, err := q.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := db.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}

	got := []models.Tag{}
	for rows.Next() {
		var r models.Tag
		err = rows.Scan(&r.ID, &r.Name)
		if err != nil {
			return nil, err
		}
		got = append(got, r)
	}

	return got, nil
}

// RemoveMany removes Tags of given tagIDs from database
func (r TagRepository) RemoveMany(ctx context.Context, tagIDs []uint) error {
	db := r.getDB(ctx)
	q := r.SB.Delete("tags")

	q = q.Where(sq.Eq{"tags.id": tagIDs})
	sql, args, err := q.ToSql()
	if err != nil {
		return err
	}
	_, err = db.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}
	return nil
}
