CREATE TABLE crons (
    cron_id       INTEGER   NOT NULL PRIMARY KEY,
    connection_id INTEGER   NOT NULL REFERENCES connections (connection_id),
    name          TEXT      NOT NULL,
    command       TEXT      NOT NULL,
    schedule      TEXT      NOT NULL,

    created_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at    TIMESTAMP,
    last_run_at   TIMESTAMP
);

CREATE TABLE cron_outputs (
    cron_id INTEGER NOT NULL REFERENCES crons (cron_id),
    name    TEXT    NOT NULL,
    type   TEXT    NOT NULL,

    CONSTRAINT cron_outputs_pk PRIMARY KEY (cron_id, name)
);