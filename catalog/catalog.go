package catalog

import (
	"github.com/toasterson/pkg6-go/packageinfo"
	"github.com/toasterson/pkg6-go/util"
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
	Packages            map[string][]packageinfo.PackageInfo
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

func (c *Catalog) LoadFromV1(location string){
	fd, err := os.Open(location+"/catalog.attrs")
	defer fd.Close()
	util.Error(err, "Opening Catalog")
	decoder := json.NewDecoder(fd)
	err2 := decoder.Decode(&c)
	util.Panic(err2, "Decoding Catalog")
	for key, value := range c.Parts {
		c.loadV1Part(location+"/"+key, Signature{value["signature-sha-1"]})
	}
}

func (c *Catalog) loadV1Part(location string, signature Signature){
	fd, err := os.Open(location)
	defer fd.Close()
	util.Error(err, "Opening Catalog Part")
	decoder := json.NewDecoder(fd)
	var anything map[string]map[string]interface{}
	err = decoder.Decode(&anything)
	util.Error(err, "Decoding Catalog Part")
	if c.Packages == nil {
		c.Packages = map[string][]packageinfo.PackageInfo{}
	}
	for k := range anything {
		switch k {
		case "_SIGNATURE": {
			util.Panic(signature.Check(anything[k]["sha-1"].(string)), "Checking Signature")
		}
		default : {
			for packname, tmp := range anything[k]{
				rawPackages := tmp.([]interface{})
				var packArr []packageinfo.PackageInfo
				alreadypresent := false
				if val, ok := c.Packages[packname]; ok {
					packArr = val
					alreadypresent = true
				}
				for pack := range rawPackages {
					rawPack := rawPackages[pack].(map[string]interface{})
					var thePackage packageinfo.PackageInfo
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
}