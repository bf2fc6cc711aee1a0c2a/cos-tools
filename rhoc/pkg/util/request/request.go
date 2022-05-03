package request

import (
	"strconv"

	"github.com/antihax/optional"
)

type ListOptions struct {
	Page     int
	Limit    int
	AllPages bool
	OrderBy  string
	Search   string
}

func OptionalString(value string) optional.String {
	if value == "" {
		return optional.EmptyString()
	}

	return optional.NewString(value)
}

func OptionalInt(value int) optional.String {
	return optional.NewString(strconv.Itoa(value))
}

func OptionalBool(value bool) optional.String {
	return optional.NewString(strconv.FormatBool(value))
}
