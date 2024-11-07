package storage

import (
	"context"
	"fmt"
	"githib.com/s4bb4t/leadgen/internal/config"
	"githib.com/s4bb4t/leadgen/internal/lib/models"
	"githib.com/s4bb4t/leadgen/internal/storage/pgsql"
)

type RepositoryI interface {
	Save(context.Context, models.Building) (models.Building, error)
	Building(context.Context, string) (models.Building, error)
	Buildings(context.Context, models.Query) (models.Buildings, error)
	Close() error
}

func MustLoad(cfg *config.Config) RepositoryI {
	repo, err := pgsql.Connect(cfg)
	if err != nil {
		panic(err)
	}

	fmt.Println("repo ok")

	return repo
}
