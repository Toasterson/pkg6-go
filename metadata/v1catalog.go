package metadata

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
)

type V1Catalog struct {
	Location            string                        `json:"-"`
	Created             string                        `json:"created"`
	LastModified        string                        `json:"last-modified"`
	PackageCount        int                           `json:"package-count"`
	PackageVersionCount int                           `json:"package-version-count"`
	Parts               map[string]CatalogPart        `json:"parts"`
	V1PartContent       map[string]*V1CatalogPartFile `json:"-"`
	Updates             map[string]CatalogPart        `json:"updates"`
	Version             int                           `json:"version"`
	Packages            map[string][]PackageInfo      `json:"-"`
	Signature           Signature                     `json:"_SIGNATURE"`
}

type CatalogPart struct {
	LastModified string `json:"last-modified"`
	Signature    string `json:"signature-sha-1"`
}

func NewV1Catalog(location string) *V1Catalog {
	return &V1Catalog{
		Location:      location,
		Parts:         make(map[string]CatalogPart),
		V1PartContent: make(map[string]*V1CatalogPartFile),
		Updates:       make(map[string]CatalogPart),
		Packages:      make(map[string][]PackageInfo),
	}
}

func (c *V1Catalog) Load() error {
	var content []byte
	var err error
	if content, err = ioutil.ReadFile(filepath.Join(c.Location, "catalog.attrs")); err != nil {
		return err
	}
	if err = json.Unmarshal(content, c); err != nil {
		return err
	}
	for part := range c.Parts {
		var partContent []byte
		if partContent, err = ioutil.ReadFile(filepath.Join(c.Location, part)); err != nil {
			return err
		}
		if c.V1PartContent[part], err = DeSerializeV1Part(partContent); err != nil {
			return err
		}
	}
	return nil
}

func (c *V1Catalog) Save() error {
	var serialized []byte
	var err error
	if serialized, err = json.Marshal(c); err != nil {
		return err
	}
	if err = ioutil.WriteFile(filepath.Join(c.Location, "catalog.attrs"), serialized, 0644); err != nil {
		return err
	}
	for part := range c.Parts {
		if serialized, err = c.SerializeV1Part(part); err != nil {
			return err
		}
		if err = ioutil.WriteFile(filepath.Join(c.Location, part), serialized, 0644); err != nil {
			return err
		}
	}
	return nil
}

func (c *V1Catalog) SerializeV1Part(part string) ([]byte, error) {
	if blob, err := json.Marshal(c.V1PartContent[part]); err != nil {
		return nil, err
	} else {
		return blob, nil
	}
}

func DeSerializeV1Part(partBlob []byte) (file *V1CatalogPartFile, err error) {
	partFile := &V1CatalogPartFile{}
	if err = json.Unmarshal(partBlob, &partFile); err != nil {
		return nil, err
	}
	return partFile, nil
}

func (c *V1Catalog) AddPackage(pkg *PackageInfo) (err error) {
	return
}

func (c *V1Catalog) UpdatePackage(pkg *PackageInfo) (err error) {
	return
}

func (c *V1Catalog) RemovePackage(pkg *PackageInfo) (err error) {
	return
}

func (c *V1Catalog) HasPackage(fmri string) bool {
	return false
}

func (c *V1Catalog) GetPackage(fmri string) (pkg PackageInfo, err error) {
	return
}

func (c *V1Catalog) GetPackages(fmris []string) (pkgs []PackageInfo, err error) {
	for _, fmri := range fmris {
		pkg, err := c.GetPackage(fmri)
		if err != nil {
			return nil, err
		}
		pkgs = append(pkgs, pkg)
	}
	return
}
