CREATE TABLE IF NOT EXISTS nodes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    display_name TEXT,
    address TEXT NOT NULL,
    ssh_port INTEGER NOT NULL,
    ssh_user TEXT DEFAULT 'root',
    ssh_password TEXT NOT NULL,
    agent_port INTEGER DEFAULT 0,
    status TEXT DEFAULT 'active',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    is_local BOOLEAN DEFAULT FALSE,
    metadata TEXT
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_nodes_name ON nodes (name);