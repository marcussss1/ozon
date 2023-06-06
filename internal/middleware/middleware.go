package middleware

import (
	"encoding/json"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"project/internal/pkg/errors"
	httpUtils "project/internal/pkg/http_utils"
)

type jsonError struct {
	Err error `json:"error"`
}

func (j jsonError) MarshalJSON() ([]byte, error) {
	return json.Marshal(j.Err.Error())
}

func LoggerMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		requestId := rand.Int63()
		log.Info("Incoming request: ", ctx.Request().URL, ", ip: ", ctx.RealIP(), ", method: ", ctx.Request().Method, ", request_id: ", requestId)

		if err := next(ctx); err != nil {
			apiErr := errors.FromGrpcErrorToOriginalError(err)
			statusCode := httpUtils.StatusCode(apiErr)

			log.Error("HTTP code: ", statusCode, ", Error: ", err, ", request_id: ", requestId)

			err = apiErr

			if statusCode == 500 {
				jsonErr, err := json.Marshal(jsonError{
					Err: errors.ErrInternal,
				})
				if err != nil {
					log.Error(err)
				}

				return ctx.JSONBlob(statusCode, jsonErr)
			}

			jsonErr, err := json.Marshal(jsonError{
				Err: err,
			})
			if err != nil {
				log.Error(err)
			}

			return ctx.JSONBlob(statusCode, jsonErr)
		}

		log.Info("HTTP code: ", ctx.Response().Status, ", request_id: ", requestId)

		return nil
	}
}
