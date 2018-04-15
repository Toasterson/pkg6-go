package action

import "strings"

type DirectoryAction struct {
	Type   string `json:"type,-"`
	Path   string `json:"path"`
	Mode   string `json:"mode"`
	Owner  string `json:"owner"`
	Group  string `json:"group"`
	Facets map[string]string `json:"facets"`
}

func (a *DirectoryAction) FromActionString(actionString string) {
	a.Facets = make(map[string]string)
	for _, value := range tokenize(actionString) {
		equalpos := strings.Index(value, "=")
		key := value[0: equalpos]
		value = value[equalpos+1:]
		switch key {
		case "action_type":
			a.Type = value
		case "path":
			a.Path = value
		case "mode":
			a.Mode = value
		case "owner":
			a.Owner = value
		case "group":
			a.Group = value
		default:
			if strings.Contains(key, "facet.") {
				a.Facets[key] = value
			}
		}
	}
}
