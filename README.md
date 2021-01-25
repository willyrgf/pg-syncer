# pg-sync

A query-based syncer between PostgreSQL databases.

## PostgreSQL version
- 9.2 or higher

## Sync modes
- fullsync: truncate the destination table and copy all result of the query to then (in transaction mode)
- onlydiff: sync only the diff data between the source and destination **(not implemented yet)**
- partialsync: copy all result from query without truncate the destination **(not implemented yet)**

*Check the syncer.toml.sample in the repo to more options and configs*

## Configuration examples

### System configuration (config.toml)
```
[config]
dev = true # enable debug
trace = false # enable the trace of each execution and errors

# repository of the pg-syncer, only to internal things
# not used yet
[repository]
url = "postgres://cryp:XkZPqxHC5h5f6koZrzap@127.0.0.1:5432/cryp?pool_max_conns=2"
```

### Sync configuration (syncer.toml)
```
# repositories to syncers services sync the data.
# we could have a lot of repositories, you'll choose on the
# syncers configuration below what you'll use

[repository_1]
url = "postgres://cryp:XkZPqxHC5h5f6koZrzap@210.253.255.252:5432/cryp?pool_max_conns=2"

[repository_2]
url = "postgres://report:zVmwGt6Pj67nXkhKNceupBjT@210.252.255.255:5432/report?pool_max_conns=2"


# the syncers configuration area.
# each syncer configuration below will start a service to
# handle with that sync
[syncers]
    # full documented configuration
    [[syncers.access]]
    enabled = true # if enabled or not
    sync_mode = "fullsync" # fullsync, onlydiff, partialsync ...
    source_repository = "repository_1" # the label of the source repository
    source_db = "cryp" # database name of the source repository
    source_query = "select   now() as updated_at,   buy_op_at,   filled::numeric,   now() as created_at,   bot_license_id as bot_id,   diff_from_now::numeric,   diff_percent_from_now::numeric,   estimate_buy_price::numeric,   current_bid_price::numeric,   buy_price::numeric,   total_if_sell_now::numeric,   buyed_amount::numeric,   estimate::numeric,   market,   'binance' as exchange,   strategy_label from   stats.open_positions" # the query will be executed on the source to get the data
    destination_repository = "repository_2" # the label of the destination repository
    destination_db = "report" # database name of the destination repository
    destination_schema = "public" # the schema name of the destination of the data
    destination_table = "open_positions" # the destination table that we'll write the data getted on source
    periodicity_value = "1" # periodicity value that this syncer will run each periodicity unit (int value > 0)
    periodicity_unit = "minute" # periodicity unit that this syncer will run (second, minute, hour, day, week) to understand more: https://github.com/cryp-com-br/pg-syncer/blob/e92559d0881d4bc2d380345d5a8f3be45dd07808/syncer/scheduler.go#L11

    # another example
    [[syncers.access]]
    enabled = false
    sync_mode = "fullsync" # fullsync, onlydiff, partialsync ...
    source_repository = "repository_1"
    source_db = "cryp"
    source_query = "select id::uuid as id,  created_at,  updated_at, (holder_info ->> 'label')::text as label from bot_licenses"
    destination_repository = "repository_2"
    destination_db = "report"
    destination_schema = "public"
    destination_table = "bots"
    periodicity_value = "2"
    periodicity_unit = "weeks"


```
