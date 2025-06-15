# Grognon

Grognon allows creating timeseries from existing databases.

## How it works

Grognon allows connecting to different databases (named connections), and defining recurrent SQL scripts to run against those connections (named crons).

The outputs are saved with the timestamp to a database table on Grognon's side, allowing us to have timeseries data without having to modify the schema of the target databases.

## Roadmap

This is not ordered and will evolve over time

- [ ] Use a PostgreSQL database for the internal storage
- [ ] Support other database types for Connections
- [ ] Allow editing Connections
- [ ] Migrate from incremental IDs to UUIDs
- [ ] Add events table for auditing purposes
