package catalog

import (
	"encoding/json"
	"github.com/toasterson/pkg6-go/packageinfo"
	"io/ioutil"
	"os"
	"path/filepath"
)

type V1Catalog struct {
	Location            string                               `json:"-"`
	Created             string                               `json:"created"`
	LastModified        string                               `json:"last-modified"`
	PackageCount        int                                  `json:"package-count"`
	PackageVersionCount int                                  `json:"package-version-count"`
	Parts               map[string]CatalogPart               `json:"parts"`
	V1PartContent       map[string]V1CatalogPartFile         `json:"-"`
	Updates             map[string]CatalogPart               `json:"updates"`
	Version             int                                  `json:"version"`
	Packages            map[string][]packageinfo.PackageInfo `json:"-"`
	Signature           Signature                            `json:"_SIGNATURE"`
}

type CatalogPart struct {
	LastModified string `json:"last-modified"`
	Signature    string `json:"signature-sha-1"`
}

type V1CatalogPartFile struct {
	Publishers map[string]V1Publisher
	Signature  Signature `json:"_SIGNATURE"`
}

func (f *V1CatalogPartFile) UnmarshalJSON(blob []byte) error {
	raw := map[string]*json.RawMessage{}
	if err := json.Unmarshal(blob, &raw); err != nil {
		return err
	}
	if len(f.Publishers) == 0 {
		f.Publishers = make(map[string]V1Publisher)
	}
	for key, val := range raw {
		if key == "_SIGNATURE" {
			sig := Signature{}
			if err := json.Unmarshal(*val, &sig); err != nil {
				return err
			}
			f.Signature = sig
		} else {
			pub := V1Publisher{}
			if err := json.Unmarshal(*val, &pub); err != nil {
				return err
			}
			f.Publishers[key] = pub
		}
	}
	return nil
}

type V1Publisher struct {
	Packages map[string]V1Packages
}

func (pub *V1Publisher) UnmarshalJSON(blob []byte) error {
	pub.Packages = make(map[string]V1Packages)
	raw := make(map[string]*json.RawMessage)
	if err := json.Unmarshal(blob, &raw); err != nil {
		return err
	}
	for key, value := range raw {
		pak := V1Packages{}
		if err := json.Unmarshal(*value, &pak); err != nil {
			return err
		}
		pub.Packages[key] = pak
	}
	return nil
}

type V1Packages []V1PackageVersion

type V1PackageVersion struct {
	Actions   V1Actions `json:"actions"`
	Signature string    `json:"signature-sha-1"`
	Version   string    `json:"version"`
}

type V1Actions []string

func NewV1Catalog(location string) *V1Catalog {
	return &V1Catalog{
		Location:      location,
		Parts:         make(map[string]CatalogPart),
		V1PartContent: make(map[string]V1CatalogPartFile),
		Updates:       make(map[string]CatalogPart),
		Packages:      make(map[string][]packageinfo.PackageInfo),
	}
}

func (c *V1Catalog) SerializeToV1() (map[string][]byte, error) {
	retVal := make(map[string][]byte)
	if attrFile, err := json.Marshal(c); err != nil {
		return nil, err
	} else {
		retVal["catalog.attrs"] = attrFile
	}
	for p := range c.V1PartContent {
		if blob, err := c.SerializeV1Part(p); err != nil {
			return nil, err
		} else {
			retVal[p] = blob
		}
	}
	return retVal, nil
}

func (c *V1Catalog) SerializeV1Part(part string) ([]byte, error) {
	if blob, err := json.Marshal(c.V1PartContent[part]); err != nil {
		return nil, err
	} else {
		return blob, nil
	}
}

func (c *V1Catalog) LoadFromV1() error {
	var err error
	var fd *os.File
	if fd, err = os.Open(c.Location + "/catalog.attrs"); err != nil {
		return err
	}
	defer fd.Close()
	decoder := json.NewDecoder(fd)
	if err = decoder.Decode(&c); err != nil {
		return err
	}
	for key, value := range c.Parts {
		partFile := V1CatalogPartFile{}
		path := filepath.Join(c.Location, key)
		var blob []byte
		if blob, err = ioutil.ReadFile(path); err != nil {
			return err
		}
		if err = json.Unmarshal(blob, &partFile); err != nil {
			return err
		}
		if err = partFile.Signature.Check(value.Signature); err != nil {
			return err
		}
		c.V1PartContent[key] = partFile
	}
	return nil
}

func (c *V1Catalog) Load() error {
	return c.LoadFromV1()
}
