package catalog

type PackageCatalog interface {
	Load() error
	Save() error
	Patch(update PackageCatalog) error
}
