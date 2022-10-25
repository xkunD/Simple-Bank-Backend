# Simple Bank Backend Project
🏦 In this backend project, I designed, developed and deployed a backend web service for a simple bank with Golang, Postgres, Gin, Kubernetes, gRPC, AWS.

This project can provide APIs for the frontend to do the following things:
- Create and manage bank accounts.
- Record all balance changes to each of the accounts.
- Perform a money transfer between 2 accounts.

## How to play:
#### 1. Clone repo
#### 2. Run docker compose from root
```
$ docker compose up
```
_UNDER THE HOOD: Creates a docker network and attaches postgres and golang containers to it. Waits until the DB is up, runs the migrations and starts the BANK API._
