package redis

import (
	"context"
	"github.com/go-redis/redismock/v9"
	"github.com/stretchr/testify/require"
	"project/internal/model"
	pkgErrors "project/internal/pkg/errors"
	"testing"
)

func TestRedis_GetOriginalUrl_OK(t *testing.T) {
	excpectedOriginalUrl := "vk.com"
	abbreviatedUrl := "AF4_ResRbL"

	mockedClient, mock := redismock.NewClientMock()
	repo := NewLinksRepository(mockedClient)

	mock.ExpectGet(abbreviatedUrl).SetVal(excpectedOriginalUrl)

	originalUrl, err := repo.GetOriginalUrl(context.TODO(), abbreviatedUrl)

	require.NoError(t, err)
	require.Equal(t, excpectedOriginalUrl, originalUrl)
}

func TestRedis_GetOriginalUrl_NotFound(t *testing.T) {
	excpectedOriginalUrl := "vk.com"
	abbreviatedUrl := "AF4_ResRbL"
	notFoundAbbreviatedUrl := "AF5_RSSSSS"

	mockedClient, mock := redismock.NewClientMock()
	repo := NewLinksRepository(mockedClient)

	mock.ExpectGet(abbreviatedUrl).SetVal(excpectedOriginalUrl)

	_, err := repo.GetOriginalUrl(context.TODO(), notFoundAbbreviatedUrl)

	require.Error(t, err, pkgErrors.ErrAbbreviatedUrlNotFound)
}

func TestRedis_SaveAbbreviatedLink_OK(t *testing.T) {
	link := model.LinkDB{
		OriginalUrl:    "vk.com",
		AbbreviatedUrl: "AF4_ResRbL",
	}

	mockedClient, mock := redismock.NewClientMock()
	repo := NewLinksRepository(mockedClient)

	mock.ExpectSet(link.AbbreviatedUrl, link.OriginalUrl, 0).SetVal(link.OriginalUrl)

	err := repo.SaveAbbreviatedLink(context.TODO(), link)

	require.NoError(t, err)
}

func TestRedis_CheckExistAbbreviatedLink_OK(t *testing.T) {
	excpectedOriginalUrl := "vk.com"
	abbreviatedUrl := "AF4_ResRbL"

	mockedClient, mock := redismock.NewClientMock()
	repo := NewLinksRepository(mockedClient)

	mock.ExpectGet(abbreviatedUrl).SetVal(excpectedOriginalUrl)

	err := repo.CheckExistAbbreviatedLink(context.TODO(), abbreviatedUrl)

	require.NoError(t, err)
}

func TestRedis_CheckExistAbbreviatedLink_NotFound(t *testing.T) {
	excpectedOriginalUrl := "vk.com"
	abbreviatedUrl := "AF4_ResRbL"
	notFoundAbbreviatedUrl := "AF5_RSSSSS"

	mockedClient, mock := redismock.NewClientMock()
	repo := NewLinksRepository(mockedClient)

	mock.ExpectGet(abbreviatedUrl).SetVal(excpectedOriginalUrl)

	err := repo.CheckExistAbbreviatedLink(context.TODO(), notFoundAbbreviatedUrl)

	require.Error(t, err, pkgErrors.ErrAbbreviatedUrlNotFound)
}
