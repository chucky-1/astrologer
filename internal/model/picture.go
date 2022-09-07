// Package model stores entities
package model

import "time"

// Picture - image structure obtained from AstronomyPictureOfTheDay
type Picture struct {
	Title string
	Date  time.Time
	Image []byte
}
