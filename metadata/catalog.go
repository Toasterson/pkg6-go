package metadata

import (
	"github.com/toasterson/pkg6-go/pkg"
	"errors"
	"fmt"
)

type Catalog struct {
	Created             string `json:"created"`
	LastModified        string `json:"last-modified"`
	PackageCount        int `json:"package-count"`
	PackageVersionCount int `json:"package-version-count"`
	Parts               map[string]map[string]string `json:"parts"`
	Updates             map[string]map[string]string `json:"updates"`
	Version             int `json:"version"`
	Packages            map[string][]pkg.PackageInfo
	Signature           Signature `json:"_SIGNATURE"`
}

type Signature struct {
	SHA1 string `json:"sha-1"`
}

func (s *Signature) Check(signature string) error{
	if s.SHA1 == signature {
		return nil
	}
	return errors.New(fmt.Sprintf("Signature check failed %s != %s", s.SHA1, signature))
}

