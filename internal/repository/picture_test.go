package repository

import (
	"astrologer/internal/model"
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/ory/dockertest"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"os"
	"os/exec"
	"testing"
	"time"
)

var (
	dbPool      *pgxpool.Pool
	pictureRepo *PicturePostgres
)

func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		logrus.Fatalf("Could not connect to docker: %s", err)
	}

	resource, err := pool.Run("postgres", "14.1-alpine", []string{"POSTGRES_PASSWORD=password123"})
	if err != nil {
		logrus.Fatalf("Could not start resource: %s", err)
	}

	var dbHostAndPort string

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = pool.Retry(func() error {
		dbHostAndPort = resource.GetHostPort("5432/tcp")

		dbPool, err = pgxpool.Connect(ctx, fmt.Sprintf("postgresql://postgres:password123@%v/postgres", dbHostAndPort))
		if err != nil {
			return err
		}

		return dbPool.Ping(ctx)
	})
	if err != nil {
		logrus.Fatalf("Could not connect to database: %s", err)
	}

	pictureRepo = NewPicturePostgres(dbPool)

	cmd := exec.Command("flyway",
		"-user=postgres",
		"-password=password123",
		"-locations=filesystem:../../migrations",
		fmt.Sprintf("-url=jdbc:postgresql://%v/postgres", dbHostAndPort),
		"migrate")

	err = cmd.Run()
	if err != nil {
		logrus.Fatalf("There are errors in migrations: %s", err)
	}

	code := m.Run()

	if err = pool.Purge(resource); err != nil {
		logrus.Fatalf("Could not purge resource: %s", err)
	}

	err = resource.Expire(1)
	if err != nil {
		logrus.Fatal(err)
	}

	os.Exit(code)
}

func TestPicturePostgres_AddAndGet(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	defer func() {
		_, err := dbPool.Exec(ctx, "TRUNCATE table picture")
		require.NoError(t, err)
	}()

	tm := time.Date(2022, 9, 15, 0, 0, 0, 0, time.UTC)
	pic := model.Picture{
		Title: "title",
		Date:  tm,
		Image: []byte{
			8, 45, 99,
		},
	}

	err := pictureRepo.Add(ctx, &pic)
	require.NoError(t, err)

	picture, err := pictureRepo.GetByDate(ctx, tm)
	require.NoError(t, err)
	require.Equal(t, &pic, picture)
}

func TestPicturePostgres_GetByDateIfPictureNil(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	picture, err := pictureRepo.GetByDate(ctx, time.Date(2021, 9, 15, 0, 0, 0, 0, time.UTC))
	require.NoError(t, err)
	require.Nil(t, picture)
}

func TestPicturePostgres_GetAllPictures(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	defer func() {
		_, err := dbPool.Exec(ctx, "TRUNCATE table picture")
		require.NoError(t, err)
	}()

	err := pictureRepo.Add(ctx, &model.Picture{
		Title: "title",
		Date:  time.Date(2022, 9, 15, 0, 0, 0, 0, time.UTC),
		Image: []byte{
			8, 45, 99,
		},
	})
	require.NoError(t, err)
	err = pictureRepo.Add(ctx, &model.Picture{
		Title: "title2",
		Date:  time.Date(2021, 9, 15, 0, 0, 0, 0, time.UTC),
		Image: []byte{
			8, 45, 99,
		},
	})
	require.NoError(t, err)
	err = pictureRepo.Add(ctx, &model.Picture{
		Title: "title3",
		Date:  time.Date(2020, 9, 15, 0, 0, 0, 0, time.UTC),
		Image: []byte{
			8, 45, 99,
		},
	})
	require.NoError(t, err)

	pictures, err := pictureRepo.GetAllPictures(ctx)
	require.NoError(t, err)
	require.Equal(t, 3, len(pictures))
}

func TestPicturePostgres_GetAllPicturesIfInPicturesNotElements(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	pictures, err := pictureRepo.GetAllPictures(ctx)
	require.NoError(t, err)
	require.Equal(t, 0, len(pictures))
}
