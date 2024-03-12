package url

import (
	"bytes"
)

// QueryParam is a key value structure used for creating a ErrorLog.Query
// with FormatQuery.
type QueryParam struct {
	Key   string
	Value string
}

// FormatQuery creates a string that could be used for ErrorLog.Query
// It combines a prefix and any non-empty queryParams that are passed in.
// The queryParams will be included in the order they are passed into
// the slice.
func FormatQuery(prefix string, queryParams []QueryParam) string {
	buff := &bytes.Buffer{}
	buff.WriteString(prefix)
	hasNonEmptyValues := false // will be set true once we hit first one

	for _, kv := range queryParams {

		if kv.Value != "" { // only include non-empty pairs

			if hasNonEmptyValues { // always except first one
				buff.WriteString(", ")
			}

			buff.WriteString(kv.Key)
			buff.WriteString(": ")
			buff.WriteString(kv.Value)

			hasNonEmptyValues = true
		}
	}

	return buff.String()
}
