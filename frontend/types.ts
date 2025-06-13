export const DB_TYPES = [
  'sqlite',
] as const
export type DbType = typeof DB_TYPES[number]

export type ConnectionCreate = {
  ConnectionUrl: string
  DbType: DbType
}

export type Connection = {
  ConnectionId: number
  CreatedAt: string
  DeletedAt: string | null
  Connected: boolean
  LastConnectedAt: string | null
  LastError: string | null
} & ConnectionCreate

export type Column = {
  ConnectionId: number
  TableName: string
  Name: string
  Type: string
  Notnull: boolean
  DfltValue: string | null
  PK: number
}

export const SCHEDULES = [
  'minute',
  'hour',
  'day',
  'week',
  'month',
  'year',
] as const
type Schedule = typeof SCHEDULES[number]

export type CronCreate = {
  ConnectionId: number
  Name: string
  Command: string
  Schedule: Schedule
}

export type Cron = {
  ConnectionId: number
  Name: string
  Command: string
  Schedule: string

  CronId: number
  CreatedAt: string
  DeletedAt: string | null
  LastRunAt: string | null
}

export type CronOutput = {
  CronId: number
  Name: string
  Type: string
}
