package http_utils

import (
	"errors"
	"net/http"
	pkgErrors "project/internal/pkg/errors"
)

func StatusCode(err error) int {
	switch {
	case errors.Is(err, pkgErrors.ErrUrlIsTooLong):
		return http.StatusBadRequest
	case errors.Is(err, pkgErrors.ErrAbbreviatedUrlIsTooLong):
		return http.StatusBadRequest
	case errors.Is(err, pkgErrors.ErrAbbreviatedUrlNotFound):
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}
