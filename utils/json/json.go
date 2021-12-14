//go:build !jsoniter
// +build !jsoniter

package json

import (
	"encoding/json"
)

var (
	// Marshal is exported by utils/json package.
	Marshal = json.Marshal
	// Unmarshal is exported by utils/json package.
	Unmarshal = json.Unmarshal
	// MarshalIndent is exported by utils/json package.
	MarshalIndent = json.MarshalIndent
	// NewDecoder is exported by utils/json package.
	NewDecoder = json.NewDecoder
	// NewEncoder is exported by utils/json package.
	NewEncoder = json.NewEncoder

	// MarshalString is exported by utils/json package.
	MarshalString = marshalString
	// UnmarshalString is exported by utils/json package.
	UnmarshalString = unmarshalString
)

func marshalString(v interface{}) (s string, err error) {
	bs, err := json.Marshal(v)
	if err == nil {
		s = BytesToString(bs)
	}
	return
}

func unmarshalString(s string, v interface{}) error {
	return json.Unmarshal(StringToBytes(s), v)
}
