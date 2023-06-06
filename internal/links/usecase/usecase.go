package usecase

import (
	"context"
	"project/internal/links"
	"project/internal/model"
	pkgErrors "project/internal/pkg/errors"
	"project/internal/pkg/hash"
	"project/internal/pkg/validation"
)

type usecase struct {
	linksRepo links.Repository
}

func NewLinksUsecase(linksRepo links.Repository) links.Usecase {
	return usecase{
		linksRepo: linksRepo,
	}
}

func (u usecase) GetOriginalLink(ctx context.Context, abbreviatedUrl string) (model.Link, error) {
	err := validation.ValidateAbbreviatedUrl(abbreviatedUrl)
	if err != nil {
		return model.Link{}, err
	}

	originalUrl, err := u.linksRepo.GetOriginalUrl(ctx, abbreviatedUrl)
	if err != nil {
		return model.Link{}, err
	}

	return model.Link{
		Url: originalUrl,
	}, nil
}

func (u usecase) SaveAbbreviatedLink(ctx context.Context, originalUrl string) (model.Link, error) {
	err := validation.ValidateOriginalUrl(originalUrl)
	if err != nil {
		return model.Link{}, err
	}

	abbreviatedUrl := hash.Hash(originalUrl)

	err = u.linksRepo.CheckExistAbbreviatedLink(ctx, abbreviatedUrl)
	if err == nil {
		return model.Link{
			Url: abbreviatedUrl,
		}, nil
	}

	if err != pkgErrors.ErrAbbreviatedUrlNotFound {
		return model.Link{}, err
	}

	err = u.linksRepo.SaveAbbreviatedLink(ctx, model.LinkDB{
		OriginalUrl:    originalUrl,
		AbbreviatedUrl: abbreviatedUrl,
	})
	if err != nil {
		return model.Link{}, err
	}

	return model.Link{
		Url: abbreviatedUrl,
	}, nil
}
