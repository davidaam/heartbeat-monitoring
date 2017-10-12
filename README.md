# Heartbeat Monitoring
Simple heartbeat monitoring system based on UDP protocol implemented in Go.

Client sents a UDP packet to server letting it know it's up. The server then registers this heartbeat in a database.
If after a configurable time period, the client doesn't send a heartbeat, the server triggers a notification to the
sysadmin.

For now, client doesn't wait for an ACK from server, if for some reason the packet is lost, the server is down, etc.
it will be up to register this downtime, later on a mechanism can be implemented to either wait for ACK, or register
heartbeat data on the client as well, so it can be resent to the server when it's up and running again.

## Requirements

- ###Go-sql
```
go get -u github.com/go-sql-driver/mysql
```
