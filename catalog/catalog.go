package catalog

import (
	"github.com/toasterson/pkg6-go/pkg"
	"os"
	"encoding/json"
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

func (c *Catalog) LoadFromV1(location string)(err error){
	fd, err := os.Open(location+"/catalog.attrs")
	defer func() {
		err = fd.Close()
	}()
	if err != nil {
		return err
	}
	decoder := json.NewDecoder(fd)
	err = decoder.Decode(&c)
	if err != nil {
		return err
	}
	for key, value := range c.Parts {
		c.loadV1Part(location+"/"+key, Signature{value["signature-sha-1"]})
	}
	return
}

func (c *Catalog) loadV1Part(location string, signature Signature)(err error){
	fd, err := os.Open(location)
	defer func() {
		err = fd.Close()
	}()
	if err != nil {
		return
	}
	decoder := json.NewDecoder(fd)
	var anything map[string]map[string]interface{}
	err = decoder.Decode(&anything)
	if err != nil {
		return err
	}
	if c.Packages == nil {
		c.Packages = map[string][]pkg.PackageInfo{}
	}
	for k := range anything {
		switch k {
		case "_SIGNATURE": {
			err = signature.Check(anything[k]["sha-1"].(string))
			if err != nil {
				panic(fmt.Errorf("signature check failed: %s", err))
			}
		}
		default : {
			for packname, tmp := range anything[k]{
				rawPackages := tmp.([]interface{})
				var packArr []pkg.PackageInfo
				alreadypresent := false
				if val, ok := c.Packages[packname]; ok {
					packArr = val
					alreadypresent = true
				}
				for pack := range rawPackages {
					rawPack := rawPackages[pack].(map[string]interface{})
					var thePackage pkg.PackageInfo
					thePackage.FromMap(rawPack)
					thePackage.Name = packname
					if alreadypresent {
						for _, p := range packArr {
							if thePackage.CompareVersion(p) == "equals" {
								thePackage.Merge(&p)
							}
						}
					} else {
						packArr = append(packArr, thePackage)
					}
				}
				c.Packages[packname] = packArr
			}
		}
		}
	}
	return
}