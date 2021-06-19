package delivery

import (
	"adv/internal/pkg/advert"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type AdvertHandler struct {
	AdvertUC advert.AdvertUsecase
}

func (h *AdvertHandler) GetAdverts(c echo.Context) error {
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil || page < 1 {
		page = 1
	}
	sortPrice := c.QueryParam("price")
	sortDate := c.QueryParam("date")

	adverts, err := h.AdvertUC.GetAdverts(page, sortPrice, sortDate)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}
	return c.JSON(http.StatusOK, adverts)
}

func (h *AdvertHandler) GetAdvert(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid advert id")
	}
	fields, err := strconv.ParseBool(c.QueryParam("fields"))
	if err != nil {
		fields = false
	}
	adv, err := h.AdvertUC.GetAdvert(id, fields)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}
	return c.JSON(http.StatusOK, adv)
}

func (h *AdvertHandler) CreateAdvert(c echo.Context) error {
	var adv advert.Advert
	defer c.Request().Body.Close()
	err := json.NewDecoder(c.Request().Body).Decode(&adv)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	id, err := h.AdvertUC.CreateAdvert(adv)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusCreated, struct {
		Id int `json:"id"`
	}{
		id,
	})
}

