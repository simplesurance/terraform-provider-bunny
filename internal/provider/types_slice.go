package provider

import (
	"sort"
	"strings"
)

// strSliceEqual returns true if all elements in the slices are equal.
func strSliceEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

// strSliceToLower converts all elements in strs to lowercase.
func strSliceToLower(strs []string) {
	for i := range strs {
		strs[i] = strings.ToLower(strs[i])
	}
}

// strSliceAsInterfaceSlice returns an []interface{} containing all elements of strs.
func strSliceAsInterfaceSlice(strs []string) []interface{} {
	res := make([]interface{}, len(strs))

	for i, elem := range strs {
		res[i] = elem
	}

	return res
}

func interfaceSlicetoStrSlice(in []interface{}) []string {
	res := make([]string, len(in))
	for i, elem := range in {
		res[i] = elem.(string)
	}

	return res
}

func strSliceAsNormalizedStr(in []string, sep string) string {
	if len(in) == 0 {
		return ""
	}

	normalizedSl := make([]string, len(in))
	for i := range in {
		normalizedSl[i] = strings.TrimSpace(in[i])
	}

	sort.Strings(normalizedSl)

	return strings.Join(normalizedSl, sep)
}
