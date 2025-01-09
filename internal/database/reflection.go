package database

import (
	"context"
	"fmt"
	"log"

	"github.com/yaitoo/sqle"
)

func reflectSqlite(db *Database, con_id int, con *sqle.DB) error {
	log.Printf("Reflecting connection %d", con_id)

	tx, err := db.BeginTx(context.TODO(), nil)
	if err != nil {
		return err
	}

	var table_names []string
	rows, err := con.Query("SELECT name FROM sqlite_schema WHERE type ='table' AND name NOT LIKE 'sqlite_%'")
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var table string
		if err := rows.Scan(&table); err != nil {
			return err
		}
		table_names = append(table_names, table)
	}

	// Wipe cache
	_, err = tx.Exec("DELETE FROM tables WHERE connection_id = ?;", con_id)
	if err != nil {
		return err
	}
	_, err = tx.Exec("DELETE FROM columns WHERE connection_id = ?;", con_id)
	if err != nil {
		return err
	}

	for _, table := range table_names {
		_, err := tx.Exec("INSERT INTO tables (connection_id, table_name) VALUES (?, ?) ;", con_id, table)
		if err != nil {
			return err
		}

		var columns []Column
		// The built-in parameter substitution fails with `near "?": syntax error`
		rows, err := con.Query(
			fmt.Sprintf("PRAGMA table_info('%s');", table),
		)
		if err != nil {
			return err
		}
		if err := rows.Bind(&columns); err != nil {
			return err
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
			return err
		}
		cnt, err := res.RowsAffected()
		if err != nil {
			return err
		}
		log.Println("Inserted", cnt, "columns for table", table)
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func ReflectDB(db *Database, connections Connections, con Connection) error {
	con_db, ok := connections[con.ID]
	if !ok {
		return fmt.Errorf("Connection %d not found", con.ID)
	}

	switch con.DbType {
	case "sqlite":
		if err := reflectSqlite(db, con.ID, con_db); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown database type %s", con.DbType)
	}
	return nil
}

func ReflectAll(db *Database, connections Connections) error {
	log.Println("Reflecting all connections...")
	var connection_list []Connection
	rows, err := db.Query("SELECT * FROM connections WHERE deleted_at IS NULL")
	if err != nil {
		return err
	}
	if err := rows.Bind(&connection_list); err != nil {
		return err
	}

	for _, con := range connection_list {
		err := ReflectDB(db, connections, con)
		if err != nil {
			fmt.Printf("Failed to reflect connection %d: %s\n", con.ID, err)
		}
	}

	return nil
}
