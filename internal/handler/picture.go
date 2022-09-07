// Package handler responsible for handling user requests
package handler

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"

	"github.com/chucky-1/astrologer/internal/service"
)

const layout = "2006-01-02"

// Picture - responsible for handling user requests
type Picture struct {
	servicePicture service.Picture
}

// NewPicture is constructor
func NewPicture(servicePicture service.Picture) *Picture {
	return &Picture{servicePicture: servicePicture}
}

// getPictureByDateRequest format: YYYY-MM-DD
type getPictureByDateRequest struct {
	Date string `json:"date"`
}

// GetPictureByDate - get a picture for a certain day
func (p *Picture) GetPictureByDate(c echo.Context) error {
	request := new(getPictureByDateRequest)
	if err := c.Bind(request); err != nil {
		logrus.Errorf("handler.GetPictureByDate, Bind: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, "can't get picture, try again.")
	}

	date, err := time.Parse(layout, request.Date)
	if err != nil {
		logrus.Errorf("handler.GetPictureByDate, Parse: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, "can't get picture, try again.")
	}

	picture, err := p.servicePicture.GetPictureByDate(c.Request().Context(), date)
	if err != nil {
		logrus.Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "something went wrong, try again.")
	}

	return c.JSON(http.StatusOK, picture)
}

// GetAllPictures - get all available pictures
func (p *Picture) GetAllPictures(c echo.Context) error {
	pictures, err := p.servicePicture.GetAllPictures(c.Request().Context())
	if err != nil {
		logrus.Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, "something went wrong, try again.")
	}
	return c.JSON(http.StatusOK, pictures)
}
