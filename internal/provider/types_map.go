package provider

import (
	"errors"
	"fmt"
	"sort"
)

// strIntMapKeysSorted returns a sorted []string containing the keys of m.
func strIntMapKeysSorted(m map[string]int) []string {
	res := make([]string, 0, len(m))

	for key := range m {
		res = append(res, key)
	}

	sort.Strings(res)

	return res
}

func reverseStrIntMap(m map[string]int) map[int]string {
	res := make(map[int]string, len(m))

	for k, v := range m {
		res[v] = k
	}

	return res
}

func strIntMapGet(m map[string]int, key string) (int, error) {
	if v, exists := m[key]; exists {
		return v, nil
	}

	return -1, fmt.Errorf("key %q not found", key)
}

func intStrMapGet(m map[int]string, key *int) (string, error) {
	if key == nil {
		return "", errors.New("key is nil")
	}

	if v, exists := m[*key]; exists {
		return v, nil
	}

	return "", fmt.Errorf("key '%d' not found", key)
}
