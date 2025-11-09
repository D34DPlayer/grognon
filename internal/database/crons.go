package database

import (
	"context"
	"fmt"
	"log/slog"
	"time"

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
	defer func() {
		if err := rows.Close(); err != nil {
			slog.Error("Error closing rows", slog.Any("error", err))
		}
	}()

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
		output := CronOutput{cron.CronId, col, "NULL"}
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
		timestamp TIMESTAMPTZ NOT NULL`, cron.CronId)
	indexQuery := fmt.Sprintf(
		"CREATE INDEX cron_%d_timestamp ON cron_%d(timestamp);",
		cron.CronId,
		cron.CronId,
	)

	for _, output := range outputs {
		tableQuery += fmt.Sprintf(",%s %s", output.Name, output.Type)
	}
	tableQuery += ");"
	slog.Debug("Cron table creation", slog.String("tableQuery", tableQuery), slog.String("indexQuery", indexQuery))

	if _, err := db.Exec(tableQuery); err != nil {
		return nil, err
	}
	if _, err := db.Exec(indexQuery); err != nil {
		return nil, err
	}

	return outputs, nil
}

func AddCron(db *Database, cons Connections, input CronCreate) (*Cron, error) {
	con, ok := cons[input.ConnectionId]
	if !ok {
		return nil, fmt.Errorf("connection %d not found", input.ConnectionId)
	}

	// Create Cron in DB
	row := db.QueryRow(
		"INSERT INTO crons (connection_id, name, command, schedule) VALUES ($1, $2, $3, $4) RETURNING cron_id;",
		input.ConnectionId,
		input.Name,
		input.Command,
		input.Schedule,
	)
	var cronId int64
	err := row.Scan(&cronId)
	if err != nil {
		return nil, errors.Wrap(err, "Error getting cron ID")
	}

	cron, err := GetCron(db, cronId)
	if err != nil {
		return nil, errors.Wrap(err, "Error getting cron")
	}

	// Create Cron table
	outputs, err := createCronTable(db, con, *cron)
	if err != nil {
		_, _ = db.Exec("DELETE FROM crons WHERE cron_id = $1;", cron.CronId)
		return nil, errors.Wrap(err, "Error creating cron table")
	}

	// We can't create the TX before as it would lock table creation
	tx, err := db.BeginTx(context.TODO(), nil)
	if err != nil {
		return nil, errors.Wrap(err, "Error starting transaction")
	}
	// Save outputs
	insertQuery := `INSERT INTO cron_outputs (cron_id, name, type) VALUES `
	var params []interface{}

	for i, output := range outputs {
		offset := 3*i + 1
		insertQuery += fmt.Sprintf("($%d,$%d,$%d),", offset, offset+1, offset+2)
		params = append(params, cron.CronId, output.Name, output.Type)
	}
	insertQuery = insertQuery[:len(insertQuery)-1] + ";"

	if _, err := tx.Exec(insertQuery, params...); err != nil {
		_ = tx.Rollback()
		return nil, errors.Wrap(err, "Error inserting cron outputs")
	}

	return cron, tx.Commit()
}

func updateCronLastRun(db *Database, cronId int64) (*time.Time, error) {
	now := time.Now()
	// Update last run time of cron
	if _, err := db.Exec("UPDATE crons SET last_run_at = $1 WHERE cron_id = $2", now, cronId); err != nil {
		return nil, errors.Wrap(err, "Error updating cron last run")
	}
	return &now, nil
}

func ExecuteCrons(db *Database, cons Connections) error {
	slog.Debug("Executing crons")
	var crons []Cron
	rows, err := db.Query("SELECT * FROM crons WHERE deleted_at IS NULL")
	if err != nil {
		return errors.Wrap(err, "Error getting crons")
	}
	if err := rows.Bind(&crons); err != nil {
		return errors.Wrap(err, "Error binding crons")
	}
	slog.Debug("Found crons", slog.Int("count", len(crons)))

	for _, cron := range crons {
		if !cron.NeedsToRun() {
			continue
		}

		con, ok := cons[cron.ConnectionId]
		if !ok {
			slog.Error("connection not found", slog.Int64("id", cron.ConnectionId))
			continue
		}

		slog.Info("Executing cron", slog.Int64("id", cron.CronId))

		now, err := updateCronLastRun(db, cron.CronId)
		if err != nil {
			slog.Error("Error updating cron last run", slog.Int64("id", cron.CronId), slog.Any("error", err))
			continue
		}

		objects, cols, err := executeCron(con, cron)
		if err != nil {
			slog.Error("Error executing cron", slog.Int64("id", cron.CronId), slog.Any("error", err))
			continue
		}
		slog.Info("Saving results for cron", slog.Int64("id", cron.CronId))
		insertQuery := fmt.Sprintf("INSERT INTO cron_%d (timestamp", cron.CronId)
		for _, col := range cols {
			insertQuery += "," + col
		}
		insertQuery += ") VALUES "

		var params []interface{}
		argCounter := 1
		for i := range objects {
			insertQuery += fmt.Sprintf("($%d,", argCounter)
			argCounter++
			params = append(params, now)
			for _, col := range cols {
				insertQuery += fmt.Sprintf("$%d,", argCounter)
				argCounter++
				params = append(params, objects[i][col])
			}
			insertQuery = insertQuery[:len(insertQuery)-1] + "),"
		}
		insertQuery = insertQuery[:len(insertQuery)-1] + ";"

		slog.Debug("Inserting cron results", slog.String("query", insertQuery), slog.Int("params_count", len(params)), slog.Any("params", params))

		if _, err := db.Exec(insertQuery, params...); err != nil {
			slog.Error("Error inserting cron", slog.Int64("id", cron.CronId), slog.Any("error", err))
			continue
		}
		slog.Info("Cron executed", slog.Int64("id", cron.CronId))
	}

	return nil
}

func GetCrons(db *Database, connectionId *int64) ([]Cron, error) {
	var crons []Cron
	query := "SELECT * FROM crons WHERE deleted_at IS NULL"
	var args []interface{}

	if connectionId != nil {
		query += " AND connection_id = $1"
		args = append(args, *connectionId)
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, errors.Wrap(err, "Error getting crons")
	}
	if err := rows.Bind(&crons); err != nil {
		return nil, errors.Wrap(err, "Error binding crons")
	}
	return crons, nil
}

func GetCron(db *Database, cronId int64) (*Cron, error) {
	var cron Cron
	row := db.QueryRow("SELECT * FROM crons WHERE cron_id = $1 AND deleted_at IS NULL", cronId)
	if err := row.Bind(&cron); err != nil {
		return nil, errors.Wrap(err, "Error binding cron")
	}
	if cron.CronId == 0 {
		return nil, fmt.Errorf("cron %d not found", cronId)
	}
	return &cron, nil
}

func GetCronOutputs(db *Database, cronId int64) ([]CronOutput, error) {
	var outputs []CronOutput
	rows, err := db.Query("SELECT * FROM cron_outputs WHERE cron_id = $1", cronId)
	if err != nil {
		return nil, errors.Wrap(err, "Error getting cron outputs")
	}
	if err := rows.Bind(&outputs); err != nil {
		return nil, errors.Wrap(err, "Error binding cron outputs")
	}
	return outputs, nil
}

func DeleteCron(db *Database, cronId int64) error {
	// Mark cron as deleted
	if _, err := db.Exec("UPDATE crons SET deleted_at = $1 WHERE cron_id = $2", time.Now(), cronId); err != nil {
		return errors.Wrap(err, "Error deleting cron")
	}

	return nil
}

func UpdateCron(db *Database, con *sqle.DB, cron Cron) error {
	// Get current state
	tx, err := db.BeginTx(context.TODO(), nil)
	if err != nil {
		return errors.Wrap(err, "Error starting transaction")
	}
	oldOutputs, err := GetCronOutputs(db, cron.CronId)
	if err != nil {
		_ = tx.Rollback()
		return errors.Wrap(err, "Error getting cron outputs")
	}

	// Update cron in DB
	if _, err := tx.Exec(
		"UPDATE crons SET name = $1, command = $2, schedule = $3 WHERE cron_id = $4",
		cron.Name,
		cron.Command,
		cron.Schedule,
		cron.CronId,
	); err != nil {
		return errors.Wrap(err, "Error updating cron")
	}

	newOutputs, err := reflectCron(con, cron)
	if err != nil {
		_ = tx.Rollback()
		return errors.Wrap(err, "Error reflecting cron")
	}

	if len(newOutputs) != len(oldOutputs) {
		_ = tx.Rollback()
		return fmt.Errorf("cannot update cron outputs, number of outputs changed: %d -> %d", len(oldOutputs), len(newOutputs))
	}

	for i := range oldOutputs {
		if oldOutputs[i].Name != newOutputs[i].Name || oldOutputs[i].Type != newOutputs[i].Type {
			_ = tx.Rollback()
			return fmt.Errorf("cannot update cron outputs, output %d changed: %s %s -> %s %s",
				i, oldOutputs[i].Name, oldOutputs[i].Type, newOutputs[i].Name, newOutputs[i].Type)
		}
	}

	return tx.Commit()
}

func GetCronData(db *Database, cronId int64) ([]CronData, error) {
	var data []CronData
	query := fmt.Sprintf("SELECT * FROM cron_%d ORDER BY timestamp DESC", cronId)
	rows, err := db.Query(query)
	if err != nil {
		return nil, errors.Wrap(err, "Error getting cron data")
	}
	if err := rows.Bind(&data); err != nil {
		return nil, errors.Wrap(err, "Error binding cron data")
	}
	return data, nil
}
