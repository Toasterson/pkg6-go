package action

import "strings"

type LinkAction struct {
	Type   string `json:"type,-"`
	Path   string `json:"path"`
	Target string `json:"target"`
}

func (a *LinkAction) FromActionString(actionString string) {
	for _, value := range tokenize(actionString) {
		equalpos := strings.Index(value, "=")
		key := value[0: equalpos]
		value = value[equalpos+1:]
		switch key {
		case "action_type":
			a.Type = value
		case "path":
			a.Path = value
		case "target":
			a.Target = value
		default:
		}
	}
}
