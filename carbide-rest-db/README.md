## Forge Cloud DB

Forge Cloud DB is the data store interface for Cloud control plane

## DB Setup

Set up a Postgres database instance (version 12 or later recommended).

## Schema Setup

Once postgres is up and running, create the database/user:

    PGPASSWORD=postgres psql -h 127.0.0.1 -U postgres -p 30432 < scripts/setup.sql

Build the migrations command:

    make build

Export postgres env vars for migration command to read (could be command line arguments in future):

    export PGHOST=localhost PGPORT=30432 PGUSER=forge PGPASSWORD=forge

Initialize the migration tables:

    ./migrations db init

Apply Forge schema:

    ./migrations db migrate

In order to create a new migration file use the following to create a basic migration template

    PGUSER=forge PGPASSWORD=forge PGDATABASE=forge ./migrations db create_go {migration_description}
