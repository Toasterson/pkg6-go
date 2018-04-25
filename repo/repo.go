package repo

import (
	"fmt"
	"github.com/toasterson/pkg6-go/metadata"
	"os"
	"strings"
)

type Repository interface {
	GetFile(publisher string, hash string) (*os.File, error)
	GetPackage(fmri string) (metadata.PackageInfo, error)
	GetPath() string
	GetPublishers() []string
	GetPackageFMRIs(publisher string, partial bool) []string
	Create() error
	Load() error
	Save() (err error)
	Destroy() error
	GetVersion() int
	GetCatalog(publisher string) *metadata.V1Catalog
	GetCatalogFile(publisher, part string) (*os.File, error)
	Search(params map[string]string, query string) string
	AddPackage(info metadata.PackageInfo) error
}

func NewRepo(url string) (Repository, error) {
	switch {
	case strings.HasPrefix(url, "file://"):
		if stat, err := os.Stat(strings.Replace(url, "file://", "", -1)); err != nil {
			return nil, err
		} else if stat.IsDir() {
			return &FileRepo{Path: strings.Replace(url, "file://", "", -1)}, nil
		} else if stat.Mode().IsRegular() {
			return nil, fmt.Errorf("p5p format not implemented yet")
		}
		return nil, fmt.Errorf("%s is nether dir, symlink or regular file", url)
	case strings.HasPrefix(url, "http://"):
		fallthrough
	case strings.HasPrefix(url, "https://"):
		return &HttpRepo{}, nil
	default:
		return nil, fmt.Errorf("can not create repo object invalid url")
	}
}
