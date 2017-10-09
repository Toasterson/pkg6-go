package catalog

import (
	"github.com/toasterson/pkg6-go/packageinfo"
	"github.com/toasterson/pkg6-go/util"
	"os"
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
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

func (c *Catalog) Load(location string) error{
	if _, err := os.Stat(location+"/catalog.attrs"); os.IsExist(err) {
		c.LoadFromV1(location)
	} else {
		file, ferr := os.OpenFile(location+"/catalog.json", os.O_RDONLY, 0666)
		if ferr != nil{
			return ferr
		}
		defer file.Close()
		var b = []byte{}
		if _, rerr := file.Read(b); rerr != nil {
			return rerr
		}
		if merr := json.Unmarshal(b, c); merr != nil {
			return merr
		}
	}
	return nil
}

func (c *Catalog) Save(location string) (err error) {
	b, err := json.Marshal(c)
	if err != nil{
		return err
	}
	path := location+"/catalog.json"
	file, ferr := os.OpenFile(path, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0666 )
	if ferr != nil {
		return ferr
	}
	defer func() {
		if cerr := file.Close(); err == nil{
			err = cerr
		}
	}()
	if _, werr := file.Write(b); werr != nil{
		return werr
	}
	return nil
}

func (c Catalog) Upgrade(location string) error {
	if err := c.Save(location); err != nil {
		return err
	}
	files := []string{"catalog.attrs", "catalog.base.C", "catalog.dependency.C", "catalog.summary.C"}
	if updates, err := filepath.Glob("update.*.C"); err != nil{
		files = append(files, updates...)
	}
	for _, f := range files {
		if err := os.Remove(filepath.Join(location, f)); err != nil{
			return err
		}
	}
	return nil
}
