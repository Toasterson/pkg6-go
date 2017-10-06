package action

import (
	"strings"
	"fmt"
	"strconv"
)

type AttributeAction struct {
	Type      string `json:"type,-"`
	Name      string `json:"name"`
	Values    []string `json:"values"`
	Optionals map[string]string `json:"optionals,omitempty"`
}

type DependAction struct {
	ActionType string `json:"action_type,-"`
	FMRI       string `json:"fmri"`
	Type       string `json:"type"`
	Predicate  string `json:"predicate"`
	Optional   map[string]string `json:"optional"`
}

type DirectoryAction struct {
	Type   string `json:"type,-"`
	Path   string `json:"path"`
	Mode   string `json:"mode"`
	Owner  string `json:"owner"`
	Group  string `json:"group"`
	Facets map[string]string `json:"facets"`
}

type FileAction struct {
	Type         string            `json:"type,-"`
	Sha1         string            `json:"sha1"`
	Path         string            `json:"path"`
	Size         int               `json:"size"`
	Csize        int               `json:"csize"`
	Chash        string            `json:"chash"`
	Owner        string            `json:"owner"`
	Group        string            `json:"group"`
	Mode         string            `json:"mode"`
	Preserve     bool              `json:"preserve"`
	Overlay      bool              `json:"overlay"`
	OriginalName string            `json:"original_name"`
	ReleaseNote  string            `json:"release_note"`
	RevertTag    string            `json:"revert_tag"`
	Elfarch      string            `json:"elfarch"`
	Elfbits      string            `json:"elfbits"`
	Elfhash      string            `json:"elfhash"`
	Attributes   map[string]string `json:"attributes"`
}

type LinkAction struct {
	Type   string `json:"type,-"`
	Path   string `json:"path"`
	Target string `json:"target"`
}

type LicenseAction struct {
	Type    string `json:"type,-"`
	Sha1    string `json:"sha1"`
	Chash   string `json:"chash"`
	License string `json:"license"`
	Csize   int `json:"csize"`
	Size    int `json:"size"`
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

func (a *FileAction) FromActionString(actionString string) {
	a.Attributes = make(map[string]string)
	for _, value := range tokenize(actionString) {
		equalpos := strings.Index(value, "=")
		key := value[0: equalpos]
		value = value[equalpos+1:]
		switch key {
		case "key":
			a.Sha1 = value
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
		case "sha1":
			a.Sha1 = value
		case "size":
		case "pkg.size":
			i, err := strconv.Atoi(value)
			if err == nil {
				a.Size = i
			}
		case "csize":
		case "pkg.csize":
			i, err := strconv.Atoi(value)
			if err == nil {
				a.Csize = i
			}
		case "chash":
			a.Chash = value
		case "preserve":
			if value == "true" {
				a.Preserve = true
			} else {
				a.Preserve = false
			}
		case "overlay":
			if value == "true" {
				a.Overlay = true
			} else {
				a.Overlay = false
			}
		case "original_name":
			a.OriginalName = value
		case "release_note":
			a.ReleaseNote = value
		case "revert_tag":
			a.RevertTag = value
		case "elfarch":
			a.Elfarch = value
		case "elfbits":
			a.Elfbits = value
		case "elfhash":
			a.Elfhash = value
		default:
			a.Attributes[key] = value
		}
	}
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

func (a *LicenseAction) FromActionString(actionString string) {
	for _, value := range tokenize(actionString) {
		equalpos := strings.Index(value, "=")
		key := value[0: equalpos]
		value = value[equalpos+1:]
		switch key {
		case "key":
			a.Sha1 = value
		case "action_type":
			a.Type = value
		case "sha1":
			a.Sha1 = value
		case "size":
		case "pkg.size":
			i, err := strconv.Atoi(value)
			if err == nil {
				a.Size = i
			}
		case "csize":
		case "pkg.csize":
			i, err := strconv.Atoi(value)
			if err == nil {
				a.Csize = i
			}
		case "chash":
			a.Chash = value
		case "license":
			a.License = value
		default:
		}
	}
}

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
