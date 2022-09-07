package main

import (
	"context"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"

	"github.com/chucky-1/astrologer/internal/dataprovider"
	"github.com/chucky-1/astrologer/internal/handler"
	"github.com/chucky-1/astrologer/internal/repository"
	"github.com/chucky-1/astrologer/internal/service"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("No .env file found")
	}

	apdClient := dataprovider.NewAstronomyPictureOfTheDayClient()
	conn, err := pgxpool.Connect(ctx, os.Getenv("POSTGRES_ENDPOINT"))
	if err != nil {
		logrus.Fatalf("couldn't connect to database: %v", err)
	}
	if err = conn.Ping(ctx); err != nil {
		logrus.Fatalf("couldn't ping database: %v", err)
	}

	pictureRepo := repository.NewPicturePostgres(conn)
	pictureService := service.NewPictureImp(apdClient, pictureRepo)
	pictureHandler := handler.NewPicture(pictureService)

	e := echo.New()
	e.GET("/picture", pictureHandler.GetPictureByDate)
	e.GET("/pictures", pictureHandler.GetAllPictures)
	e.Logger.Fatal(e.Start(":1323"))
}
