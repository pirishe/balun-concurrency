package engine

import "fmt"

var storage = make(map[string]string)

func Set(key, value string) error {
	if key == "" {
		return fmt.Errorf("key cannot be empty")
	}
	storage[key] = value
	return nil
}

func Get(key string) (string, bool, error) {
	if key == "" {
		return "", false, fmt.Errorf("key cannot be empty")
	}
	value, ok := storage[key]
	return value, ok, nil
}

func Del(key string) error {
	if key == "" {
		return fmt.Errorf("key cannot be empty")
	}
	delete(storage, key)
	return nil
}
