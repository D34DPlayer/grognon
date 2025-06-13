package database

import "time"

type ConnectionCreate struct {
	DbType        string
	ConnectionUrl string
}

type Connection struct {
	ConnectionId    int64
	DbType          string
	ConnectionUrl   string
	CreatedAt       EpochTime
	DeletedAt       *EpochTime
	Connected       bool
	LastConnectedAt *EpochTime
	LastError       *string
}

type Table struct {
	ConnectionId int64
	TableName    string
}

type Column struct {
	ConnectionId int64
	TableName    string
	Name         string
	Type         string
	Notnull      bool
	DfltValue    *string
	PK           int
}

type CronCreate struct {
	ConnectionId int64
	Name         string
	Command      string
	Schedule     string
}

type Cron struct {
	ConnectionId int64
	Name         string
	Command      string
	Schedule     string

	CronId    int64
	CreatedAt EpochTime
	DeletedAt *EpochTime
	LastRunAt *EpochTime
}

func (c *Cron) NeedsToRun() bool {
	if c.LastRunAt == nil {
		return true
	}
	switch c.Schedule {
	case "minute":
		return c.LastRunAt.Add(1 * time.Minute).Before(time.Now())
	case "hour":
		return c.LastRunAt.Add(1 * time.Hour).Before(time.Now())
	case "day":
		return c.LastRunAt.Add(24 * time.Hour).Before(time.Now())
	case "week":
		return c.LastRunAt.Add(7 * 24 * time.Hour).Before(time.Now())
	case "month":
		return c.LastRunAt.Add(30 * 24 * time.Hour).Before(time.Now())
	case "year":
		return c.LastRunAt.Add(365 * 24 * time.Hour).Before(time.Now())
	default:
		return false
	}
}

type CronOutput struct {
	CronId int64
	Name   string
	Type   string
}

type CronData map[string]interface{}
