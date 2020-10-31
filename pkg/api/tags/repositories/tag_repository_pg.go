package repositories

import (
	"context"

	sq "github.com/Masterminds/squirrel"

	"github.com/DavudSafarli/Critique/pkg/database/postgres"
	"github.com/DavudSafarli/Critique/pkg/domain"
)

// TagRepository is TagRepository
type TagRepository struct {
	storage *postgres.Storage
}

// CreateMany persists new Tags into the database
func (r *TagRepository) CreateMany(ctx context.Context, tags []domain.Tag) ([]domain.Tag, error) {
	q := r.storage.SB.Insert("tags")

	for _, tag := range tags {
		q = q.SetMap(map[string]interface{}{
			"name": tag.Name,
		})
	}
	q = q.Suffix("RETURNING id, name")

	sql, args, err := q.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.storage.DB.Query(ctx, sql, args...)

	got := []domain.Tag{}
	for rows.Next() {
		var r domain.Tag
		err = rows.Scan(&r.ID, &r.Name)
		if err != nil {
			return nil, err
		}
		got = append(got, r)
	}

	return got, nil
}

// Get returns all Tags
func (r *TagRepository) Get(ctx context.Context) ([]domain.Tag, error) {
	q := r.storage.SB.Select("*").From("tags")

	sql, args, err := q.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.storage.DB.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}

	got := []domain.Tag{}
	for rows.Next() {
		var r domain.Tag
		err = rows.Scan(&r.ID, &r.Name)
		if err != nil {
			return nil, err
		}
		got = append(got, r)
	}

	return got, nil
}

// RemoveMany removes Tags of given tagIDs from database
func (r *TagRepository) RemoveMany(ctx context.Context, tagIDs []uint) error {
	q := r.storage.SB.Delete("tags")

	q = q.Where(sq.Eq{"tags.id": tagIDs})
	sql, args, err := q.ToSql()
	if err != nil {
		return err
	}
	_, err = r.storage.DB.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}
	return nil
}

// RemoveAll removes Tags of given tagIDs from database
func (r *TagRepository) RemoveAll(ctx context.Context) error {
	q := r.storage.SB.Delete("tags")

	sql, args, err := q.ToSql()
	if err != nil {
		return err
	}
	_, err = r.storage.DB.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}
	return nil
}
