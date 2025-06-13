CREATE TABLE connections (
    connection_id     INTEGER NOT NULL PRIMARY KEY,
    db_type           TEXT    NOT NULL,
    connection_url    TEXT    NOT NULL,
    created_at        INTEGER NOT NULL DEFAULT (unixepoch()),
    deleted_at        INTEGER,

    connected         INTEGER NOT NULL DEFAULT FALSE,
    last_connected_at INTEGER,
    last_error        TEXT
) STRICT;

-- INSERT INTO connections (db_type, connection_url) VALUES
-- ('sqlite', './data/giveaway.db');