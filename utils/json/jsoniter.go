//go:build jsoniter
// +build jsoniter

package json

import (
	jsoniter "github.com/json-iterator/go"
)

var (
	json = jsoniter.ConfigCompatibleWithStandardLibrary
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
	MarshalString = json.MarshalToString
	// UnmarshalString is exported by utils/json package.
	UnmarshalString = json.UnmarshalFromString
)
