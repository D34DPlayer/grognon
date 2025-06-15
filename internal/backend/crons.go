package backend

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"d34d.one/grognon/internal/database"
	"github.com/gorilla/mux"
	inertia "github.com/romsar/gonertia"
)

func GetCrons(i *inertia.Inertia, db *database.Database) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		errs := NewErrors(r)
		vars := mux.Vars(r)

		props := inertia.Props{
			"connectionId": nil,
			"connection":   nil,
			"crons":        nil,
		}

		connectionIdStr, ok := vars["connection_id"]
		var connectionId *int64
		if ok {
			id, err := strconv.ParseInt(connectionIdStr, 10, 64)
			if err != nil {
				slog.Error("Failed to parse connection id", slog.Any("error", err))
				errs.Add("input", err)

				Render(w, errs.Request(r), i, "Home/Crons", nil)
				return
			}
			props["connectionId"] = id
			connectionId = &id
		}

		if connectionId != nil {
			connection, err := database.GetConnection(db, *connectionId)
			if err != nil {
				slog.Error("Failed to get connection", slog.Any("error", err))
				errs.Add("connection", err)

				Render(w, errs.Request(r), i, "Home/Crons", props)
				return
			}
			props["connection"] = connection
		}

		crons, err := database.GetCrons(db, connectionId)
		if err != nil {
			slog.Error("Failed to get crons", slog.Any("error", err))
			errs.Add("crons", err)
		}
		props["crons"] = crons

		Render(w, errs.Request(r), i, "Home/Crons", props)
	}

	return i.Middleware(http.HandlerFunc(fn))
}

func GetNewCrons(i *inertia.Inertia, db *database.Database) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		errs := NewErrors(r)
		vars := mux.Vars(r)

		props := inertia.Props{
			"connectionId": nil,
			"connections":  nil,
			"columns":      nil,
		}

		connectionIdStr, ok := vars["connection_id"]
		var connectionId *int64
		if ok {
			id, err := strconv.ParseInt(connectionIdStr, 10, 64)
			if err != nil {
				slog.Error("Failed to parse connection id", slog.Any("error", err))
				errs.Add("input", err)

				Render(w, errs.Request(r), i, "Home/CronsCreate", nil)
				return
			}
			props["connectionId"] = id
			connectionId = &id
		}

		if connectionId != nil {
			columns, err := database.GetColumns(db, *connectionId)
			if err != nil {
				slog.Error("Failed to get columns", slog.Any("error", err))
				errs.Add("columns", err)
			}
			props["columns"] = columns
		}

		connections, err := database.GetConnections(db)
		if err != nil {
			slog.Error("Failed to get connections", slog.Any("error", err))
			errs.Add("connections", err)
			Render(w, errs.Request(r), i, "Home/CronsCreate", nil)
			return
		}
		props["connections"] = connections

		Render(w, errs.Request(r), i, "Home/CronsCreate", props)
	}

	return i.Middleware(http.HandlerFunc(fn))
}

func PostNewCrons(i *inertia.Inertia, db *database.Database, cons database.Connections) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		errs := NewErrors(r)

		var body database.CronCreate
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			slog.Error("Failed to decode request body", slog.Any("error", err))
			errs.Add("body", err)
		}

		cron, err := database.AddCron(db, cons, body)
		if err != nil {
			slog.Error("Failed to create cron", slog.Any("error", err))
			errs.Add("creation", errors.Unwrap(err))
		}

		if errs.HasErrors() {
			errs.Save(w, r)
			i.Back(w, r)
		} else {
			url := fmt.Sprintf("/crons/%d", cron.CronId)
			i.Redirect(w, r, url)
		}
		SaveSession(w, r)
	}

	return i.Middleware(http.HandlerFunc(fn))
}

func GetCron(i *inertia.Inertia, db *database.Database) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		errs := NewErrors(r)
		vars := mux.Vars(r)

		cronId, err := strconv.ParseInt(vars["cron_id"], 10, 64)
		if err != nil {
			slog.Error("Failed to parse cron id", slog.Any("error", err))
			errs.Add("input", err)

			Render(w, errs.Request(r), i, "Home/Cron", nil)
			return
		}

		cron, err := database.GetCron(db, cronId)
		if err != nil {
			slog.Error("Failed to get cron", slog.Any("error", err))
			errs.Add("cron", err)

			Render(w, errs.Request(r), i, "Home/Cron", nil)
			return
		}

		connection, err := database.GetConnection(db, cron.ConnectionId)
		if err != nil {
			slog.Error("Failed to get connection for cron", slog.Any("error", err))
			errs.Add("connection", err)
			Render(w, errs.Request(r), i, "Home/Cron", nil)
			return
		}

		outputs, err := database.GetCronOutputs(db, cronId)
		if err != nil {
			slog.Error("Failed to get cron outputs", slog.Any("error", err))
			errs.Add("outputs", err)
			Render(w, errs.Request(r), i, "Home/Cron", nil)
			return
		}

		props := inertia.Props{
			"cron":        cron,
			"connection":  connection,
			"cronOutputs": outputs,
		}

		Render(w, errs.Request(r), i, "Home/Cron", props)
	}

	return i.Middleware(http.HandlerFunc(fn))
}

func GetCronData(i *inertia.Inertia, db *database.Database) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		errs := NewErrors(r)
		vars := mux.Vars(r)

		props := inertia.Props{
			"cronId":      nil,
			"cronOutputs": nil,
			"cron":        nil,
			"data":        nil,
		}

		cronId, err := strconv.ParseInt(vars["cron_id"], 10, 64)
		if err != nil {
			slog.Error("Failed to parse cron id", slog.Any("error", err))
			errs.Add("input", err)

			Render(w, errs.Request(r), i, "Home/CronData", props)
			return
		}
		props["cronId"] = cronId

		cron, err := database.GetCron(db, cronId)
		if err != nil {
			slog.Error("Failed to get cron", slog.Any("error", err))
			errs.Add("cron", err)
			Render(w, errs.Request(r), i, "Home/CronData", props)
			return
		}
		props["cron"] = cron

		outputs, err := database.GetCronOutputs(db, cronId)
		if err != nil {
			slog.Error("Failed to get cron outputs", slog.Any("error", err))
			errs.Add("outputs", err)
			Render(w, errs.Request(r), i, "Home/CronData", props)
			return
		}
		props["cronOutputs"] = outputs

		data, err := database.GetCronData(db, cronId)
		if err != nil {
			slog.Error("Failed to get cron data", slog.Any("error", err))
			errs.Add("data", err)

			Render(w, errs.Request(r), i, "Home/CronData", props)
			return
		}
		props["data"] = data

		Render(w, errs.Request(r), i, "Home/CronData", props)
	}

	return i.Middleware(http.HandlerFunc(fn))
}

func DeleteCrons(i *inertia.Inertia, db *database.Database) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		errs := NewErrors(r)
		vars := mux.Vars(r)

		cronId, err := strconv.ParseInt(vars["cron_id"], 10, 64)
		if err != nil {
			slog.Error("Failed to parse cron id", slog.Any("error", err))
			errs.Add("input", err)

			Render(w, errs.Request(r), i, "Home/Cron", nil)
			return
		}

		err = database.DeleteCron(db, cronId)
		if err != nil {
			slog.Error("Failed to delete cron", slog.Any("error", err))
			errs.Add("deletion", err)
		}

		if errs.HasErrors() {
			errs.Save(w, r)
			i.Back(w, r)
		} else {
			i.Redirect(w, r, "/crons")
		}
		SaveSession(w, r)
	}

	return i.Middleware(http.HandlerFunc(fn))
}
