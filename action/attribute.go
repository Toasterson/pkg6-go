package action

import "strings"

type AttributeAction struct {
	Type      string `json:"type,-"`
	Name      string `json:"name"`
	Values    []string `json:"values"`
	Optionals map[string]string `json:"optionals,omitempty"`
}

func (a *AttributeAction) FromActionString(actionString string) {
	a.Optionals = make(map[string]string)
	for _, value := range tokenize(actionString) {
		equalpos := strings.Index(value, "=")
		key := value[0: equalpos]
		value = value[equalpos+1:]
		switch key {
		case "action_type":
			a.Type = value
		case "name":
			a.Name = value
		case "value":
			a.Values = append(a.Values, value)
		default:
			a.Optionals[key] = value
		}
	}
}
