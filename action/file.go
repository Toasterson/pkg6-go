package action

import (
	"strconv"
	"strings"
)

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
