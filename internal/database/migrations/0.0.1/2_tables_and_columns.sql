CREATE TABLE tables (
    connection_id INTEGER NOT NULL REFERENCES connections (connection_id),
    table_name    TEXT    NOT NULL,

    CONSTRAINT tables_pk PRIMARY KEY (connection_id, table_name)
);

CREATE TABLE columns (
    connection_id INTEGER NOT NULL REFERENCES connections (connection_id),
    table_name    TEXT    NOT NULL REFERENCES tables (table_name),
    name          TEXT    NOT NULL,
    type          TEXT    NOT NULL,
    "notnull"     BOOLEAN NOT NULL,
    dflt_value    TEXT,
    pk            INTEGER NOT NULL DEFAULT 0,

    CONSTRAINT columns_pk PRIMARY KEY (connection_id, table_name, name)
)
