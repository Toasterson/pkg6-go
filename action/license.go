package action

import (
	"strings"
	"strconv"
)

type LicenseAction struct {
	Type    string `json:"type,-"`
	Sha1    string `json:"sha1"`
	Chash   string `json:"chash"`
	License string `json:"license"`
	Csize   int `json:"csize"`
	Size    int `json:"size"`
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
