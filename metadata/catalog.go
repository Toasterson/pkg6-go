package metadata

type Catalog struct {
	Created             string                 `json:"created"`
	LastModified        string                 `json:"last-modified"`
	PackageCount        int                    `json:"package-count"`
	PackageVersionCount int                    `json:"package-version-count"`
	Parts               map[string]CatalogPart `json:"parts"`
	Updates             map[string]CatalogPart `json:"updates"`
	Version             int                    `json:"version"`
	Packages            map[string]PackageInfo `json:"-"`
	Signature           Signature              `json:"_SIGNATURE"`
}

type CatalogPart struct {
	LastModified string `json:"last-modified"`
	Signature    string `json:"signature-sha-1"`
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
