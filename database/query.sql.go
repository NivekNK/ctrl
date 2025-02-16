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
    app_status, 
    app_version, 
    app_available
) 
VALUES (
    ?, 
    ?, 
    ?, 
    ?, 
    datetime('now', 'utc'), 
    'not-installed', 
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

const findInProgressApps = `-- name: FindInProgressApps :many
SELECT 
    a.app_index AS "index", 
    a.app_id AS "id", 
    r.registry_name AS "name",
    a.app_os AS "os", 
    a.app_version AS "version", 
    a.app_available AS "available",
    a.app_last_updated AS "last_updated"
FROM app a
JOIN registry r ON a.app_registry_id = r.registry_id
WHERE a.app_status = 'in-progress'
`

type FindInProgressAppsRow struct {
	Index       int64
	ID          string
	Name        string
	Os          string
	Version     sql.NullString
	Available   sql.NullString
	LastUpdated string
}

func (q *Queries) FindInProgressApps(ctx context.Context) ([]FindInProgressAppsRow, error) {
	rows, err := q.db.QueryContext(ctx, findInProgressApps)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []FindInProgressAppsRow
	for rows.Next() {
		var i FindInProgressAppsRow
		if err := rows.Scan(
			&i.Index,
			&i.ID,
			&i.Name,
			&i.Os,
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

const findInProgressAppsOS = `-- name: FindInProgressAppsOS :many
SELECT 
    a.app_index AS "index", 
    a.app_id AS "id", 
    r.registry_name AS "name", 
    a.app_version AS "version", 
    a.app_available AS "available",
    a.app_last_updated AS "last_updated"
FROM app a
JOIN registry r ON a.app_registry_id = r.registry_id
WHERE a.app_status = 'in-progress'
AND a.app_os = ?
`

type FindInProgressAppsOSRow struct {
	Index       int64
	ID          string
	Name        string
	Version     sql.NullString
	Available   sql.NullString
	LastUpdated string
}

func (q *Queries) FindInProgressAppsOS(ctx context.Context, appOs string) ([]FindInProgressAppsOSRow, error) {
	rows, err := q.db.QueryContext(ctx, findInProgressAppsOS, appOs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []FindInProgressAppsOSRow
	for rows.Next() {
		var i FindInProgressAppsOSRow
		if err := rows.Scan(
			&i.Index,
			&i.ID,
			&i.Name,
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

const findInstalledApps = `-- name: FindInstalledApps :many
SELECT 
    a.app_index AS "index",
    a.app_id AS "id",
    r.registry_name AS "name",
    a.app_os AS "os",
    a.app_version AS "version",
    a.app_available AS "available",
    a.app_last_updated AS "last_updated"
FROM app a
JOIN registry r ON a.app_registry_id = r.registry_id
WHERE a.app_status = 'installed'
`

type FindInstalledAppsRow struct {
	Index       int64
	ID          string
	Name        string
	Os          string
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
    a.app_version AS "version", 
    a.app_available AS "available",
    a.app_last_updated AS "last_updated"
FROM app a
JOIN registry r ON a.app_registry_id = r.registry_id
WHERE a.app_status = 'installed'
AND a.app_os = ?
`

type FindInstalledAppsOSRow struct {
	Index       int64
	ID          string
	Name        string
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
    a.app_version AS "version", 
    a.app_available AS "available",
    a.app_last_updated AS "last_updated"
FROM app a
JOIN registry r ON a.app_registry_id = r.registry_id
WHERE a.app_status = 'not-installed'
`

type FindNotInstalledAppsRow struct {
	Index       int64
	ID          string
	Name        string
	Os          string
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
    a.app_version AS "version", 
    a.app_available AS "available",
    a.app_last_updated AS "last_updated"
FROM app a
JOIN registry r ON a.app_registry_id = r.registry_id
WHERE a.app_status = 'not-installed'
AND a.app_os = ?
`

type FindNotInstalledAppsOSRow struct {
	Index       int64
	ID          string
	Name        string
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

const listApps = `-- name: ListApps :many
SELECT 
    a.app_index AS "index", 
    a.app_id AS "id", 
    r.registry_name AS "name", 
    a.app_os AS "os",
    a.app_status AS "status",
    a.app_version AS "version", 
    a.app_available AS "available"
FROM app a
JOIN registry r ON a.app_registry_id = r.registry_id
`

type ListAppsRow struct {
	Index     int64
	ID        string
	Name      string
	Os        string
	Status    string
	Version   sql.NullString
	Available sql.NullString
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
			&i.Status,
			&i.Version,
			&i.Available,
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
    a.app_status AS "status",
    a.app_version AS "version", 
    a.app_available AS "available"
FROM app a
JOIN registry r ON a.app_registry_id = r.registry_id
WHERE a.app_os = ?
`

type ListAppsOSRow struct {
	Index     int64
	ID        string
	Name      string
	Status    string
	Version   sql.NullString
	Available sql.NullString
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
			&i.Status,
			&i.Version,
			&i.Available,
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
    a.app_version AS "version", 
    a.app_available AS "available",
    a.app_last_updated AS "last_updated"
FROM app a
JOIN registry r ON a.app_registry_id = r.registry_id
WHERE a.app_status = 'installed'
AND a.app_available IS NOT NULL
`

type ListUpgradableAppsRow struct {
	Index       int64
	ID          string
	Name        string
	Os          string
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
    a.app_version AS "version", 
    a.app_available AS "available",
    a.app_last_updated AS "last_updated"
FROM app a
JOIN registry r ON a.app_registry_id = r.registry_id
WHERE a.app_status = 'installed' 
AND a.app_os = ?
AND a.app_available IS NOT NULL
`

type ListUpgradableAppsOSRow struct {
	Index       int64
	ID          string
	Name        string
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
