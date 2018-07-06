package godash

import "strings"

// DifferenceStr - returns a slice of elements in items that are not in vals
func DifferenceStr(items, vals []string) []string {
	diff := []string{}
	dstMap := ToMap(diff)

	for _, item := range items {
		if _, ok := dstMap[item]; !ok {
			diff = append(diff, item)
		}
	}

	return diff
}

// DifferenceSubstr - return a slice of elements in items excluding elements where any element from vals is a substring
func DifferenceSubstr(items, vals []string) []string {
	diff := []string{}

	for _, item := range items {
		if !IncludesStr(vals, item) {
			diff = append(diff, item)
		}
	}

	return diff
}

// IncludesStr - returns a boolean indicating if items contains key
func IncludesStr(items []string, key string) bool {
	for _, item := range items {
		if key == item || strings.Contains(key, item) {
			return true
		}
	}

	return false
}

// ToMap - convert a slice of strings to a map for fast lookups
func ToMap(items []string) map[string]string {
	itemMap := make(map[string]string)

	for _, item := range items {
		itemMap[item] = ""
	}

	return itemMap
}
