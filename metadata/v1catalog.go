package metadata

import (
	"os"
	"encoding/json"
	"github.com/toasterson/pkg6-go/pkg"
	"fmt"
)

func LoadCatalogV1(location string, catalog *Catalog)(err error){
	fd, err := os.Open(location+"/catalog.attrs")
	defer func() {
		err = fd.Close()
	}()
	if err != nil {
		return err
	}
	decoder := json.NewDecoder(fd)
	err = decoder.Decode(&catalog)
	if err != nil {
		return err
	}
	for key, value := range catalog.Parts {
		loadV1Part(location+"/"+key, catalog, Signature{value["signature-sha-1"]})
	}
	return
}

func loadV1Part(location string, catalog *Catalog, signature Signature)(err error){
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
	if catalog.Packages == nil {
		catalog.Packages = map[string][]pkg.PackageInfo{}
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
				if val, ok := catalog.Packages[packname]; ok {
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
				catalog.Packages[packname] = packArr
			}
		}
		}
	}
	return
}