package assert

import "fmt"

func Assert(mustBeTrue bool, formatStringAndArgs ...interface{}) {
	if !mustBeTrue {
		var message string
		if len(formatStringAndArgs) > 0 {
			switch v := formatStringAndArgs[0].(type) {
			case string:
				message = fmt.Sprintf(v, formatStringAndArgs[1:])
			default:
				message = "Invalid message type in call to Assert"
			}
		} else {
			message = "No message provided"
		}
		panic("Assert failure: " + message)
	}
}
