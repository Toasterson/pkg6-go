package metadata

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type catalogV1Part map[string]map[string]interface{}

func LoadCatalogV1(location string, catalog *Catalog) (err error) {
	fd, err := os.Open(location + "/catalog.attrs")
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
		loadV1Part(filepath.Join(location, key), catalog, Signature{SHA1: value.Signature})
	}
	catalog.Parts = nil
	for key, value := range catalog.Updates {
		loadV1Part(filepath.Join(location, key), catalog, Signature{SHA1: value.Signature})
	}
	catalog.Updates = nil
	return
}

func loadV1Part(location string, catalog *Catalog, signature Signature) (err error) {
	fd, err := os.Open(location)
	defer func() {
		err = fd.Close()
	}()
	if err != nil {
		return
	}
	decoder := json.NewDecoder(fd)
	var v1Parts catalogV1Part
	err = decoder.Decode(&v1Parts)
	if err != nil {
		return err
	}
	if catalog.Packages == nil {
		catalog.Packages = make(map[string]PackageInfo)
	}
	for k := range v1Parts {
		switch k {
		case "_SIGNATURE":
			{
				err = signature.Check(v1Parts[k]["sha-1"].(string))
				if err != nil {
					panic(fmt.Errorf("signature check failed: %s", err))
				}
			}
		default:
			{
				for packBase, packageVersionsRaw := range v1Parts[k] {
					packageVersions := packageVersionsRaw.([]interface{})
					for pack := range packageVersions {
						rawPack := packageVersions[pack].(map[string]interface{})
						var packageInfo PackageInfo
						packageInfo.FromMap(rawPack)
						packageInfo.Name = packBase
						FMRI := packageInfo.GetFmri()
						if _, ok := catalog.Packages[FMRI]; !ok {
							catalog.Packages[FMRI] = packageInfo
						}
					}
				}
			}
		}
	}
	return
}
