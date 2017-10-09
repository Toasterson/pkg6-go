package repo

import (
	"github.com/toasterson/pkg6-go/packageinfo"
	"os"
	"github.com/toasterson/pkg6-go/catalog"
	"strings"
	"fmt"
)

type Repository interface {
	GetFile(publisher string, hash string) (*os.File, error)
	GetPackage(fmri string) (packageinfo.PackageInfo, error)
	GetPath() string
	GetPublishers() []string
	GetPackageFMRIs(publisher string, partial bool) []string
	Create() error
	Load() error
	Destroy() error
	Upgrade() error
	GetVersion() int
	GetCatalog(publisher string) catalog.Catalog
	AddPackage(info packageinfo.PackageInfo) error
}


func NewRepo(url string) (Repository, error){
	switch {
	case strings.HasPrefix(url, "file://"):
		return &FileRepo{Path: strings.Replace(url, "file://", "", -1)}, nil
	case strings.HasPrefix(url, "http://"):
		fallthrough
	case strings.HasPrefix(url, "https://"):
		return nil, fmt.Errorf("not implemented")
	default:
		return nil, fmt.Errorf("can not create repo object invalid url")
	}
}