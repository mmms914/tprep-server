package internal

import (
	"github.com/google/uuid"
	"github.com/gookit/slog"
)

func GenerateUUID() string {
	gen, err := uuid.NewRandom()
	if err != nil {
		slog.FatalErr(err)
	}
	return gen.String()
}
