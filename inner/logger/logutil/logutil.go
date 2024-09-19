package logutil

import (
	"log/slog"
)

// Using to add error to the slog as attribute
func Err(err error) slog.Attr {
	return slog.Attr{
		Key:   "error",
		Value: slog.StringValue(err.Error()),
	}
}
