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

type Table struct {
	TableId int
	Name    string
}

type Column struct {
	Name      string
	Type      string
	Notnull   bool
	DfltValue *string
	PK        int
}
