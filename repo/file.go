package repository

import (
	"github.com/go-ini/ini"
	"github.com/toasterson/pkg6-go/util"
	"os"
	"github.com/toasterson/pkg6-go/packageinfo"
	"io/ioutil"
	"github.com/toasterson/pkg6-go/catalog"
	"strings"
)

type FileRepo struct {
	Path string
	Version int
	TrustAnchorDirectory string
	CheckCertificateRevocation bool
	SignatureRequiredNames []string
	Publishers []string
	Catalogs map[string]catalog.Catalog
}

func (r *FileRepo) GetPath() string {
	return r.Path
}

func (r *FileRepo) Create() error {
	return nil
}

func (r *FileRepo)Load() error {
	r.Publishers = r.getAllPublishersFromDisk()
	inifile, err := ini.Load(r.GetPath()+"/pkg5.repository")
	util.Error(err, "Error Loading FileRepo Configuration")
	repoCFG, _ := inifile.GetSection("repository")
	r.Version = repoCFG.Key("version").MustInt()
	r.CheckCertificateRevocation = repoCFG.Key("check-certificate-revocation").MustBool()
	r.TrustAnchorDirectory = repoCFG.Key("trust-anchor-directory").MustString("/etc/certs/CA/")
	//TODO Full Load of Config as described in Documentation
	r.Catalogs = make(map[string]catalog.Catalog)
	for _, pub := range r.Publishers {
		cat := catalog.Catalog{}
		cat.LoadFromV1(r.Path + "/publisher/"+pub+ "/catalog")
		r.Catalogs[pub] = cat
	}
	return nil
}

func (r *FileRepo) Destroy() error{
	return os.RemoveAll(r.GetPath())
}

func (r *FileRepo)Upgrade() error{
	for _, pub := range r.Publishers{
		pkgPath := r.Path + "/publisher/"+pub+"/pkg"
		for _, pkgFMRI := range r.GetPackageFMRIs(pub, false){
			pkg := packageinfo.FromFMRI(pkgFMRI)
			pkg.ReadManifest(pkgPath)
			if err := pkg.UpgradeFormat(pkgPath); err != nil {
				return err
			}

		}
		cat := r.Catalogs[pub]
		cat.Save(r.Path+"/publisher/"+pub+"/catalog")
	}
	return nil
}

func (r *FileRepo)GetPackageFMRIs(publisher string, partial bool) []string{
	var FMRIS = []string{}
	pkgPath := r.Path+"/publisher/"+publisher+"/pkg"
	packages, _ := ioutil.ReadDir(pkgPath)
	for _, pkg := range packages{
		manifests, _ := ioutil.ReadDir(pkgPath+"/"+pkg.Name())
		for _, manifest := range manifests{
			if partial {
				FMRIS = append(FMRIS, "pkg:/"+packageinfo.Unicode2FMRI(pkg.Name()+"@"+manifest.Name()))
			} else {
				FMRIS = append(FMRIS, "pkg://"+publisher+"/"+packageinfo.Unicode2FMRI(pkg.Name()+"@"+manifest.Name()))
			}
		}
	}
	return FMRIS
}

func (r *FileRepo) GetPublishers() []string {
	return r.Publishers
}

func (r *FileRepo)getAllPublishersFromDisk() []string{
	var publishers = []string{}
	files, _ := ioutil.ReadDir(r.Path+"/publisher")
	for _, f := range files {
		if f.IsDir() {
			publishers = append(publishers, f.Name())
		}
	}
	return publishers
}

func (r *FileRepo)GetFile(publisher string, hash string) *os.File{
	file, err := os.OpenFile(r.Path+"/publisher/"+publisher+"/file/"+hash[0:2]+"/"+hash, os.O_RDONLY, 0666)
	if err != nil{
		return nil
	}
	return file
}

func (r *FileRepo)GetPackageInfo(fmri string) packageinfo.PackageInfo {
	//TODO Return Package Matching of Highest Priority Publisher if Asked
	if !strings.Contains(fmri, "pkg://"){
		return packageinfo.PackageInfo{}
	}
	pkg := packageinfo.PackageInfo{}
	pkg.SetFmri(fmri)
	pkg.Load(r.Path+"/publisher/"+pkg.Publisher+"/pkg")
	return pkg
}

func (r *FileRepo)GetCatalog(publisher string) catalog.Catalog{
	return r.Catalogs[publisher]
}