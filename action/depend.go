package action

import "strings"

type DependAction struct {
	ActionType string `json:"action_type,-"`
	FMRI       string `json:"fmri"`
	Type       string `json:"type"`
	Predicate  string `json:"predicate"`
	Optional   map[string]string `json:"optional"`
}

func (d *DependAction) FromActionString(actionString string) {
	d.Optional = make(map[string]string)
	for _, value := range tokenize(actionString) {
		equalpos := strings.Index(value, "=")
		key := value[0: equalpos]
		value = value[equalpos+1:]
		switch key {
		case "action_type":
			d.ActionType = value
		case "type":
			d.Type = value
		case "fmri":
			d.FMRI = value
		case "predicate":
			d.Predicate = value
		default:
			d.Optional[key] = value
		}
	}
}