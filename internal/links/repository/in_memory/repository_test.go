package in_memory

import (
	"context"
	"errors"
	"github.com/stretchr/testify/require"
	"project/internal/model"
	pkgErrors "project/internal/pkg/errors"
	"testing"
	"time"
)

type testCase struct {
	abbreviatedUrl      string
	expectedOriginalUrl string
	expectedError       error
	name                string
}

func TestInMemory_GetOriginalUrl_OK(t *testing.T) {
	test := testCase{
		abbreviatedUrl:      "877ca9b8dd",
		expectedOriginalUrl: "vk.com",
		expectedError:       nil,
		name:                "Оригинальная ссылка успешно получена",
	}

	repo := NewLinksRepository()

	err := repo.SaveAbbreviatedLink(context.TODO(), model.LinkDB{
		AbbreviatedUrl: test.abbreviatedUrl,
		OriginalUrl:    test.expectedOriginalUrl,
	})
	require.NoError(t, err)

	originalUrl, err := repo.GetOriginalUrl(context.TODO(), test.abbreviatedUrl)

	require.NoError(t, err)
	require.Equal(t, test.expectedOriginalUrl, originalUrl)
}

func TestInMemory_GetOriginalUrl_NotFound(t *testing.T) {
	test := testCase{
		abbreviatedUrl:      "877ca9b8dd",
		expectedOriginalUrl: "",
		expectedError:       pkgErrors.ErrAbbreviatedUrlNotFound,
		name:                "Оригинальная ссылка не найдена",
	}

	repo := NewLinksRepository()

	originalUrl, err := repo.GetOriginalUrl(context.TODO(), test.abbreviatedUrl)

	require.Error(t, test.expectedError, err)
	require.Equal(t, test.expectedOriginalUrl, originalUrl)
}

func TestInMemory_GetOriginalUrl_ContextCancel(t *testing.T) {
	test := testCase{
		abbreviatedUrl:      "877ca9b8dd",
		expectedOriginalUrl: "",
		expectedError:       errors.New("Функция завершилась по контексту"),
		name:                "Функция завершилась по контексту",
	}

	repo := NewLinksRepository()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	cancel()

	originalUrl, err := repo.GetOriginalUrl(ctx, test.abbreviatedUrl)

	require.Error(t, test.expectedError, err)
	require.Equal(t, test.expectedOriginalUrl, originalUrl)
}

func TestInMemory_SaveAbbreviatedLink_OK(t *testing.T) {
	test := testCase{
		abbreviatedUrl:      "877ca9b8dd",
		expectedOriginalUrl: "vk.com",
		expectedError:       nil,
		name:                "Сокращенная ссылка успешно сохранена",
	}

	repo := NewLinksRepository()

	err := repo.SaveAbbreviatedLink(context.TODO(), model.LinkDB{
		AbbreviatedUrl: test.abbreviatedUrl,
		OriginalUrl:    test.expectedOriginalUrl,
	})

	require.NoError(t, err)
}

func TestInMemory_SaveAbbreviatedLink_ContextCancel(t *testing.T) {
	test := testCase{
		abbreviatedUrl:      "877ca9b8dd",
		expectedOriginalUrl: "vk.com",
		expectedError:       errors.New("Функция завершилась по контексту"),
		name:                "Функция завершилась по контексту",
	}

	repo := NewLinksRepository()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	cancel()

	err := repo.SaveAbbreviatedLink(ctx, model.LinkDB{
		AbbreviatedUrl: test.abbreviatedUrl,
		OriginalUrl:    test.expectedOriginalUrl,
	})

	require.Error(t, test.expectedError, err)
}

func TestInMemory_CheckExistAbbreviatedLink_OK(t *testing.T) {
	test := testCase{
		abbreviatedUrl:      "877ca9b8dd",
		expectedOriginalUrl: "vk.com",
		expectedError:       nil,
		name:                "Оригинальная ссылка по такой сокращённой ссылке существует",
	}

	repo := NewLinksRepository()

	err := repo.SaveAbbreviatedLink(context.TODO(), model.LinkDB{
		AbbreviatedUrl: test.abbreviatedUrl,
		OriginalUrl:    test.expectedOriginalUrl,
	})
	require.NoError(t, err)

	err = repo.CheckExistAbbreviatedLink(context.TODO(), test.abbreviatedUrl)

	require.NoError(t, err)
}

func TestInMemory_CheckExistAbbreviatedLink_NotFound(t *testing.T) {
	test := testCase{
		abbreviatedUrl:      "877ca9b8dd",
		expectedOriginalUrl: "",
		expectedError:       pkgErrors.ErrAbbreviatedUrlNotFound,
		name:                "Сокращённая ссылка не найдена",
	}

	repo := NewLinksRepository()

	err := repo.CheckExistAbbreviatedLink(context.TODO(), test.abbreviatedUrl)

	require.Error(t, test.expectedError, err)
}

func TestInMemory_CheckExistAbbreviatedLink_ContextCancel(t *testing.T) {
	test := testCase{
		abbreviatedUrl:      "877ca9b8dd",
		expectedOriginalUrl: "",
		expectedError:       errors.New("Функция завершилась по контексту"),
		name:                "Функция завершилась по контексту",
	}

	repo := NewLinksRepository()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	cancel()

	err := repo.CheckExistAbbreviatedLink(ctx, test.abbreviatedUrl)

	require.Error(t, test.expectedError, err)
}
