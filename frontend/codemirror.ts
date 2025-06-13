import type { Completion } from '@codemirror/autocomplete'
import type { Extension } from '@codemirror/state'
import type { Column } from './types'
import { sql, SQLite } from '@codemirror/lang-sql'

export function getExtensions(columns?: Column[]): Extension[] {
  const schema: Record<string, Completion[]> = {}

  if (columns) {
    for (const column of columns) {
      const table = schema[column.TableName] || []
      table.push({
        label: column.Name,
        type: 'variable',
        detail: column.Type,
      })
      schema[column.TableName] = table
    }
  }

  const sqlExtension = sql({
    schema,
    upperCaseKeywords: true,
    dialect: SQLite,
  })
  return [sqlExtension]
}
