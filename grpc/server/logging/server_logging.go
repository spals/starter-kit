package logging

import (
	"context"
	"fmt"
	nativelog "log"
	"os"
	"path"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spals/starter-kit/grpc/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/status"
)

// An implementation of grpclog.LoggerV2 for Grpc core logging
// which delegates to zerolog
// See https://github.com/cheapRoc/grpc-zerolog
type zerologGrpcLoggerV2 struct {
	grpclog.LoggerV2

	delegate zerolog.Logger
}

func (l zerologGrpcLoggerV2) Fatal(args ...interface{}) {
	l.delegate.Fatal().Msg(fmt.Sprint(args...))
}

func (l zerologGrpcLoggerV2) Fatalf(format string, args ...interface{}) {
	l.delegate.Fatal().Msg(fmt.Sprintf(format, args...))
}

func (l zerologGrpcLoggerV2) Fatalln(args ...interface{}) {
	l.Fatal(args...)
}

func (l zerologGrpcLoggerV2) Error(args ...interface{}) {
	l.delegate.Error().Msg(fmt.Sprint(args...))
}

func (l zerologGrpcLoggerV2) Errorf(format string, args ...interface{}) {
	l.delegate.Error().Msg(fmt.Sprintf(format, args...))
}

func (l zerologGrpcLoggerV2) Errorln(args ...interface{}) {
	l.Error(args...)
}

func (l zerologGrpcLoggerV2) Info(args ...interface{}) {
	l.delegate.Debug().Msg(fmt.Sprint(args...))
}

func (l zerologGrpcLoggerV2) Infof(format string, args ...interface{}) {
	l.delegate.Debug().Msg(fmt.Sprintf(format, args...))
}

func (l zerologGrpcLoggerV2) Infoln(args ...interface{}) {
	l.Info(args...)
}

func (l zerologGrpcLoggerV2) Warning(args ...interface{}) {
	l.delegate.Warn().Msg(fmt.Sprint(args...))
}

func (l zerologGrpcLoggerV2) Warningf(format string, args ...interface{}) {
	l.delegate.Warn().Msg(fmt.Sprintf(format, args...))
}

func (l zerologGrpcLoggerV2) Warningln(args ...interface{}) {
	l.Warning(args...)
}

func (l zerologGrpcLoggerV2) Print(args ...interface{}) {
	l.delegate.Info().Msg(fmt.Sprint(args...))
}

func (l zerologGrpcLoggerV2) Printf(format string, args ...interface{}) {
	l.delegate.Info().Msg(fmt.Sprintf(format, args...))
}

func (l zerologGrpcLoggerV2) Println(args ...interface{}) {
	l.Print(args...)
}

func (l zerologGrpcLoggerV2) V(level int) bool {
	return true
}

var StreamRequestLogMiddleware grpc.StreamServerInterceptor
var UnaryRequestLogMiddleware grpc.UnaryServerInterceptor

func ConfigureLogging(config *proto.GrpcServerConfig) {
	// Set the default logger as the application logger
	log.Logger = newLogger(config).With().Str("system", "starter-kit-grpc").Logger()

	// Configure a logger to handle Grpc core logging
	grpcCoreLogger := newLogger(config).With().Str("system", "grpc-core").Logger()
	grpclog.SetLoggerV2(zerologGrpcLoggerV2{delegate: grpcCoreLogger})

	grpcUnaryRequestLogger := newLogger(config).With().Str("system", "grpc-request").Str("request-type", "unary").Logger()
	UnaryRequestLogMiddleware = newUnaryRequestLogMiddleware(grpcUnaryRequestLogger)
}

// ========== Private Helpers ==========

func newLogger(config *proto.GrpcServerConfig) zerolog.Logger {
	logLevel, err := zerolog.ParseLevel(config.GetLogLevel())
	if err != nil {
		nativelog.Fatalf("Error while parsing log level: %s. Available log levels are (trace|debug|info|warn|error|fatal|panic)", err)
	} else if logLevel == zerolog.NoLevel {
		nativelog.Fatalf("No log level configured. Please specify a log level (trace|debug|info|warn|error|fatal|panic)")
	}
	zerolog.SetGlobalLevel(logLevel)

	if config.GetDev() {
		output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
		output.FormatFieldName = func(i interface{}) string {
			return fmt.Sprintf("[%s]:", i)
		}

		return zerolog.New(output).With().Timestamp().Caller().Logger()
	} else {
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
		return zerolog.New(os.Stderr).With().Timestamp().Logger()
	}
}

// See https://www.gitmemory.com/issue/rs/zerolog/211/774897456
// See https://github.com/grpc-ecosystem/go-grpc-middleware/blob/master/logging/zap/server_interceptors.go
func newUnaryRequestLogMiddleware(reqLogger zerolog.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		startTime := time.Now()
		resp, err := handler(ctx, req)
		duration := time.Since(startTime)

		code := status.Code(err)
		level := reqLogLevel(code)
		reqLogEvent := reqLogger.WithLevel(level).
			Err(err).
			Str("grpc.code", code.String()).
			Str("grpc.method", path.Base(info.FullMethod)).
			Str("grpc.service", path.Dir(info.FullMethod)).
			Str("grpc.start_time", startTime.Format(time.RFC3339)).
			Dur("grpc.time_ms", duration)

		if d, ok := ctx.Deadline(); ok {
			reqLogEvent = reqLogEvent.Str("grpc.request.deadline", d.Format(time.RFC3339))
		}

		reqLogEvent.Msg("Finished unary call")
		return resp, err
	}
}

func reqLogLevel(code codes.Code) zerolog.Level {
	switch code {
	case codes.OK:
		return zerolog.InfoLevel
	case codes.Canceled:
		return zerolog.InfoLevel
	case codes.Unknown:
		return zerolog.InfoLevel
	case codes.InvalidArgument:
		return zerolog.InfoLevel
	case codes.DeadlineExceeded:
		return zerolog.InfoLevel
	case codes.NotFound:
		return zerolog.WarnLevel
	case codes.AlreadyExists:
		return zerolog.WarnLevel
	case codes.PermissionDenied:
		return zerolog.InfoLevel
	case codes.Unauthenticated:
		return zerolog.WarnLevel
	case codes.ResourceExhausted:
		return zerolog.WarnLevel
	case codes.FailedPrecondition:
		return zerolog.WarnLevel
	case codes.Aborted:
		return zerolog.DebugLevel
	case codes.OutOfRange:
		return zerolog.DebugLevel
	case codes.Unimplemented:
		return zerolog.WarnLevel
	case codes.Internal:
		return zerolog.WarnLevel
	case codes.Unavailable:
		return zerolog.WarnLevel
	case codes.DataLoss:
		return zerolog.WarnLevel
	default:
		return zerolog.InfoLevel
	}
}
