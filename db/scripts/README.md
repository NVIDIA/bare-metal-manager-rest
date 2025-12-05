## Executing Setup Script

The setup script creates forge database and role.

**Note**: Please update the password for `forge` DB/role in the script before executing.

It also needs to install default extensions for all databases. Hence the script must be run as postgres user and connect to template1 DB:

    PGPASSWORD=postgres psql -U postgres -p 30432 -d template1 < scripts/setup.sql
