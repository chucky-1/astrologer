// Package service responsible for business logic
package service

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/chucky-1/astrologer/internal/dataprovider"
	"github.com/chucky-1/astrologer/internal/model"
	"github.com/chucky-1/astrologer/internal/repository"
)

// Picture - available actions for working with pictures
type Picture interface {
	GetPictureByDate(ctx context.Context, date time.Time) (*model.Picture, error)
	GetAllPictures(ctx context.Context) ([]*model.Picture, error)
}

// PictureImp - implements Picture
type PictureImp struct {
	apdClient   dataprovider.AstronomyPictureOfTheDay
	pictureRepo repository.Picture
}

// NewPictureImp is constructor
func NewPictureImp(apdClient dataprovider.AstronomyPictureOfTheDay, pictureRepo repository.Picture) *PictureImp {
	return &PictureImp{
		apdClient:   apdClient,
		pictureRepo: pictureRepo,
	}
}

// GetPictureByDate - get a picture for a certain date.
// First, we check the database, if we don’t find it, we turn to the api and add it to the database
func (p *PictureImp) GetPictureByDate(ctx context.Context, date time.Time) (*model.Picture, error) {
	picture, err := p.pictureRepo.GetByDate(ctx, date)
	if err != nil {
		return nil, err
	}
	if picture != nil {
		return picture, nil
	}

	picture, err = p.apdClient.GetPictureByDate(date)
	if err != nil {
		return nil, err
	}

	err = p.pictureRepo.Add(ctx, picture)
	if err != nil {
		logrus.Error(err)
	}

	return picture, nil
}

// GetAllPictures - get all pictures from database
func (p *PictureImp) GetAllPictures(ctx context.Context) ([]*model.Picture, error) {
	return p.pictureRepo.GetAllPictures(ctx)
}
