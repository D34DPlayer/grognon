CREATE TABLE connections (
    connection_id     INTEGER     NOT NULL PRIMARY KEY,
    db_type           VARCHAR(32) NOT NULL,
    connection_url    TEXT        NOT NULL,
    created_at        INTEGER   NOT NULL DEFAULT (unixepoch()),
    deleted_at        INTEGER,

    connected         BOOLEAN     NOT NULL DEFAULT FALSE,
    last_connected_at INTEGER,
    last_error        TEXT
);

INSERT INTO connections (db_type, connection_url) VALUES
('sqlite', './data/giveaway.db');