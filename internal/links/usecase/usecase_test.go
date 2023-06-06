package usecase

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	mockLinksRepository "project/internal/links/repository/mocks"
	"project/internal/model"
	pkgErrors "project/internal/pkg/errors"
	"testing"
)

type testGetCase struct {
	abbreviatedUrl string
	originalUrl    string
	expectedLink   model.Link
	expectedError  error
	name           string
}

func Test_GetOriginalLink(t *testing.T) {
	tests := []testGetCase{
		{
			abbreviatedUrl: "AF4_ResRbL",
			originalUrl:    "vk.com",
			expectedLink: model.Link{
				Url: "vk.com",
			},
			expectedError: nil,
			name:          "Успешно возвращена оригинальная ссылка",
		},
		{
			abbreviatedUrl: "AF4_ResRbL111",
			originalUrl:    "vk.com",
			expectedLink:   model.Link{},
			expectedError:  pkgErrors.ErrAbbreviatedUrlIsTooLong,
			name:           "Слишком длинная ссылка(больше 10 символов)",
		},
		{
			abbreviatedUrl: "AF4_ResRbL111",
			originalUrl:    "vk.com",
			expectedLink:   model.Link{},
			expectedError:  pkgErrors.ErrAbbreviatedUrlNotFound,
			name:           "Ссылка не найдена",
		},
		{
			abbreviatedUrl: "AF4_ResRbL111",
			originalUrl:    "vk.com",
			expectedLink:   model.Link{},
			expectedError:  pkgErrors.ErrInternal,
			name:           "Внутренняя ошибка сервера",
		},
	}

	ctl := gomock.NewController(t)
	defer ctl.Finish()

	linksRepository := mockLinksRepository.NewMockRepository(ctl)
	linksUsecase := NewLinksUsecase(linksRepository)

	for _, test := range tests {
		linksRepository.EXPECT().GetOriginalUrl(context.TODO(), test.abbreviatedUrl).Return(test.originalUrl, nil).AnyTimes()
		link, err := linksUsecase.GetOriginalLink(context.TODO(), test.abbreviatedUrl)

		if test.expectedError == nil {
			require.NoError(t, err)
		} else {
			require.Error(t, test.expectedError, err)
		}

		require.Equal(t, test.expectedLink, link)
	}
}

type testSaveCase struct {
	abbreviatedUrl     string
	originalUrl        string
	expectedLink       model.Link
	expectedExistError error
	expectedSaveError  error
	expectedError      error
	existCall          int
	saveCall           int
	name               string
}

func Test_SaveAbbreviatedLink(t *testing.T) {
	tests := []testSaveCase{
		{
			abbreviatedUrl: "OQ4g14JATq",
			originalUrl:    "vk.com",
			expectedLink: model.Link{
				Url: "OQ4g14JATq",
			},
			expectedExistError: pkgErrors.ErrAbbreviatedUrlNotFound,
			expectedSaveError:  nil,
			expectedError:      nil,
			existCall:          1,
			saveCall:           1,
			name:               "Успешно сохранена оригинальная ссылка",
		},
		{
			abbreviatedUrl:     "AF4_ResRbL111",
			originalUrl:        "vk.commmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmm",
			expectedLink:       model.Link{},
			expectedExistError: nil,
			expectedSaveError:  nil,
			expectedError:      pkgErrors.ErrUrlIsTooLong,
			existCall:          0,
			saveCall:           0,
			name:               "Слишком длинная ссылка(больше 512 символов)",
		},
		{
			abbreviatedUrl: "OQ4g14JATq",
			originalUrl:    "vk.com",
			expectedLink: model.Link{
				Url: "OQ4g14JATq",
			},
			expectedExistError: nil,
			expectedSaveError:  nil,
			expectedError:      nil,
			existCall:          1,
			saveCall:           0,
			name:               "Такая ссылка уже существует",
		},
		{
			abbreviatedUrl:     "OQ4g14JATq",
			originalUrl:        "vk.com",
			expectedLink:       model.Link{},
			expectedExistError: pkgErrors.ErrAbbreviatedUrlNotFound,
			expectedSaveError:  pkgErrors.ErrInternal,
			expectedError:      pkgErrors.ErrInternal,
			existCall:          1,
			saveCall:           1,
			name:               "Внутренняя ошибка сервера",
		},
		{
			abbreviatedUrl:     "OQ4g14JATq",
			originalUrl:        "vk.com",
			expectedLink:       model.Link{},
			expectedExistError: pkgErrors.ErrInternal,
			expectedSaveError:  nil,
			expectedError:      pkgErrors.ErrInternal,
			existCall:          1,
			saveCall:           0,
			name:               "Внутренняя ошибка сервера",
		},
	}

	ctl := gomock.NewController(t)
	defer ctl.Finish()

	linksRepository := mockLinksRepository.NewMockRepository(ctl)
	linksUsecase := NewLinksUsecase(linksRepository)

	for _, test := range tests {
		linksRepository.EXPECT().CheckExistAbbreviatedLink(context.TODO(), test.abbreviatedUrl).Return(test.expectedExistError).Times(test.existCall)

		linksRepository.EXPECT().SaveAbbreviatedLink(context.TODO(), model.LinkDB{
			OriginalUrl:    test.originalUrl,
			AbbreviatedUrl: test.abbreviatedUrl,
		}).Return(test.expectedSaveError).Times(test.saveCall)

		link, err := linksUsecase.SaveAbbreviatedLink(context.TODO(), test.originalUrl)
		if test.expectedError == nil {
			require.NoError(t, err, test.name)
		} else {
			require.Error(t, test.expectedError, err, test.name)
		}

		require.Equal(t, test.expectedLink, link, test.name)
	}
}
