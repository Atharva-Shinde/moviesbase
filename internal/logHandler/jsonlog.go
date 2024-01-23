package loghandler

import (
	"io"
	"sync"
)

type Logger struct {
	logLevel int
	mu       sync.Mutex
	w        io.Writer
}

func New(level int, w io.Writer) *Logger {
	return &Logger{
		logLevel: level,
		w:        w,
	}
}
