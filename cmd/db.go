package cmd

import (
	"context"
	"ctrl/database"
	"ctrl/util"
	"database/sql"
	"errors"
	"fmt"
	"path/filepath"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

var Schema string

type Instance struct {
	db    *sql.DB
	ctx   context.Context
	query *database.Queries
}

func tablesExists(db *sql.DB) bool {
	if _, err := db.Query("SELECT * FROM registry"); err == nil {
		return true
	}
	return false
}

var ErrCreatingDatabase = errors.New("couldnt create database")

func InitializeInstance() (*Instance, error) {
	dataPath, err := util.DataPath()
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite3", fmt.Sprintf("file:%s", filepath.Join(dataPath, "ctrl.db")))
	if err != nil {
		return nil, fmt.Errorf("%w :: %s", ErrCreatingDatabase, err.Error())
	}

	ctx := context.Background()
	if !tablesExists(db) {
		if _, err := db.ExecContext(ctx, Schema); err != nil {
			return nil, fmt.Errorf("%w :: %s", ErrCreatingDatabase, err)
		}
		fmt.Println("ctrl sqlite tables created!")
	}

	query := database.New(db)

	instance := Instance{db, ctx, query}
	return &instance, nil
}

var ErrAddingApp = errors.New("couldnt add app to database")

func (instance *Instance) AddApp(id string, name string, source string) error {
	tx, err := instance.db.Begin()
	if err != nil {
		return fmt.Errorf("%w :: %s", ErrAddingApp, err.Error())
	}
	defer tx.Rollback()

	qtx := instance.query.WithTx(tx)
	registryId := uuid.New().String()

	err = qtx.AddRegistryApp(instance.ctx, database.AddRegistryAppParams{
		RegistryID:   registryId,
		RegistryName: name,
	})
	if err != nil {
		switch err.Error() {
		case "UNIQUE constraint failed: registry.registry_name":
			return fmt.Errorf("%w :: %s", ErrAddingApp, "app already exists in the database")
		default:
			return fmt.Errorf("%w :: %s", ErrAddingApp, err.Error())
		}
	}

	err = qtx.SyncRegistrySearchApps(instance.ctx)
	if err != nil {
		return fmt.Errorf("%w :: %s", ErrAddingApp, err.Error())
	}

	err = qtx.AddApp(instance.ctx, database.AddAppParams{
		AppID:         id,
		AppSource:     source,
		AppOs:         util.GetOS(),
		AppRegistryID: registryId,
	})
	if err != nil {
		return fmt.Errorf("%w :: %s", ErrAddingApp, err.Error())
	}

	return tx.Commit()
}
