package metadata

import (
	"errors"
	"fmt"
)

type Catalog struct {
	Created             string                       `json:"created"`
	LastModified        string                       `json:"last-modified"`
	PackageCount        int                          `json:"package-count"`
	PackageVersionCount int                          `json:"package-version-count"`
	Parts               map[string]map[string]string `json:"parts"`
	Updates             map[string]map[string]string `json:"updates"`
	Version             int                          `json:"version"`
	Packages            map[string]PackageInfo
	Signature           Signature `json:"_SIGNATURE"`
}

type Signature struct {
	SHA1 string `json:"sha-1"`
}

func (s *Signature) Check(signature string) error {
	if s.SHA1 == signature {
		return nil
	}
	return errors.New(fmt.Sprintf("Signature check failed %s != %s", s.SHA1, signature))
}

func (c Catalog) AddPackage(pkg *PackageInfo) (err error) {
	return
}

func (c Catalog) UpdatePackage(pkg *PackageInfo) (err error) {
	return
}

func (c Catalog) RemovePackage(pkg *PackageInfo) (err error) {
	return
}

func (c Catalog) HasPackage(fmri string) bool {
	return false
}

func (c Catalog) GetPackage(fmri string) (pkg PackageInfo, err error) {
	return
}

func (c Catalog) GetPackages(fmris []string) (pkgs []PackageInfo, err error) {
	for _, fmri := range fmris {
		pkg, err := c.GetPackage(fmri)
		if err != nil {
			return
		}
		pkgs = append(pkgs, pkg)
	}
	return
}
