## Carbide REST API DB

Carbide REST API DB is the data store interface for Cloud control plane

## DB Setup

Set up a Postgres database instance (version 12 or later recommended).

## Schema Setup

Once postgres is up and running, create the database/user:

    PGPASSWORD=postgres psql -h 127.0.0.1 -U postgres -p 30432 < scripts/setup.sql
