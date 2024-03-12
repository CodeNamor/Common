package logging

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/ascarter/requestid"
	"github.com/stretchr/testify/assert"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func Test_NewLogger(t *testing.T) {
	l := New(TraceLevel)
	buffer, resetFn := redirectLoggerOutputToNewBuffer(l.(*Logger))
	defer resetFn()
	l.Error("ERR123")
	l.Trace("TRACE456")
	l.SetLevel(ErrorLevel)
	l.Trace("TRACE789")
	l.Error("ERRABC")
	expected := "level=error msg=ERR123\nlevel=trace msg=TRACE456\nlevel=error msg=ERRABC\n"
	require.Equal(t, expected, buffer.String())
}

func Test_PackageLoggingMethods(t *testing.T) {
	// package logging methods use the default logger
	buffer, resetFn := redirectLoggerOutputToNewBuffer(DefaultLogger())
	defer resetFn()
	// default level is error
	Error("ERR123")
	Trace("TRACE456")
	SetLevel(TraceLevel)
	defer SetLevel(ErrorLevel)
	Trace("TRACE789")
	expected := "level=error msg=ERR123\nlevel=trace msg=TRACE789\n"
	require.Equal(t, expected, buffer.String())
}

func Test_LoggerInfoLevelProvidesMsgs(t *testing.T) {
	l := New(InfoLevel)
	buffer, resetFn := redirectLoggerOutputToNewBuffer(l.(*Logger))
	defer resetFn()

	printMessageList(l)

	// Assert only 3 messages logged
	expected := "level=info msg=INFO456\nlevel=warning msg=WARN789\nlevel=error msg=ERRABC\n"
	require.Equal(t, expected, buffer.String())
}

func Test_LoggerWarningLevelSupressesInfoMsgs(t *testing.T) {
	l := New(WarningLevel)
	buffer, resetFn := redirectLoggerOutputToNewBuffer(l.(*Logger))
	defer resetFn()

	printMessageList(l)

	// Assert only 2 messages logged
	expected := "level=warning msg=WARN789\nlevel=error msg=ERRABC\n"
	require.Equal(t, expected, buffer.String())
}

func Test_LoggerErrorLevelSupressesInfoAndWarningMsgs(t *testing.T) {
	l := New(ErrorLevel)
	buffer, resetFn := redirectLoggerOutputToNewBuffer(l.(*Logger))
	defer resetFn()

	printMessageList(l)

	// Assert only 1 message logged
	expected := "level=error msg=ERRABC\n"
	require.Equal(t, expected, buffer.String())
}

func printMessageList(l Logging) {
	l.Trace("TRACE123")
	l.Info("INFO456")
	l.Warning("WARN789")
	l.Error("ERRABC")
}

func redirectLoggerOutputToNewBuffer(logger *Logger) (*bytes.Buffer, func()) {
	buffer := &bytes.Buffer{}

	origDefaultOutput := logger.LoggerImpl.Out
	logger.LoggerImpl.SetOutput(buffer)

	origDefaultFormatter := logger.LoggerImpl.Formatter
	logger.LoggerImpl.SetFormatter(&logrus.TextFormatter{DisableTimestamp: true})

	resetFn := func() {
		logger.LoggerImpl.SetOutput(origDefaultOutput)
		logger.LoggerImpl.SetFormatter(origDefaultFormatter)
	}
	return buffer, resetFn
}

func Test_DefaultLogger(t *testing.T) {
	require.Equal(t, defaultLogger, DefaultLogger())
}

func Test_ConfigureDefaultLoggingFromString(t *testing.T) {
	testcases := []struct {
		name            string
		env             string
		loggingLevelStr string
		expectedError   string
	}{
		{
			name:            "1 - invalid logging level",
			env:             "local",
			loggingLevelStr: "badlevel",
			expectedError:   "invalid level specified:badlevel",
		},
		{
			name:            "2 - valid logging level",
			env:             "test",
			loggingLevelStr: "info",
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			err := ConfigureDefaultLoggingFromString(tc.env, tc.loggingLevelStr)
			if tc.expectedError != "" {
				require.EqualError(t, err, tc.expectedError)
			} else { // not expecting error
				require.NoError(t, err)
			}
		})
	}
}

func Test_converLoggingLevelToConst(t *testing.T) {
	testcases := []struct {
		name          string
		loggingLevel  Level
		expectedLevel logrus.Level
		expectedError string
	}{
		{
			name:          "2 - valid logging level trace",
			loggingLevel:  TraceLevel,
			expectedLevel: logrus.TraceLevel,
		},
		{
			name:          "3 - valid logging level info",
			loggingLevel:  InfoLevel,
			expectedLevel: logrus.InfoLevel,
		},
		{
			name:          "3 - valid logging level warn",
			loggingLevel:  WarningLevel,
			expectedLevel: logrus.WarnLevel,
		},
		{
			name:          "4 - valid logging level error",
			loggingLevel:  ErrorLevel,
			expectedLevel: logrus.ErrorLevel,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			levelConst := convertLoggingLevelToConst(tc.loggingLevel)
			require.Equal(t, tc.expectedLevel, levelConst)
		})
	}

}

func Test_convertLoggingLevelStrToLoggingLevel(t *testing.T) {
	testcases := []struct {
		name            string
		loggingLevelStr string
		expectedLevel   Level
		expectedError   string
	}{
		{
			name:            "1 - invalid logging level",
			loggingLevelStr: "badlevel",
			expectedError:   "invalid level specified:badlevel",
		},
		{
			name:            "2 - valid logging level trace",
			loggingLevelStr: "trace",
			expectedLevel:   TraceLevel,
		},
		{
			name:            "3 - valid logging level info",
			loggingLevelStr: "info",
			expectedLevel:   InfoLevel,
		},
		{
			name:            "3 - valid logging level warn",
			loggingLevelStr: "warn",
			expectedLevel:   WarningLevel,
		},
		{
			name:            "4 - valid logging level error",
			loggingLevelStr: "error",
			expectedLevel:   ErrorLevel,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			level, err := convertLoggingLevelStrToLoggingLevel(tc.loggingLevelStr)
			if tc.expectedError != "" {
				require.EqualError(t, err, tc.expectedError)
			} else { // not expecting error
				require.NoError(t, err)
				require.Equal(t, tc.expectedLevel, level)
			}
		})
	}
}

func Test_WithRequestID(t *testing.T) {

	testcases := []struct {
		name        string
		mockContext context.Context
		expectedID  string
	}{
		{
			name:        "Request ID set on context should produce log entry with same ID",
			mockContext: requestid.NewContext(mockContext{}, "1234"),
			expectedID:  "1234",
		},
		{
			name:        "Request ID not set should produce log entry with empty string",
			mockContext: mockContext{},
			expectedID:  "",
		},
	}

	for _, test := range testcases {
		t.Run(test.name, func(t *testing.T) {

			entry := WithRequestID(test.mockContext)

			require.IsType(t, &logrus.Entry{}, entry)
			assert.Equal(t, test.expectedID, entry.Data["requestId"])
		})
	}
}

func Test_WithField(t *testing.T) {
	entry := WithField("mykey", "myvalue")
	expected := &logrus.Entry{}
	require.IsType(t, expected, entry)
}

func Test_WithFields(t *testing.T) {
	entry := WithFields(map[string]interface{}{"one": 1, "two": 2})
	expected := &logrus.Entry{}
	require.IsType(t, expected, entry)
}

type mockContext struct {
	Values map[interface{}]interface{}
}

func (c mockContext) Deadline() (deadline time.Time, ok bool) {
	return
}

func (c mockContext) Done() <-chan struct{} {
	return nil
}

func (c mockContext) Err() error {
	return nil
}

func (c mockContext) Value(key interface{}) interface{} {
	return c.Values[key]
}
