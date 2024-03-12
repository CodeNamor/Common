/*
Below are keys for fields used in logging using the "github.com/sirupsen/logrus" structured logger

# For example, when trace logging the function called, use the instruction

logging.WithFields( logfields.Function, "aFnName" ).Trace()

In this way the logger will include "function=aFnName" in the logging statement. The goal of having these constants for
field names is so that the fields are consistently labeled across the APIs.
*/
package logfields

// Common fields across all APIs
const (
	RequestId     = "requestId"
	RemoteAddr    = "remoteAddr"
	URI           = "uri"
	Elapsed       = "elapsed"
	Request       = "request"
	StatusCode    = "statusCode"
	Function      = "function"
	RequestParams = "requestParams"
	ResultsCount  = "resultsCount"
	ErrorsCount   = "errorsCount"
	RequestURL    = "requestURL"
	ServiceName   = "serviceName"
)

// Common fields for error logs
const (
	RootCause   = "rootCause"
	ErrorSource = "errorSource"
)

// Common fields testing logs
const (
	IsTest   = "isTest"
	TestName = "testName"
)
