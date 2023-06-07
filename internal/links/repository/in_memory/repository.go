package in_memory

import (
	"context"
	"errors"
	"project/internal/model"
	pkgErrors "project/internal/pkg/errors"
	"sync"
)

type repository struct {
	storage map[string]string // ключ - сокращенная ссылка, значение - оригинальная ссылка
	mu      sync.Mutex
}

func NewLinksRepository() *repository {
	return &repository{
		storage: make(map[string]string),
	}
}

func (r *repository) GetOriginalUrl(ctx context.Context, abbreviatedUrl string) (string, error) {
	select {
	case <-ctx.Done():
		return "", errors.New("Функция завершилась по контексту")
	default:
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	originalUrl, ok := r.storage[abbreviatedUrl]
	if !ok {
		return "", pkgErrors.ErrAbbreviatedUrlNotFound
	}

	return originalUrl, nil
}

func (r *repository) SaveAbbreviatedLink(ctx context.Context, link model.LinkDB) error {
	select {
	case <-ctx.Done():
		return errors.New("Функция завершилась по контексту")
	default:
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	r.storage[link.AbbreviatedUrl] = link.OriginalUrl

	return nil
}

func (r *repository) CheckExistAbbreviatedLink(ctx context.Context, abbreviatedUrl string) error {
	select {
	case <-ctx.Done():
		return errors.New("Функция завершилась по контексту")
	default:
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	_, ok := r.storage[abbreviatedUrl]
	if !ok {
		return pkgErrors.ErrAbbreviatedUrlNotFound
	}

	return nil
}
