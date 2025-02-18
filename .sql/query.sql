-- name: ListApps :many
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
JOIN registry r ON a.app_registry_id = r.registry_id;

-- name: ListAppsOS :many
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
WHERE a.app_os = ?;

-- name: FindNotInstalledApps :many
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
WHERE a.app_installed = 0;

-- name: FindNotInstalledAppsOS :many
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
AND a.app_os = ?;

-- name: FindInstalledApps :many
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
WHERE a.app_installed = 1;

-- name: FindInstalledAppsOS :many
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
AND a.app_os = ?;

-- name: ListUpgradableApps :many
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
AND a.app_available IS NOT NULL;

-- name: ListUpgradableAppsOS :many
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
AND a.app_available IS NOT NULL;

-- name: AddRegistryApp :exec
INSERT INTO registry (registry_id, registry_name) VALUES (?, ?);

-- name: SyncRegistrySearchApps :exec
INSERT INTO registry_search (registry_name) SELECT registry_name FROM registry;

-- name: AddApp :exec
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
);

-- name: FindAppBySourceAndId :one
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
LIMIT 1;

-- name: FindAppsByText :many
SELECT registry_name FROM registry_search WHERE registry_name MATCH ?;

-- name: FindAppsByNamesGivenText :many
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
AND a.app_os = ?;

-- name: InstalledApp :exec
UPDATE app
SET 
    app_installed = 1,
    app_version = ?,
    app_available = NULL,
    app_last_updated = datetime('now', 'utc')
WHERE app_index = ?;

-- name: UninstalledApp :exec
UPDATE app
SET 
    app_installed = 0,
    app_version = NULL,
    app_available = ?,
    app_last_updated = datetime('now', 'utc')
WHERE app_index = ?;
