CREATE TABLE IF NOT EXISTS apps (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    display TEXT, -- display name
    description TEXT, -- description
    version TEXT NOT NULL,
    icon TEXT,
    qa TEXT,
    docker_compose TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    metadata TEXT
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_apps_name ON apps (name);


CREATE TABLE IF NOT EXISTS services (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    app_id INTEGER NOT NULL,
    node_id INTEGER NOT NULL,
    status TEXT DEFAULT 'active',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    metadata TEXT
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_services_name ON services (name);
CREATE INDEX IF NOT EXISTS idx_services_app_id_node_id ON services (app_id, node_id);