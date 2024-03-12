package errors

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	pkgerrors "github.com/pkg/errors"
)

//TODO: need to move this out to the errors module

// ErrorLog struct is an structure for capturing detailed errors
// that implements the error interface so *ErrorLog can be passed
// anywhere errors are allowed.
type ErrorLog struct {
	RootCause             string `json:"RootCause,omitempty"` // like context msg
	Trace                 string `json:"Trace,omitempty"`
	StatusCode            string `json:"StatusCode,omitempty"`
	Source                string `json:"Source,omitempty"`
	Scope                 string `json:"Scope,omitempty"`
	Query                 string `json:"Query,omitempty"`
	AdditionalInformation string `json:"AdditionalInformation,omitempty"`
	ExceptionType         string `json:"ExceptionType,omitempty"`
	Err                   error  `json:"-"`
}

// New creates an ErrorLog structure that populates the Err
// field by creating an error from pkg/errors that contains the
// current stack
func New(errorString string) *ErrorLog {
	return FromError(pkgerrors.New(errorString))
}

// NewRootMsgStatusCode creates an ErrorLog structure that populates the RootCause,
// StatusCode, and Err field by creating an error from pkg/errors
// that contains the current stack
func NewRootMsgStatusCode(rootCause string, errmsg string, statusCode string) *ErrorLog {
	return &ErrorLog{
		RootCause:  rootCause,
		StatusCode: statusCode,
		Err:        pkgerrors.New(errmsg),
	}
}

// Errorf creates an ErrorLog structure that populates the Err
// field by creating an error from pkg/errors that contains the
// current stack using the format and args
func Errorf(format string, args ...interface{}) *ErrorLog {
	return FromError(pkgerrors.Errorf(format, args...))
}

// FromError creates an ErrorLog structure populating the Err
// field. If nil is passed as the err then a nil *ErrorLog is
// returned
func FromError(err error) *ErrorLog {
	if err == nil {
		return nil
	}
	// if it already is a ErrorLog simply return itself
	if errorLog, ok := err.(*ErrorLog); ok {
		return errorLog
	}
	return &ErrorLog{
		Err: err,
	}
}

// WithErrorAndCause creates an ErrorLog from an error and
// a RootCause, if err is nil then a nil is returned making
// it easy to use this to wrap whatever comes back from a function
func WithErrorAndCause(err error, rootCause string) *ErrorLog {
	if err == nil {
		return nil
	}
	return &ErrorLog{
		RootCause: rootCause,
		Err:       err,
	}
}

// String implements the stringer interface where objects
// can coerce themselves into a string
func (el ErrorLog) String() string {
	// include the error message in the output
	return errorLogToErrorString(true, &el)
}

// Error implements the error interface so this struct
// can be used anywhere an error can be used
func (el ErrorLog) Error() string {
	// include the error message in the output
	return errorLogToErrorString(true, &el)
}

// errorLogToErrorString is the implementation that converts
// data from ErrorLog to an error string. includeErr determines
// wheter the errorLog.err msg is included in the string, since
// in the format case it will be done separately.
func errorLogToErrorString(includeErr bool, el *ErrorLog) string {
	segments := []string{}
	if el.RootCause != "" {
		segments = append(segments, el.RootCause)
	}
	if el.Trace != "" {
		segments = append(segments, el.Trace)
	}
	if el.Err != nil && includeErr && el.Err.Error() != "" {
		segments = append(segments, el.Err.Error())
	}
	if el.StatusCode != "" {
		segments = append(segments, "StatusCode:"+el.StatusCode)
	}
	if el.Source != "" {
		segments = append(segments, "Source:"+el.Source)
	}
	if el.Scope != "" {
		segments = append(segments, "Scope:"+el.Scope)
	}

	if el.Query != "" {
		segments = append(segments, "Query:"+el.Query)
	}
	if el.AdditionalInformation != "" {
		segments = append(segments, "AdditionalInformation:"+el.AdditionalInformation)
	}
	if el.ExceptionType != "" {
		segments = append(segments, "ExceptionType:"+el.ExceptionType)
	}
	return strings.Join(segments, " ")
}

// Format implements the Formatter interface so an ErrorLog can be formatted.
// %v and %s outputs errorLog.Error()
// %q outputs a quoted errorLog.Error()
// %+v outputs the errorLog.Error() without the error and appends a full output
// of formatted Err including stack trace
func (el ErrorLog) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') { // %+v provides a more detailed error
			// not including err in string since is added below
			firstPart := errorLogToErrorString(false, &el)
			io.WriteString(s, firstPart)
			if el.Err != nil {
				if firstPart != "" {
					io.WriteString(s, " ")
				}
				if f, ok := (el.Err).(fmt.Formatter); ok { // is a formatter
					f.Format(s, verb)
				} else {
					io.WriteString(s, fmt.Sprintf("%v", el.Err))
				}
			}
			return
		}
		fallthrough
	case 's':
		io.WriteString(s, el.Error())
	case 'q':
		fmt.Fprintf(s, "%q", el.Error())
	}
}

// MarshalJSON implements custom marshalling for ErrorLog
// such that Trace and Err fields are merged
func (el ErrorLog) MarshalJSON() ([]byte, error) {
	type Alias ErrorLog
	return json.Marshal(&struct {
		Trace string
		Alias // embed ErrorLog fields, but not methods
	}{
		Trace: combineStringAndError(el.Trace, el.Err),
		Alias: (Alias)(el),
	})
}

// combineStringAndError combines the string and
// error string values inserting a space separator
// if needed
func combineStringAndError(str string, err error) string {
	// if both str and err then combine
	// if just str use it
	// if just err then use err.Error()
	arr := make([]string, 0, 2)
	if str != "" {
		arr = append(arr, str)
	}
	if err != nil {
		s := err.Error()
		if s != "" {
			arr = append(arr, s)
		}
	}
	return strings.Join(arr, " ")
}
