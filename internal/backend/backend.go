package backend

import (
	"log/slog"
	"net/http"
	"os"

	"d34d.one/grognon/internal/database"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	inertia "github.com/romsar/gonertia"
)

func Setup(db *database.Database, cons database.Connections, ssrHost string) error {
	sessionKey := os.Getenv("SESSION_KEY")
	if sessionKey == "" {
		slog.Warn("SESSION_KEY not set")
		sessionKey = "grognon"
	}
	store := sessions.NewCookieStore([]byte(sessionKey))
	i, err := initInertia(ssrHost)
	if err != nil {
		return err
	}

	router := mux.NewRouter()

	router.PathPrefix("/build").Handler(http.StripPrefix("/build/", http.FileServer(http.Dir("./public/build"))))

	router.Methods("GET").Path("/connections/create").
		Handler(GetNewConnections(i))
	router.Methods("GET").Path("/connections/{connection_id}").
		Handler(GetConnection(i, db))
	router.Methods("DELETE").Path("/connections/{connection_id}").
		Handler(DeleteConnection(i, db, cons))
	router.Methods("GET").Path("/connections/{connection_id}/crons").
		Handler(GetCrons(i, db))
	router.Methods("GET").Path("/connections/{connection_id}/crons/create").
		Handler(GetNewCrons(i, db))
	router.Methods("GET").Path("/connections").
		Handler(GetConnections(i, db))
	router.Methods("POST").Path("/connections").
		Handler(PostNewConnections(i, db, cons))
	router.Methods("GET").Path("/crons/create").
		Handler(GetNewCrons(i, db))
	router.Methods("GET").Path("/crons/{cron_id}/data").
		Handler(GetCronData(i, db))
	router.Methods("GET").Path("/crons/{cron_id}").
		Handler(GetCron(i, db))
	router.Methods("DELETE").Path("/crons/{cron_id}").
		Handler(DeleteCrons(i, db))
	router.Methods("GET").Path("/crons").
		Handler(GetCrons(i, db))
	router.Methods("POST").Path("/crons").
		Handler(PostNewCrons(i, db, cons))
	router.Methods("GET").Path("/").
		Handler(http.RedirectHandler("/connections", http.StatusTemporaryRedirect))
	router.PathPrefix("/").
		Handler(NotFound(i))

	return http.ListenAndServe(":3000", SessionMiddleware(store, router))
}

func GetHome(i *inertia.Inertia) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		Render(w, r, i, "Home", nil)
	}

	return i.Middleware(http.HandlerFunc(fn))
}

func NotFound(i *inertia.Inertia) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		props := inertia.Props{
			"url": r.URL.Path,
		}

		Render(w, r, i, "NotFound", props)
	}

	return i.Middleware(http.HandlerFunc(fn))
}

type Errors struct {
	session *sessions.Session
	errs    inertia.ValidationErrors
}

func NewErrors(r *http.Request) *Errors {
	session := GetSession(r)
	return &Errors{session: session, errs: make(inertia.ValidationErrors)}
}

func (e *Errors) Count() int {
	return len(e.errs)
}

func (e *Errors) HasErrors() bool {
	return e.Count() > 0
}

func (e *Errors) Add(key string, err error) {
	e.errs[key] = err.Error()
}

func (e *Errors) Save(w http.ResponseWriter, r *http.Request) {
	e.session.Values["errors"] = e.errs
	err := e.session.Save(r, w)
	if err != nil {
		slog.Error("Failed to save session with errors", "error", err)
		return
	}
}

func (e *Errors) Request(r *http.Request) *http.Request {
	ctx := r.Context()
	ctx = inertia.SetProps(ctx, inertia.Props{
		"errors": e.errs,
	})
	return r.WithContext(ctx)
}
