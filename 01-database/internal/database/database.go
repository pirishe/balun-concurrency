package database

import (
	"context"
	"fmt"
	"raptor/internal/compute/parser"
	"raptor/internal/storage/engine"
)

type Database struct {
}

func New() Database {
	return Database{}
}

func (d Database) HandleQuery(ctx context.Context, input string) (string, error) {
	q, err := parser.Parse(input)
	if err != nil {
		return "", fmt.Errorf("query parse failed: %w", err)
	}
	switch q.Command() {
	case "GET":
		value, _, err := engine.Get(q.Key())
		if err != nil {
			return "", fmt.Errorf("get command failed: %w", err)
		}
		return value, nil
	case "SET":
		err = engine.Set(q.Key(), q.Value())
		if err != nil {
			return "", fmt.Errorf("set command failed: %w", err)
		}
		return "", nil
	case "DEL":
		err = engine.Del(q.Key())
		if err != nil {
			return "", fmt.Errorf("del command failed: %w", err)
		}
		return "", nil
	}
	panic("unreachable code")
}
