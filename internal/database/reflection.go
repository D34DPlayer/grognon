package database

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/pkg/errors"
	"github.com/yaitoo/sqle"
)

func reflectSqlite(db *Database, con_id int64, con *sqle.DB) error {
	slog.Info("Reflecting connection", slog.Int64("id", con_id))

	tx, err := db.BeginTx(context.TODO(), nil)
	if err != nil {
		return errors.Wrap(err, "failed to start transaction")
	}

	var table_names []string
	rows, err := con.Query("SELECT name FROM sqlite_schema WHERE type ='table' AND name NOT LIKE 'sqlite_%'")
	if err != nil {
		return errors.Wrap(err, "failed to query tables")
	}
	defer func() {
		if err := rows.Close(); err != nil {
			slog.Error("Error closing rows", slog.Any("error", err))
		}
	}()

	for rows.Next() {
		var table string
		if err := rows.Scan(&table); err != nil {
			return errors.Wrap(err, "failed to scan table name")
		}
		table_names = append(table_names, table)
	}

	// Wipe cache
	_, err = tx.Exec("DELETE FROM tables WHERE connection_id = ?;", con_id)
	if err != nil {
		return errors.Wrap(err, "failed to delete tables")
	}
	_, err = tx.Exec("DELETE FROM columns WHERE connection_id = ?;", con_id)
	if err != nil {
		return errors.Wrap(err, "failed to delete columns")
	}

	for _, table := range table_names {
		_, err := tx.Exec("INSERT INTO tables (connection_id, table_name) VALUES (?, ?) ;", con_id, table)
		if err != nil {
			return errors.Wrap(err, "failed to insert table")
		}

		var columns []Column
		// The built-in parameter substitution fails with `near "?": syntax error`
		rows, err := con.Query(
			fmt.Sprintf("PRAGMA table_info('%s');", table),
		)
		if err != nil {
			return errors.Wrap(err, "failed to query columns")
		}
		if err := rows.Bind(&columns); err != nil {
			return errors.Wrap(err, "failed to bind columns")
		}

		cols_query := "INSERT INTO columns(connection_id,table_name,name,type,\"notnull\",dflt_value,pk) VALUES "
		params := []interface{}{}
		for _, column := range columns {
			//  VALUES(?,?,?,?,?,?,?)
			cols_query += "(?,?,?,?,?,?,?),"
			params = append(params, con_id, table, column.Name, column.Type, column.Notnull, column.DfltValue, column.PK)
		}
		cols_query = cols_query[:len(cols_query)-1] + ";"
		res, err := tx.Exec(cols_query, params...)
		if err != nil {
			return errors.Wrap(err, "failed to insert columns")
		}
		cnt, err := res.RowsAffected()
		if err != nil {
			return errors.Wrap(err, "failed to get rows affected")
		}
		slog.Info("Inserted columns for table", slog.String("table", table), slog.Int("count", int(cnt)))
	}

	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, "failed to commit transaction")
	}
	return nil
}

func ReflectDB(db *Database, connections Connections, con Connection) error {
	con_db, ok := connections[con.ConnectionId]
	if !ok {
		return fmt.Errorf("Connection %d not found", con.ConnectionId)
	}

	switch con.DbType {
	case "sqlite":
		if err := reflectSqlite(db, con.ConnectionId, con_db); err != nil {
			return errors.Wrap(err, "failed to reflect sqlite database")
		}
	default:
		return fmt.Errorf("unknown database type %s", con.DbType)
	}
	return nil
}

func ReflectAll(db *Database, connections Connections) error {
	slog.Info("Reflecting all connections...")
	var connection_list []Connection
	rows, err := db.Query("SELECT * FROM connections WHERE deleted_at IS NULL")
	if err != nil {
		return errors.Wrap(err, "failed to query connections")
	}
	if err := rows.Bind(&connection_list); err != nil {
		return errors.Wrap(err, "failed to bind connections")
	}

	for _, con := range connection_list {
		err := ReflectDB(db, connections, con)
		if err != nil {
			fmt.Printf("Failed to reflect connection %d: %s\n", con.ConnectionId, err)
		}
	}

	return nil
}

func GetColumns(db *Database, con_id int64) ([]Column, error) {
	var cols []Column
	rows, err := db.Query("SELECT * FROM columns WHERE connection_id = ? ", con_id)
	if err != nil {
		return nil, err
	}
	if err := rows.Bind(&cols); err != nil {
		return nil, err
	}
	return cols, nil
}
