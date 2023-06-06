package links

import (
	"context"
	"project/internal/model"
)

type Repository interface {
	GetOriginalUrl(ctx context.Context, abbreviatedUrl string) (string, error)
	SaveAbbreviatedLink(ctx context.Context, link model.LinkDB) error
	CheckExistAbbreviatedLink(ctx context.Context, abbreviatedUrl string) error
}
