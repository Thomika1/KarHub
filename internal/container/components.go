package container

import (
	"context"

	"github.com/Thomika1/KarHub/internal/config"
)

func setupComponents(ctx context.Context) (*components, error) {
	cmp := &components{}

	if err := setupAPIComponents(ctx, cmp); err != nil {
		return nil, err
	}

	return cmp, nil

}

func setupAPIComponents(ctx context.Context, cmp *components) error {
	// Init logger
	cmp.Logger = config.InitLogger()

	// Init database
	_, dbInstance, err := config.InitDB(ctx)
	if err != nil {
		return err
	}
	cmp.Database = dbInstance

	return nil
}
