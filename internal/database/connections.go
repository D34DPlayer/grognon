package database

import (
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/pkg/errors"
	"github.com/yaitoo/sqle"
)

type Connections map[int64]*sqle.DB

func pushToConnections(db *Database, con Connection, connections Connections) error {
	var con_db *sqle.DB
	var sqldb *sql.DB
	var err error

	switch con.DbType {
	case "sqlite":
		sqldb, err = sql.Open("sqlite3", "file:"+con.ConnectionUrl)
		if err == nil {
			con_db = sqle.Open(sqldb)
			err = con_db.Ping()
		}

		if err == nil {
			_, err = db.Exec("UPDATE connections SET connected = true, last_error = NULL, last_connected_at = unixepoch() WHERE connection_id = ?", con.ConnectionId)
		} else {
			_, saveErr := db.Exec("UPDATE connections SET connected = false, last_error = ? WHERE connection_id = ?", err.Error(), con.ConnectionId)
			if saveErr != nil {
				err = errors.Wrap(saveErr, "failed to save connection error")
			}
		}
	default:
		err = fmt.Errorf("unknown database type %s", con.DbType)
	}
	if err != nil {
		return errors.Wrap(err, "Error during the connection process")
	}

	slog.Info("Connected to connection", slog.Int64("id", con.ConnectionId), slog.String("type", con.DbType), slog.String("url", con.ConnectionUrl))
	connections[con.ConnectionId] = con_db
	return nil
}

func GetConnection(db *Database, id int64) (*Connection, error) {
	var con Connection
	err := db.QueryRow("SELECT * FROM connections WHERE connection_id = ? AND deleted_at IS NULL", id).
		Bind(&con)
	if err != nil {
		return nil, errors.Wrap(err, "Error getting connection")
	}
	return &con, nil
}

func AddConnection(db *Database, connections Connections, input ConnectionCreate) (*Connection, error) {
	res, err := db.Exec("INSERT INTO connections (db_type, connection_url) VALUES (?, ?)",
		input.DbType, input.ConnectionUrl)
	if err != nil {
		return nil, errors.Wrap(err, "Error during the connection process")
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, errors.Wrap(err, "Error getting last insert id")
	}

	con, err := GetConnection(db, id)
	if err != nil {
		return nil, errors.Wrap(err, "Error getting connection")
	}
	err = pushToConnections(db, *con, connections)
	if err != nil {
		RemoveConnection(db, connections, id)
		return nil, errors.Wrap(err, "Error pushing connection")
	}
	return con, nil
}

func RemoveConnection(db *Database, connections Connections, id int64) error {
	_, err := db.Exec("UPDATE connections SET connected = false, last_error = 'Connection removed', deleted_at = unixepoch() WHERE connection_id = ?", id)
	if err != nil {
		return errors.Wrap(err, "Error saving the disconnection")
	}

	con, ok := connections[id]
	if !ok {
		slog.Error("Connection not found", slog.Int64("id", id))
		return nil
	}
	delete(connections, id)

	err = con.Close()
	if err != nil {
		slog.Error("Failed to close connection", slog.Int64("id", id), slog.Any("error", err))
	}
	return nil
}

func GetConnections(db *Database) ([]Connection, error) {
	var connection_list []Connection
	rows, err := db.Query("SELECT * FROM connections WHERE deleted_at IS NULL")
	if err != nil {
		return nil, err
	}
	if err := rows.Bind(&connection_list); err != nil {
		return nil, err
	}
	return connection_list, nil
}

func SetupConnections(db *Database) (Connections, error) {
	slog.Info("Setting up connections...")
	connections := make(Connections)

	connection_list, err := GetConnections(db)
	if err != nil {
		return nil, errors.Wrap(err, "Error getting connections")
	}

	for _, con := range connection_list {
		err := pushToConnections(db, con, connections)
		if err != nil {
			slog.Error("Failed to connect to connection", slog.Int64("id", con.ConnectionId), slog.Any("error", err))
		}
	}

	return connections, nil
}

func RefreshConnections(db *Database, connections Connections) error {
	slog.Info("Refreshing connections...")
	connection_list, err := GetConnections(db)
	if err != nil {
		return errors.Wrap(err, "Error getting connections")
	}

	for _, con := range connection_list {
		dbCon, ok := connections[con.ConnectionId]
		if ok {
			err := dbCon.Ping()
			if err != nil {
				slog.Error("Connection ping failed", slog.Int64("id", con.ConnectionId), slog.Any("error", err))
				ok = false
			}
		}
		if !ok {
			err := pushToConnections(db, con, connections)
			if err != nil {
				slog.Error("Failed to connect to connection", slog.Int64("id", con.ConnectionId), slog.Any("error", err))
			}
		}
	}

	return nil
}
