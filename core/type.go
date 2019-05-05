package core

import "github.com/tokopedia/sweep-log/core/enum"

const (
	NOTIFY_SUCCESS = 1
	VALIDATE_USE   = 2
)

type Filter struct {
	GrepType enum.GrepType
	Value    string
}
