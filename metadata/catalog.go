package metadata

type PackageCatalog interface {
	Load() error
	Save() error
	AddPackage(pkg *PackageInfo) (err error)
	UpdatePackage(pkg *PackageInfo) (err error)
	RemovePackage(pkg *PackageInfo) (err error)
	HasPackage(fmri string) bool
	GetPackage(fmri string) (pkg PackageInfo, err error)
	GetPackages(fmris []string) (pkgs []PackageInfo, err error)
}
