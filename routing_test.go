package grom

import (
	"bytes"
	"sort"
)

type Ctx struct{}

type routeTest struct {
	route string
	get   string
	vars  map[string]string
}

// Converts the map into a consistent, string-comparable string (to compare with another map)
// Eg, stringifyMap({"foo": "bar"}) == stringifyMap({"foo": "bar"})
func stringifyMap(m map[string]string) string {
	if m == nil {
		return ""
	}

	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	keysLenMinusOne := len(keys) - 1

	var b bytes.Buffer
	b.WriteString("[")
	for i, k := range keys {
		b.WriteString(k)
		b.WriteRune(':')
		b.WriteString(m[k])

		if i != keysLenMinusOne {
			b.WriteRune(' ')
		}
	}
	b.WriteRune(']')

	return b.String()
}
