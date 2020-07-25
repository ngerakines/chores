package chores

import (
	"strings"
)

// Config is neat.
type Config struct {
	Address  string
	Database string
	Init     bool
	People   StringsFlag
}

// StringsFlag is neat.
type StringsFlag []string

func (f *StringsFlag) String() string {
	return strings.Join(*f, ",")
}

// Set is neat.
func (f *StringsFlag) Set(value string) error {
	for _, field := range strings.Split(value, ",") {
		field := strings.TrimSpace(field)
		if len(field) > 0 {
			*f = append(*f, strings.TrimSpace(field))
		}
	}
	return nil
}
