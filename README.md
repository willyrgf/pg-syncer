# pg-sync

A query-based syncer between PostgreSQL databases.

## PostgreSQL version

- 9.2 or higher

## Sync modes
- fullsync: truncate the destination table and copy all result of the query to then
- onlydiff: sync only the diff data between the source and destination **(not implemented yet)**
- partialsync: copy all result from query without truncate the destination **(not implemented yet)**

*Check the syncer.toml.sample in the repo to more options and configs*

