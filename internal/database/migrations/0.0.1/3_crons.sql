CREATE TABLE crons (
    cron_id       INTEGER   NOT NULL PRIMARY KEY,
    connection_id INTEGER   NOT NULL REFERENCES connections (connection_id),
    name          TEXT      NOT NULL,
    command       TEXT      NOT NULL,
    schedule      TEXT      NOT NULL,

    created_at    INTEGER   NOT NULL DEFAULT (unixepoch()),
    deleted_at    INTEGER,
    last_run_at   INTEGER
) STRICT;

CREATE TABLE cron_outputs (
    cron_id INTEGER NOT NULL REFERENCES crons (cron_id),
    name    TEXT    NOT NULL,
    type   TEXT    NOT NULL,

    CONSTRAINT cron_outputs_pk PRIMARY KEY (cron_id, name)
) STRICT;