# raptor

### Introduction
Raptor - is an in-memory key-value database with support for asynchronous replication (physical) and WAL (Write-Ahead Logging).

### Server launch

```bash
go run cmd/server.go
```

### CLI launch

```bash
go run cmd/client.go
```

### See logs

```bash
tail -f raptor_client.log raptor_server.log
```

### Run tests

```bash
go test -race ./internal/...
```
