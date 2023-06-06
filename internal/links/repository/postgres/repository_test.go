package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/gommon/log"
	"github.com/stretchr/testify/require"
	"project/internal/model"
	pkgErrors "project/internal/pkg/errors"
	"regexp"
	"testing"
)

func TestPostgres_GetOriginalUrl_OK(t *testing.T) {
	excpectedOriginalUrl := "vk.com"
	abbreviatedUrl := "AF4_ResRbL"

	db, mock, err := sqlmock.New()
	require.Nil(t, err, fmt.Errorf("cant create mock: %s", err))
	defer func() {
		err := db.Close()
		if err != nil {
			log.Error(err)
		}
	}()

	row := sqlmock.NewRows([]string{"original_url", "abbreviated_url"}).
		AddRow(excpectedOriginalUrl, abbreviatedUrl)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM link WHERE abbreviated_url=$1`)).
		WithArgs(abbreviatedUrl).
		WillReturnRows(row)

	dbx := sqlx.NewDb(db, "sqlmock")
	repo := NewLinksRepository(dbx)

	originalUrl, err := repo.GetOriginalUrl(context.TODO(), abbreviatedUrl)
	require.NoError(t, err)
	require.Equal(t, excpectedOriginalUrl, originalUrl)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestPostgres_GetOriginalUrl_NotFoud(t *testing.T) {
	abbreviatedUrl := "AF4_ResRbL"

	db, mock, err := sqlmock.New()
	require.Nil(t, err, fmt.Errorf("cant create mock: %s", err))
	defer func() {
		err := db.Close()
		if err != nil {
			log.Error(err)
		}
	}()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM link WHERE abbreviated_url=$1`)).
		WithArgs(abbreviatedUrl).
		WillReturnError(sql.ErrNoRows)

	dbx := sqlx.NewDb(db, "sqlmock")
	repo := NewLinksRepository(dbx)

	_, err = repo.GetOriginalUrl(context.TODO(), abbreviatedUrl)
	require.Error(t, err, pkgErrors.ErrAbbreviatedUrlNotFound)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestPostgres_SaveAbbreviatedLink_OK(t *testing.T) {
	link := model.LinkDB{
		AbbreviatedUrl: "AF4_ResRbL",
		OriginalUrl:    "vk.com",
	}

	db, mock, err := sqlmock.New()
	require.Nil(t, err, fmt.Errorf("cant create mock: %s", err))
	defer func() {
		err := db.Close()
		if err != nil {
			log.Error(err)
		}
	}()

	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO link(original_url, abbreviated_url) VALUES(?, ?)`)).
		WithArgs(link.OriginalUrl, link.AbbreviatedUrl).
		WillReturnResult(sqlmock.NewResult(1, 1))

	dbx := sqlx.NewDb(db, "sqlmock")
	repo := NewLinksRepository(dbx)

	err = repo.SaveAbbreviatedLink(context.TODO(), link)
	require.NoError(t, err)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestPostgres_CheckExistAbbreviatedLink_OK(t *testing.T) {
	abbreviatedUrl := "AF4_ResRbL"

	db, mock, err := sqlmock.New()
	require.Nil(t, err, fmt.Errorf("cant create mock: %s", err))
	defer func() {
		err := db.Close()
		if err != nil {
			log.Error(err)
		}
	}()

	row := sqlmock.NewRows([]string{"exists"}).
		AddRow(true)

	mock.
		ExpectQuery(regexp.QuoteMeta(`SELECT EXISTS(SELECT 1 FROM link WHERE abbreviated_url=$1)`)).
		WithArgs(abbreviatedUrl).
		WillReturnRows(row)

	dbx := sqlx.NewDb(db, "sqlmock")
	repo := NewLinksRepository(dbx)

	err = repo.CheckExistAbbreviatedLink(context.TODO(), abbreviatedUrl)
	require.NoError(t, err)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestPostgres_CheckExistAbbreviatedLink_False(t *testing.T) {
	abbreviatedUrl := "AF4_ResRbL"

	db, mock, err := sqlmock.New()
	require.Nil(t, err, fmt.Errorf("cant create mock: %s", err))
	defer func() {
		err := db.Close()
		if err != nil {
			log.Error(err)
		}
	}()

	row := sqlmock.NewRows([]string{"exists"}).
		AddRow(false)

	mock.
		ExpectQuery(regexp.QuoteMeta(`SELECT EXISTS(SELECT 1 FROM link WHERE abbreviated_url=$1)`)).
		WithArgs(abbreviatedUrl).
		WillReturnRows(row)

	dbx := sqlx.NewDb(db, "sqlmock")
	repo := NewLinksRepository(dbx)

	err = repo.CheckExistAbbreviatedLink(context.TODO(), abbreviatedUrl)
	require.Error(t, err, pkgErrors.ErrAbbreviatedUrlNotFound)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}
