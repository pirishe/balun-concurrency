package engine

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStorageInMemory(t *testing.T) {
	value, ok, err := Get("")
	require.ErrorContains(t, err, "key cannot be empty")
	value, ok, err = Get("asd")
	require.NoError(t, err)
	require.False(t, ok)
	require.Empty(t, value)
	err = Del("")
	require.ErrorContains(t, err, "key cannot be empty")
	err = Del("asd")
	require.NoError(t, err)
	err = Set("", "")
	require.ErrorContains(t, err, "key cannot be empty")
	err = Set("key1", "")
	value, ok, err = Get("key1")
	require.NoError(t, err)
	require.True(t, ok)
	require.Empty(t, value)
	err = Set("key1", "111")
	value, ok, err = Get("key1")
	require.NoError(t, err)
	require.True(t, ok)
	require.Equal(t, "111", value)
	err = Del("key1")
	require.NoError(t, err)
	value, ok, err = Get("key1")
	require.NoError(t, err)
	require.False(t, ok)
	require.Empty(t, value)
}
