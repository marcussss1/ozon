package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"project/internal/links"
	"project/internal/model"
	pkgErrors "project/internal/pkg/errors"
)

func NewLinksRepository(db *redis.Client) links.Repository {
	return &repository{db: db}
}

type repository struct {
	db *redis.Client
}

func (r repository) GetOriginalUrl(ctx context.Context, abbreviatedUrl string) (string, error) {
	originalUrl, err := r.db.Get(ctx, abbreviatedUrl).Result()
	if err != nil {
		if err == redis.Nil {
			return "", pkgErrors.ErrAbbreviatedUrlNotFound
		}

		return "", err
	}

	return originalUrl, nil
}

func (r repository) SaveAbbreviatedLink(ctx context.Context, link model.LinkDB) error {
	err := r.db.Set(ctx, link.AbbreviatedUrl, link.OriginalUrl, 0).Err()
	return err
}

func (r repository) CheckExistAbbreviatedLink(ctx context.Context, abbreviatedUrl string) error {
	_, err := r.db.Get(ctx, abbreviatedUrl).Result()
	if err != nil {
		if err == redis.Nil {
			return pkgErrors.ErrAbbreviatedUrlNotFound
		}

		return err
	}

	return nil
}
