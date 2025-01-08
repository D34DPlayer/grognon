package database

import (
	"database/sql"
	"log"

	"github.com/yaitoo/sqle"
)

type Connections map[int]*sqle.DB

func OnConnection(db *Database, id int, og_error error) error {
	if og_error == nil {
		log.Printf("Connected to connection %d", id)

		_, err := db.Exec("UPDATE connections SET connected = true, last_error = NULL WHERE id = ?", id)
		return err
	} else {
		log.Printf("Failed to connect to connection %d: %s", id, og_error)

		_, err := db.Exec("UPDATE connections SET connected = false, last_error = ? WHERE id = ?", og_error.Error(), id)
		return err
	}
}

func SetupConnections(db *Database) (Connections, error) {
	connections := make(Connections)

	var connection_list []Connection
	rows, err := db.Query("SELECT * FROM connections")
	if err != nil {
		return nil, err
	}
	if err := rows.Bind(&connection_list); err != nil {
		return nil, err
	}

	for _, con := range connection_list {
		var con_db *sqle.DB
		switch con.DbType {
		case "sqlite":
			sqldb, err := sql.Open("sqlite3", "file:"+con.ConnectionUrl)
			OnConnection(db, con.ID, err)
			if err == nil {
				con_db = sqle.Open(sqldb)
			}
		}
		if con_db != nil {
			connections[con.ID] = con_db
		}
	}

	return connections, nil
}
