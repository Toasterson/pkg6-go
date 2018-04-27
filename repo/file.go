package repo

import (
	"fmt"
	"github.com/go-ini/ini"
	"github.com/toasterson/pkg6-go/metadata"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type FileRepo struct {
	Path                       string `json:"-"`
	IsMirror                   bool   `json:"-"`
	Version                    int
	TrustAnchorDirectory       string
	CheckCertificateRevocation bool
	SignatureRequiredNames     []string
	Publishers                 []string `json:"-"`
}

// This is a helper method to add  a Catalog to the Repository in the right location
// only works if FileRepo is the Backing Store of a Mirror
func (r *FileRepo) AddCatalogFile(publisher, part string, file *os.File) error {
	if !r.IsMirror {
		return fmt.Errorf("only mirrors can add catalog files directly use AddPackage instead")
	}
	path := filepath.Join(r.Path, "publisher", publisher, "catalog", part)
	t, err := os.Create(path)
	if err != nil {
		return err
	}
	defer t.Close()
	_, err = io.Copy(t, file)
	return err
}

// Add a Manifest to the Repo.
// Helper function for Mirror Operations
func (r *FileRepo) AddPackageManifest(fmri string, file *os.File) error {
	var publisher, name, version string
	ok := true
	fmriMap := metadata.SplitFmri(fmri)
	if publisher, ok = fmriMap["publisher"]; !ok {
		return fmt.Errorf("no publisher in frmi: %s", fmri)
	}
	if name, ok = fmriMap["name"]; !ok {
		return fmt.Errorf("no name in fmri: %s", fmri)
	}
	if version, ok = fmriMap["version"]; !ok {
		return fmt.Errorf("no version component in fmri: %s", fmri)
	}
	version += "," + fmriMap["build_release"] + "-" + fmriMap["branch"] + ":" + fmriMap["packaging_date"]
	folderName := metadata.FMRI2Unicode(name)
	fileName := metadata.FMRI2Unicode(version)
	path := filepath.Join(r.Path, "publisher", publisher, "pkg", folderName, fileName)
	t, err := os.Create(path)
	if err != nil {
		return err
	}
	defer t.Close()
	if _, err = io.Copy(t, file); err != nil {
		return err
	}
	return nil
}

// Helper Function for Mirror Operations adds the a File to the correct Path
func (r *FileRepo) AddFile(publisher, hash string, file *os.File) error {
	path := filepath.Join(r.Path, "publisher", publisher, "file", hash[0:2], hash)
	t, err := os.Create(path)
	if err != nil {
		return err
	}
	defer t.Close()
	if _, err = io.Copy(t, file); err != nil {
		return err
	}
	return nil
}

func (r *FileRepo) HasPublisher(publisher string) bool {
	for _, pub := range r.Publishers {
		if pub == publisher {
			return true
		}
	}
	return false
}

func (r *FileRepo) Exists() bool {
	if _, err := os.Stat(r.Path); os.IsNotExist(err) {
		return false
	}
	return true
}

func (r *FileRepo) GetPath() string {
	return r.Path
}

func (r *FileRepo) Create() error {
	if r.Exists() {
		return fmt.Errorf("repository at %s already exists", r.Path)
	}
	if err := os.MkdirAll(r.Path, 0755); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Join(r.Path, "publisher"), 0755); err != nil {
		return err
	}
	for _, pub := range r.Publishers {
		pubPath := filepath.Join(r.Path, "publisher", pub)
		if err := os.MkdirAll(pubPath, 0755); err != nil {
			return err
		}
		if err := os.MkdirAll(filepath.Join(pubPath, "catalog"), 0755); err != nil {
			return err
		}
		if err := os.MkdirAll(filepath.Join(pubPath, "file"), 0755); err != nil {
			return err
		}
		if err := os.MkdirAll(filepath.Join(pubPath, "pkg"), 0755); err != nil {
			return err
		}
		if err := os.MkdirAll(filepath.Join(pubPath, "tmp"), 0755); err != nil {
			return err
		}
	}
	return nil
}

func (r *FileRepo) Load() error {
	r.Publishers = r.getAllPublishersFromDisk()
	inifile, err := ini.Load(r.GetPath() + "/pkg5.repository")
	if err != nil {
		return fmt.Errorf("can not load configuration %s: %s", r.Path, err)
	}
	if pubCFG, err := inifile.GetSection("publisher"); err != nil {
		return err
	} else {
		pub := pubCFG.Key("prefix").MustString("")
		if !r.HasPublisher(pub) {
			r.Publishers = append(r.Publishers, pub)
		}
	}
	repoCFG, _ := inifile.GetSection("repository")
	r.Version = repoCFG.Key("version").MustInt()
	r.CheckCertificateRevocation = repoCFG.Key("check-certificate-revocation").MustBool()
	r.TrustAnchorDirectory = repoCFG.Key("trust-anchor-directory").MustString("/etc/certs/CA/")
	//TODO Full Load of Config as described in Documentation
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
	if len(r.Publishers) == 0 {
		return r.getAllPublishersFromDisk()
	}
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
	return os.Open(filepath.Join(r.Path, "publisher", publisher, "file", hash[0:2], hash))
}

func (r *FileRepo) GetPackage(fmri string) (*metadata.PackageInfo, error) {
	if !strings.Contains(fmri, "pkg://") {
		return nil, fmt.Errorf("package needs to be with publisher to retrieve from repo")
	}
	pkg := &metadata.PackageInfo{}
	if err := pkg.SetFmri(fmri); err != nil {
		return nil, err
	}
	pkgPath := filepath.Join(r.Path, "publisher", pkg.Publisher, "pkg")
	switch r.GetVersion() {
	case 4:
		pkg.ReadManifest(pkgPath)
	case 5:
		pkg.Load(pkgPath)
	default:
		return pkg, fmt.Errorf("can not load Package from a Repository with version: %d", r.Version)
	}
	return pkg, nil
}

func (r *FileRepo) GetCatalog(publisher string) (*metadata.V1Catalog, error) {
	catalogPath := filepath.Join(r.Path, "publisher", publisher, "catalog")
	cat := metadata.NewV1Catalog(catalogPath)
	if err := cat.Load(); err != nil {
		return nil, err
	}
	return cat, nil
}

func (r *FileRepo) GetCatalogFile(publisher, part string) (*os.File, error) {
	return os.Open(filepath.Join(r.Path, "publisher", publisher, "catalog", part))
}

func (r *FileRepo) GetVersion() int {
	return 4
}

func (r *FileRepo) AddPackage(info metadata.PackageInfo) error {
	if r.IsMirror {
		return fmt.Errorf("repository is a local Component of a mirror and thus readonly")
	}
	panic("implement me")
}

func (r *FileRepo) Search(params map[string]string, query string) string {
	return ""
}

func (r *FileRepo) Save() (err error) {
	if r.IsMirror {
		return fmt.Errorf("repository is a local Component of a mirror and thus readonly")
	}
	panic("implement me")
}
