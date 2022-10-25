# go-simple-bank
Following Tech School Backend Masterclass to design, develop, and deploy a complete backend system from scratch using Golang, PostgreSQL and Docker (by the moment).

## How to play:
#### 1. Clone repo
#### 2. Run docker compose from root
```
$ docker compose up
```
_UNDER THE HOOD: Creates a docker network and attaches postgres and golang containers to it. Waits until the DB is up, runs the migrations and starts the BANK API._
