CREATE TABLE registry (
    registry_id   TEXT PRIMARY KEY,
    registry_name TEXT UNIQUE NOT NULL
);

CREATE VIRTUAL TABLE registry_search USING fts5(registry_name);

CREATE TABLE app (
    app_index        INTEGER PRIMARY KEY,
    app_id           TEXT NOT NULL,
    app_source       TEXT NOT NULL,
    app_os           TEXT CHECK(app_os IN ('windows', 'linux')) NOT NULL,
    app_registry_id  TEXT NOT NULL,
    app_last_updated TEXT DEFAULT (datetime('now', 'utc')) NOT NULL,
    app_installed    BOOLEAN DEFAULT 0 NOT NULL,
    app_version      TEXT,
    app_available    TEXT,
    FOREIGN KEY (app_registry_id) REFERENCES registry(registry_id)
);
