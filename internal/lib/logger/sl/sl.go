package sl

import "golang.org/x/exp/slog"

func Err(err error) slog.Attr {
	return slog.String(
		"error",
		err.Error(),
	)
}
