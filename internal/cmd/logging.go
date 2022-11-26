package cmd

import (
	"io"

	"github.com/rs/zerolog"
)

// LevelWriter implements zerolog.LevelWriter.
type LevelWriter struct {
	io.Writer
	Level zerolog.Level
}

// NewLevelWriter returns a new implementation of zerolog.LevelWriter.
func NewLevelWriter(level zerolog.Level, w io.Writer) *LevelWriter {
	return &LevelWriter{
		Level:  level,
		Writer: w,
	}
}

// WriteLevel implements zerolog.LevelWriter.WriteLevel.
func (lw *LevelWriter) WriteLevel(l zerolog.Level, p []byte) (n int, err error) {
	if l >= lw.Level { // Notice that it's ">=", not ">"
		return lw.Writer.Write(p)
	}
	return len(p), nil
}
