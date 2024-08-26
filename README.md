# NATS Consumer

NATS consumer example with JWT auth and account jetstream subject export/import

Run with `./run.sh` or do the following:

1. Run setup with `go run ./setup.go`
1. Start NATS server with `docker-compose up -d`
1. Run application with `go run ./main.go`

See related issue: https://github.com/nats-io/nats.go/issues/1703

Possibly related issue: https://github.com/nats-io/nats.go/issues/1622
