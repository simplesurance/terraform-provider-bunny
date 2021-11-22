package provider

import (
	"sort"
	"strings"
)

func normalizeStrList(s string, sep rune) []string {
	res := strings.FieldsFunc(s, func(r rune) bool { return r == sep })

	for i := range res {
		res[i] = strings.TrimSpace(res[i])
	}

	sort.Strings(res)
	return res
}
