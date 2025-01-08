package database

import "time"

type Connection struct {
	ID              int
	DbType          string
	ConnectionUrl   string
	CreatedAt       time.Time
	DeletedAt       *time.Time
	Connected       bool
	LastConnectedAt *time.Time
	LastError       *string
}
