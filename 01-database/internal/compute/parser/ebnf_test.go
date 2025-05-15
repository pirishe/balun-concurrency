package parser

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseQuery(t *testing.T) {
	cases := []struct {
		desc        string
		input       string
		expectedErr string
		expectedRes Query
	}{
		{
			desc:        "empty input",
			input:       "",
			expectedErr: "failed to parse query",
		},
		{
			desc:        "invalid command",
			input:       "asd asd",
			expectedErr: "unexpected command",
		},
		{
			desc:  "DEL command",
			input: "DEL asd",
			expectedRes: Query{
				command: "DEL",
				key:     "asd",
				value:   "",
			},
		},
		{
			desc:        "DEL command must empty value",
			input:       "DEL asd qwe",
			expectedErr: "DEL command must have 1 argument, got 2",
		},
		{
			desc:  "GET command",
			input: "GET qwe",
			expectedRes: Query{
				command: "GET",
				key:     "qwe",
				value:   "",
			},
		},
		{
			desc:        "GET command must empty value",
			input:       "GET asd qwe",
			expectedErr: "GET command must have 1 argument, got 2",
		},
		{
			desc:  "SET command",
			input: "SET qwe oiu",
			expectedRes: Query{
				command: "SET",
				key:     "qwe",
				value:   "oiu",
			},
		},
		{
			desc:        "SET command must have value",
			input:       "SET asd",
			expectedErr: "SET command must have 2 arguments, got 1",
		},
		{
			desc:        "SET command must have value",
			input:       "SET asd ",
			expectedErr: "failed to parse query, invalid string format",
		},
	}

	for _, tc := range cases {
		res, err := Parse(tc.input)
		if tc.expectedErr != "" {
			require.ErrorContains(t, err, tc.expectedErr)
		} else if res != tc.expectedRes {
			t.Fatalf("test failed %s: expected: %v, got: %v", tc.desc, tc.expectedRes, res)
		}
	}
}
