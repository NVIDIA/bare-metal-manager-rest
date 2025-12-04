# Intro

Elektra Site Agent is the site manager who manages Carbide Data center sites and provisions DPU by interfacing with Site Controller

Elektra Site Agent Components
* Site Agent
* Carbide gRPC Library
* Site Agent Data Store

Elektra Site Agent interfaces
- Northbound - Temporal interface with Carbide REST API
- Southbound - Site Controller

We need the Elektra Site Agent Framework to be:
* Extensible
* Highly-Available
* Scalable

Elektra Site Agent manages Site workflows and store Site workflow information in DB

# Components

* Elektra Site agent
* NoSql Database
* Local Temporal
