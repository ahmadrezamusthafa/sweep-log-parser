package core

import "github.com/tokopedia/sweep-log/core/enum"

const (
	NOTIFY_SUCCESS = 1
	VALIDATE_USE   = 2

	MODE_ONE_BY_ONE = 1
	MODE_ALL_IN_ONE = 2
)

type Filter struct {
	GrepType enum.GrepType
	Value    string
}
