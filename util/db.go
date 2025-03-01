package util

import (
	"context"
	"ctrl/database"
	"database/sql"
	"errors"
	"fmt"
	"path/filepath"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

var Schema *string

type Instance struct {
	DB    *sql.DB
	Ctx   context.Context
	Query *database.Queries
}

func tablesExists(db *sql.DB) bool {
	if _, err := db.Query("SELECT * FROM registry"); err == nil {
		return true
	}
	return false
}

var ErrCreatingDatabase = errors.New("couldnt create database")

func InitializeInstance() (*Instance, error) {
	dataPath, err := DataPath()
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite3", fmt.Sprintf("file:%s", filepath.Join(dataPath, "ctrl.db")))
	if err != nil {
		return nil, fmt.Errorf("%w :: %s", ErrCreatingDatabase, err.Error())
	}

	ctx := context.Background()
	if !tablesExists(db) {
		if _, err := db.ExecContext(ctx, *Schema); err != nil {
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
	tx, err := instance.DB.Begin()
	if err != nil {
		return fmt.Errorf("%w :: %s", ErrAddingApp, err.Error())
	}
	defer tx.Rollback()

	qtx := instance.Query.WithTx(tx)
	registryId := uuid.New().String()

	err = qtx.AddRegistryApp(instance.Ctx, database.AddRegistryAppParams{
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

	err = qtx.SyncRegistrySearchApps(instance.Ctx)
	if err != nil {
		return fmt.Errorf("%w :: %s", ErrAddingApp, err.Error())
	}

	err = qtx.AddApp(instance.Ctx, database.AddAppParams{
		AppID:         id,
		AppSource:     source,
		AppOs:         GetOS(),
		AppRegistryID: registryId,
	})
	if err != nil {
		return fmt.Errorf("%w :: %s", ErrAddingApp, err.Error())
	}

	return tx.Commit()
}

func (instance *Instance) InstallApp(id string, name string, source string, version string, newVersion sql.NullString) error {
	tx, err := instance.DB.Begin()
	if err != nil {
		return fmt.Errorf("%w :: %s", ErrAddingApp, err.Error())
	}
	defer tx.Rollback()

	qtx := instance.Query.WithTx(tx)
	registryId := uuid.New().String()

	err = qtx.AddRegistryApp(instance.Ctx, database.AddRegistryAppParams{
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

	err = qtx.SyncRegistrySearchApps(instance.Ctx)
	if err != nil {
		return fmt.Errorf("%w :: %s", ErrAddingApp, err.Error())
	}

	err = qtx.InstallApp(instance.Ctx, database.InstallAppParams{
		AppID:         id,
		AppSource:     source,
		AppOs:         GetOS(),
		AppRegistryID: registryId,
		AppVersion:    sql.NullString{String: version, Valid: true},
		AppAvailable:  newVersion,
	})
	if err != nil {
		return fmt.Errorf("%w :: %s", ErrAddingApp, err.Error())
	}

	return tx.Commit()
}
