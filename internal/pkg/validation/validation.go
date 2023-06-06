package validation

import (
	"project/internal/pkg/errors"
)

func ValidateAbbreviatedUrl(url string) error {
	if len(url) > 10 {
		return errors.ErrAbbreviatedUrlIsTooLong
	}

	return nil
}

func ValidateOriginalUrl(url string) error {
	if len(url) > 512 {
		return errors.ErrUrlIsTooLong
	}

	return nil
}
