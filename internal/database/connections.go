package database

import (
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/pkg/errors"
	"github.com/yaitoo/sqle"
)

type Connections map[int]*sqle.DB

func pushToConnections(db *Database, con Connection, connections Connections) error {
	var con_db *sqle.DB
	var sqldb *sql.DB
	var err error

	switch con.DbType {
	case "sqlite":
		sqldb, err = sql.Open("sqlite3", "file:"+con.ConnectionUrl)
		if err == nil {
			con_db = sqle.Open(sqldb)
			_, err = db.Exec("UPDATE connections SET connected = true, last_error = NULL WHERE connection_id = ?", con.ID)
		} else {
			_, err = db.Exec("UPDATE connections SET connected = false, last_error = ? WHERE connection_id = ?", err.Error(), con.ID)
			if err == nil {
				err = errors.Errorf("failed to connect to connection %d: %s", con.ID, err)
			}
		}
	default:
		err = fmt.Errorf("unknown database type %s", con.DbType)
	}
	if err != nil {
		return errors.Wrap(err, "Error during the connection process")
	}

	slog.Info("Connected to connection", slog.Int("id", con.ID), slog.String("type", con.DbType), slog.String("url", con.ConnectionUrl))
	connections[con.ID] = con_db
	return nil
}

func AddConnection(db *Database, connections Connections, id int) error {
	var con Connection
	err := db.QueryRow("SELECT * FROM connections WHERE connection_id = ?", id).Scan(&con)
	if err != nil {
		return errors.Wrap(err, "Error during the connection process")
	}

	return pushToConnections(db, con, connections)
}

func RemoveConnection(db *Database, connections Connections, id int) error {
	con, ok := connections[id]
	if !ok {
		return errors.Errorf("connection %d not found", id)
	}
	err := con.Close()
	if err != nil {
		slog.Error("Failed to close connection", slog.Int("id", id), slog.Any("error", err))
	}

	delete(connections, id)

	_, err = db.Exec("UPDATE connections SET connected = false, last_error = 'Connection removed', deleted_at = CURRENT_TIMESTAMP WHERE connection_id = ?", id)
	return errors.Wrap(err, "Error saving the disconnection")
}

func SetupConnections(db *Database) (Connections, error) {
	slog.Info("Setting up connections...")
	connections := make(Connections)

	var connection_list []Connection
	rows, err := db.Query("SELECT * FROM connections WHERE deleted_at IS NULL")
	if err != nil {
		return nil, err
	}
	if err := rows.Bind(&connection_list); err != nil {
		return nil, err
	}

	for _, con := range connection_list {
		err := pushToConnections(db, con, connections)
		if err != nil {
			slog.Error("Failed to connect to connection", slog.Int("id", con.ID), slog.Any("error", err))
		}
	}

	return connections, nil
}
