package logger

import "io"

type writerAdapter struct {
	logger Logger
	level  Level
}

func NewWriterAdapter(logger Logger, level Level) io.Writer {
	return &writerAdapter{
		logger: logger,
		level:  level,
	}
}

func (w *writerAdapter) Write(p []byte) (n int, err error) {
	w.logger.Print(w.level, string(p))
	return len(p), nil
}
