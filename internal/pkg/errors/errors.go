package errors

import (
	"errors"
)

var (
	// 400 HTTP code
	GrpcErrUrlIsTooLong = errors.New("rpc error: code = Unknown desc = слишком длинный URL(больше 512 символов)")
	ErrUrlIsTooLong     = errors.New("слишком длинный URL(больше 512 символов)")

	GrpcErrAbbreviatedUrlIsTooLong = errors.New("rpc error: code = Unknown desc = слишком длинный URL(больше 10 символов)")
	ErrAbbreviatedUrlIsTooLong     = errors.New("слишком длинный URL(больше 10 символов)")

	// 404 HTTP code
	GrpcErrAbbreviatedUrlNotFound = errors.New("rpc error: code = Unknown desc = сокращённый URL не найден")
	ErrAbbreviatedUrlNotFound     = errors.New("сокращённый URL не найден")

	// 500 HTTP code
	GrpcErrInternal = errors.New("rpc error: code = Unknown desc = внутренняя ошибка сервера")
	ErrInternal     = errors.New("внутренняя ошибка сервера")
)

func FromGrpcErrorToOriginalError(err error) error {
	switch err.Error() {
	case GrpcErrUrlIsTooLong.Error():
		return ErrUrlIsTooLong
	case GrpcErrAbbreviatedUrlIsTooLong.Error():
		return ErrAbbreviatedUrlIsTooLong
	case GrpcErrAbbreviatedUrlNotFound.Error():
		return ErrAbbreviatedUrlNotFound
	default:
		return ErrInternal
	}
}
