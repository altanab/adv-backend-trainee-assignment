package repository

import (
	"adv/internal/pkg/advert"
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"regexp"
	"testing"
)

type Suite struct {
	suite.Suite
	DB *gorm.DB
	mock sqlmock.Sqlmock
	testRep AdvertRep
}

func (s *Suite) SetupSuite() {
	var (
		db  *sql.DB
		err error
	)

	db, s.mock, err = sqlmock.New()
	require.NoError(s.T(), err)

	s.DB, err = gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	require.NoError(s.T(), err)


	s.testRep = AdvertRep{
		s.DB,
	}
}

func TestInit(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (s *Suite) TestGetAdverts() {
	adverts := []advert.Advert{
		{
			Id: 1,
			Name: "New test advert",
			Price: 100,
			MainPhoto: "media/1/photo",
		},

		{
			Id: 2,
			Name: "New test advert 2",
			Price: 10,
			MainPhoto: "media/2/photo",
		},
	}

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT id, adv_name, photos[1] as photo, price FROM "adverts" ORDER BY price DESC, created ASC LIMIT 10 OFFSET 10`)).
		WillReturnRows(sqlmock.NewRows([]string{
			"id",
			"adv_name",
			"photo",
			"price",
		}).AddRow(
			adverts[0].Id,
			adverts[0].Name,
			adverts[0].MainPhoto,
			adverts[0].Price,
		).AddRow(
			adverts[1].Id,
			adverts[1].Name,
			adverts[1].MainPhoto,
			adverts[1].Price,
	))

	responseAdverts, err := s.testRep.GetAdverts(2, "DESC", "ASC", 10)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), adverts, responseAdverts)

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT id, adv_name, photos[1] as photo, price FROM "adverts" ORDER BY created DESC LIMIT 10 OFFSET 10`)).
		WillReturnRows(sqlmock.NewRows([]string{
			"id",
			"adv_name",
			"photo",
			"price",
		}).AddRow(
			adverts[0].Id,
			adverts[0].Name,
			adverts[0].MainPhoto,
			adverts[0].Price,
		).AddRow(
			adverts[1].Id,
			adverts[1].Name,
			adverts[1].MainPhoto,
			adverts[1].Price,
		))

	responseAdverts, err = s.testRep.GetAdverts(2, "", "", 10)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), adverts, responseAdverts)
}

func (s *Suite) TestGetAdvert() {
	adv := advert.Advert{
			Name: "New test advert",
			About: "About new test",
			Price: 100,
			Photos: pq.StringArray([]string{"media/1/photo", "media/2/photo", "media/3/photo"}),
	}

	s.mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT adv_name, photos, price, about FROM "adverts" WHERE id=$1`)).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{
			"adv_name",
			"photos",
			"price",
			"about",
		}).AddRow(
			adv.Name,
			adv.Photos,
			adv.Price,
			adv.About,
		))

	responseAdvert, err := s.testRep.GetAdvert(1)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), adv, responseAdvert)

}

func (s *Suite) TestCreateAdvert() {
	adv := advert.Advert{
		Name: "New test advert",
		About: "About new test",
		Price: 100,
		Photos: pq.StringArray([]string{"media/1/photo", "media/2/photo", "media/3/photo"}),
	}

	s.mock.MatchExpectationsInOrder(false)
	s.mock.ExpectBegin()
	s.mock.ExpectQuery("INSERT INTO").
		WithArgs(
			adv.Name,
			adv.About,
			adv.Photos,
			adv.Price,
		).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).
			AddRow(1))
	s.mock.ExpectCommit()
	id, err := s.testRep.CreateAdvert(adv)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), 1, id)
}
