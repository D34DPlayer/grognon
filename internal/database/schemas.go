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

type CronCreate struct {
	ConnectionId int
	Name         string
	Command      string
	Schedule     string
}

type Cron struct {
	*CronCreate
	CronId    int
	CreatedAt time.Time
	DeletedAt *time.Time
	LastRunAt *time.Time
}

type CronOutput struct {
	Name string
	Type string
}
