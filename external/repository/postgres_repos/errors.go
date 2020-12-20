package postgres_repos

import (
	"github.com/DavudSafarli/Critique/domain"
	"github.com/jackc/pgx/v4"
)

func convertPgErrorToDomainError(err error) error {
	switch err {
	case nil:
		break
	case pgx.ErrNoRows:
		err = domain.ErrNotFound
	}
	return err
}
