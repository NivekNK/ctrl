-- name: ListApps :many
SELECT 
    a.app_index AS "index", 
    a.app_id AS "id", 
    r.registry_name AS "name", 
    a.app_os AS "os",
    a.app_status AS "status",
    a.app_version AS "version", 
    a.app_available AS "available"
FROM app a
JOIN registry r ON a.app_registry_id = r.registry_id;

-- name: ListAppsOS :many
SELECT 
    a.app_index AS "index",
    a.app_id AS "id", 
    r.registry_name AS "name",
    a.app_status AS "status",
    a.app_version AS "version", 
    a.app_available AS "available"
FROM app a
JOIN registry r ON a.app_registry_id = r.registry_id
WHERE a.app_os = ?;

-- name: FindNotInstalledApps :many
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
WHERE a.app_status = 'not-installed';

-- name: FindNotInstalledAppsOS :many
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
AND a.app_os = ?;

-- name: FindInProgressApps :many
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
WHERE a.app_status = 'in-progress';

-- name: FindInProgressAppsOS :many
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
AND a.app_os = ?;

-- name: FindInstalledApps :many
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
WHERE a.app_status = 'installed';

-- name: FindInstalledAppsOS :many
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
AND a.app_os = ?;

-- name: ListUpgradableApps :many
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
AND a.app_available IS NOT NULL;

-- name: ListUpgradableAppsOS :many
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
AND a.app_available IS NOT NULL;

-- name: AddRegistryApp :exec
INSERT INTO registry (registry_id, registry_name) VALUES (?, ?);

-- name: AddApp :exec
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
);
