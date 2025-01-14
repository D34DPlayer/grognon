package database

type Connection struct {
	ID              int
	DbType          string
	ConnectionUrl   string
	CreatedAt       EpochTime
	DeletedAt       *EpochTime
	Connected       bool
	LastConnectedAt *EpochTime
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

type Cron struct {
	ConnectionId int
	Name         string
	Command      string
	Schedule     string

	CronId    int
	CreatedAt EpochTime
	DeletedAt *EpochTime
	LastRunAt *EpochTime
}

type CronOutput struct {
	Name string
	Type string
}
