// Package repository responsible for database operations
package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/chucky-1/astrologer/internal/model"
)

// Picture - all possible database operations
type Picture interface {
	Add(ctx context.Context, picture *model.Picture) error
	GetByDate(ctx context.Context, date time.Time) (*model.Picture, error)
	GetAllPictures(ctx context.Context) ([]*model.Picture, error)
}

// PicturePostgres - implements operations with Postgres
type PicturePostgres struct {
	conn *pgxpool.Pool
}

// NewPicturePostgres is constructor
func NewPicturePostgres(conn *pgxpool.Pool) *PicturePostgres {
	return &PicturePostgres{conn: conn}
}

// Add - adds title, date and byte array
func (p *PicturePostgres) Add(ctx context.Context, picture *model.Picture) error {
	_, err := p.conn.Exec(ctx, `INSERT INTO picture (title, date, image) VALUES ($1, $2, $3)`, picture.Title,
		picture.Date, picture.Image)
	if err != nil {
		return fmt.Errorf("repository.PicturePostgres.Add, Exec: %v", err)
	}
	return nil
}

// GetByDate - returns a record by date
func (p *PicturePostgres) GetByDate(ctx context.Context, date time.Time) (*model.Picture, error) {
	var picture model.Picture
	err := p.conn.QueryRow(ctx, `SELECT title, date, image FROM picture WHERE date = $1`, date).Scan(&picture.Title,
		&picture.Date, &picture.Image)
	if err != nil && err != pgx.ErrNoRows {
		return nil, fmt.Errorf("repository.PicturePostgres.GetByDate, QueryRow.Scan: %v", err)
	} else if err == pgx.ErrNoRows {
		return nil, nil
	}
	return &picture, nil
}

// GetAllPictures - returns all records from the database
func (p *PicturePostgres) GetAllPictures(ctx context.Context) ([]*model.Picture, error) {
	rows, err := p.conn.Query(ctx, `SELECT title, date, image FROM picture`)
	if err != nil {
		return nil, fmt.Errorf("repository.PicturePostgres.GetAllPictures, Query: %v", err)
	}
	defer rows.Close()

	var pictures []*model.Picture

	for rows.Next() {
		var picture model.Picture

		errScan := rows.Scan(&picture.Title, &picture.Date, &picture.Image)
		if errScan != nil {
			return nil, fmt.Errorf("repository.PicturePostgres.GetAllPictures, Scan: %v", errScan)
		}

		pictures = append(pictures, &picture)
	}

	return pictures, nil
}
