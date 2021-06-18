package usecase

import (
	"adv/internal/pkg/advert"
	"adv/internal/pkg/advert/mocks"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetAdverts(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockRep := mocks.NewMockAdvertRepository(mockCtrl)
	testUc := AdvertUC{
		mockRep,
		10,
		3,
		200,
		1000,
	}

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

	mockRep.EXPECT().GetAdverts(1, "", "", 10).Return(adverts, nil).Times(1)
	responseAdverts, err := testUc.GetAdverts(1, "", "")
	if err != nil {
		t.Errorf("cannot get adverts: %v", err.Error())
	}
	assert.Equal(t, adverts, responseAdverts)

	mockRep.EXPECT().GetAdverts(1, "DESC", "ASC", 10).Return(adverts, nil).Times(1)
	responseAdverts, err = testUc.GetAdverts(1, "desc", "asc")
	if err != nil {
		t.Errorf("cannot get adverts: %v", err.Error())
	}
	assert.Equal(t, adverts, responseAdverts)

	mockRep.EXPECT().GetAdverts(1, "", "DESC", 10).Return(adverts, nil).Times(1)
	responseAdverts, err = testUc.GetAdverts(1, "invalidParam", "DESC")
	if err != nil {
		t.Errorf("cannot get adverts: %v", err.Error())
	}
	assert.Equal(t, adverts, responseAdverts)

	newError := errors.New("cannot get adverts")
	mockRep.EXPECT().GetAdverts(1, "", "", 10).Return(nil, newError).Times(1)
	_, err = testUc.GetAdverts(1, "", "")
	if err != newError {
		t.Errorf("didn't pass error")
	}
}

func TestGetAdvert(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockRep := mocks.NewMockAdvertRepository(mockCtrl)
	testUc := AdvertUC{
		mockRep,
		10,
		3,
		200,
		1000,
	}

	adv := advert.Advert{
			Name: "New test advert",
			About: "About new test",
			Price: 100,
			Photos: []string{"media/1/photo", "media/2/photo"},
	}

	mockRep.EXPECT().GetAdvert(1).Return(adv, nil).Times(1)
	responseAdv, err := testUc.GetAdvert(1, true)
	if err != nil {
		t.Errorf("cannot get advert: %v", err.Error())
	}
	assert.Equal(t, adv, responseAdv)

	mockRep.EXPECT().GetAdvert(1).Return(adv, nil).Times(1)
	responseAdv, err = testUc.GetAdvert(1, false)
	if err != nil {
		t.Errorf("cannot get advert: %v", err.Error())
	}
	expectedAdv := advert.Advert{
		Name: "New test advert",
		Price: 100,
		MainPhoto: "media/1/photo",
	}
	assert.Equal(t, expectedAdv, responseAdv)

	newError := errors.New("cannot get advert")
	mockRep.EXPECT().GetAdvert(1).Return(advert.Advert{}, newError).Times(1)
	_, err = testUc.GetAdvert(1, false)
	assert.Error(t, newError)
}

func TestCreateAdvert(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockRep := mocks.NewMockAdvertRepository(mockCtrl)
	testUc := AdvertUC{
		mockRep,
		10,
		3,
		200,
		1000,
	}

	adv := advert.Advert{
		Name: "New test advert",
		About: "About new test",
		Price: 100,
		Photos: []string{"media/1/photo", "media/2/photo"},
	}

	mockRep.EXPECT().CreateAdvert(adv).Return(1, nil).Times(1)
	id, err := testUc.CreateAdvert(adv)
	if err != nil {
		t.Errorf("cannot create advert: %v", err.Error())
	}
	assert.Equal(t, 1, id)

	newError := errors.New("cannot create advert")
	mockRep.EXPECT().CreateAdvert(adv).Return(0, newError).Times(1)
	_, err = testUc.CreateAdvert(adv)
	assert.Error(t, newError)

	adv.Photos = append(adv.Photos, adv.Photos...)
	_, err = testUc.CreateAdvert(adv)
	assert.Error(t, errors.New("incorrect data"))
}