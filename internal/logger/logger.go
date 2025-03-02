// Package logger provides a function Get that returns a zerolog.Logger
// instance for global logging, allowing consistent logging throughout
// the application.
package logger

import (
	"io"
	"os"
	"runtime/debug"
	"slices"
	"sync"
	"time"

	"github.com/KengoWada/meetup-clone/internal"
	"github.com/KengoWada/meetup-clone/internal/config"
	"github.com/KengoWada/meetup-clone/internal/utils"
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

		appEnv := config.AppEnv(utils.EnvGetString("SERVER_ENVIRONMENT", string(config.AppEnvProd)))
		if !slices.Contains(config.Environments, appEnv) {
			appEnv = config.AppEnvProd
		}

		if appEnv == config.AppEnvProd {
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

		logLevel := utils.EnvGetInt("LOG_LEVEL", int(zerolog.InfoLevel))
		if appEnv == config.AppEnvTest {
			logLevel = int(zerolog.Disabled)
		}

		log = zerolog.New(output).
			Level(zerolog.Level(logLevel)).
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
