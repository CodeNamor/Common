package logging

import (
	"context"
	"fmt"
	"os"

	"github.com/CodeNamor/Common/logging/logfields"
	"github.com/ascarter/requestid"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type Level int

const (
	_ Level = iota // no longer supporting disabled level 0
	TraceLevel
	InfoLevel
	WarningLevel
	ErrorLevel
)

var levelToConst map[Level]logrus.Level
var levelStrToLoggingLevel map[string]Level
var defaultLogger *Logger

func init() {
	levelToConst = map[Level]logrus.Level{
		TraceLevel:   logrus.TraceLevel,
		InfoLevel:    logrus.InfoLevel,
		WarningLevel: logrus.WarnLevel,
		ErrorLevel:   logrus.ErrorLevel,
	}

	levelStrToLoggingLevel = map[string]Level{
		"trace": TraceLevel,
		"info":  InfoLevel,
		"warn":  WarningLevel,
		"error": ErrorLevel,
	}

	loggerImpl := logrus.StandardLogger()
	// default logger is error level before configuration
	loggerImpl.SetLevel(logrus.ErrorLevel)
	defaultLogger = &Logger{
		LoggerImpl: loggerImpl,
	}
}

func DefaultLogger() *Logger {
	return defaultLogger
}

func ConfigureDefaultLoggingFromString(env string, loggingLevelStr string) error {
	loggingLevel, err := convertLoggingLevelStrToLoggingLevel(loggingLevelStr)
	if err != nil {
		return err
	}
	return ConfigureDefaultLogging(env, loggingLevel)
}

func ConfigureDefaultLogging(env string, level Level) error {
	defaultLogger.SetLevel(level)

	// Default is os.StdErr
	defaultLogger.LoggerImpl.SetOutput(os.Stdout)

	// For any other environment we use the JSON log format to enable logstash auto-parsing
	if env != "local" {
		defaultLogger.LoggerImpl.SetFormatter(&logrus.JSONFormatter{})
	}

	return nil
}

type Logging interface {
	SetLevel(loggingLevel Level)
	Trace(args ...interface{})
	Info(args ...interface{})
	Warning(args ...interface{})
	Error(args ...interface{})
	WithField(key string, value interface{}) *logrus.Entry
	WithFields(fields logrus.Fields) *logrus.Entry
}

type Logger struct {
	LoggerImpl *logrus.Logger
}

func convertLoggingLevelToConst(loggingLevel Level) logrus.Level {
	logLevelConst, ok := levelToConst[loggingLevel]
	if !ok {
		defaultLogger.Warning(fmt.Sprintf("invalid level specified to SetLevel: %v", loggingLevel))
		return logrus.ErrorLevel
	}
	return logLevelConst
}

func convertLoggingLevelStrToLoggingLevel(levelStr string) (Level, error) {
	loggingLevel, ok := levelStrToLoggingLevel[levelStr]
	if !ok {
		return ErrorLevel, errors.Errorf("invalid level specified:%v", levelStr)
	}
	return loggingLevel, nil
}

func (l *Logger) SetLevel(loggingLevel Level) {
	l.LoggerImpl.SetLevel(convertLoggingLevelToConst(loggingLevel))
}

func (l *Logger) Trace(args ...interface{}) {
	l.LoggerImpl.Trace(args...)
}

func (l *Logger) Info(args ...interface{}) {
	l.LoggerImpl.Info(args...)
}

func (l *Logger) Warning(args ...interface{}) {
	l.LoggerImpl.Warning(args...)
}

func (l *Logger) Error(args ...interface{}) {
	l.LoggerImpl.Error(args...)
}

func (l *Logger) WithField(key string, value interface{}) *logrus.Entry {
	return l.LoggerImpl.WithField(key, value)
}

func (l *Logger) WithFields(fields logrus.Fields) *logrus.Entry {
	return l.LoggerImpl.WithFields(fields)
}

// Used when using the package directly

func New(level Level) Logging {
	logger := &Logger{
		LoggerImpl: logrus.New(),
	}
	logger.SetLevel(level)
	return logger
}

func SetLevel(loggingLevel Level) {
	defaultLogger.SetLevel(loggingLevel)
}

func Trace(args ...interface{}) {
	defaultLogger.Trace(args...)
}

func Info(args ...interface{}) {
	defaultLogger.Info(args...)
}

func Warning(args ...interface{}) {
	defaultLogger.Warning(args...)
}

func Error(args ...interface{}) {
	defaultLogger.Error(args...)
}

func WithRequestID(requestContext context.Context) *logrus.Entry {
	requestID := ""
	if idFromContext, ok := requestid.FromContext(requestContext); ok {
		requestID = idFromContext
	}
	return WithField(logfields.RequestId, requestID)
}

func WithField(key string, value interface{}) *logrus.Entry {
	return defaultLogger.LoggerImpl.WithField(key, value)
}

func WithFields(fields logrus.Fields) *logrus.Entry {
	return defaultLogger.LoggerImpl.WithFields(fields)
}
