/*
Package errors contains a detailed error structure ErrorLog which implements the error, formatter, and stringer interfaces.

# Creating ErrorLog structures

There are a variety of helper factory functions to create *ErrorLog structs.

	  errorLog := New("myerror")
	  errorLog := NewRootMsgStatusCode("rootcause", "myerror", "206")
	  errorLog := Errorf("myerror %d %v", 10, "hey")
		errorLog := FromError(fmt.Errorf("myerror")
		errorLog := WithErrorAndCause(fmt.Errorf("myerror"), "MyRootCause")
	  errorLog := &ErrorLog{
	    RootCause:             "myroot",
	    Trace:                 "mytrace",
	    StatusCode:            "206",
	    Source:                "abs",
	    Scope:                 "GeneralInfo",
	    AdditionalInformation: "myinfo",
	    ExceptionType:         "myexc",
	    Err:                   pkgerrors.New("myerror"),
	  }

# Interfaces implemented

ErrorLog implements the error, formatter, and stringer interfaces. It also provides a custom JSON marshaller which merges the Trace and Err fields in the JSON output.

	fmt.Println(errorLog) // stringer interface allows errorLog.Error() be coerced to a string
	errorLog.Error()  // error interface outputs the fields as an error string
	fmt.Printf("%v", errorLog)  // formats the errorLog.Error()
	fmt.Printf("%+v", errorLog) // provides extended detailed output including callstack if err contains it
	bytes, err := json.Marshal(errorLog)   // custom marshalled JSON with Trace and Err merged
*/
package errors
