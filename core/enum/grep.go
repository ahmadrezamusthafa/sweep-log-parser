package enum

import (
	"strings"
)

type GrepType int

const (
	GrepStandard GrepType = 1 + iota
	GrepCombined
)

var format = [...]string{
	"grep '%s'",
	"grep -E '%s'",
}

var formatV2 = [...][]string{
	{"{{value}}"},
	{"-E", "{{value}}"},
}

func (c GrepType) Format() string {
	return format[c-1]
}

func (c GrepType) FormatV2(value string) []string {
	val := formatV2[c-1]
	for i, v := range val {
		val[i] = strings.Replace(v, "{{value}}", value, -1)
	}

	return val
}
