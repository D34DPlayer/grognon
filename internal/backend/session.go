package backend

import (
	"context"
	"encoding/gob"
	"log/slog"
	"net/http"

	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	inertia "github.com/romsar/gonertia"
)

type SessionFlashProvider struct {
}

type SessionKey struct{}

func SessionMiddleware(store *sessions.CookieStore, next http.Handler) http.Handler {
	gob.Register(&inertia.ValidationErrors{})

	fn := func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "session")
		if session.IsNew {
			session.Values["sessionID"] = securecookie.GenerateRandomKey(16)
		}

		ctx := inertia.SetProps(r.Context(), inertia.Props{
			"errors": session.Values["errors"],
			"flash":  session.Flashes("flash"),
		})
		delete(session.Values, "errors")
		err := session.Save(r, w)
		if err != nil {
			slog.Error("Failed to save session", slog.Any("error", err))
		}

		ctx = context.WithValue(ctx, SessionKey{}, session)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}

func GetSession(r *http.Request) *sessions.Session {
	session := r.Context().Value(SessionKey{}).(*sessions.Session)
	if session == nil {
		slog.Error("Session not found")
	}

	return session
}

func SaveSession(w http.ResponseWriter, r *http.Request) {
	session := GetSession(r)
	err := session.Save(r, w)
	if err != nil {
		slog.Error("Failed to save session", slog.Any("error", err))
	}
}
