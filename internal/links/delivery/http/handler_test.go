package http

import (
	"bytes"
	"context"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	mockLinksUsecase "project/internal/links/usecase/mocks"
	"project/internal/model"
	pkgErrors "project/internal/pkg/errors"
	"testing"
)

type getTestCase struct {
	name           string
	abbreviatedUrl string
	usecaseLink    model.Link
	usecaseError   error
	status         int
}

func TestHandlers_GetAbbreviatedUrl(t *testing.T) {
	tests := []getTestCase{
		{
			name:           "Успешное получение ссылки",
			abbreviatedUrl: "AF4_ResRbL",
			usecaseLink: model.Link{
				Url: "vk.com",
			},
			usecaseError: nil,
			status:       http.StatusOK,
		},
		{
			name:           "Такой ссылки не существует",
			abbreviatedUrl: "AF4_ResRbL",
			usecaseLink: model.Link{
				Url: "vk.com",
			},
			usecaseError: pkgErrors.ErrAbbreviatedUrlNotFound,
			status:       http.StatusNotFound,
		},
		{
			name:           "Введена сокращенная ссылка больше 10 символов",
			abbreviatedUrl: "AF4_ResRbL",
			usecaseLink: model.Link{
				Url: "vk.com",
			},
			usecaseError: pkgErrors.ErrAbbreviatedUrlIsTooLong,
			status:       http.StatusBadRequest,
		},
		{
			name:           "Внутренняя ошибка сервера",
			abbreviatedUrl: "AF4_ResRbL",
			usecaseLink: model.Link{
				Url: "vk.com",
			},
			usecaseError: pkgErrors.ErrInternal,
			status:       http.StatusInternalServerError,
		},
	}

	ctl := gomock.NewController(t)
	defer ctl.Finish()

	e := echo.New()

	linksUsecase := mockLinksUsecase.NewMockUsecase(ctl)
	handler := NewLinksHandler(e, linksUsecase)

	for _, test := range tests {
		url := "/api/v1/get/" + test.abbreviatedUrl

		linksUsecase.EXPECT().GetOriginalLink(context.TODO(), test.abbreviatedUrl).Return(test.usecaseLink, test.usecaseError).Times(1)

		r := httptest.NewRequest("GET", url, nil)

		w := httptest.NewRecorder()

		ctx := e.NewContext(r, w)
		ctx.SetParamNames("url")
		ctx.SetParamValues(test.abbreviatedUrl)

		err := handler.GetAbbreviatedUrlHandler(ctx)

		if err == nil {
			require.NoError(t, test.usecaseError)
		} else {
			require.Error(t, err, test.usecaseError)
			continue
		}

		require.Equal(t, test.status, w.Code)
	}
}

type saveTestCase struct {
	name            string
	originalUrlJson []byte
	originalUrl     string
	usecaseLink     model.Link
	usecaseError    error
	status          int
}

func TestHandlers_SaveOriginalLink(t *testing.T) {
	tests := []saveTestCase{
		{
			name:            "Успешное сохранение ссылки",
			originalUrlJson: []byte(`{"url":"vk.com"}`),
			originalUrl:     "vk.com",
			usecaseLink: model.Link{
				Url: "AF4_ResRbL",
			},
			usecaseError: nil,
			status:       http.StatusCreated,
		},
		{
			name:            "Введена ссылка больше 512 символов ",
			originalUrlJson: []byte(`{"url":"vk.commmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmm..."}`),
			originalUrl:     "vk.commmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmmm...",
			usecaseLink: model.Link{
				Url: "AF4_ResRbL",
			},
			usecaseError: pkgErrors.ErrUrlIsTooLong,
			status:       http.StatusBadRequest,
		},
		{
			name:            "Внутренняя ошибка сервера",
			originalUrlJson: []byte(`{"url":"vk.com"}`),
			originalUrl:     "vk.com",
			usecaseLink: model.Link{
				Url: "AF4_ResRbL",
			},
			usecaseError: pkgErrors.ErrInternal,
			status:       http.StatusInternalServerError,
		},
		{
			name:            "Битое тело",
			originalUrlJson: []byte(`s}{d}sff{"url":"vk.com"SD{{{}}SD`),
			originalUrl:     "vk.com",
			usecaseLink: model.Link{
				Url: "AF4_ResRbL",
			},
			usecaseError: pkgErrors.ErrInternal,
			status:       http.StatusInternalServerError,
		},
	}

	ctl := gomock.NewController(t)
	defer ctl.Finish()

	e := echo.New()

	linksUsecase := mockLinksUsecase.NewMockUsecase(ctl)
	handler := NewLinksHandler(e, linksUsecase)

	for _, test := range tests {
		url := "/api/v1/save"

		if test.name != "Битое тело" {
			linksUsecase.EXPECT().SaveAbbreviatedLink(context.TODO(), test.originalUrl).Return(test.usecaseLink, test.usecaseError).Times(1)
		}

		r := httptest.NewRequest("POST", url, bytes.NewReader(test.originalUrlJson))
		r.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		w := httptest.NewRecorder()

		ctx := e.NewContext(r, w)

		err := handler.SaveOriginalLinkHandler(ctx)

		if err == nil {
			require.NoError(t, test.usecaseError)
		} else {
			require.Error(t, err, test.usecaseError)
			continue
		}

		require.Equal(t, test.status, w.Code)
	}
}
