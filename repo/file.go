package repo

import (
	"encoding/json"
	"fmt"
	"github.com/go-ini/ini"
	"github.com/toasterson/pkg6-go/metadata"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type FileRepo struct {
	Path                       string `json:"-"`
	Version                    int
	TrustAnchorDirectory       string
	CheckCertificateRevocation bool
	SignatureRequiredNames     []string
	Publishers                 []string                       `json:"-"`
	Catalogs                   map[string]*metadata.V1Catalog `json:"-"`
}

func (r *FileRepo) GetPath() string {
	return r.Path
}

func (r *FileRepo) Create() error {
	return nil
}

func (r *FileRepo) Load() error {
	r.Publishers = r.getAllPublishersFromDisk()
	inifile, err := ini.Load(r.GetPath() + "/pkg5.repository")
	if err != nil {
		return fmt.Errorf("can not load configuration %s: %s", r.Path, err)
	}
	repoCFG, _ := inifile.GetSection("repository")
	r.Version = repoCFG.Key("version").MustInt()
	r.CheckCertificateRevocation = repoCFG.Key("check-certificate-revocation").MustBool()
	r.TrustAnchorDirectory = repoCFG.Key("trust-anchor-directory").MustString("/etc/certs/CA/")
	//TODO Full Load of Config as described in Documentation
	r.Catalogs = make(map[string]*metadata.V1Catalog)
	for _, pub := range r.Publishers {
		catalogPath := filepath.Join(r.Path, "publisher", pub, "catalog")
		cat := metadata.NewV1Catalog(catalogPath)
		cat.Load()
		r.Catalogs[pub] = cat
	}
	return nil
}

func (r *FileRepo) Destroy() error {
	return os.RemoveAll(r.GetPath())
}

func (r *FileRepo) GetPackageFMRIs(publisher string, partial bool) []string {
	var FMRIS []string
	pkgPath := filepath.Join(r.Path, "publisher", publisher, "pkg")
	packages, _ := ioutil.ReadDir(pkgPath)
	for _, pkg := range packages {
		manifests, _ := ioutil.ReadDir(filepath.Join(pkgPath, pkg.Name()))
		for _, manifest := range manifests {
			if partial {
				FMRIS = append(FMRIS, "pkg:/"+metadata.Unicode2FMRI(pkg.Name()+"@"+manifest.Name()))
			} else {
				FMRIS = append(FMRIS, "pkg://"+publisher+"/"+metadata.Unicode2FMRI(pkg.Name()+"@"+manifest.Name()))
			}
		}
	}
	return FMRIS
}

func (r *FileRepo) GetPublishers() []string {
	return r.Publishers
}

func (r *FileRepo) getAllPublishersFromDisk() []string {
	var publishers []string
	files, _ := ioutil.ReadDir(filepath.Join(r.Path, "publisher"))
	for _, f := range files {
		if f.IsDir() {
			publishers = append(publishers, f.Name())
		}
	}
	return publishers
}

func (r *FileRepo) GetFile(publisher string, hash string) (*os.File, error) {
	path := filepath.Join(r.Path, "publisher", publisher, "file", hash[0:2], hash)
	file, err := os.OpenFile(path, os.O_RDONLY, 0666)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func (r *FileRepo) GetPackage(fmri string) (metadata.PackageInfo, error) {
	if !strings.Contains(fmri, "pkg://") {
		return metadata.PackageInfo{}, fmt.Errorf("package needs to be with publisher to retrieve from repo")
	}
	pkg := metadata.PackageInfo{}
	pkg.SetFmri(fmri)
	pkgPath := filepath.Join(r.Path, "publisher", pkg.Publisher, "pkg")
	switch r.Version {
	case 4:
		pkg.ReadManifest(pkgPath)
	case 5:
		pkg.Load(pkgPath)
	default:
		return pkg, fmt.Errorf("can not load Package from a Repository with version: %s", r.Version)
	}
	return pkg, nil
}

func (r *FileRepo) GetCatalog(publisher string) *metadata.V1Catalog {
	return r.Catalogs[publisher]
}

func (r *FileRepo) GetVersion() int {
	return r.Version
}

func (r *FileRepo) AddPackage(info metadata.PackageInfo) error {
	return nil
}

func (r *FileRepo) Search(params map[string]string, query string) string {
	return ""
}

func (r *FileRepo) Save() (err error) {
	file, err := os.Create(filepath.Join(r.Path, "repository.json"))
	defer func() {
		if cerr := file.Close(); err == nil {
			err = cerr
		}
	}()
	if err != nil {
		return err
	}
	if err := json.NewEncoder(file).Encode(r); err != nil {
		return err
	}
	return err
}
