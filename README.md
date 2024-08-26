# NATS Consumer

NATS consumer example with JWT auth and account jetstream subject export/import

Run with `./run.sh` or do the following:

1. Run setup with `go run ./setup.go`
1. Start NATS server with `docker-compose up -d`
1. Run application with `go run ./main.go`

See related issue: https://github.com/nats-io/nats.go/issues/1703

Possibly related issue: https://github.com/nats-io/nats.go/issues/1622

Here's the output from `main.go`:

```
u1 connected to nats
u2 connected to nats
publishing 2024-08-26T11:27:19-06:00 1
publishing 2024-08-26T11:27:19-06:00 2
publishing 2024-08-26T11:27:19-06:00 3
publishing 2024-08-26T11:27:19-06:00 4
publishing 2024-08-26T11:27:19-06:00 5
publishing 2024-08-26T11:27:19-06:00 6
publishing 2024-08-26T11:27:19-06:00 7
publishing 2024-08-26T11:27:19-06:00 8
publishing 2024-08-26T11:27:19-06:00 9
publishing 2024-08-26T11:27:19-06:00 10
consuming 2024-08-26T11:27:19-06:00 order 1
__consumer__ delivered-sequence: 0 0 ackfloor-sequence: 0 0 sent: 0 redelivered: 0 pending: 1 ack-pending: 0 pulls-waiting: 1
__consumer__ delivered-sequence: 10 10 ackfloor-sequence: 10 10 sent: 1 redelivered: 0 pending: 0 ack-pending: 9 pulls-waiting: 1
__consumer__ delivered-sequence: 10 10 ackfloor-sequence: 10 10 sent: 1 redelivered: 0 pending: 0 ack-pending: 9 pulls-waiting: 1
consume error nats: no heartbeat received
consuming 2024-08-26T11:27:49-06:00 order 2
__consumer__ delivered-sequence: 19 10 ackfloor-sequence: 19 10 sent: 2 redelivered: 8 pending: 0 ack-pending: 8 pulls-waiting: 1
__consumer__ delivered-sequence: 19 10 ackfloor-sequence: 19 10 sent: 2 redelivered: 8 pending: 0 ack-pending: 8 pulls-waiting: 1
__consumer__ delivered-sequence: 19 10 ackfloor-sequence: 19 10 sent: 2 redelivered: 8 pending: 0 ack-pending: 8 pulls-waiting: 1
consume error nats: no heartbeat received
consuming 2024-08-26T11:28:19-06:00 order 3
__consumer__ delivered-sequence: 27 10 ackfloor-sequence: 27 10 sent: 3 redelivered: 7 pending: 0 ack-pending: 7 pulls-waiting: 1
__consumer__ delivered-sequence: 27 10 ackfloor-sequence: 27 10 sent: 3 redelivered: 7 pending: 0 ack-pending: 7 pulls-waiting: 1
__consumer__ delivered-sequence: 27 10 ackfloor-sequence: 27 10 sent: 3 redelivered: 7 pending: 0 ack-pending: 7 pulls-waiting: 1
consume error nats: no heartbeat received
consuming 2024-08-26T11:28:49-06:00 order 4
__consumer__ delivered-sequence: 34 10 ackfloor-sequence: 34 10 sent: 4 redelivered: 6 pending: 0 ack-pending: 6 pulls-waiting: 1
__consumer__ delivered-sequence: 34 10 ackfloor-sequence: 34 10 sent: 4 redelivered: 6 pending: 0 ack-pending: 6 pulls-waiting: 1
__consumer__ delivered-sequence: 34 10 ackfloor-sequence: 34 10 sent: 4 redelivered: 6 pending: 0 ack-pending: 6 pulls-waiting: 1
consume error nats: no heartbeat received
consuming 2024-08-26T11:29:19-06:00 order 6
__consumer__ delivered-sequence: 40 10 ackfloor-sequence: 40 10 sent: 4 redelivered: 5 pending: 0 ack-pending: 5 pulls-waiting: 1
__consumer__ delivered-sequence: 40 10 ackfloor-sequence: 40 10 sent: 4 redelivered: 5 pending: 0 ack-pending: 5 pulls-waiting: 1
__consumer__ delivered-sequence: 40 10 ackfloor-sequence: 40 10 sent: 4 redelivered: 5 pending: 0 ack-pending: 5 pulls-waiting: 1
consume error nats: no heartbeat received
consuming 2024-08-26T11:29:49-06:00 order 5
__consumer__ delivered-sequence: 45 10 ackfloor-sequence: 45 10 sent: 4 redelivered: 5 pending: 0 ack-pending: 5 pulls-waiting: 1
__consumer__ delivered-sequence: 45 10 ackfloor-sequence: 45 10 sent: 6 redelivered: 4 pending: 0 ack-pending: 4 pulls-waiting: 1
__consumer__ delivered-sequence: 45 10 ackfloor-sequence: 45 10 sent: 6 redelivered: 4 pending: 0 ack-pending: 4 pulls-waiting: 1
consume error nats: no heartbeat received
consuming 2024-08-26T11:30:19-06:00 order 7
__consumer__ delivered-sequence: 49 10 ackfloor-sequence: 49 10 sent: 7 redelivered: 3 pending: 0 ack-pending: 3 pulls-waiting: 1
__consumer__ delivered-sequence: 49 10 ackfloor-sequence: 49 10 sent: 7 redelivered: 3 pending: 0 ack-pending: 3 pulls-waiting: 1
__consumer__ delivered-sequence: 49 10 ackfloor-sequence: 49 10 sent: 7 redelivered: 3 pending: 0 ack-pending: 3 pulls-waiting: 1
consume error nats: no heartbeat received
consuming 2024-08-26T11:30:49-06:00 order 8
__consumer__ delivered-sequence: 52 10 ackfloor-sequence: 52 10 sent: 8 redelivered: 2 pending: 0 ack-pending: 2 pulls-waiting: 1
__consumer__ delivered-sequence: 52 10 ackfloor-sequence: 52 10 sent: 8 redelivered: 2 pending: 0 ack-pending: 2 pulls-waiting: 1
__consumer__ delivered-sequence: 52 10 ackfloor-sequence: 52 10 sent: 8 redelivered: 2 pending: 0 ack-pending: 2 pulls-waiting: 1
consume error nats: no heartbeat received
consuming 2024-08-26T11:31:19-06:00 order 9
__consumer__ delivered-sequence: 54 10 ackfloor-sequence: 54 10 sent: 9 redelivered: 1 pending: 0 ack-pending: 1 pulls-waiting: 1
__consumer__ delivered-sequence: 54 10 ackfloor-sequence: 54 10 sent: 9 redelivered: 1 pending: 0 ack-pending: 1 pulls-waiting: 1
__consumer__ delivered-sequence: 54 10 ackfloor-sequence: 54 10 sent: 9 redelivered: 1 pending: 0 ack-pending: 1 pulls-waiting: 1
consume error nats: no heartbeat received
consuming 2024-08-26T11:31:49-06:00 order 10
```
