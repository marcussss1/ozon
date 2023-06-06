package links

import (
	"context"
	"project/internal/model"
)

type Usecase interface {
	GetOriginalLink(ctx context.Context, abbreviatedUrl string) (model.Link, error)
	SaveAbbreviatedLink(ctx context.Context, originalUrl string) (model.Link, error)
}
