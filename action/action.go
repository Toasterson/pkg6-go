package action

import (
	"strings"
)

type AttributeAction struct {
	Type string `json:"type,-"`
	Name string `json:"name"`
	Values []string `json:"values"`
	Optionals map[string]string `json:"optionals,omitempty"`
}

type DependAction struct {
	ActionType string `json:"action_type,-"`
	FMRI string `json:"fmri"`
	Type string `json:"type"`
}

func (a *AttributeAction) FromActionString(action_string string){
	for key, value := range tokenize(action_string){
		switch key {
		case "action_type" : a.Type = value
		case "name" : a.Name = value
		case "value" : a.Values = append(a.Values, value)
		default: a.Optionals[key] = value
		}
	}
}

func (d *DependAction) FromActionString(action_string string){
	for key, value := range tokenize(action_string){
		switch key {
		case "action_type": d.ActionType = value
		case "type": d.Type = value
		case "fmri": d.FMRI = value
		default:
		}
	}
}

func tokenize(action_string string) map[string]string {
	retVal := map[string]string{}
	typespacepos := strings.Index(action_string, " ")
	retVal["action_type"] = action_string[0:typespacepos]
	action_string = action_string[typespacepos+1:]
	for strings.Contains(action_string, "="){
		var key, value string
		equalpos := strings.Index(action_string, "=")
		key = action_string[0:equalpos]
		action_string = action_string[equalpos+1:]
		if strings.Contains(action_string, "="){
			secondequalpos := strings.Index(action_string, "=")
			spacepos := strings.Index(action_string[0:secondequalpos], " ")
			value = action_string[0:spacepos]
			action_string = action_string[spacepos+1:]
		} else {
			value = action_string
			action_string = ""
		}
		retVal[key] = value
	}
	return retVal
}
