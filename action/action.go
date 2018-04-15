package action

import (
	"strings"
	"fmt"
)

func tokenize(actionString string) []string {
	retVal := []string{}
	typespacepos := strings.Index(actionString, " ")
	retVal = append(retVal, fmt.Sprintf("%s=%s", "action_type", actionString[0:typespacepos]))
	actionString = actionString[typespacepos+1:]
	for strings.Contains(actionString, "=") {
		var key, value string
		equalpos := strings.Index(actionString, "=")
		key = actionString[0:equalpos]
		actionString = actionString[equalpos+1:]
		if strings.Contains(key, " ") {
			keyspacepos := strings.LastIndex(key, " ")
			keyval := key[0:keyspacepos]
			keyval = cleanFromChars(keyval)
			key = key[keyspacepos+1:]
			retVal = append(retVal, fmt.Sprintf("key=%s", keyval))
		}
		if strings.Contains(actionString, "=") && strings.Contains(actionString, " ") {
			secondequalpos := strings.Index(actionString, "=")
			spacepos := strings.LastIndex(actionString[0:secondequalpos], " ")
			value = actionString[0:spacepos]
			actionString = actionString[spacepos+1:]
		} else {
			value = actionString
			actionString = ""
		}
		value = cleanFromChars(value)
		retVal = append(retVal, fmt.Sprintf("%s=%s", key, value))
	}
	return retVal
}

func cleanFromChars(input string) string {
	input = strings.Replace(input, "\"", "", -1)
	input = strings.Replace(input, "\\\"'", "", -1)
	input = strings.Replace(input, "\\\"", "", -1)
	input = strings.Replace(input, "\\'", "", -1)
	input = strings.Replace(input, "'\\", "", -1)
	input = strings.Replace(input, "\\", "", -1)
	return input
}
