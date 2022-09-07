// Package dataprovider Package data provider responsible for working with third-party APIs
package dataprovider

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/chucky-1/astrologer/internal/model"
)

const (
	defaultTimeoutForRequest = 5 * time.Second
	queryWithoutParameters   = "https://api.nasa.gov/planetary/apod?api_key=DEMO_KEY"
)

// AstronomyPictureOfTheDay - the ability to get a picture for a specific date
type AstronomyPictureOfTheDay interface {
	GetPictureByDate(date time.Time) (*model.Picture, error)
}

// AstronomyPictureOfTheDayClient - works with api https://api.nasa.gov/
type AstronomyPictureOfTheDayClient struct {
	client *http.Client
}

// NewAstronomyPictureOfTheDayClient is constructor
func NewAstronomyPictureOfTheDayClient() *AstronomyPictureOfTheDayClient {
	return &AstronomyPictureOfTheDayClient{
		client: &http.Client{
			Timeout: defaultTimeoutForRequest,
		},
	}
}

type getPictureByDateResponse struct {
	Title string `json:"title"`
	URL   string `json:"url"`
}

// GetPictureByDate - get an image for a specific date
func (a *AstronomyPictureOfTheDayClient) GetPictureByDate(date time.Time) (*model.Picture, error) {
	path := fmt.Sprintf("%v&date=%d-%02d-%02d", queryWithoutParameters, date.Year(), int(date.Month()), date.Day())

	response, err := a.getRequest(path)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		errClose := Body.Close()
		if errClose != nil {
			logrus.Errorf("dataprovider.GetPictureByDate, Close: %v", errClose)
		}
	}(response.Body)

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("dataprovider.AstronomyPictureOfTheDayClient.getRequestAndRead, ReadAll: %v", err)
	}

	var responseObject getPictureByDateResponse
	err = json.Unmarshal(responseData, &responseObject)
	if err != nil {
		return nil, fmt.Errorf("dataprovider.AstronomyPictureOfTheDayClient.GetPictureByDate, Unmarshal: %v", err)
	}

	image, err := a.getImageByURL(responseObject.URL)
	if err != nil {
		return nil, err
	}

	return &model.Picture{
		Title: responseObject.Title,
		Date:  date,
		Image: image,
	}, nil
}

func (a *AstronomyPictureOfTheDayClient) getImageByURL(url string) ([]byte, error) {
	response, err := a.getRequest(url)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		errClose := Body.Close()
		if errClose != nil {
			logrus.Errorf("dataprovider.getImageByURL, Close: %v", errClose)
		}
	}(response.Body)

	image, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("dataprovider.AstronomyPictureOfTheDayClient.getImageByURL, ReadAll: %v", err)
	}

	return image, nil
}

func (a *AstronomyPictureOfTheDayClient) getRequest(url string) (*http.Response, error) {
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("dataprovider.AstronomyPictureOfTheDayClient.getRequest, NewRequest: %v", err)
	}

	response, err := a.client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("dataprovider.AstronomyPictureOfTheDayClient.getRequest, Do: %v", err)
	}

	return response, nil
}
