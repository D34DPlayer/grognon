CREATE TABLE connections (
    connection_id     INTEGER     NOT NULL PRIMARY KEY,
    db_type           VARCHAR(32) NOT NULL,
    connection_url    TEXT        NOT NULL,

    created_at        TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at        TIMESTAMP,

    connected         BOOLEAN     NOT NULL DEFAULT FALSE,
    last_connected_at TIMESTAMP,
    last_error        TEXT
);
