# replicated-state-machine

A minimal example on using the [hashicorp/raft](https://github.com/hashicorp/raft) package with go inspired by [this talk](https://www.youtube.com/watch?v=EGRmmxVFOfE).

## Requirements

- Go >= 1.18

## Usage

- Follower

```bash
go run cmd/server/main.go -config static/config/follower.toml
```

- Bootstrap

```bash
go run cmd/server/main.go -config static/config/bootstrap.toml
```

## Requests

Change the state of the application

```bash
curl http://localhost:3001/update -H 'Content-Type: application/json' -d '{"value": 10}'
```
