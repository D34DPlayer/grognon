package background

import (
	"context"
	"time"

	"d34d.one/grognon/internal/database"
)

func backgroundTask(ctx context.Context, duration time.Duration, task func()) {
	ticker := time.NewTicker(duration)
	go func() {
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

func SetupCronJobs(ctx context.Context, db *database.Database, cons database.Connections) error {
	backgroundTask(ctx, 30*time.Second, func() {
		database.ExecuteCrons(db, cons)
	})
	return nil
}

func SetupReflection(ctx context.Context, db *database.Database, cons database.Connections) error {
	backgroundTask(ctx, 30*time.Minute, func() {
		database.ReflectAll(db, cons)
	})
	return nil
}

func SetupDBRefresh(ctx context.Context, db *database.Database, cons database.Connections) error {
	backgroundTask(ctx, 5*time.Minute, func() {
		database.RefreshConnections(db, cons)
	})
	return nil
}
