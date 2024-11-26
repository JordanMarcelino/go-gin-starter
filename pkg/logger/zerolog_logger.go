package logger

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/JordanMarcelino/go-gin-starter/internal/config"
	"github.com/rs/zerolog"
)

type zerologLogger struct {
	log *zerolog.Logger
}

func SetZerologLogger(cfg *config.Config) {
	log := zerolog.
		New(zerolog.ConsoleWriter{Out: os.Stdout, TimeLocation: time.UTC, TimeFormat: zerolog.TimeFormatUnix}).
		Level(zerolog.Level(cfg.Logger.Level)).
		With().
		Timestamp().
		Logger()

	log = log.With().Caller().Int("pid", os.Getpid()).Logger()
	zerologLogger := &zerologLogger{log: &log}

	SetLogger(zerologLogger)
}

func (l *zerologLogger) GetWriter() io.Writer {
	return l.log
}

func (l *zerologLogger) Printf(format string, args ...any) {
	l.log.Printf(format, args...)
}

func (l *zerologLogger) Error(args ...any) {
	l.log.Error().Msg(argsToString(args...))
}

func (l *zerologLogger) Errorf(format string, args ...any) {
	l.log.Error().Msgf(format, args...)
}

func (l *zerologLogger) Fatal(args ...any) {
	l.log.Fatal().Msg(argsToString(args...))
}

func (l *zerologLogger) Fatalf(format string, args ...any) {
	l.log.Fatal().Msgf(format, args...)
}

func (l *zerologLogger) Info(args ...any) {
	l.log.Info().Msg(argsToString(args...))
}

func (l *zerologLogger) Infof(format string, args ...any) {
	l.log.Info().Msgf(format, args...)
}

func (l *zerologLogger) Warn(args ...any) {
	l.log.Warn().Msg(argsToString(args...))
}

func (l *zerologLogger) Warnf(format string, args ...any) {
	l.log.Warn().Msgf(format, args...)
}

func (l *zerologLogger) Debug(args ...any) {
	l.log.Debug().Msg(argsToString(args...))
}

func (l *zerologLogger) Debugf(format string, args ...any) {
	l.log.Debug().Msgf(format, args...)
}

func (l *zerologLogger) WithField(key string, value any) Logger {
	var log zerolog.Logger
	if err, ok := value.(error); ok {
		log = l.log.With().Err(err).Logger()
	} else {
		log = l.log.With().Any(key, value).Logger()
	}

	return &zerologLogger{
		log: &log,
	}
}

func (l *zerologLogger) WithFields(fields map[string]any) Logger {
	logCtx := l.log.With()
	for k, v := range fields {
		if errs, ok := v.([]error); ok {
			logCtx = logCtx.Errs(k, errs)
		} else if err, ok := v.(error); ok {
			logCtx = logCtx.AnErr(k, err)
		} else {
			logCtx = logCtx.Any(k, v)
		}
	}

	log := logCtx.Logger()
	return &zerologLogger{
		log: &log,
	}
}

func argsToString(args ...any) string {
	var sb strings.Builder
	for _, arg := range args {
		sb.WriteString(fmt.Sprintf("%v ", arg))
	}

	return strings.TrimSpace(sb.String())
}
