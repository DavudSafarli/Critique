package postgres_repos

import (
	"context"

	"github.com/DavudSafarli/Critique/domain/models"

	sq "github.com/Masterminds/squirrel"
)

// TagRepository is TagRepository
type TagRepository struct {
	storage *Storage
}

// NewPGTagRepository ..
func NewPGTagRepository(connstr string) TagRepository {
	storage, err := NewDbConnection(connstr)
	if err != nil {
		panic("db could not be initialized")
	}
	return TagRepository{storage: storage}
}

// CreateMany persists new Tags into the database
func (r TagRepository) CreateMany(ctx context.Context, tags []models.Tag) ([]models.Tag, error) {
	q := r.storage.SB.Insert("tags").Columns("name")

	for _, tag := range tags {
		q = q.Values(tag.Name)
	}
	q = q.Suffix("RETURNING id, name")

	sql, args, err := q.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.storage.DB.Query(ctx, sql, args...)

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

// Get returns all Tags
func (r TagRepository) Get(ctx context.Context) ([]models.Tag, error) {
	q := r.storage.SB.Select("*").From("tags")

	sql, args, err := q.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.storage.DB.Query(ctx, sql, args...)
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
