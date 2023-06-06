package postgres

import (
	"context"
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
	"project/internal/links"
	"project/internal/model"
	pkgErrors "project/internal/pkg/errors"
)

type repository struct {
	db *sqlx.DB
}

func NewLinksRepository(db *sqlx.DB) links.Repository {
	return &repository{
		db: db,
	}
}

func (r repository) GetOriginalUrl(ctx context.Context, abbreviatedUrl string) (string, error) {
	var link model.LinkDB

	err := r.db.GetContext(ctx, &link, `SELECT * FROM link WHERE abbreviated_url=$1`, abbreviatedUrl)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", pkgErrors.ErrAbbreviatedUrlNotFound
		}

		return "", err
	}

	return link.OriginalUrl, nil
}

func (r repository) SaveAbbreviatedLink(ctx context.Context, link model.LinkDB) error {
	_, err := r.db.NamedExecContext(ctx, `INSERT INTO link(original_url, abbreviated_url) VALUES(:original_url, :abbreviated_url)`, link)
	if err != nil {
		return err
	}

	return nil
}

func (r repository) CheckExistAbbreviatedLink(ctx context.Context, abbreviatedUrl string) error {
	var exists bool

	err := r.db.GetContext(ctx, &exists, `SELECT EXISTS(SELECT 1 FROM link WHERE abbreviated_url=$1)`, abbreviatedUrl)
	if err != nil {
		return err
	}

	if !exists {
		return pkgErrors.ErrAbbreviatedUrlNotFound
	}

	return nil
}
