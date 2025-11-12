package background

import (
	"context"
	"log/slog"
	"time"

	"d34d.one/grognon/internal/database"
)

func backgroundTask(ctx context.Context, duration time.Duration, eager bool, task func()) {
	ticker := time.NewTicker(duration)
	go func() {
		if eager {
			task()
		}
		for {
			select {
			case <-ticker.C:
				task()
			case <-ctx.Done():
				ticker.Stop()
				return
			}
		}
	}()
}

func SetupCronJobs(ctx context.Context, db *database.Database, cons database.Connections) {
	backgroundTask(ctx, 30*time.Second, true, func() {
		err := database.ExecuteCrons(db, cons)
		if err != nil {
			slog.Error("Failed to execute crons", "error", err)
		}
	})
}

func SetupReflection(ctx context.Context, db *database.Database, cons database.Connections) {
	backgroundTask(ctx, 30*time.Minute, true, func() {
		err := database.ReflectAll(db, cons)
		if err != nil {
			slog.Error("Failed to reflect database", "error", err)
		}
	})
}

func SetupDBRefresh(ctx context.Context, db *database.Database, cons database.Connections) {
	backgroundTask(ctx, 5*time.Minute, false, func() {
		err := database.RefreshConnections(db, cons)
		if err != nil {
			slog.Error("Failed to refresh database connections", "error", err)
		}
	})
}
