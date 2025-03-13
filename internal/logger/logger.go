// Package logger provides a function Get that returns a zerolog.Logger
// instance for global logging, allowing consistent logging throughout
// the application.
package logger

import (
	"io"
	"net/http"
	"os"
	"runtime/debug"
	"sync"
	"time"

	"github.com/KengoWada/meetup-clone/internal"
	"github.com/KengoWada/meetup-clone/internal/config"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	once sync.Once
	log  zerolog.Logger
)

// Get returns a singleton instance of the zerolog.Logger.
// It uses the sync.Once mechanism to ensure that the logger is created only once,
// regardless of how many times the function is called, providing a global logger
// instance for consistent logging throughout the application.
func Get() zerolog.Logger {
	once.Do(func() {
		cfg := config.Get()
		zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
		zerolog.TimeFieldFormat = time.RFC3339Nano

		var output io.Writer = zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: internal.DateTimeFormat,
			FieldsExclude: []string{
				"user_agent",
				"git_revision",
				"go_version",
			},
		}

		if cfg.Environment == config.AppEnvProd {
			fileLogger := &lumberjack.Logger{
				Filename:   "meetup_clone.log",
				MaxSize:    5,
				MaxBackups: 10,
				MaxAge:     14,
				Compress:   true,
			}

			output = zerolog.MultiLevelWriter(os.Stderr, fileLogger)
		}

		var gitRevision string

		buildInfo, ok := debug.ReadBuildInfo()
		if ok {
			for _, v := range buildInfo.Settings {
				if v.Key == "vcs.revision" {
					gitRevision = v.Value
					break
				}
			}
		}

		log = zerolog.New(output).
			Level(zerolog.Level(cfg.LogLevel)).
			With().
			Stack().
			Caller().
			Timestamp().
			Str("git_revision", gitRevision).
			Str("go_version", buildInfo.GoVersion).
			Logger()

		zerolog.DefaultContextLogger = &log
	})

	return log
}

func ErrLoggerCache(r *http.Request, err error) {
	logger := Get()

	reqIDRaw := middleware.GetReqID(r.Context())
	logger.Error().
		Str("requestID", reqIDRaw).
		Str("method", r.Method).
		Str("url", r.URL.Path).
		Err(errors.Wrap(err, "cache error")).
		Msg("Cache Error")
}
