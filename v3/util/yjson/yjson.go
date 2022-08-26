package yjson

import jsoniter "github.com/json-iterator/go"

var YJson = jsoniter.Config{
	EscapeHTML:             true,
	ValidateJsonRawMessage: true,
	TagKey:                 "yjson",
}.Froze()

func SetDebug() {
	YJson = jsoniter.Config{
		EscapeHTML:             true,
		ValidateJsonRawMessage: true,
		TagKey:                 "yjson",
		SortMapKeys:            true,
	}.Froze()
}
