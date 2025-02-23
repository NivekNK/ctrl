// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: query.sql

package database

import (
	"context"
	"database/sql"
)

const addApp = `-- name: AddApp :exec
INSERT INTO app (
    app_id, 
    app_source, 
    app_os, 
    app_registry_id, 
    app_last_updated, 
    app_installed, 
    app_version, 
    app_available
) 
VALUES (
    ?, 
    ?, 
    ?, 
    ?, 
    datetime('now', 'utc'), 
    0, 
    NULL, 
    NULL
)
`

type AddAppParams struct {
	AppID         string
	AppSource     string
	AppOs         string
	AppRegistryID string
}

func (q *Queries) AddApp(ctx context.Context, arg AddAppParams) error {
	_, err := q.db.ExecContext(ctx, addApp,
		arg.AppID,
		arg.AppSource,
		arg.AppOs,
		arg.AppRegistryID,
	)
	return err
}

const addRegistryApp = `-- name: AddRegistryApp :exec
INSERT INTO registry (registry_id, registry_name) VALUES (?, ?)
`

type AddRegistryAppParams struct {
	RegistryID   string
	RegistryName string
}

func (q *Queries) AddRegistryApp(ctx context.Context, arg AddRegistryAppParams) error {
	_, err := q.db.ExecContext(ctx, addRegistryApp, arg.RegistryID, arg.RegistryName)
	return err
}

const findAppBySourceAndId = `-- name: FindAppBySourceAndId :one
SELECT 
    app.app_index AS "index",
    app.app_id AS "id",
    registry.registry_name AS "name",
    app.app_installed AS "installed",
    app.app_version AS "version",
    app.app_available AS "available",
    app.app_last_updated AS "last_updated"
FROM app
JOIN registry ON app.app_registry_id = registry.registry_id
WHERE app.app_source = ? 
    AND app.app_id = ? 
    AND app.app_os = ?
LIMIT 1
`

type FindAppBySourceAndIdParams struct {
	AppSource string
	AppID     string
	AppOs     string
}

type FindAppBySourceAndIdRow struct {
	Index       int64
	ID          string
	Name        string
	Installed   bool
	Version     sql.NullString
	Available   sql.NullString
	LastUpdated string
}

func (q *Queries) FindAppBySourceAndId(ctx context.Context, arg FindAppBySourceAndIdParams) (FindAppBySourceAndIdRow, error) {
	row := q.db.QueryRowContext(ctx, findAppBySourceAndId, arg.AppSource, arg.AppID, arg.AppOs)
	var i FindAppBySourceAndIdRow
	err := row.Scan(
		&i.Index,
		&i.ID,
		&i.Name,
		&i.Installed,
		&i.Version,
		&i.Available,
		&i.LastUpdated,
	)
	return i, err
}

const findAppsByNamesGivenText = `-- name: FindAppsByNamesGivenText :many
SELECT 
    a.app_index AS "index",
    a.app_id AS "id",
    r.registry_name AS "name",
    a.app_installed AS "installed",
    a.app_source AS "source",
    a.app_version AS "version",
    a.app_available AS "available",
    a.app_last_updated AS "last_updated"
FROM app a
JOIN registry r ON a.app_registry_id = r.registry_id
WHERE r.registry_name IN (
    SELECT rs.registry_name FROM registry_search rs WHERE rs.registry_name MATCH ?
)
AND a.app_os = ?
`

type FindAppsByNamesGivenTextParams struct {
	RegistryName string
	AppOs        string
}

type FindAppsByNamesGivenTextRow struct {
	Index       int64
	ID          string
	Name        string
	Installed   bool
	Source      string
	Version     sql.NullString
	Available   sql.NullString
	LastUpdated string
}

func (q *Queries) FindAppsByNamesGivenText(ctx context.Context, arg FindAppsByNamesGivenTextParams) ([]FindAppsByNamesGivenTextRow, error) {
	rows, err := q.db.QueryContext(ctx, findAppsByNamesGivenText, arg.RegistryName, arg.AppOs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []FindAppsByNamesGivenTextRow
	for rows.Next() {
		var i FindAppsByNamesGivenTextRow
		if err := rows.Scan(
			&i.Index,
			&i.ID,
			&i.Name,
			&i.Installed,
			&i.Source,
			&i.Version,
			&i.Available,
			&i.LastUpdated,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const findAppsByText = `-- name: FindAppsByText :many
SELECT registry_name FROM registry_search WHERE registry_name MATCH ?
`

func (q *Queries) FindAppsByText(ctx context.Context, registryName string) ([]string, error) {
	rows, err := q.db.QueryContext(ctx, findAppsByText, registryName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []string
	for rows.Next() {
		var registry_name string
		if err := rows.Scan(&registry_name); err != nil {
			return nil, err
		}
		items = append(items, registry_name)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const findInstalledApps = `-- name: FindInstalledApps :many
SELECT 
    a.app_index AS "index",
    a.app_id AS "id",
    r.registry_name AS "name",
    a.app_os AS "os",
    a.app_source AS "source",
    a.app_version AS "version",
    a.app_available AS "available",
    a.app_last_updated AS "last_updated"
FROM app a
JOIN registry r ON a.app_registry_id = r.registry_id
WHERE a.app_installed = 1
`

type FindInstalledAppsRow struct {
	Index       int64
	ID          string
	Name        string
	Os          string
	Source      string
	Version     sql.NullString
	Available   sql.NullString
	LastUpdated string
}

func (q *Queries) FindInstalledApps(ctx context.Context) ([]FindInstalledAppsRow, error) {
	rows, err := q.db.QueryContext(ctx, findInstalledApps)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []FindInstalledAppsRow
	for rows.Next() {
		var i FindInstalledAppsRow
		if err := rows.Scan(
			&i.Index,
			&i.ID,
			&i.Name,
			&i.Os,
			&i.Source,
			&i.Version,
			&i.Available,
			&i.LastUpdated,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const findInstalledAppsOS = `-- name: FindInstalledAppsOS :many
SELECT 
    a.app_index AS "index",
    a.app_id AS "id",
    r.registry_name AS "name",
    a.app_source AS "source",
    a.app_version AS "version",
    a.app_available AS "available",
    a.app_last_updated AS "last_updated"
FROM app a
JOIN registry r ON a.app_registry_id = r.registry_id
WHERE a.app_installed = 1
AND a.app_os = ?
`

type FindInstalledAppsOSRow struct {
	Index       int64
	ID          string
	Name        string
	Source      string
	Version     sql.NullString
	Available   sql.NullString
	LastUpdated string
}

func (q *Queries) FindInstalledAppsOS(ctx context.Context, appOs string) ([]FindInstalledAppsOSRow, error) {
	rows, err := q.db.QueryContext(ctx, findInstalledAppsOS, appOs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []FindInstalledAppsOSRow
	for rows.Next() {
		var i FindInstalledAppsOSRow
		if err := rows.Scan(
			&i.Index,
			&i.ID,
			&i.Name,
			&i.Source,
			&i.Version,
			&i.Available,
			&i.LastUpdated,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const findNotInstalledApps = `-- name: FindNotInstalledApps :many
SELECT 
    a.app_index AS "index", 
    a.app_id AS "id", 
    r.registry_name AS "name", 
    a.app_os AS "os",
    a.app_source AS "source",
    a.app_version AS "version", 
    a.app_available AS "available",
    a.app_last_updated AS "last_updated"
FROM app a
JOIN registry r ON a.app_registry_id = r.registry_id
WHERE a.app_installed = 0
`

type FindNotInstalledAppsRow struct {
	Index       int64
	ID          string
	Name        string
	Os          string
	Source      string
	Version     sql.NullString
	Available   sql.NullString
	LastUpdated string
}

func (q *Queries) FindNotInstalledApps(ctx context.Context) ([]FindNotInstalledAppsRow, error) {
	rows, err := q.db.QueryContext(ctx, findNotInstalledApps)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []FindNotInstalledAppsRow
	for rows.Next() {
		var i FindNotInstalledAppsRow
		if err := rows.Scan(
			&i.Index,
			&i.ID,
			&i.Name,
			&i.Os,
			&i.Source,
			&i.Version,
			&i.Available,
			&i.LastUpdated,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const findNotInstalledAppsOS = `-- name: FindNotInstalledAppsOS :many
SELECT 
    a.app_index AS "index", 
    a.app_id AS "id", 
    r.registry_name AS "name",
    a.app_source AS "source",
    a.app_version AS "version", 
    a.app_available AS "available",
    a.app_last_updated AS "last_updated"
FROM app a
JOIN registry r ON a.app_registry_id = r.registry_id
WHERE a.app_installed = 0
AND a.app_os = ?
`

type FindNotInstalledAppsOSRow struct {
	Index       int64
	ID          string
	Name        string
	Source      string
	Version     sql.NullString
	Available   sql.NullString
	LastUpdated string
}

func (q *Queries) FindNotInstalledAppsOS(ctx context.Context, appOs string) ([]FindNotInstalledAppsOSRow, error) {
	rows, err := q.db.QueryContext(ctx, findNotInstalledAppsOS, appOs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []FindNotInstalledAppsOSRow
	for rows.Next() {
		var i FindNotInstalledAppsOSRow
		if err := rows.Scan(
			&i.Index,
			&i.ID,
			&i.Name,
			&i.Source,
			&i.Version,
			&i.Available,
			&i.LastUpdated,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const installApp = `-- name: InstallApp :exec
INSERT INTO app (
    app_id, 
    app_source, 
    app_os, 
    app_registry_id, 
    app_last_updated, 
    app_installed, 
    app_version, 
    app_available
) 
VALUES (
    ?, 
    ?, 
    ?, 
    ?, 
    datetime('now', 'utc'), 
    1, 
    ?, 
    NULL
)
`

type InstallAppParams struct {
	AppID         string
	AppSource     string
	AppOs         string
	AppRegistryID string
	AppVersion    sql.NullString
}

func (q *Queries) InstallApp(ctx context.Context, arg InstallAppParams) error {
	_, err := q.db.ExecContext(ctx, installApp,
		arg.AppID,
		arg.AppSource,
		arg.AppOs,
		arg.AppRegistryID,
		arg.AppVersion,
	)
	return err
}

const listApps = `-- name: ListApps :many
SELECT 
    a.app_index AS "index", 
    a.app_id AS "id", 
    r.registry_name AS "name", 
    a.app_os AS "os",
    a.app_installed AS "installed",
    a.app_source AS "source",
    a.app_version AS "version", 
    a.app_available AS "available",
    a.app_last_updated AS "last_updated"
FROM app a
JOIN registry r ON a.app_registry_id = r.registry_id
`

type ListAppsRow struct {
	Index       int64
	ID          string
	Name        string
	Os          string
	Installed   bool
	Source      string
	Version     sql.NullString
	Available   sql.NullString
	LastUpdated string
}

func (q *Queries) ListApps(ctx context.Context) ([]ListAppsRow, error) {
	rows, err := q.db.QueryContext(ctx, listApps)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListAppsRow
	for rows.Next() {
		var i ListAppsRow
		if err := rows.Scan(
			&i.Index,
			&i.ID,
			&i.Name,
			&i.Os,
			&i.Installed,
			&i.Source,
			&i.Version,
			&i.Available,
			&i.LastUpdated,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listAppsOS = `-- name: ListAppsOS :many
SELECT 
    a.app_index AS "index",
    a.app_id AS "id", 
    r.registry_name AS "name",
    a.app_installed AS "installed",
    a.app_source AS "source",
    a.app_version AS "version", 
    a.app_available AS "available",
    a.app_last_updated AS "last_updated"
FROM app a
JOIN registry r ON a.app_registry_id = r.registry_id
WHERE a.app_os = ?
`

type ListAppsOSRow struct {
	Index       int64
	ID          string
	Name        string
	Installed   bool
	Source      string
	Version     sql.NullString
	Available   sql.NullString
	LastUpdated string
}

func (q *Queries) ListAppsOS(ctx context.Context, appOs string) ([]ListAppsOSRow, error) {
	rows, err := q.db.QueryContext(ctx, listAppsOS, appOs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListAppsOSRow
	for rows.Next() {
		var i ListAppsOSRow
		if err := rows.Scan(
			&i.Index,
			&i.ID,
			&i.Name,
			&i.Installed,
			&i.Source,
			&i.Version,
			&i.Available,
			&i.LastUpdated,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listUpgradableApps = `-- name: ListUpgradableApps :many
SELECT 
    a.app_index AS "index", 
    a.app_id AS "id", 
    r.registry_name AS "name",
    a.app_os AS "os",
    a.app_source AS "source",
    a.app_version AS "version", 
    a.app_available AS "available",
    a.app_last_updated AS "last_updated"
FROM app a
JOIN registry r ON a.app_registry_id = r.registry_id
WHERE a.app_installed = 1
AND a.app_available IS NOT NULL
`

type ListUpgradableAppsRow struct {
	Index       int64
	ID          string
	Name        string
	Os          string
	Source      string
	Version     sql.NullString
	Available   sql.NullString
	LastUpdated string
}

func (q *Queries) ListUpgradableApps(ctx context.Context) ([]ListUpgradableAppsRow, error) {
	rows, err := q.db.QueryContext(ctx, listUpgradableApps)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListUpgradableAppsRow
	for rows.Next() {
		var i ListUpgradableAppsRow
		if err := rows.Scan(
			&i.Index,
			&i.ID,
			&i.Name,
			&i.Os,
			&i.Source,
			&i.Version,
			&i.Available,
			&i.LastUpdated,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listUpgradableAppsOS = `-- name: ListUpgradableAppsOS :many
SELECT 
    a.app_index AS "index",
    a.app_id AS "id",
    r.registry_name AS "name",
    a.app_source AS "source",
    a.app_version AS "version",
    a.app_available AS "available",
    a.app_last_updated AS "last_updated"
FROM app a
JOIN registry r ON a.app_registry_id = r.registry_id
WHERE a.app_installed = 1 
AND a.app_os = ?
AND a.app_available IS NOT NULL
`

type ListUpgradableAppsOSRow struct {
	Index       int64
	ID          string
	Name        string
	Source      string
	Version     sql.NullString
	Available   sql.NullString
	LastUpdated string
}

func (q *Queries) ListUpgradableAppsOS(ctx context.Context, appOs string) ([]ListUpgradableAppsOSRow, error) {
	rows, err := q.db.QueryContext(ctx, listUpgradableAppsOS, appOs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListUpgradableAppsOSRow
	for rows.Next() {
		var i ListUpgradableAppsOSRow
		if err := rows.Scan(
			&i.Index,
			&i.ID,
			&i.Name,
			&i.Source,
			&i.Version,
			&i.Available,
			&i.LastUpdated,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const syncRegistrySearchApps = `-- name: SyncRegistrySearchApps :exec
INSERT INTO registry_search (registry_name) SELECT registry_name FROM registry
`

func (q *Queries) SyncRegistrySearchApps(ctx context.Context) error {
	_, err := q.db.ExecContext(ctx, syncRegistrySearchApps)
	return err
}

const updateAvailable = `-- name: UpdateAvailable :exec
UPDATE app
SET
    app_available = ?
WHERE app_index = ?
`

type UpdateAvailableParams struct {
	AppAvailable sql.NullString
	AppIndex     int64
}

func (q *Queries) UpdateAvailable(ctx context.Context, arg UpdateAvailableParams) error {
	_, err := q.db.ExecContext(ctx, updateAvailable, arg.AppAvailable, arg.AppIndex)
	return err
}

const updateInstalledApp = `-- name: UpdateInstalledApp :exec
UPDATE app
SET 
    app_installed = 1,
    app_version = ?,
    app_available = ?,
    app_last_updated = datetime('now', 'utc')
WHERE app_index = ?
`

type UpdateInstalledAppParams struct {
	AppVersion   sql.NullString
	AppAvailable sql.NullString
	AppIndex     int64
}

func (q *Queries) UpdateInstalledApp(ctx context.Context, arg UpdateInstalledAppParams) error {
	_, err := q.db.ExecContext(ctx, updateInstalledApp, arg.AppVersion, arg.AppAvailable, arg.AppIndex)
	return err
}

const updateName = `-- name: UpdateName :exec
UPDATE registry
SET registry_name = ?
WHERE registry_id = (
    SELECT app_registry_id FROM app WHERE app_index = ?
)
`

type UpdateNameParams struct {
	RegistryName string
	AppIndex     int64
}

func (q *Queries) UpdateName(ctx context.Context, arg UpdateNameParams) error {
	_, err := q.db.ExecContext(ctx, updateName, arg.RegistryName, arg.AppIndex)
	return err
}

const updateSource = `-- name: UpdateSource :exec
UPDATE app
SET
    app_id = ?,
    app_source = ?,
    app_available = ?
WHERE app_index = ? AND app_installed = 0
`

type UpdateSourceParams struct {
	AppID        string
	AppSource    string
	AppAvailable sql.NullString
	AppIndex     int64
}

func (q *Queries) UpdateSource(ctx context.Context, arg UpdateSourceParams) error {
	_, err := q.db.ExecContext(ctx, updateSource,
		arg.AppID,
		arg.AppSource,
		arg.AppAvailable,
		arg.AppIndex,
	)
	return err
}

const updateUninstalledApp = `-- name: UpdateUninstalledApp :exec
UPDATE app
SET 
    app_installed = 0,
    app_version = NULL,
    app_available = ?,
    app_last_updated = datetime('now', 'utc')
WHERE app_index = ?
`

type UpdateUninstalledAppParams struct {
	AppAvailable sql.NullString
	AppIndex     int64
}

func (q *Queries) UpdateUninstalledApp(ctx context.Context, arg UpdateUninstalledAppParams) error {
	_, err := q.db.ExecContext(ctx, updateUninstalledApp, arg.AppAvailable, arg.AppIndex)
	return err
}
