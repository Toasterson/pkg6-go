package repository

import (
	"github.com/toasterson/pkg6-go/packageinfo"
	"os"
	"github.com/toasterson/pkg6-go/catalog"
)

type Repository interface {
	GetFile(publisher string, hash string) *os.File
	GetPackageInfo(fmri string) packageinfo.PackageInfo
	GetPath() string
	GetPublishers() []string
	GetPackageFMRIs(publisher string, partial bool) []string
	Create() error
	Load() error
	Destroy() error
	Upgrade() error
	GetCatalog(publisher string) catalog.Catalog
}
