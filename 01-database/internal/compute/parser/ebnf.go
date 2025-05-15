package parser

import (
	"fmt"
	"regexp"
	"strings"
)

type Query struct {
	command string
	key     string
	value   string
}

func (q Query) Command() string {
	return q.command
}

func (q Query) Key() string {
	return q.key
}

func (q Query) Value() string {
	return q.value
}

func Parse(input string) (Query, error) {
	re := regexp.MustCompile(`^([^\s]+)[\s]+([^\s]+)(?:[\s]+([^\s]+))?$`)
	res := re.FindStringSubmatch(strings.TrimRight(input, "\r\n"))
	if len(res) != 4 {
		return Query{}, fmt.Errorf("failed to parse query, invalid string format: %s", input)
	}

	command, key, value := res[1], res[2], res[3]
	switch command {
	case "GET":
		fallthrough
	case "DEL":
		if value != "" {
			return Query{}, fmt.Errorf("%s command must have 1 argument, got 2: %s", command, input)
		}
	case "SET":
		if value == "" {
			return Query{}, fmt.Errorf("SET command must have 2 arguments, got 1: %s", input)
		}
	default:
		return Query{}, fmt.Errorf("unexpected command: %s %s %s", command, key, value)
	}
	return Query{command: command, key: key, value: value}, nil
}
