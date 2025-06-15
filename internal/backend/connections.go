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

func GetNewConnections(i *inertia.Inertia) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		Render(w, r, i, "Home/ConnectionsCreate", nil)
	}

	return i.Middleware(http.HandlerFunc(fn))
}

func PostNewConnections(i *inertia.Inertia, db *database.Database, cons database.Connections) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		errs := NewErrors(r)

		var body database.ConnectionCreate
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			slog.Error("Failed to decode request body", slog.Any("error", err))
			errs.Add("body", err)
		}

		con, err := database.AddConnection(db, cons, body)
		if err != nil {
			slog.Error("Failed to create connection", slog.Any("error", err))
			errs.Add("creation", errors.Unwrap(err))
		}

		if errs.HasErrors() {
			errs.Save(w, r)
			i.Back(w, r)
		} else {
			url := fmt.Sprintf("/connections/%d", con.ConnectionId)
			i.Redirect(w, r, url)
		}
		SaveSession(w, r)
	}

	return i.Middleware(http.HandlerFunc(fn))
}

func GetConnections(i *inertia.Inertia, db *database.Database) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		errs := NewErrors(r)

		connections, err := database.GetConnections(db)
		if err != nil {
			slog.Error("Failed to get connections", slog.Any("error", err))
			errs.Add("connections", err)
		}

		props := inertia.Props{
			"connections": connections,
		}

		Render(w, errs.Request(r), i, "Home/Connections", props)
	}

	return i.Middleware(http.HandlerFunc(fn))
}

func GetConnection(i *inertia.Inertia, db *database.Database) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		errs := NewErrors(r)
		vars := mux.Vars(r)

		connectionId, err := strconv.ParseInt(vars["connection_id"], 10, 64)
		if err != nil {
			slog.Error("Failed to parse connection id", slog.Any("error", err))
			errs.Add("input", err)

			Render(w, errs.Request(r), i, "Home/Connection", nil)
			return
		}
		props := inertia.Props{
			"connectionId": connectionId,
		}

		connection, err := database.GetConnection(db, connectionId)
		if err != nil {
			slog.Error("Failed to get connection", slog.Any("error", err))
			errs.Add("connection", err)

			Render(w, errs.Request(r), i, "Home/Connection", props)
			return
		}
		props["connection"] = connection

		cols, err := database.GetColumns(db, connectionId)
		if err != nil {
			slog.Error("Failed to get columns", slog.Any("error", err))
			errs.Add("columns", err)

			Render(w, errs.Request(r), i, "Home/Connection", props)
			return
		}
		props["columns"] = cols

		Render(w, errs.Request(r), i, "Home/Connection", props)
	}

	return i.Middleware(http.HandlerFunc(fn))
}

func DeleteConnection(i *inertia.Inertia, db *database.Database, cons database.Connections) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		errs := NewErrors(r)
		vars := mux.Vars(r)

		connectionId, err := strconv.ParseInt(vars["connection_id"], 10, 64)
		if err != nil {
			slog.Error("Failed to parse connection id", slog.Any("error", err))
			errs.Add("input", err)

			Render(w, errs.Request(r), i, "Home/Connection", nil)
			return
		}

		err = database.DeleteConnection(db, cons, connectionId)
		if err != nil {
			slog.Error("Failed to delete connection", slog.Any("error", err))
			errs.Add("deletion", err)
		}

		if errs.HasErrors() {
			errs.Save(w, r)
			i.Back(w, r)
		} else {
			i.Redirect(w, r, "/connections")
		}
	}

	return i.Middleware(http.HandlerFunc(fn))
}
