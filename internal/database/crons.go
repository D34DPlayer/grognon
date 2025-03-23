package database

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/pkg/errors"
	"github.com/yaitoo/sqle"
)

type Object map[string]interface{}

func executeCron(con *sqle.DB, cron Cron) ([]Object, []string, error) {
	var output []Object

	rows, err := con.Query(cron.Command)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return nil, nil, err
	}

	scanSlots := make([]interface{}, len(cols))
	for rows.Next() {
		object := make(Object, len(cols))
		// Fill scanSlots with pointers to the values in object
		for i := range cols {
			scanSlots[i] = new(interface{})
		}

		if err := rows.Scan(scanSlots...); err != nil {
			return nil, nil, err
		}

		for i, col := range cols {
			object[col] = *scanSlots[i].(*interface{})
		}
		output = append(output, object)
	}

	return output, cols, nil
}

func reflectCron(con *sqle.DB, cron Cron) ([]CronOutput, error) {
	var outputs []CronOutput

	objects, cols, err := executeCron(con, cron)
	if err != nil {
		return nil, err
	}
	if len(objects) == 0 {
		return nil, fmt.Errorf("no rows returned")
	}

	for _, col := range cols {
		output := CronOutput{col, "NULL"}
		for _, object := range objects {
			colValue := object[col]
			if colValue == nil {
				return nil, fmt.Errorf("column %s is null", col)
			}
			var valueType string
			switch colValue.(type) {
			case string:
				valueType = "TEXT"
			case int:
			case int64:
				valueType = "INTEGER"
			case float64:
				valueType = "REAL"
			default:
				return nil, fmt.Errorf("unknown type %T", colValue)
			}
			if output.Type == "NULL" {
				output.Type = valueType
			} else if output.Type != valueType {
				return nil, fmt.Errorf("column %s has mixed types", col)
			}
		}
		outputs = append(outputs, output)
	}

	return outputs, nil
}

func createCronTable(db *Database, con *sqle.DB, cron Cron) ([]CronOutput, error) {
	outputs, err := reflectCron(con, cron)
	if err != nil {
		return nil, err
	}

	tableQuery := fmt.Sprintf(`CREATE TABLE cron_%d (
		timestamp INTEGER NOT NULL`, cron.CronId)
	indexQuery := fmt.Sprintf(
		"CREATE INDEX cron_%d_timestamp ON cron_%d(timestamp);",
		cron.CronId,
		cron.CronId,
	)

	for _, output := range outputs {
		tableQuery += fmt.Sprintf(",%s %s", output.Name, output.Type)
	}
	tableQuery += ");"

	if _, err := db.Exec(tableQuery); err != nil {
		return nil, err
	}
	if _, err := db.Exec(indexQuery); err != nil {
		return nil, err
	}

	return outputs, nil
}

func AddCron(db *Database, cons Connections, cron Cron) error {
	con, ok := cons[cron.ConnectionId]
	if !ok {
		return fmt.Errorf("connection %d not found", cron.ConnectionId)
	}

	// Create Cron in DB
	res, err := db.Exec(
		"INSERT INTO crons (connection_id, name, command, schedule) VALUES (?, ?, ?, ?);",
		cron.ConnectionId,
		cron.Name,
		cron.Command,
		cron.Schedule,
	)
	if err != nil {
		return errors.Wrap(err, "Error inserting cron")
	}
	cronId, err := res.LastInsertId()
	if err != nil {
		return errors.Wrap(err, "Error getting cron ID")
	}
	cron.CronId = int(cronId)

	// Create Cron table
	outputs, err := createCronTable(db, con, cron)
	if err != nil {
		db.Exec("DELETE FROM crons WHERE cron_id = ?", cron.CronId)
		return errors.Wrap(err, "Error creating cron table")
	}

	// We can't create the TX before as it would lock table creation
	tx, err := db.BeginTx(context.TODO(), nil)
	if err != nil {
		return errors.Wrap(err, "Error starting transaction")
	}
	// Save outputs
	insertQuery := `INSERT INTO cron_outputs (cron_id, name, type) VALUES `
	var params []interface{}

	for _, output := range outputs {
		insertQuery += "(?,?,?),"
		params = append(params, cron.CronId, output.Name, output.Type)
	}
	insertQuery = insertQuery[:len(insertQuery)-1] + ";"

	if _, err := tx.Exec(insertQuery, params...); err != nil {
		tx.Rollback()
		return errors.Wrap(err, "Error inserting cron outputs")
	}

	return tx.Commit()
}

func ExecuteCrons(db *Database, cons Connections, schedule string) error {
	slog.Info("Executing crons for schedule", slog.String("schedule", schedule))
	var crons []Cron
	rows, err := db.Query("SELECT * FROM crons WHERE schedule = ? AND deleted_at IS NULL", schedule)
	if err != nil {
		return errors.Wrap(err, "Error getting crons")
	}
	if err := rows.Bind(&crons); err != nil {
		return errors.Wrap(err, "Error binding crons")
	}
	slog.Info("Found crons", slog.Int("count", len(crons)))

	for _, cron := range crons {
		con, ok := cons[cron.ConnectionId]
		if !ok {
			slog.Error("connection not found", slog.Int("id", cron.ConnectionId))
			continue
		}

		slog.Info("Executing cron", slog.Int("id", cron.CronId))
		now := Now()
		objects, cols, err := executeCron(con, cron)
		if err != nil {
			slog.Error("Error executing cron", slog.Int("id", cron.CronId), slog.Any("error", err))
			continue
		}
		slog.Info("Saving results for cron", slog.Int("id", cron.CronId))
		insertQuery := fmt.Sprintf("INSERT INTO cron_%d (timestamp", cron.CronId)
		for _, col := range cols {
			insertQuery += "," + col
		}
		insertQuery += ") VALUES "

		var params []interface{}
		for i := range objects {
			insertQuery += "(?,"
			params = append(params, now)
			for _, col := range cols {
				insertQuery += "?,"
				params = append(params, objects[i][col])
			}
			insertQuery = insertQuery[:len(insertQuery)-1] + "),"
		}
		insertQuery = insertQuery[:len(insertQuery)-1] + ";"

		if _, err := db.Exec(insertQuery, params...); err != nil {
			slog.Error("Error inserting cron", slog.Int("id", cron.CronId), slog.Any("error", err))
			continue
		}
		slog.Info("Cron executed", slog.Int("id", cron.CronId))
	}

	return nil
}
