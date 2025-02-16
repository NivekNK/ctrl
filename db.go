package main

import (
	"context"
	"ctrl/database"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"

	_ "embed"

	"github.com/google/uuid"
	gap "github.com/muesli/go-app-paths"
	_ "modernc.org/sqlite"
)

//go:embed .sql/schema.sql
var ddl string

type CtrlDB struct {
	instance *sql.DB
	ctx      context.Context
	db       *database.Queries
}

func initAppsDir(path string) error {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return os.Mkdir(path, 0o770)
		}
		return err
	}
	return nil
}

func setupPath() string {
	scope := gap.NewScope(gap.User, "ctrl")

	dirs, err := scope.DataDirs()
	if err != nil {
		log.Fatal(err)
	}

	var appsDir string
	if len(dirs) > 0 {
		appsDir = dirs[0]
	} else {
		appsDir, _ = os.UserHomeDir()
	}

	if err := initAppsDir(appsDir); err != nil {
		log.Fatal(err)
	}

	return appsDir
}

func tablesExists(instance *sql.DB) bool {
	if _, err := instance.Query("SELECT * FROM registry"); err == nil {
		return true
	}
	return false
}

func initDB() (*CtrlDB, error) {
	ctx := context.Background()
	path := setupPath()

	instance, err := sql.Open("sqlite", fmt.Sprintf("file:%s", filepath.Join(path, "ctrl.db")))
	if err != nil {
		return nil, err
	}

	// create tables
	if !tablesExists(instance) {
		if _, err := instance.ExecContext(ctx, ddl); err != nil {
			return nil, err
		}
	}

	db := database.New(instance)

	ctrlDB := CtrlDB{instance, ctx, db}
	return &ctrlDB, nil
}

func getOS() string {
	if runtime.GOOS == "windows" {
		return "windows"
	}
	return "linux"
}

func addApp(ctx context.Context, db *database.Queries, id string, name string, source string) error {
	registryId := uuid.New().String()

	err := db.AddRegistryApp(ctx, database.AddRegistryAppParams{
		RegistryID:   registryId,
		RegistryName: name,
	})
	if err != nil {
		return err
	}

	err = db.AddApp(ctx, database.AddAppParams{
		AppID:         id,
		AppSource:     source,
		AppOs:         getOS(),
		AppRegistryID: registryId,
	})
	if err != nil {
		return err
	}

	return nil
}
