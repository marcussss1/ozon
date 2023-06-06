package http

import (
	"context"
	"github.com/labstack/echo/v4"
	"net/http"
	"project/internal/links"
	"project/internal/model"
)

type linksHandler struct {
	linksUsecase links.Usecase
}

func (u linksHandler) GetAbbreviatedUrlHandler(ctx echo.Context) error {
	abbreviatedUrl := ctx.Param("url")

	originalLink, err := u.linksUsecase.GetOriginalLink(context.TODO(), abbreviatedUrl)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, originalLink)
}

func (u linksHandler) SaveOriginalLinkHandler(ctx echo.Context) error {
	var originalLink model.Link

	err := ctx.Bind(&originalLink)
	if err != nil {
		return err
	}

	abbreviatedLink, err := u.linksUsecase.SaveAbbreviatedLink(context.TODO(), originalLink.Url)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusCreated, abbreviatedLink)
}

func NewLinksHandler(e *echo.Echo, linksUsecase links.Usecase) linksHandler {
	handler := linksHandler{
		linksUsecase: linksUsecase,
	}

	getUrl := "/get/:url"
	saveUrl := "/save"

	api := e.Group("api/v1")

	get := api.Group(getUrl)
	save := api.Group(saveUrl)

	get.GET("", handler.GetAbbreviatedUrlHandler)
	save.POST("", handler.SaveOriginalLinkHandler)

	return handler
}
