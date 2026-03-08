package factory

import (
	"database/sql"
	"goflow/config"
	"goflow/database"
)

type Factory struct {
	DB     *sql.DB
	Config *config.Config
}

func New(cfg *config.Config) (*Factory, error) {
	db, err := database.OpenSQLite(cfg.DBPath)
	if err != nil {
		return nil, err
	}

	return &Factory{
		DB:     db,
		Config: cfg,
	}, nil
}
