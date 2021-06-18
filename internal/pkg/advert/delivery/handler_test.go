package delivery

import (
	"adv/internal/pkg/advert"
	"adv/internal/pkg/advert/mocks"
	"bytes"
	"encoding/json"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestGetAdverts(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockUc := mocks.NewMockAdvertUsecase(mockCtrl)
	testHandler := AdvertHandler{mockUc}

	e := echo.New()
	q := make(url.Values)
	q.Set("page", "10")
	q.Set("price", "DESC")
	q.Set("date", "ASC")
	req := httptest.NewRequest(http.MethodGet, "/advert?"+q.Encode(), nil)
	response := httptest.NewRecorder()
	echoContext := e.NewContext(req, response)

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

	mockUc.EXPECT().GetAdverts(10, "DESC", "ASC").Return(adverts, nil).Times(1)

	err := testHandler.GetAdverts(echoContext)
	if err != nil {
		t.Errorf("couldn't get adverts: %s", err.Error())
	}

	responseAdverts := make([]advert.Advert, 0)
	err = json.NewDecoder(response.Body).Decode(&responseAdverts)
	if err != nil {
		t.Errorf("couldn't unmarshal adverts: %s", err.Error())
	}
	assert.Equal(t, adverts, responseAdverts)

	req = httptest.NewRequest(http.MethodGet, "/advert", nil)
	response = httptest.NewRecorder()
	echoContext = e.NewContext(req, response)
	mockUc.EXPECT().GetAdverts(1, "", "").Return(nil, errors.New("cannot get adverts")).Times(1)
	err = testHandler.GetAdverts(echoContext)
	if httperr, ok := err.(*echo.HTTPError); ok {
		if httperr.Code != http.StatusNotFound {
			t.Errorf("didn't pass error: %v\n", err)
		}
	} else {
		t.Errorf("didn't pass error: %v\n", err)
	}
}

func TestGetAdvert(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockUc := mocks.NewMockAdvertUsecase(mockCtrl)
	testHandler := AdvertHandler{mockUc}

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/advert/1", nil)
	response := httptest.NewRecorder()
	echoContext := e.NewContext(req, response)
	echoContext.SetPath("/advert/:id")
	echoContext.SetParamNames("id")
	echoContext.SetParamValues("1")

	adv := advert.Advert{
			Name: "New test advert",
			Price: 100,
			MainPhoto: "media/1/photo",
	}
	mockUc.EXPECT().GetAdvert(1, false).Return(adv, nil).Times(1)
	err := testHandler.GetAdvert(echoContext)
	if err != nil {
		t.Errorf("couldn't get advert: %s", err.Error())
	}
	var responseAdvert advert.Advert
	err = json.NewDecoder(response.Body).Decode(&responseAdvert)
	if err != nil {
		t.Errorf("couldn't unmarshal advert: %s", err.Error())
	}
	assert.Equal(t, adv, responseAdvert)


	req = httptest.NewRequest(http.MethodGet, "/advert/1?fields=true", nil)
	response = httptest.NewRecorder()
	echoContext = e.NewContext(req, response)
	echoContext.SetPath("/advert/:id")
	echoContext.SetParamNames("id")
	echoContext.SetParamValues("1")

	advWithFields := advert.Advert{
		Name: "New test advert",
		Price: 100,
		About : "About test",
		Photos: []string{"media/1", "media/2", "media/3"},
	}
	mockUc.EXPECT().GetAdvert(1, true).Return(advWithFields, nil).Times(1)
	err = testHandler.GetAdvert(echoContext)
	if err != nil {
		t.Errorf("couldn't get advert with fields: %s", err.Error())
	}
	var responseAdvertWithFields advert.Advert
	err = json.NewDecoder(response.Body).Decode(&responseAdvertWithFields)
	if err != nil {
		t.Errorf("couldn't unmarshal advert with field: %s", err.Error())
	}
	assert.Equal(t, advWithFields, responseAdvertWithFields)

	req = httptest.NewRequest(http.MethodGet, "/advert/1", nil)
	response = httptest.NewRecorder()
	echoContext = e.NewContext(req, response)
	echoContext.SetPath("/advert/:id")
	echoContext.SetParamNames("id")
	echoContext.SetParamValues("1")
	mockUc.EXPECT().GetAdvert(1, false).Return(advert.Advert{}, errors.New("cannot get advert")).Times(1)
	err = testHandler.GetAdvert(echoContext)
	if httperr, ok := err.(*echo.HTTPError); ok {
		if httperr.Code != http.StatusNotFound {
			t.Errorf("didn't pass error: %v\n", err)
		}
	} else {
		t.Errorf("didn't pass error: %v\n", err)
	}

	req = httptest.NewRequest(http.MethodGet, "/advert/one", nil)
	response = httptest.NewRecorder()
	echoContext = e.NewContext(req, response)
	echoContext.SetPath("/advert/:id")
	echoContext.SetParamNames("id")
	echoContext.SetParamValues("one")
	err = testHandler.GetAdvert(echoContext)
	if httperr, ok := err.(*echo.HTTPError); ok {
		if httperr.Code != http.StatusBadRequest {
			t.Errorf("didn't pass invalid id: %v\n", err)
		}
	} else {
		t.Errorf("didn't pass invalid id: %v\n", err)
	}
}

func TestCreateAdvert(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockUc := mocks.NewMockAdvertUsecase(mockCtrl)
	testHandler := AdvertHandler{mockUc}

	adv := advert.Advert{
		Name: "New test advert",
		About: "About new test",
		Price: 100,
		Photos: []string{"media/1/photo", "media/2/photo"},
	}
	e := echo.New()
	body, _ := json.Marshal(adv)
	req := httptest.NewRequest(http.MethodPost, "/advert", bytes.NewReader(body))
	response := httptest.NewRecorder()
	echoContext := e.NewContext(req, response)
	mockUc.EXPECT().CreateAdvert(adv).Return(1, nil).Times(1)
	err := testHandler.CreateAdvert(echoContext)
	if err != nil {
		t.Errorf("couldn't create advert: %s", err.Error())
	}
	id := struct {
		Id int `json:"id"`
	}{}
	err = json.NewDecoder(response.Body).Decode(&id)
	if err != nil {
		t.Errorf("couldn't unmarshal id: %s", err.Error())
	}
	assert.Equal(t, 1, id.Id)

	req = httptest.NewRequest(http.MethodPost, "/advert", bytes.NewReader(body))
	response = httptest.NewRecorder()
	echoContext = e.NewContext(req, response)
	mockUc.EXPECT().CreateAdvert(adv).Return(0, errors.New("cannot create advert")).Times(1)
	err = testHandler.CreateAdvert(echoContext)
	if httperr, ok := err.(*echo.HTTPError); ok {
		if httperr.Code != http.StatusInternalServerError {
			t.Errorf("didn't pass error: %v\n", err)
		}
	} else {
		t.Errorf("didn't pass error: %v\n", err)
	}
}